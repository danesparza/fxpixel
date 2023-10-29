package leds

import "github.com/Jon-Bright/ledctl/pixarray"

func lerp(c1, c2 pixarray.Pixel, t float32) pixarray.Pixel {
	return pixarray.Pixel{
		R: c1.R + int(t*float32(c2.R-c1.R)),
		G: c1.G + int(t*float32(c2.G-c1.G)),
		B: c1.B + int(t*float32(c2.B-c1.B)),
		W: c1.W + int(t*float32(c2.W-c1.W)),
	}
}

type Artist interface {
	Draw(*pixarray.PixArray)
}

type Gradient struct {
	Colors [][]int `json:"colors"` // Pass 3 colors (or 4 when using RGBW strips)
}

func (g *Gradient) Draw(arr *pixarray.PixArray) {
	c1 := color(g.Colors[0][:])
	c2 := color(g.Colors[1][:])
	for i := 0; i < arr.NumPixels(); i++ {
		t := float32(i) / float32(arr.NumPixels())
		arr.SetOne(i, lerp(c1, c2, t))
	}
}

type Sequence struct {
	Colors [][]int `json:"colors"` // Pass 3 colors (or 4 when using RGBW strips)
}

func (seq *Sequence) Draw(arr *pixarray.PixArray) {
	for i := 0; i < arr.NumPixels(); i++ {
		c := seq.Colors[i%len(seq.Colors)][:]
		arr.SetOne(i, color(c))
	}
}
