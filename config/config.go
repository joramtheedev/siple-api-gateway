package config

import (
	"encoding/json"
	"os"
)

// Route maps a public path prefix (e.g. "/users") to a backend service URL.
type Route struct {
	Prefix string `json:"prefix"`
	Target string `json:"target"`
}

// Load reads routes from a JSON file.
func Load(path string) ([]Route, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var routes []Route
	if err := json.Unmarshal(data, &routes); err != nil {
		return nil, err
	}
	return routes, nil
}
