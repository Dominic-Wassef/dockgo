package analysis

import "sort"

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
