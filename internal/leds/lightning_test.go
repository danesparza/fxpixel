package leds

import (
	"context"
	"reflect"
	"testing"

	"github.com/danesparza/fxpixel/internal/data"
	"github.com/danesparza/fxpixel/internal/ledctl/pixarray"
)

func TestLightningRestoresAmbientColors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		colors int
		pixels []pixarray.Pixel
	}{
		{"solid RGB", 3, []pixarray.Pixel{{R: 20, G: 30, B: 40}, {R: 20, G: 30, B: 40}}},
		{"RGBW pattern", 4, []pixarray.Pixel{{R: 10, W: 15}, {G: 20, W: 25}, {B: 30, W: 35}}},
		{"black", 4, []pixarray.Pixel{{}, {}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			strip := &recordingStrip{stubStrip: newStubStrip(len(tc.pixels), tc.colors)}
			copy(strip.pixels, tc.pixels)
			sp := StepProcessor{PixArray: pixarray.NewPixArray(len(tc.pixels), tc.colors, strip)}
			// The module supports Go 1.23, before testing.T.Context was introduced.
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			err := sp.ProcessLightningEffect(ctx, data.TimelineStep{MetaInfo: data.LightningMeta{Bursts: 3, BurstLength: 1, BurstSpacing: 1, BurstBrightness: 100}})
			if err != nil {
				t.Fatal(err)
			}
			if len(strip.frames) != 6 {
				t.Fatalf("got %d frames, want 6", len(strip.frames))
			}
			for i, frame := range strip.frames {
				if i%2 == 0 {
					for _, p := range frame {
						if p != (pixarray.Pixel{B: 50, W: 100}) {
							t.Fatalf("flash %d: %+v", i, p)
						}
					}
				} else if !reflect.DeepEqual(frame, tc.pixels) {
					t.Fatalf("ambient frame %d = %+v, want %+v", i, frame, tc.pixels)
				}
			}
			if !reflect.DeepEqual(strip.pixels, tc.pixels) {
				t.Fatal("final ambient colors not restored")
			}
		})
	}
}

type cancelLightningStrip struct {
	*recordingStrip
	cancel   context.CancelFunc
	cancelAt int
}

func (s *cancelLightningStrip) Write() error {
	err := s.recordingStrip.Write()
	if len(s.frames) == s.cancelAt {
		s.cancel()
	}
	return err
}

func TestLightningCancellationClearsStrip(t *testing.T) {
	for _, tc := range []struct {
		name     string
		cancelAt int
	}{
		{"before starting", 0}, {"during flash", 1}, {"between flashes", 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			strip := &cancelLightningStrip{recordingStrip: &recordingStrip{stubStrip: newStubStrip(2, 4)}, cancel: cancel, cancelAt: tc.cancelAt}
			strip.pixels[0] = pixarray.Pixel{R: 20, W: 30}
			strip.pixels[1] = pixarray.Pixel{G: 10}
			if tc.cancelAt == 0 {
				cancel()
			}
			sp := StepProcessor{PixArray: pixarray.NewPixArray(2, 4, strip)}
			if err := sp.ProcessLightningEffect(ctx, data.TimelineStep{MetaInfo: data.LightningMeta{Bursts: 3, BurstLength: 1, BurstSpacing: 1, BurstBrightness: 100}}); err != nil {
				t.Fatal(err)
			}
			for _, p := range strip.pixels {
				if p != (pixarray.Pixel{}) {
					t.Fatalf("cancel left pixel lit: %+v", p)
				}
			}
		})
	}
}
