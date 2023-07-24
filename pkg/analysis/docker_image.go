package analysis

import (
	"time"
)

// DockerLayer holds information about a Docker layer.
type DockerLayer struct {
	ID        string
	Size      int64 // in bytes
	Command   string
	Author    string
	Created   time.Time
	CreatedBy string
	Tags      []string
	Parent    *DockerLayer
}

// DockerImage holds information about a docker image
type DockerImage struct {
	Name   string
	Layers []DockerLayer
	Size   int64 // Total size in bytes
}
