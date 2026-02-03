package date

import "time"

const (
	ISO8601Format = "2006-01-02T15:04:05Z07:00"
	ISO8601Date   = "2006-01-02"
	RFC3339Format = time.RFC3339
)

func Format(t time.Time) string {
	return t.Format(RFC3339Format)
}
