package step

type StepType int

//go:generate stringer -type=StepType
const (
	Effect StepType = iota + 1
	Sleep
	RandomSleep
	Trigger
	Loop
)
