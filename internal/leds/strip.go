package leds

import (
	"github.com/Jon-Bright/ledctl/pixarray"
	"strings"
)

type StripOptions struct {
	NumPixels    int
	Order        int
	OscFrequency uint
	DMAChannel   int
	PWMPins      []int
	Brightness   float32
	NumColors    int
}

type option func(*StripOptions)

func color(c []int) pixarray.Pixel {

	retval := pixarray.Pixel{}

	if len(c) == 3 {
		retval = pixarray.Pixel{
			R: c[0],
			G: c[1],
			B: c[2],
		}
	}

	if len(c) == 4 {
		retval = pixarray.Pixel{
			R: c[0],
			G: c[1],
			B: c[2],
			W: c[3],
		}
	}

	return retval
}

func WithPixelOrder(order string) option {
	return func(opts *StripOptions) {
		opts.Order = pixarray.StringOrders[strings.ToUpper(order)]
	}
}

func WithOscFreq(freq uint) option {
	return func(opts *StripOptions) {
		opts.OscFrequency = freq
	}
}

func WithGPIOPIn(pin int) option {
	return func(opts *StripOptions) {
		opts.PWMPins = []int{pin}
	}
}

func WithNumberOfColors(num int) option {
	return func(opts *StripOptions) {
		opts.NumColors = num
	}
}

func WithDMAChannel(channel int) option {
	return func(opts *StripOptions) {
		opts.DMAChannel = channel
	}
}

func Scale(c1 pixarray.Pixel, t float32) pixarray.Pixel {
	return pixarray.Pixel{
		R: int(t * float32(c1.R)),
		G: int(t * float32(c1.G)),
		B: int(t * float32(c1.B)),
		W: int(t * float32(c1.W)),
	}
}

func WithBrightness(b float32) option {
	return func(opts *StripOptions) {
		opts.Brightness = b
	}
}

type LEDStripWithBrightness struct {
	pixarray.LEDStrip
	Brightness float32
}

func (s *LEDStripWithBrightness) SetPixel(i int, p pixarray.Pixel) {
	s.LEDStrip.SetPixel(i, GammaCorrect(Scale(p, s.Brightness)))
}

func NewStrip(numPixels int, options ...option) (pixarray.LEDStrip, error) {
	opts := StripOptions{
		NumColors:    3,
		NumPixels:    numPixels,
		Order:        pixarray.GRB,
		OscFrequency: 800000,
		DMAChannel:   10,
		PWMPins:      []int{18},
		Brightness:   0,
	}

	for _, o := range options {
		o(&opts)
	}

	strip, err := pixarray.NewWS281x(
		opts.NumPixels,
		opts.NumColors,
		opts.Order,
		uint(opts.OscFrequency),
		opts.DMAChannel,
		opts.PWMPins,
	)
	if err != nil {
		return nil, err
	}

	if opts.Brightness != 0. {
		return &LEDStripWithBrightness{
			LEDStrip:   strip,
			Brightness: opts.Brightness,
		}, nil
	} else {
		return strip, nil
	}
}
