// Code generated by "stringer -type=EffectType"; DO NOT EDIT.

package effect

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[Solid-1]
	_ = x[Fade-2]
	_ = x[Gradient-3]
	_ = x[Sequence-4]
	_ = x[Rainbow-5]
	_ = x[Zip-6]
	_ = x[KnightRider-7]
	_ = x[Lightning-8]
}

const _EffectType_name = "UnknownSolidFadeGradientSequenceRainbowZipKnightRiderLightning"

var _EffectType_index = [...]uint8{0, 7, 12, 16, 24, 32, 39, 42, 53, 62}

func (i EffectType) String() string {
	if i < 0 || i >= EffectType(len(_EffectType_index)-1) {
		return "EffectType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _EffectType_name[_EffectType_index[i]:_EffectType_index[i+1]]
}
