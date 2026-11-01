package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
)

var (
	doCpuProf, doMemProf, doTrace       bool
	cpuProfFile, memProfFile, traceFile string
	fhCpu, fhMem, fhTrace               *os.File
)

func init() {
	cpuProfFile = fmt.Sprintf("%s.cpuprof", path.Base(os.Args[0]))
	memProfFile = fmt.Sprintf("%s.memprof", path.Base(os.Args[0]))
	traceFile = fmt.Sprintf("%s.trace", path.Base(os.Args[0]))

	flag.BoolVar(&doCpuProf, "cpuprof", false,
		"write cpu profile to "+cpuProfFile)
	flag.BoolVar(&doMemProf, "memprof", false,
		"write memory profile to "+memProfFile)
	flag.BoolVar(&doTrace, "trace", false,
		"write trace data to "+traceFile)
}

func StartProfiling() {
	var err error

	if doCpuProf {
		fhCpu, err = os.Create(cpuProfFile)
		if err != nil {
			log.Fatal("couldn't create cpu profile: ", err)
		}
		err = pprof.StartCPUProfile(fhCpu)
		if err != nil {
			log.Fatal("couldn't start cpu profiling: ", err)
		}
	}

	if doMemProf {
		fhMem, err = os.Create(memProfFile)
		if err != nil {
			log.Fatal("couldn't create memory profile: ", err)
		}
	}

	if doTrace {
		fhTrace, err = os.Create(traceFile)
		if err != nil {
			log.Fatal("couldn't create tracefile: ", err)
		}
		trace.Start(fhTrace)
	}
}

func StopProfiling() {
	if fhCpu != nil {
		pprof.StopCPUProfile()
		fhCpu.Close()
	}

	if fhMem != nil {
		runtime.GC()
		err := pprof.WriteHeapProfile(fhMem)
		if err != nil {
			log.Fatal("couldn't write memory profile: ", err)
		}
		fhMem.Close()
	}

	if fhTrace != nil {
		trace.Stop()
		fhTrace.Close()
	}
}
