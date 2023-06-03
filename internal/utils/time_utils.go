package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"time"
)

var (
	ageRegexp = regexp.MustCompile(`(\d+)([hdw])`)
)

func AgeToTime(clock Clock, value string) (time.Time, error) {
	submatches := ageRegexp.FindAllStringSubmatch(value, -1)
	if submatches == nil {
		return time.Time{}, errors.New(
			fmt.Sprintf(
				`failed to parse age '%v'. 
Valid arguments are: 
* 'h' - hours
* 'd' - days (24 hours)
* 'w' - weeks (7 days)`,
				value,
			),
		)
	}

	durationNanos := int64(0)

	for _, submatch := range submatches {
		count, err := strconv.Atoi(submatch[1])
		if err != nil {
			return time.Time{}, errors.Wrap(err, "failed to parse number of units")
		}

		unit := submatch[2]

		switch unit {
		case "h":
			durationNanos += (time.Hour * time.Duration(count)).Nanoseconds()
		case "d":
			durationNanos += 24 * (time.Hour * time.Duration(count)).Nanoseconds()
		case "w":
			durationNanos += 7 * 24 * (time.Hour * time.Duration(count)).Nanoseconds()
		default:
			return time.Time{}, errors.Errorf("unsupported age unit %v", unit)
		}
	}

	return clock.Now().Add(-time.Duration(durationNanos)), nil
}
