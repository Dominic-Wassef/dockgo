package analysis

import (
	"fmt"
	"strconv"
	"strings"
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

// NewDockerLayer creates a new DockerLayer from a line of output from `docker history`.
func NewDockerLayer(line string, parent *DockerLayer) (*DockerLayer, error) {
	fields := strings.Fields(line)

	if len(fields) < 6 {
		return nil, fmt.Errorf("invalid line: %s", line)
	}

	size, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid size: %w", err)
	}

	created, err := time.Parse(time.RFC3339, fields[4])
	if err != nil {
		return nil, fmt.Errorf("invalid creation time: %w", err)
	}

	tags := strings.Split(fields[5], ",")

	layer := DockerLayer{
		ID:        fields[0],
		Size:      size,
		Command:   fields[2],
		Author:    fields[3],
		Created:   created,
		CreatedBy: fields[6],
		Tags:      tags,
		Parent:    parent,
	}

	return &layer, nil
}

// ParentLayer returns the parent layer of the given Docker layer, or nil if it has no parent.
func ParentLayer(layer *DockerLayer) *DockerLayer {
	return layer.Parent
}
