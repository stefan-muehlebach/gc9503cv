package event

type InputEvent interface{}

type EventQueue chan InputEvent
