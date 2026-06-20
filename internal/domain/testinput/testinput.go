package testinput

type TestInputType string

// TODO consider multiple format types, like single input vs select, maybe as a stretch goal

const (
	// TestInputTypeConstant cannot be edited when test is started
	TestInputTypeConstant TestInputType = "constant"
	// TestInputTypeVariable can be edited during test
	TestInputTypeVariable TestInputType = "variable"
)

type TestInput struct {
	ID     string
	TestID string
	Type   TestInputType
	Value  string // TODO this can be multiple types
}
