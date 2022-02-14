package traceid

import "shorturl/pkg/helpers"

var TraceID string

func Boot(t string) {
	if len(t) <= 0 {
		TraceID = helpers.GenerateTraceID()
	} else {
		TraceID = t
	}
}
