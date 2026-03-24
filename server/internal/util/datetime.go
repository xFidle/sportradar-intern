package util

import "time"

const DateLayout = "2006-01-02"
const TimestampLayout = time.RFC3339

const (
	ParsingDateMessage      = "failed to parse '%s' (YYYY-MM-DD)"
	ParsingTimestampMessage = "failed to parse '%s' (" + time.RFC3339 + ")"
)
