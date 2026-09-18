package leds

import (
	"fmt"
	"sync"

	"github.com/danesparza/fxpixel/internal/data"
	"github.com/danesparza/fxpixel/internal/ledctl/pixarray"
)

func (bp *BackgroundProcess) timelinePixels(config data.SystemConfig) (*pixarray.PixArray, error) {
	bp.stripMu.Lock()
	defer bp.stripMu.Unlock()
	if bp.pixels != nil {
		if bp.stripConfig != config {
			return nil, fmt.Errorf("strip configuration changed; restart fxpixel before playing timelines with the new configuration")
		}
		return bp.pixels, nil
	}
	create := bp.newStrip
	if create == nil {
		create = NewStrip
	}
	strip, err := create(config.LEDs, WithGPIOPIn(config.GPIO), WithPixelOrder(config.PixelOrder), WithNumberOfColors(config.NumberOfColors))
	if err != nil {
		return nil, err
	}
	bp.pixels = pixarray.NewPixArray(config.LEDs, config.NumberOfColors, &synchronizedStrip{LEDStrip: strip})
	bp.stripConfig = config
	return bp.pixels, nil
}

// Timelines may run concurrently. Protect the shared driver's pixel and DMA
// buffers from concurrent reads/writes while retaining existing playback behavior.
type synchronizedStrip struct {
	pixarray.LEDStrip
	mu sync.Mutex
}

func (s *synchronizedStrip) GetPixel(i int) pixarray.Pixel {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.LEDStrip.GetPixel(i)
}
func (s *synchronizedStrip) SetPixel(i int, p pixarray.Pixel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LEDStrip.SetPixel(i, p)
}
func (s *synchronizedStrip) Write() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.LEDStrip.Write()
}
