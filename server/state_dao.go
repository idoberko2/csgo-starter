package server

import (
	"csgo-starter/types"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/pkg/errors"
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
func (dao *StateDAO) GetState() (*types.State, error) {
	dao.mutex.Lock()
	defer dao.mutex.Unlock()

	return unprotectedGetState()
}

func unprotectedGetState() (*types.State, error) {
	jsonFile, err := os.Open("data/state.json")
	if err != nil {
		if errors.As(err, os.ErrNotExist) {
			return &types.State{}, nil
		}
		return nil, err
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	state := types.State{}
	if err := json.Unmarshal(bytes, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// SetState writes the server's state
func (dao *StateDAO) SetState(state *types.State) (*types.State, error) {
	dao.mutex.Lock()
	defer dao.mutex.Unlock()

	return unprotectedSetState(state)
}

func unprotectedSetState(state *types.State) (*types.State, error) {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile("data/state.json", data, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return state, nil
}

// SetStartingState checks server's mode and starts if not already started atomically
func (dao *StateDAO) SetStartingState() (*types.State, error) {
	dao.mutex.Lock()
	defer dao.mutex.Unlock()

	state, err := unprotectedGetState()
	if err != nil {
		return nil, err
	}

	if state.Mode > types.ModeIdle {
		return nil, types.ErrServerStarted{
			IP: state.DropletIP,
		}
	}

	starting := types.State{
		Mode: types.ModeStarting,
	}

	return unprotectedSetState(&starting)
}
