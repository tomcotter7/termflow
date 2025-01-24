package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Task struct {
	Status   string `json:"status"`
	Desc     string `json:"desc"`
	FullDesc string `json:"fulldesc"`
	Created  string `json:"created"`
	Due      string `json:"due"`
	Blocked  bool   `json:"blocked"`
}

type Handler struct {
	dataFolder string
}

func New() (*Handler, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home directory: %w", err)
	}

	dataFolder := filepath.Join(home, ".termflow")
	if err := os.MkdirAll(dataFolder, 0o755); err != nil {
		return nil, fmt.Errorf("creating data directory: %w", err)
	}

	return &Handler{
		dataFolder: dataFolder,
	}, nil
}

func (h *Handler) LoadTasks(file string) (map[string]Task, error) {
	path := filepath.Join(h.dataFolder, file)
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	defer f.Close()

	data := make(map[string]Task)

	if err := json.NewDecoder(f).Decode(&data); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("decoding json: %w", err)
	}

	return data, nil
}

func (h *Handler) SaveTasks(file string, tasks map[string]Task) error {
	path := filepath.Join(h.dataFolder, file)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	if err := enc.Encode(tasks); err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}

	return nil
}
