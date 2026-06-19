package labtest

import "github.com/google/uuid"

type Status string

const (
	StatusOpen      Status = "open"
	StatusCompleted Status = "completed"
)

type Test struct {
	ID       string
	ReportID string
	Name     string
	Status   Status
}

func (t Test) EnsureCanModify() error {
	if t.Status == StatusCompleted {
		return ErrTestCompleted
	}
	return nil
}

type TestFields struct {
	Name   string
	Status Status
}

type CreateTestInput struct {
	ID       string
	ReportID string
	TestFields
}

// New creates a new lab test
func New(input CreateTestInput) (*Test, error) {
	// TODO domain service to ensure report exists
	if input.ReportID == "" {
		return nil, ErrReportIDBlank
	}

	if input.Name == "" {
		return nil, ErrNameBlank
	}

	return &Test{
		ID:       uuid.NewString(),
		ReportID: input.ReportID,
		Name:     input.Name,
		Status:   input.Status,
	}, nil
}

func (t *Test) Update(fields TestFields) error {
	if t == nil {
		return ErrNilEntity
	}

	// validate update ok
	if fields.Name == "" {
		return ErrNameBlank
	}

	// update
	t.Name = fields.Name
	t.Status = fields.Status

	return nil
}
