package mytime

import "time"

var Now = func() time.Time {
	return time.Now()
}
