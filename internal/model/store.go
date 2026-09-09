package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type MapperStore struct {
	path string
}

func NewMapperStore(path string) *MapperStore { return &MapperStore{path: path} }

func (s *MapperStore) Load() (Mapper, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return Mapper{}, fmt.Errorf("reading mapper: %w", err)
	}
	var m Mapper
	if err := json.Unmarshal(data, &m); err != nil {
		return Mapper{}, fmt.Errorf("parsing mapper: %w", err)
	}
	return m, nil
}

func (s *MapperStore) Save(m Mapper) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding mapper: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("writing mapper: %w", err)
	}
	return nil
}
