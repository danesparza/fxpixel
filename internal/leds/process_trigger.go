package leds

import (
	"bytes"
	"context"
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

// ProcessTrigger processes the passed trigger meta
func (sp StepProcessor) ProcessTrigger(step data.TimelineStep) error {

	//	Convert the meta information:
	meta := step.MetaInfo.(data.TriggerMeta)

	//	Set our defaults:
	if meta.Verb == "" {
		meta.Verb = http.MethodPost
	}

	//	Log the meta information we have:
	log.Debug().
		Str("stepid", step.ID).
		Int32("time", step.Time.Int32).
		Any("verb", meta.Verb).
		Any("url", meta.URL).
		Any("headers", meta.Headers).
		Msg("Processing trigger")

	//	First, build the initial request with the verb, url and body (if the body exists)
	req, err := http.NewRequestWithContext(context.Background(), meta.Verb, meta.URL, bytes.NewBuffer(meta.Body))
	if err != nil {
		log.Err(err).Str("stepid", step.ID).Str("url", meta.URL).Msg("Error creating request for trigger")
		return err
	}

	//	Then, set our initial content-type header
	req.Header.Set("Content-Type", "application/json")

	//	Next, set any custom headers
	for _, h := range meta.Headers {

		//	Split the key and value
		splits := strings.Split(h, ":")

		req.Header.Set(splits[0], splits[1])
	}

	//	Finally, send the request
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Str("stepid", step.ID).Str("url", meta.URL).Msg("Error with response for trigger")
		return err
	}
	defer resp.Body.Close()

	return nil
}
