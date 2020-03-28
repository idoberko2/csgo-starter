package mocks

import (
	"csgo-starter/types"

	"github.com/stretchr/testify/mock"
)

// StateDAO mock
type StateDAO struct {
	mock.Mock
}

// GetState mock
func (s *StateDAO) GetState() (*types.State, error) {
	args := s.Called()

	var arg0 *types.State

	if args.Get(0) != nil {
		arg0 = args.Get(0).(*types.State)
	}

	return arg0, args.Error(1)
}

// SetState mock
func (s *StateDAO) SetState(state *types.State) (*types.State, error) {
	args := s.Called(state)

	var arg0 *types.State

	if args.Get(0) != nil {
		arg0 = args.Get(0).(*types.State)
	}

	return arg0, args.Error(1)
}

// SetStartingState mock
func (s *StateDAO) SetStartingState() (*types.State, error) {
	args := s.Called()

	var arg0 *types.State

	if args.Get(0) != nil {
		arg0 = args.Get(0).(*types.State)
	}

	return arg0, args.Error(1)
}
