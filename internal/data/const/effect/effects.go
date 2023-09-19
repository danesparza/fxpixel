package effect

import "strings"

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

// FromString converts a string representation of an effect type to a EffectType
func FromString(stringtype string) EffectType {
	retval := Unknown

	switch strings.ToLower(stringtype) {
	case "solid":
		retval = Solid
	case "fade":
		retval = Fade
	case "gradient":
		retval = Gradient
	case "sequence":
		retval = Sequence
	case "rainbow":
		retval = Rainbow
	case "zip":
		retval = Zip
	case "knightrider":
		retval = KnightRider
	case "lightning":
		retval = Lightning
	}

	return retval
}
