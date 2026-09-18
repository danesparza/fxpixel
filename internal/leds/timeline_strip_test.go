package leds

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"

	"github.com/danesparza/fxpixel/internal/data"
	"github.com/danesparza/fxpixel/internal/data/const/effect"
	stepType "github.com/danesparza/fxpixel/internal/data/const/step"
	"github.com/danesparza/fxpixel/internal/ledctl/pixarray"
)

type configStub struct {
	data.AppDataService
	config data.SystemConfig
}

func (s configStub) GetSystemConfig(context.Context) (data.SystemConfig, error) { return s.config, nil }

type recordingStrip struct {
	*stubStrip
	frames [][]pixarray.Pixel
}

func (s *recordingStrip) Write() error {
	s.frames = append(s.frames, append([]pixarray.Pixel(nil), s.pixels...))
	return s.stubStrip.Write()
}

func TestFadeContinuesPreviousTimeline(t *testing.T) {
	for _, tc := range []struct {
		name   string
		colors int
		color  data.MetaColor
	}{
		{"RGB", 3, data.MetaColor{R: 200, G: 100, B: 50}},
		{"white", 4, data.MetaColor{W: 200}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			strip := &recordingStrip{stubStrip: newStubStrip(3, tc.colors)}
			creates := 0
			bp := BackgroundProcess{
				DB:       configStub{config: data.SystemConfig{GPIO: 18, LEDs: 3, PixelOrder: "GRB", NumberOfColors: tc.colors}},
				newStrip: func(int, ...option) (pixarray.LEDStrip, error) { creates++; return strip, nil },
			}
			// Go 1.23 (the module minimum) predates testing.T.Context.
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			bp.StartTimelinePlay(ctx, PlayTimelineRequest{ProcessID: "solid", RequestedTimeline: data.Timeline{Steps: []data.TimelineStep{{Type: stepType.Effect, Effect: effect.Solid, MetaInfo: data.SolidMeta{Color: tc.color}}}}})
			strip.frames = nil
			bp.StartTimelinePlay(ctx, PlayTimelineRequest{ProcessID: "fade", RequestedTimeline: data.Timeline{Steps: []data.TimelineStep{{Type: stepType.Effect, Effect: effect.Fade, Time: sql.NullInt32{Int32: 100, Valid: true}, MetaInfo: data.FadeMeta{}}}}})
			if creates != 1 {
				t.Fatalf("created %d strips, want 1", creates)
			}
			intermediate := false
			for _, frame := range strip.frames {
				for _, p := range frame {
					if tc.colors == 4 {
						intermediate = intermediate || p.W > 0 && p.W < tc.color.W
					} else {
						intermediate = intermediate || p.R > 0 && p.R < tc.color.R
					}
				}
			}
			if !intermediate {
				t.Fatal("second timeline did not fade from the previous color")
			}
			for _, p := range strip.pixels {
				if p != (pixarray.Pixel{}) {
					t.Fatalf("fade did not finish at black: %+v", p)
				}
			}
		})
	}
}

func TestTimelinePixelsInitialization(t *testing.T) {
	config := data.SystemConfig{LEDs: 3, NumberOfColors: 4}
	t.Run("concurrent reuse", func(t *testing.T) {
		creates := 0
		bp := BackgroundProcess{newStrip: func(int, ...option) (pixarray.LEDStrip, error) { creates++; return newStubStrip(3, 4), nil }}
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				pa, err := bp.timelinePixels(config)
				if err != nil {
					t.Error(err)
					return
				}
				pa.SetAll(pixarray.Pixel{W: 20})
				pa.GetPixels()
				pa.Write()
			}()
		}
		wg.Wait()
		if creates != 1 {
			t.Fatalf("created %d strips", creates)
		}
		changed := config
		changed.LEDs++
		if _, err := bp.timelinePixels(changed); err == nil {
			t.Fatal("accepted incompatible configuration")
		}
	})
	t.Run("retry failure", func(t *testing.T) {
		attempts := 0
		bp := BackgroundProcess{newStrip: func(int, ...option) (pixarray.LEDStrip, error) {
			attempts++
			if attempts == 1 {
				return nil, errors.New("unavailable")
			}
			return newStubStrip(3, 4), nil
		}}
		if _, err := bp.timelinePixels(config); err == nil {
			t.Fatal("expected initialization error")
		}
		if _, err := bp.timelinePixels(config); err != nil {
			t.Fatal(err)
		}
	})
}
