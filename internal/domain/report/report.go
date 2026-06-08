package report

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

func New(id string, name string) *Report {
	return &Report{
		ID:     id,
		Name:   name,
		Status: StatusOpen,
	}
}

func (r Report) EnsureCanModify() error {
	if r.Status == StatusClosed {
		return ErrReportClosed
	}
	return nil
}
