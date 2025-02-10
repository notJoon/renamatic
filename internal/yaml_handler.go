package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

// MappingLoader defines the interface for loading mappings from files.
type MappingLoader[K comparable, V any] interface {
	Load(path string) (Mapping[K, V], error)
}

// YAMLMappingLoader implements MappingLoader for YAML files.
type YAMLMappingLoader[K comparable, V any] struct{}

func NewYAMLMappingLoader[K comparable, V any]() *YAMLMappingLoader[K, V] {
	return &YAMLMappingLoader[K, V]{}
}

// Load loads a mapping from a YAML file.
func (l *YAMLMappingLoader[K, V]) Load(path string) (Mapping[K, V], error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var mapping Mapping[K, V]
	if err := yaml.Unmarshal(data, &mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}
