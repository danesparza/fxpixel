package effects

import (
	"github.com/danesparza/fxpixel/internal/ledctl/pixarray"
	"testing"
	"time"
)

// Follow the returned delays to check actual frame timing as well as colors.
func TestWhiteOnlyFadeTiming(t *testing.T) {
	for _, tc := range []struct {
		name             string
		pixels, from, to int
	}{
		{"single fade in", 1, 0, 255}, {"single fade out", 1, 255, 0},
		{"strip fade in", 10, 0, 255}, {"strip fade out", 10, 255, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pa := pixarray.NewPixArray(tc.pixels, 4, newTestLeds(tc.pixels))
			pa.SetAll(pixarray.Pixel{R: 12, G: 23, B: 34, W: tc.from})
			dest := pixarray.Pixel{R: 12, G: 23, B: 34, W: tc.to}
			f := NewFade(time.Second, dest)
			start := time.Unix(0, 0)
			f.Start(pa, start)
			now := start.Add(time.Millisecond)
			frames := 0
			intermediate := false
			previous := tc.from * tc.pixels
			for {
				next := f.NextStep(pa, now)
				total := 0
				for _, p := range pa.GetPixels() {
					if p.R != 12 || p.G != 23 || p.B != 34 {
						t.Fatalf("RGB changed: %+v", p)
					}
					if p.W < 0 || p.W > 255 {
						t.Fatalf("white out of range: %d", p.W)
					}
					total += p.W
					if p.W > 0 && p.W < 255 {
						intermediate = true
					}
				}
				if tc.to > tc.from && total < previous || tc.to < tc.from && total > previous {
					t.Fatal("fade reversed direction")
				}
				previous = total
				frames++
				if next == 0 {
					break
				}
				if next < 0 || next > 4*time.Millisecond {
					t.Fatalf("fade update delay = %v, want 0 < delay <= 4ms", next)
				}
				if frames > 10000 {
					t.Fatal("fade never completed")
				}
				now = now.Add(next)
			}
			if !intermediate {
				t.Fatal("no intermediate white brightness displayed")
			}
			if now.Before(start.Add(time.Second)) || now.After(start.Add(time.Second+4*time.Millisecond)) {
				t.Fatalf("incorrect completion time: %v", now.Sub(start))
			}
			for _, p := range pa.GetPixels() {
				if p != dest {
					t.Fatalf("final color = %+v, want %+v", p, dest)
				}
			}
		})
	}
}

func TestFadeIntervalStaysPositive(t *testing.T) {
	for _, pixels := range []int{1, 10} {
		pa := pixarray.NewPixArray(pixels, 4, newTestLeds(pixels))
		// These coprime channel differences exceed a 32-bit signed LCM and
		// produce a sub-nanosecond interval, which must not signal completion.
		dest := pixarray.Pixel{R: 251, G: 253, B: 254, W: 255}
		f := NewFade(time.Second, dest)
		start := time.Unix(0, 0)
		f.Start(pa, start)
		if delay := f.NextStep(pa, start.Add(time.Millisecond)); delay <= 0 {
			t.Fatalf("%d pixels: unfinished fade returned %v", pixels, delay)
		}
		if delay := f.NextStep(pa, start.Add(time.Second)); delay != 0 {
			t.Fatalf("finished fade returned %v", delay)
		}
		for _, p := range pa.GetPixels() {
			if p != dest {
				t.Fatalf("final pixel = %+v", p)
			}
		}
	}
}
