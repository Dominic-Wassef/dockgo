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