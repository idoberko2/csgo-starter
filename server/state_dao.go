package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

// StateDAO is in charge of persisting the server's state
type StateDAO struct {
	mutex *sync.Mutex
}

// NewStateDAO returns a new instance of StateDAO
func NewStateDAO() *StateDAO {
	mutex := sync.Mutex{}

	return &StateDAO{
		mutex: &mutex,
	}
}

// GetState reads the server's state
func (dao *StateDAO) GetState() (*State, error) {
	dao.mutex.Lock()
	defer dao.mutex.Unlock()

	return unprotectedGetState()
}

func unprotectedGetState() (*State, error) {
	jsonFile, err := os.Open("data/state.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &State{}, nil
		}
		return nil, err
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	state := State{}
	if err := json.Unmarshal(bytes, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// SetState writes the server's state
func (dao *StateDAO) SetState(state *State) error {
	dao.mutex.Lock()
	defer dao.mutex.Unlock()

	return unprotectedSetState(state)
}

func unprotectedSetState(state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("data/state.json", data, os.ModePerm)
}

// SetStartingState checks server's mode and starts if not already started atomically
func (dao *StateDAO) SetStartingState() error {
	dao.mutex.Lock()
	defer dao.mutex.Unlock()

	state, err := unprotectedGetState()
	if err != nil {
		return err
	}

	if state.Mode > ModeIdle {
		return ErrServerStarted{}
	}

	starting := State{
		Mode: ModeStarting,
	}

	return unprotectedSetState(&starting)
}
