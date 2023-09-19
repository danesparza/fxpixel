package step

import "strings"

type StepType int

//go:generate stringer -type=StepType
const (
	Unknown StepType = iota
	Effect
	Sleep
	RandomSleep
	Trigger
	Loop
)

// FromString converts a string representation of a step type to a StepType
func FromString(stringtype string) StepType {
	retval := Unknown

	switch strings.ToLower(stringtype) {
	case "effect":
		retval = Effect
	case "sleep":
		retval = Sleep
	case "randomsleep":
		retval = RandomSleep
	case "trigger":
		retval = Trigger
	case "loop":
		retval = Loop
	}

	return retval
}
