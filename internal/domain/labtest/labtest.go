package labtest

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

func (t *Test) AssignToReport(reportID string) {
	t.ReportID = reportID
}

func (t *Test) EnsureValid() error {
	if t.Name == "" {
		return ErrNameInvalid
	}

	return nil
}
