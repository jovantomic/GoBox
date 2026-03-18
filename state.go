package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// status could be created, running, stopped
// (we can make enum but i dont know how to do that in go :D)

type ContainerState struct {
	Id      string    `json:"id"`
	Status  string    `json:"status"`
	Command string    `json:"command"`
	Created time.Time `json:"created"`
	Pid     int       `json:"pid,omitempty"`
}

// charset is in const.go
func generateId() string {
	for {
		b := make([]byte, 8)
		for i := range b {
			b[i] = charset[rand.Int63()%int64(len(charset))]
		}
		id := string(b)
		if getContainerById(id) == nil {
			return id
		}
	}
}

func newContainerState(command string) *ContainerState {
	return &ContainerState{
		Id:      generateId(),
		Status:  "created",
		Command: command,
		Created: time.Now(),
	}
}

func getContainerById(id string) *ContainerState {
	data, err := os.ReadFile(filepath.Join(stateDir, id, "state.json"))
	if err != nil {
		return nil
	}
	var state ContainerState
	json.Unmarshal(data, &state)
	return &state
}

func saveJSON(state *ContainerState) {
	dir := filepath.Join(stateDir, state.Id)
	os.MkdirAll(dir, 0755)
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		panic(err)
	}
	must(os.WriteFile(filepath.Join(dir, "state.json"), data, 0644))
}

func deleteContainerState(id string) {
	os.RemoveAll(filepath.Join(stateDir, id))
}
