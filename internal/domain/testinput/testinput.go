package testinput

import "github.com/google/uuid"

type TestInputType string

// TODO consider multiple format types, like single input vs select, maybe as a stretch goal

const (
	// TestInputTypeConstant cannot be edited when test is started
	TestInputTypeConstant TestInputType = "constant"
	// TestInputTypeVariable can be edited during test
	TestInputTypeVariable TestInputType = "variable"
)

type TestInputCreateInput struct {
	TestID string
	Type   TestInputType
	Name   string
	Value  any
}

type TestInput struct {
	ID     string
	TestID string
	Type   TestInputType
	Name   string
	Value  any
}

func New(input TestInputCreateInput) (*TestInput, error) {
	if input.TestID == "" {
		return nil, ErrTestIDBlank
	}

	if input.Type == "" {
		return nil, ErrTypeBlank
	}

	if input.Name == "" {
		return nil, ErrNameBlank
	}

	return &TestInput{
		ID:     uuid.NewString(),
		TestID: input.TestID,
		Type:   input.Type,
		Value:  input.Value,
	}, nil
}
