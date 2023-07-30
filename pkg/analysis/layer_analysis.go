package analysis

import (
	"sort"
	"time"
)

// A general function for getting the most common elements
func mostCommon(mapWithCount map[string]int, n int) []string {
	type frequency struct {
		Value string
		Count int
	}

	frequencies := make([]frequency, 0, len(mapWithCount))
	for value, count := range mapWithCount {
		frequencies = append(frequencies, frequency{Value: value, Count: count})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].Count > frequencies[j].Count
	})

	values := make([]string, n)
	for i := 0; i < n && i < len(frequencies); i++ {
		values[i] = frequencies[i].Value
	}
	return values
}

// MostCommonCommands returns the most common commands used to create layers
func MostCommonCommands(layers []DockerLayer, n int) []string {
	commandFrequency := make(map[string]int)
	for _, layer := range layers {
		commandFrequency[layer.Command]++
	}
	return mostCommon(commandFrequency, n)
}

// MostProlificAuthors returns the authors who created the most layers.
func MostProlificAuthors(layers []DockerLayer, n int) []string {
	authorFrequency := make(map[string]int)
	for _, layer := range layers {
		authorFrequency[layer.Author]++
	}
	return mostCommon(authorFrequency, n)
}

// MostCommonTags returns the most common tags.
func MostCommonTags(layers []DockerLayer, n int) []string {
	tagFrequency := make(map[string]int)
	for _, layer := range layers {
		for _, tag := range layer.Tags {
			tagFrequency[tag]++
		}
	}
	return mostCommon(tagFrequency, n)
}

// General function for sorting layers
func sortLayers(layers []DockerLayer, comparison func(layer1, layer2 DockerLayer) bool, n int) []DockerLayer {
	copiedLayers := append([]DockerLayer(nil), layers...)
	sort.Slice(copiedLayers, func(i, j int) bool {
		return comparison(copiedLayers[i], copiedLayers[j])
	})
	if n > len(copiedLayers) {
		n = len(copiedLayers)
	}
	return copiedLayers[:n]
}

// LargestLayers returns the layers with the largest sizes.
func LargestLayers(layers []DockerLayer, n int) []DockerLayer {
	return sortLayers(layers, func(layer1, layer2 DockerLayer) bool {
		return layer1.Size > layer2.Size
	}, n)
}

// SmallestLayers returns the layers with the smallest sizes.
func SmallestLayers(layers []DockerLayer, n int) []DockerLayer {
	return sortLayers(layers, func(layer1, layer2 DockerLayer) bool {
		return layer1.Size < layer2.Size
	}, n)
}

// OldestLayers returns the oldest layers based on creation date.
func OldestLayers(layers []DockerLayer, n int) []DockerLayer {
	return sortLayers(layers, func(layer1, layer2 DockerLayer) bool {
		return layer1.Created.Before(layer2.Created)
	}, n)
}

// NewestLayers returns the newest layers based on creation date.
func NewestLayers(layers []DockerLayer, n int) []DockerLayer {
	return sortLayers(layers, func(layer1, layer2 DockerLayer) bool {
		return layer1.Created.After(layer2.Created)
	}, n)
}

// LayerSizeDistribution returns a distribution of layer sizes.
func LayerSizeDistribution(layers []DockerLayer) map[int64]int {
	distribution := make(map[int64]int)
	for _, layer := range layers {
		distribution[layer.Size]++
	}
	return distribution
}

// LayersInDateRange returns all layers created in a specific date range.
func LayersInDateRange(layers []DockerLayer, start, end time.Time) []DockerLayer {
	var result []DockerLayer
	for _, layer := range layers {
		if layer.Created.After(start) && layer.Created.Before(end) {
			result = append(result, layer)
		}
	}
	return result
}

// LayersWithTags returns all layers that have one or more tags.
func LayersWithTags(layers []DockerLayer) []DockerLayer {
	var result []DockerLayer
	for _, layer := range layers {
		if len(layer.Tags) > 0 {
			result = append(result, layer)
		}
	}
	return result
}

// LayersWithoutTags returns all layers that have no tags.
func LayersWithoutTags(layers []DockerLayer) []DockerLayer {
	var result []DockerLayer
	for _, layer := range layers {
		if len(layer.Tags) == 0 {
			result = append(result, layer)
		}
	}
	return result
}

// LayerWithTag returns all layers that contain a specific tag.
func LayerWithTag(layers []DockerLayer, tag string) []DockerLayer {
	var result []DockerLayer
	for _, layer := range layers {
		for _, t := range layer.Tags {
			if t == tag {
				result = append(result, layer)
				break
			}
		}
	}
	return result
}

// LayerCountByAuthor returns a map with authors as keys and the number of layers they have created as values.
func LayerCountByAuthor(layers []DockerLayer) map[string]int {
	result := make(map[string]int)
	for _, layer := range layers {
		result[layer.Author]++
	}
	return result
}

// LayerSizeByAuthor returns a map with authors as keys and the total size of all layers they have created as values.
func LayerSizeByAuthor(layers []DockerLayer) map[string]int64 {
	result := make(map[string]int64)
	for _, layer := range layers {
		result[layer.Author] += layer.Size
	}
	return result
}

// TotalSize returns the total size of all layers.
func TotalSize(layers []DockerLayer) int64 {
	var total int64
	for _, layer := range layers {
		total += layer.Size
	}
	return total
}

// AverageSize returns the average size of all layers
func AverageSize(layers []DockerLayer) float64 {
	if len(layers) == 0 {
		return 0
	}
	return float64(TotalSize(layers)) / float64(len(layers))
}

// MedianSize returns the median size of all layers
func MedianSize(layers []DockerLayer) int64 {
	layer := append([]DockerLayer(nil), layers...)
	sort.Slice(layer, func(i, j int) bool {
		return layer[i].Size < layer[j].Size
	})
	if len(layer) == 0 {
		return 0
	}
	middle := len(layer) / 2
	if len(layer)%2 == 0 {
		return (layer[middle-1].Size + layer[middle].Size) / 2
	} else {
		return layer[middle].Size
	}
}
