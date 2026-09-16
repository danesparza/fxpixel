package leds

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/danesparza/fxpixel/internal/ledctl/effects"
	"github.com/danesparza/fxpixel/internal/ledctl/pixarray"
	rpi "github.com/danesparza/fxpixel/internal/ledctl/rpi"
)

// stubStrip implements pixarray.LEDStrip without touching hardware.
type stubStrip struct {
	mu       sync.Mutex
	pixels   []pixarray.Pixel
	writes   int
	maxValue int
}

func newStubStrip(numPixels, numColors int) *stubStrip {
	return &stubStrip{
		pixels:   make([]pixarray.Pixel, numPixels),
		maxValue: 255,
	}
}

func (s *stubStrip) RPi() *rpi.RPi      { return nil }
func (s *stubStrip) MaxPerChannel() int { return s.maxValue }
func (s *stubStrip) GetPixel(i int) pixarray.Pixel {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pixels[i]
}
func (s *stubStrip) SetPixel(i int, p pixarray.Pixel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pixels[i] = p
}
func (s *stubStrip) Write() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writes++
	return nil
}

// fakeEffect lets us control NextStep timings for testing runEffect.
type fakeEffect struct {
	mu        sync.Mutex
	calls     int
	durations []time.Duration
	started   bool
}

func (f *fakeEffect) Start(_ *pixarray.PixArray, _ time.Time) {
	f.mu.Lock()
	f.started = true
	f.mu.Unlock()
}

func (f *fakeEffect) NextStep(_ *pixarray.PixArray, _ time.Time) time.Duration {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	if f.calls-1 < len(f.durations) {
		return f.durations[f.calls-1]
	}
	return 0
}

func (f *fakeEffect) Name() string { return "fake" }

func TestRunEffectStopsOnCompletion(t *testing.T) {
	strip := newStubStrip(5, 3)
	pa := pixarray.NewPixArray(5, 3, strip)
	eff := &fakeEffect{durations: []time.Duration{5 * time.Millisecond, 5 * time.Millisecond, 0}}

	sp := StepProcessor{PixArray: pa}
	start := time.Now()
	if err := sp.runEffect(context.Background(), eff); err != nil {
		t.Fatalf("runEffect returned error: %v", err)
	}
	elapsed := time.Since(start)

	if !eff.started {
		t.Fatal("effect Start was not called")
	}
	if eff.calls != 3 {
		t.Fatalf("expected 3 NextStep calls, got %d", eff.calls)
	}
	if strip.writes != 3 {
		t.Fatalf("expected 3 writes, got %d", strip.writes)
	}
	if elapsed < 10*time.Millisecond {
		t.Fatalf("effect finished too quickly, timer durations may have been ignored; elapsed %v", elapsed)
	}
}

func TestRunEffectCancelsAndClears(t *testing.T) {
	strip := newStubStrip(3, 3)
	pa := pixarray.NewPixArray(3, 3, strip)
	eff := &fakeEffect{durations: []time.Duration{50 * time.Millisecond, 50 * time.Millisecond}}

	sp := StepProcessor{PixArray: pa}
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := sp.runEffect(ctx, eff); err != nil {
			t.Errorf("runEffect returned error: %v", err)
		}
	}()

	// Let the first tick fire, then cancel.
	time.Sleep(5 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("runEffect did not exit after cancellation")
	}

	strip.mu.Lock()
	defer strip.mu.Unlock()

	for i, p := range strip.pixels {
		if p != (pixarray.Pixel{}) {
			t.Fatalf("pixel %d not cleared on cancel: %+v", i, p)
		}
	}
	if strip.writes == 0 {
		t.Fatal("expected at least one write on cancel path")
	}
}

// Ensure stubStrip satisfies the LEDStrip interface expected by pixarray.PixArray.
var _ pixarray.LEDStrip = (*stubStrip)(nil)
var _ effects.Effect = (*fakeEffect)(nil)
