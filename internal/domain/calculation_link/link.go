package calculationlink

import "github.com/google/uuid"

// TestEntityRef is a minimal struct to reference a calculation or test input
type TestEntityRef struct {
	ID     string
	TestID string
}

type OutputType string

const (
	// a single returned value from another calculation
	OutputTypeSingle OutputType = "single"
	// an array response from another calculation
	OutputTypeArray OutputType = "array"
	// a key from a map return from another calculation
	OutputTypeKeyed OutputType = "keyed"
	// a test input or constant
	OutputTypeInput OutputType = "input"
)

type Source struct {
	TestEntityRef
	// single, array, or keyed (based on domain calculation closure details)
	// or input (based on a test input)
	OutputType string
	// optional output name if the output is a map with keys
	OutputName string
}

type Target struct {
	TestEntityRef
	// parameter name from closure details
	InputName string
}

type CreateLinkInput struct {
	ReportID string
	Source   Source
	Target   Target
}

// Link is a link between two calculations that enables the output of one to
// be the input of another. Calculations can be from different tests but
// must be in the same report.
type Link struct {
	ID       string
	ReportID string
	Source   Source
	Target   Target
}

func New(input CreateLinkInput) (*Link, error) {
	if input.ReportID == "" {
		return nil, ErrReportIDBlank
	}

	if input.Source.ID == "" {
		return nil, ErrSourceIDBlank
	}

	if input.Source.TestID == "" {
		return nil, ErrSourceTestIDBlank
	}

	if input.Source.OutputType == "" {
		return nil, ErrSourceOutputTypeBlank
	}

	if input.Target.ID == "" {
		return nil, ErrTargetIDBlank
	}

	if input.Target.TestID == "" {
		return nil, ErrTargetTestIDBlank
	}

	if input.Target.InputName == "" {
		return nil, ErrTargetInputNameBlank
	}

	return &Link{
		ID:       uuid.NewString(),
		ReportID: input.ReportID,
		Source:   input.Source,
		Target:   input.Target,
	}, nil
}
