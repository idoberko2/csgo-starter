package types

// StateDAO represents entities that can manage the server's state
type StateDAO interface {
	GetState() (*State, error)
	SetState(state *State) (*State, error)
	SetStartingState() (*State, error)
}
