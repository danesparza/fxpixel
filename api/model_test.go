package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/danesparza/fxpixel/api"
	"github.com/sanity-io/litter"
	"testing"
)

func TestModelFromJSON(t *testing.T) {

	tests := []struct {
		name   string
		source string
		want   api.Timeline
	}{
		{
			name: "Simple JSON to ?",
			source: `{
    "id": "123abc",
    "enabled": true,
    "name": "Simple fade in and out",
    "gpio": 18,
    "steps": [
        {
            "type": "effect",
			"effect": "fade",
            "time": 1000, 
            "meta-info": { 
                "color": {
                    "R": 128,
                    "G": 0,
                    "B": 128,
                    "W": 0
                }
            }
        },
        {
            "type": "sleep",
            "time": 10000
        },
        {
            "type": "effect",
			"effect": "fade",
            "time": 1000, 
            "meta-info": {
                "color": {
                    "R": 0,
                    "G": 0,
                    "B": 0,
                    "W": 0
                }
            }
        }
    ]}`,
			want: api.Timeline{
				ID:      "123abc",
				Enabled: true,
				Created: "",
				Name:    "Simple fade in and out",
				Steps: []api.TimelineStep{
					{
						Type:   "effect",
						Effect: "fade",
						Time:   5000,
					},
					{
						Type: "sleep",
						Time: 10000,
					},
					{
						Type:   "effect",
						Effect: "fade",
						Time:   1000,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//	Unmarshal the JSON
			out := api.Timeline{}
			err := json.Unmarshal([]byte(tt.source), &out)
			if err != nil {
				t.Errorf("Problem unmarshalling JSON: %v", err)
			}

			//	Dump what we have:
			// litter.Dump(out)

			//	What is the first step meta-info?
			litter.Dump(out.Steps[0].MetaInfo)

			//	Convert it to a json string:
			//	(Method, here: https://stackoverflow.com/a/53289976/19020)
			jsonString, _ := json.Marshal(out.Steps[0].MetaInfo)
			fmt.Println(string(jsonString))

			// convert json to struct
			s := api.FadeMeta{}
			err = json.Unmarshal(jsonString, &s)
			if err != nil {
				t.Errorf("Problem unmarshalling meta: %v", err)
			}

			litter.Dump(s)

			////	Do we have the expected type?
			//got := strangeslice.StrRangeToInts(tt.source)
			//
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("StrRangeToInts() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
