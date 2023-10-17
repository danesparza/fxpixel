package data

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

// GetSystemConfig gets the system configuration information
func (a appDataService) GetSystemConfig(ctx context.Context) (SystemConfig, error) {

	retval := SystemConfig{}

	query := `select gpio, leds, pixel_order, number_of_colors from system_config limit 1;`

	stmt, err := a.DB.PreparexContext(ctx, query)
	if err != nil {
		return retval, err
	}

	rows, err := stmt.QueryxContext(ctx)
	if err != nil {
		return retval, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Err(closeErr).Msg("unable to close rows")
		}
	}()

	for rows.Next() {

		if err := rows.Scan(&retval.GPIO, &retval.LEDs, &retval.PixelOrder, &retval.NumberOfColors); err != nil {
			return retval, fmt.Errorf("problem reading into struct: %v", err)
		}

	}

	return retval, nil
}

func (a appDataService) SetSystemConfig(ctx context.Context, config SystemConfig) error {
	//TODO implement me
	panic("implement me")
}
