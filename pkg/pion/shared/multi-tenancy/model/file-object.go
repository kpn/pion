package model

import "time"

// FileObject struct object is used for configuring public paths
type FileObject struct {
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"createdAt"`
	// TODO add allowed HTTP methods for this path
}
