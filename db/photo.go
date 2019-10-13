package db

import (
	//"errors"
	"time"
	//"github.com/lib/pq"
)

// A Photo is a view into a photo
type Photo struct {
	ID, Caption, Location string
	Added                 time.Time
}
