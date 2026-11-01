//go:build ignore

// monitor prints events from evdev devices. Run with privileges:
//
//	sudo go run ./examples/monitor                       # all devices
//	sudo go run ./examples/monitor /dev/input/eventX     # one or more nodes
//	sudo go run ./examples/monitor -grab /dev/input/eventX
//
// With no path arguments it monitors every readable device, which is the easiest
// way to discover which node carries a given input (e.g. a gaming mouse exposes
// several nodes — movement and clicks are on the EV_REL node, not its
// keyboard-emulation node). Each line is prefixed with the source device.
//
// With -grab each device is taken exclusively, so its events go only to this
// program. Be careful grabbing all devices: it captures the keyboard you would
// use to stop the program. Press Ctrl-C to stop.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	evdev "github.com/mikegio27/go-evdev"

	"github.com/stefan-muehlebach/gc9503cv/event"
	"github.com/stefan-muehlebach/gc9503cv/event/key"
	"github.com/stefan-muehlebach/gc9503cv/event/mouse"
)

type MouseKeyState struct {
	X, Y      float64
	Modifiers key.Modifiers
	Button    mouse.Button
}

var (
	state MouseKeyState
)

func main() {
	grab := flag.Bool("grab", false, "grab each device exclusively (EVIOCGRAB)")
	flag.Parse()

	queue := make(chan event.InputEvent, 64)

	devices, err := openDevices(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(devices) == 0 {
		fmt.Fprintln(os.Stderr, "no readable devices (try sudo, or pass a path)")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, d := range devices {
		if *grab {
			if err := d.Grab(); err != nil {
				fmt.Fprintln(os.Stderr, "grab:", err)
				d.Close()
				continue
			}
			defer d.Ungrab()
		}
		name, _ := d.Name()
		path := d.Path()
		log.Printf("%s (%q)\n", path, name)
		log.Printf("Supported Types and Codes:")
		types, _ := d.CapableTypes()
		for _, typ := range types {
			log.Printf("  %v", typ)
			codes, _ := d.CapableCodes(typ)
			for _, code := range codes {
				log.Printf("    : %s", evdev.CodeName(typ, code))
			}
		}

		wg.Add(1)
		go func(d *evdev.Device) {
			defer wg.Done()
			defer d.Close()
			enqueue(d, queue)
			//monitor(d)
		}(d)
	}
	fmt.Println("— press Ctrl-C to stop —")
	queueMonitor(queue)
	wg.Wait()
}

// openDevices opens the given paths, or every readable device when paths is empty.
func openDevices(paths []string) ([]*evdev.Device, error) {
	if len(paths) == 0 {
		return evdev.ListDevices()
	}
	var devices []*evdev.Device
	for _, p := range paths {
		d, err := evdev.Open(p)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", p, err)
		}
		devices = append(devices, d)
	}
	return devices, nil
}

func enqueue(d *evdev.Device, queue chan event.InputEvent) {
	for {
		ev, err := d.ReadOne()
		if err != nil {
			log.Fatalf("reading on %s: %v", d.Path(), err)
		}
		switch ev.Type {
		// All event types we know of, but we don't support them
		case evdev.EV_SYN, evdev.EV_MSC:
			//
		// All events that have something to do with moving things (e.g. mouse)
		case evdev.EV_REL:
			event := mouse.Event{}
			switch ev.Code {
			case evdev.REL_X:
				event.X = float64(ev.Value)
			case evdev.REL_Y:
				event.Y = float64(ev.Value)
			default:
				log.Printf("EV_REL: unsupported event code: %(d) (%+v)",
					ev.Code, ev)
				continue
			}
			event.Button = state.Button
			event.Modifiers = state.Modifiers
			event.Direction = mouse.DirNone
			queue <- event
		// All events that have something to do with keys or buttons
		case evdev.EV_KEY:
			switch ev.Code {
			case evdev.BTN_LEFT, evdev.BTN_RIGHT, evdev.BTN_MIDDLE:
				var btn mouse.Button

				switch ev.Code {
				case evdev.BTN_LEFT:
					btn = mouse.ButtonLeft
				case evdev.BTN_RIGHT:
					btn = mouse.ButtonRight
				case evdev.BTN_MIDDLE:
					btn = mouse.ButtonMiddle
				}
				event := mouse.Event{
					X: state.X, Y: state.Y,
					Button:    btn,
					Modifiers: state.Modifiers,
				}
				switch ev.Value {
				case 0:
					event.Direction = mouse.DirRelease
					state.Button = mouse.ButtonNone
				case 1:
					event.Direction = mouse.DirPress
					state.Button = btn
				}
				queue <- event
				continue
			case evdev.KEY_LEFTSHIFT, evdev.KEY_RIGHTSHIFT:
				switch ev.Value {
				case 0:
					state.Modifiers = state.Modifiers &^ key.ModShift
				case 1:
					state.Modifiers = state.Modifiers | key.ModShift
				}
				continue
			case evdev.KEY_LEFTCTRL, evdev.KEY_RIGHTCTRL:
				switch ev.Value {
				case 0:
					state.Modifiers = state.Modifiers &^ key.ModControl
				case 1:
					state.Modifiers = state.Modifiers | key.ModControl
				}
				continue
			case evdev.KEY_LEFTALT, evdev.KEY_RIGHTALT:
				switch ev.Value {
				case 0:
					state.Modifiers = state.Modifiers &^ key.ModAlt
				case 1:
					state.Modifiers = state.Modifiers | key.ModAlt
				}
				continue
			default:
				event := key.Event{}
				event.Code = key.Code(ev.Code)
				event.Modifiers = state.Modifiers
				switch ev.Value {
				case 0:
					event.Direction = key.DirRelease
				case 1:
					event.Direction = key.DirPress
				}
				queue <- event
			}
		// If there are other events, we'd like to know them in order to
		// handle them correctly.
		default:
			log.Printf("unsupported event type: %d (%+v)", ev.Type, ev)
			continue
		}
	}
}

func queueMonitor(queue chan event.InputEvent) {
	for ev := range queue {
		switch event := ev.(type) {
		case key.Event, mouse.Event:
			fmt.Printf("%+v\n", event)
		}
	}
}

// monitor reads events from one device until it errors or disappears, printing
// each prefixed with the device's node name.
func monitor(d *evdev.Device) {
	tag := filepath.Base(d.Path())
	for {
		ev, err := d.ReadOne()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Fprintf(os.Stderr, "%s: read: %v\n", tag, err)
			}
			return
		}
		fmt.Printf("%-8s %d.%06d  %s\n", tag, ev.Time.Sec, ev.Time.Usec, ev)
	}
}
