package leds

import (
	"github.com/Jon-Bright/ledctl/pixarray"
	"time"
)

func abs(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}

func round(f float64) int {
	if f < 0 {
		return int(f - 0.5)
	}
	return int(f + 0.5)
}

func maxP(p pixarray.Pixel) int {
	if p.R < p.G {
		if p.G < p.B {
			if p.B < p.W {
				return p.W
			}
			return p.B
		}
		if p.G < p.W {
			return p.W
		}
		return p.G
	}
	if p.R < p.B {
		if p.B < p.W {
			return p.W
		}
		return p.B
	}
	if p.R < p.W {
		return p.W
	}
	return p.R
}

func max(a, b, c float64) float64 {
	if a < b {
		if b < c {
			return c
		}
		return b
	}
	if a < c {
		return c
	}
	return a
}

func min(a, b, c float64) float64 {
	if a > b {
		if b > c {
			return c
		}
		return b
	}
	if a > c {
		return c
	}
	return a
}

func lcm(p pixarray.Pixel) int {
	// TODO: no white support
	if p.R == 0 {
		p.R = 1
	}
	if p.G == 0 {
		p.G = 1
	}
	if p.B == 0 {
		p.B = 1
	}
	m := p.R * p.G
	for p.G != 0 {
		t := p.R % p.G
		p.R = p.G
		p.G = t
	}
	p.G = m / p.R
	m = p.G * p.B
	for p.B != 0 {
		t := p.G % p.B
		p.G = p.B
		p.B = t
	}
	return m / p.G
}

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

type KnightRider struct {
	pulseTime time.Duration
	pulseLen  int
	start     time.Time
}

func NewKnightRider(pulseTime time.Duration, pulseLen int) *KnightRider {
	kr := KnightRider{}
	kr.pulseTime = pulseTime
	kr.pulseLen = pulseLen
	return &kr
}

func (kr *KnightRider) Start(pa *pixarray.PixArray, now time.Time) {
	kr.start = now
	pa.SetAll(pixarray.Pixel{0, 0, 0, 0})
}

func (kr *KnightRider) NextStep(pa *pixarray.PixArray, now time.Time) time.Duration {
	pulse := now.Sub(kr.start).Nanoseconds() / kr.pulseTime.Nanoseconds()
	pulseProgress := float64(now.Sub(kr.start).Nanoseconds()-(pulse*kr.pulseTime.Nanoseconds())) / float64(kr.pulseTime.Nanoseconds())
	pulseHead := int(float64(pa.NumPixels()+kr.pulseLen) * pulseProgress)
	pulseDir := 0
	if pulse%2 == 0 {
		pulseDir = 1
	} else {
		pulseDir = -1
		pulseHead = pa.NumPixels() - pulseHead
	}
	pulseTail := pulseHead + (pulseDir * kr.pulseLen * -1)
	if pulseTail < 0 {
		pulseTail = 0
	} else if pulseTail >= pa.NumPixels() {
		pulseTail = pa.NumPixels() - 1
	}
	rangeHead := 0
	if pulseHead < 0 {
		rangeHead = 0
	} else if pulseHead >= pa.NumPixels() {
		rangeHead = pa.NumPixels() - 1
	} else {
		rangeHead = pulseHead
	}
	for i := pulseTail; i != rangeHead; i = i + pulseDir {
		v := int((float64(kr.pulseLen-abs(pulseHead-i))/float64(kr.pulseLen))*126.0) + 1
		pa.SetOne(i, pixarray.Pixel{v, 0, 0, 0})
	}
	return time.Millisecond
}
