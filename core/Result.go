package schmokin

import "time"

type Timings struct {
	ConnectDone      time.Duration
	DnsDone          time.Duration
	TlsHandshakeDone time.Duration
	FirstByteDone    time.Duration
	Complete         time.Duration
}

type Result struct {
	Headers map[string][]string
	Timings Timings
}

func NewResult() Result {
	return Result{
		Headers: map[string][]string{},
	}
}
