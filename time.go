package pixoo64

import "time"

type Time struct {
	UTCTime   int64
	LocalTime string
}

func (t *Time) Time() time.Time {
	return time.Unix(t.UTCTime, 0)
}
