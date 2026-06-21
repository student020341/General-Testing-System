package report

import "github.com/google/uuid"

type Status string

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "closed"
)

type Report struct {
	ID     string
	Name   string
	Status Status
}

type CreateReportInput struct {
	Name string
}

func New(input CreateReportInput) (*Report, error) {
	if input.Name == "" {
		return nil, ErrNameBlank
	}

	return &Report{
		ID:     uuid.New().String(),
		Name:   input.Name,
		Status: StatusOpen,
	}, nil
}

func (r Report) EnsureCanModify() error {
	if r.Status == StatusClosed {
		return ErrReportClosed
	}
	return nil
}
