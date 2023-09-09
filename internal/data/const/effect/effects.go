package effect

type EffectType int

//go:generate stringer -type=EffectType
const (
	Unknown EffectType = iota
	Solid
	Fade
	Gradient
	Sequence
	Rainbow
	Zip
	KnightRider
	Lightning
)
