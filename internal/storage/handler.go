package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	dataDir = ".termflow"
	planDir = "plans"
	fileExt = ".json"
)

type Task struct {
	ID             string `json:"id,omitempty"`
	Status         string `json:"status"`
	Desc           string `json:"desc"`
	FullDesc       string `json:"fulldesc"`
	Created        string `json:"created"`
	Due            string `json:"due"`
	Blocked        bool   `json:"blocked"`
	IgnoreFromPlan bool   `json:"ignorefromplan"`
	Priority       int    `json:"priority"`
}

type Handler struct {
	dataFolder string
}

func New() (*Handler, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataFolder := filepath.Join(home, dataDir)
	if err := os.MkdirAll(dataFolder, 0o755); err != nil {
		return nil, err
	}

	planFolder := filepath.Join(dataFolder, planDir)
	if err := os.MkdirAll(planFolder, 0o755); err != nil {
		return nil, err
	}

	return &Handler{
		dataFolder: dataFolder,
	}, nil
}

func (h *Handler) LoadTasks(file string) (map[string]Task, error) {
	path := h.GetTaskPath(file)
	data := make(map[string]Task)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(path)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			return data, nil
		} else {
			return nil, err
		}
	}

	defer f.Close()

	if err := json.NewDecoder(f).Decode(&data); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	finalData := make(map[string]Task, len(data))
	for id, taskValue := range data {
		taskValue.ID = id         // Set the ID field on the copy
		finalData[id] = taskValue // Store the task with its ID populated
	}
	return finalData, nil
}

func (h *Handler) SaveTasks(file string, tasks map[string]Task) error {
	path := h.GetTaskPath(file)

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	if err := enc.Encode(tasks); err != nil {
		return err
	}
	return nil
}

func (h *Handler) SavePlanFile(file string, content string) error {
	path := filepath.Join(h.dataFolder, planDir, file)
	return os.WriteFile(path, []byte(content), 0o644)
}

func (h Handler) ReadPlanFile(file string) ([]byte, error) {
	path := filepath.Join(h.dataFolder, planDir, file)
	return os.ReadFile(path)
}

func (h *Handler) SaveBragFile(content string) error {
	path := filepath.Join(h.dataFolder, "brag.md")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("===\n" + content + "\n")
	return err
}

func (h *Handler) ListAllProjects() ([]string, error) {
	files, err := os.ReadDir(h.dataFolder)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), fileExt) {
			filenames = append(filenames, strings.TrimSuffix(file.Name(), fileExt))
		}
	}

	return filenames, nil
}

func (h *Handler) GetTaskPath(file string) string {
	if !strings.HasSuffix(file, fileExt) {
		file = file + fileExt
	}
	return filepath.Join(h.dataFolder, file)
}

func (h *Handler) GetCurrent() string {
	path := filepath.Join(h.dataFolder, ".current")
	file, err := os.ReadFile(path)
	if err != nil {
		return "default"
	}

	return string(file)
}

func (h *Handler) SaveCurrent(current string) error {
	path := filepath.Join(h.dataFolder, ".current")
	return os.WriteFile(path, []byte(current), 0o644)
}
