package main

import (
	"math"
	"time"
)

// Dieser Typ dient der Zeitmessung.
type Stopwatch struct {
	isRunning, isPaused bool
	t                   time.Time
	d, sum, min, max    time.Duration
	n                   int
}

func NewStopwatch() *Stopwatch {
	s := &Stopwatch{}
	s.Reset()
	return s
}

// Mit Start wird eine neue Messung begonnen.
func (s *Stopwatch) Start() {
	if s.isRunning {
		return
	}
	s.t = time.Now()
	s.d = 0
	s.isRunning = true
}

// Stop beendet die Messung und aktualisiert die Variablen, welche die totale
// Messdauer als auch die Anzahl Messungen enthalten.
func (s *Stopwatch) Stop() {
	var d time.Duration

	if !s.isRunning {
		return
	}
	if !s.isPaused {
		d = time.Since(s.t)
	}
	d += s.d
	s.isRunning = false
	s.isPaused = false

	if d > s.max {
		s.max = d
	}
	if d < s.min {
		s.min = d
	}
	s.sum += d
	s.n += 1
}

// Pausiert die Zeitmessung. Danach kann die Messung entweder mit Cont()
// fortgesetzt oder mit Stop() beendet werden.
func (s *Stopwatch) Pause() {
	if !s.isRunning || s.isPaused {
		return
	}
	s.d += time.Since(s.t)
	s.isPaused = true
}

func (s *Stopwatch) Cont() {
	if !s.isRunning || !s.isPaused {
		return
	}
	s.t = time.Now()
	s.isPaused = false
}

// Setzt die gemessene Dauer auf 0 und die Anzahl Messungen ebenfalls.
func (s *Stopwatch) Reset() {
	s.isRunning = false
	s.isPaused = false
	s.sum = 0
	s.min = time.Duration(math.MaxInt64)
	s.max = time.Duration(math.MinInt64)
	s.n = 0
}

// Retourniert die totale Messdauer.
func (s *Stopwatch) Total() time.Duration {
	return s.sum
}

func (s *Stopwatch) Min() time.Duration {
	return s.min
}

func (s *Stopwatch) Max() time.Duration {
	return s.max
}

// Retourniert die Anzahl Messungen.
func (s *Stopwatch) Num() int {
	return s.n
}

// Berechnet die durchschnittliche Messdauer (also den Quotienten von
// Total() / Num()).
func (s *Stopwatch) Avg() time.Duration {
	if s.n == 0 {
		return 0
	}
	return s.sum / time.Duration(s.n)
}
