package analysis

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
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

// LayerDetails returns a string with detailed information about a DockerLayer.
func (layer *DockerLayer) LayerDetails() string {
	return fmt.Sprintf("ID: %s, Size: %d bytes, Command: %s, Author: %s, Created: %s, CreatedBy: %s, Tags: %v",
		layer.ID, layer.Size, layer.Command, layer.Author, layer.Created, layer.CreatedBy, layer.Tags)
}

// Hierarchy returns a string representing the full hierarchy of a DockerLayer.
func (layer *DockerLayer) Hierarchy() string {
	if layer.Parent == nil {
		return layer.ID
	}
	return layer.Parent.Hierarchy() + " -> " + layer.ID
}

// CumulativeSize returns the cumulative size of a DockerLayer and all its ancestors.
func (layer *DockerLayer) CumulativeSize() int64 {
	if layer.Parent == nil {
		return layer.Size
	}
	return layer.Size + layer.Parent.CumulativeSize()
}

// LayerToString returns a human-readable string representation of a DockerLayer.
func (layer *DockerLayer) LayerToString() string {
	return fmt.Sprintf("ID: %s, Size %d bytes, Command: %s, Author: %s", layer.ID, layer.Size, layer.Command, layer.Author)
}

// ImageToString returns a human-readable string representation of a DockerImage.
func (image *DockerImage) ImageToString() string {
	return fmt.Sprintf("Name: %s, Size: %d bytes, Layers: %d", image.Name, image.Size, len(image.Layers))
}

// Inspect gets detailed information about the docker image using `docker inspect`.
func (image *DockerImage) Inspect() (string, error) {
	output, err := exec.Command("docker", "insepct", image.Name).Output()
	if err != nil {
		return "", fmt.Errorf("failed to inspect image: %w", err)
	}

	var inspectOutput []map[string]interface{}
	err = json.Unmarshal(output, &inspectOutput)
	if err != nil {
		return "", fmt.Errorf("failed to parse inspect output: %w", err)
	}
	return fmt.Sprintf("%v", inspectOutput), nil
}

// LayersByAuthor returns all layers created by a specific author.
func (image *DockerImage) LayersByAuthor(author string) []DockerLayer {
	var layers []DockerLayer
	for _, layer := range image.Layers {
		if layer.Author == author {
			layers = append(layers, layer)
		}
	}
	return layers
}

// LayersByCommand returns all layers created with a specific command.
func (image *DockerImage) LayersByCommand(command string) []DockerLayer {
	var layers []DockerLayer
	for _, layer := range image.Layers {
		if layer.Command == command {
			layers = append(layers, layer)
		}
	}
	return layers
}

// LayersInTimeRange returns all layers created in a specific time range.
func (image *DockerImage) LayersInTimeRange(start, end time.Time) []DockerLayer {
	var layers []DockerLayer
	for _, layer := range image.Layers {
		if layer.Created.After(start) && layer.Created.Before(end) {
			layers = append(layers, layer)
		}
	}
	return layers
}

// LastNLayers returns the last N layers
func (image *DockerImage) LastNLayers(n int) []DockerLayer {
	if n > len(image.Layers) {
		n = len(image.Layers)
	}
	return image.Layers[len(image.Layers)-n:]
}

// LargestNLayers returns the largest N layers based on size.
func (image *DockerImage) LargestNLayers(n int) []DockerLayer {
	// Copy the slice to aviod modifiying the original
	copiedLayers := append([]DockerLayer(nil), image.Layers...)
	sort.Slice(copiedLayers, func(i, j int) bool {
		return copiedLayers[i].Size > copiedLayers[j].Size
	})
	if n > len(copiedLayers) {
		n = len(copiedLayers)
	}
	return copiedLayers[:n]
}

// TotalTags return the total number of tags in all layers
func (image *DockerImage) TotalTags() int {
	total := 0
	for _, layer := range image.Layers {
		total += len(layer.Tags)
	}
	return total
}

// UniqueAuthors returns a list of unique authors in all layers.
func (image *DockerImage) UniqueAuthors() []string {
	authorMap := make(map[string]struct{})
	for _, layer := range image.Layers {
		authorMap[layer.Author] = struct{}{}
	}

	authors := make([]string, 0, len(authorMap))
	for author := range authorMap {
		authors = append(authors, author)
	}
	return authors
}

// UniqueCommands returns a list of unique commands used in all layers
func (image *DockerImage) UniqueCommands() []string {
	commandMap := make(map[string]struct{})
	for _, layer := range image.Layers {
		commandMap[layer.Command] = struct{}{}
	}

	commands := make([]string, 0, len(commandMap))
	for command := range commandMap {
		commands = append(commands, command)
	}
	return commands
}

// UniqueTags returns a list of unique tags used in all layers
func (image *DockerImage) UniqueTags() []string {
	tagMap := make(map[string]struct{})
	for _, layer := range image.Layers {
		for _, tag := range layer.Tags {
			tagMap[tag] = struct{}{}
		}
	}
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	return tags
}
