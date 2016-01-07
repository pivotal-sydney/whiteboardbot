package model
import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct {}

func (clock RealClock) Now() time.Time {
	return time.Now()
}