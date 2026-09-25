package main

import (
	"fmt"
	"math"
	"time"
)

const (
	maxLaps = 10
)

//----------------------------------------------------------------------------

type timerValues struct {
	sum, min, max time.Duration
}

func (t *timerValues) reset() {
	t.sum = 0
	t.min = time.Duration(math.MaxInt64)
	t.max = time.Duration(math.MinInt64)
}

func (t *timerValues) update(d time.Duration) {
	if d > t.max {
		t.max = d
	}
	if d < t.min {
		t.min = d
	}
	t.sum += d
}

//----------------------------------------------------------------------------

// Dieser Typ dient der Zeitmessung.
type Stopwatch struct {
	Name                string
	isRunning, isPaused bool
	n                   int
	t0, t1              time.Time
	all                 timerValues
	timers              [maxLaps]timerValues
	laps				[maxLaps]time.Duration
	lapNames			[maxLaps]string
	NumLaps             int
}

func NewStopwatch(name string) *Stopwatch {
	s := &Stopwatch{}
	s.Name = name
	s.Reset()
	return s
}

// Mit Start wird eine neue Messung begonnen.
func (s *Stopwatch) Start() {
	if s.isRunning {
		return
	}
	for i, _ := range s.laps {
		s.laps[i] = 0
	}
	s.isRunning = true
	s.NumLaps = 0
	s.t0 = time.Now()
}

func (s *Stopwatch) SetLapNames(names ...string) {
	for i, name := range names {
		s.lapNames[i] = name
	}
}

func (s *Stopwatch) LapName(id int) string {
	if id == 0 {
		return "Total"
	} else if id <= maxLaps {
		return s.lapNames[id-1]
	} else {
		return ""
	}
}

// Stellt eine neue Rundenzeit in die Messreihe.
func (s *Stopwatch) Lap() {
	if !s.isRunning || s.isPaused {
		return
	}
	s.laps[s.NumLaps] = time.Since(s.t0)
	s.NumLaps += 1
}

// Stop beendet die Messung und aktualisiert alle abhaengigen Werte.
func (s *Stopwatch) Stop() {
	var d time.Duration

	if !s.isRunning {
		return
	}
	if s.isPaused {
		d = s.t1.Sub(s.t0)
	} else {
		d = time.Since(s.t0)
	}

	s.isRunning = false
	s.isPaused = false
	s.n += 1

	s.all.update(d)

	if s.NumLaps > 0 {
		s.laps[s.NumLaps] = d
		s.NumLaps += 1

		for i := s.NumLaps-1; i >= 0; i-- {
			if i > 0 {
				d = s.laps[i] - s.laps[i-1]
			} else {
				d = s.laps[i]
			}
			s.timers[i].update(d)
		}
	}
}

// Pausiert die Zeitmessung. Danach kann die Messung entweder mit Cont()
// fortgesetzt oder mit Stop() beendet werden.
func (s *Stopwatch) Pause() {
	if !s.isRunning || s.isPaused {
		return
	}
	s.t1 = time.Now()
	s.isPaused = true
}

// Setzt die Zeitmessung fort.
func (s *Stopwatch) Cont() {
	if !s.isRunning || !s.isPaused {
		return
	}
	s.t0 = s.t0.Add(time.Since(s.t1))
	s.isPaused = false
}

// Loescht alle gemessenen und berechneten Werte. Entspricht dem Zustand
// direkt nach NewStopwatch.
func (s *Stopwatch) Reset() {
	s.isRunning = false
	s.isPaused = false
	s.n = 0
	s.NumLaps = 0
	for i, _ := range s.timers {
		s.timers[i].reset()
	}
	s.all.reset()
	for i := range s.lapNames {
		s.lapNames[i] = fmt.Sprintf("%d", i+1)
	}
}

// Retourniert die aufkumulierte Messdauer. Sind auch Zwischenzeiten erfasst
// worden, kann mit i (in [1..n]) die Zwischenzeit gewaehlt werden. Mit
// i=0 wird die totale Messzeit ausgewaehlt.
func (s *Stopwatch) Total(i int) time.Duration {
	if i == 0 {
		return s.all.sum
	} else if i <= s.NumLaps {
		return s.timers[i-1].sum
	} else {
		return 0
	}
}

// Retourniert die kuerzeste Messdauer.
// Fuer die Bedeutung von i: siehe [Total]
func (s *Stopwatch) Min(i int) time.Duration {
	if i == 0 {
		return s.all.min
	} else if i <= s.NumLaps {
		return s.timers[i-1].min
	} else {
		return 0
	}
}

// Retourniert die laengeste Messdauer.
// Fuer die Bedeutung von i: siehe [Total]
func (s *Stopwatch) Max(i int) time.Duration {
	if i == 0 {
		return s.all.max
	} else if i <= s.NumLaps {
		return s.timers[i-1].max
	} else {
		return 0
	}
}

// Retourniert die Anzahl Messungen.
func (s *Stopwatch) Num() int {
	return s.n
}

// Berechnet die durchschnittliche Messdauer (also den Quotienten von
// Total() / Num()).
// Fuer die Bedeutung von i: siehe [Total]
func (s *Stopwatch) Avg(i int) time.Duration {
	if s.n == 0 {
		return 0
	}
	if i == 0 {
		return s.all.sum / time.Duration(s.n)
	} else if i <= s.NumLaps {
		return s.timers[i-1].sum / time.Duration(s.n)
	} else {
		return 0
	}
}
