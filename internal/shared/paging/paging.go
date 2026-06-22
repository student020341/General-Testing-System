package paging

type PageRequest struct {
	Page     uint
	PageSize uint
}

// NewPageRequest creates a new PageRequest with the given page and page size.
// Enforces default values for page and page size if invalid.
func NewPageRequest(page, pageSize uint) PageRequest {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	return PageRequest{
		Page:     page,
		PageSize: pageSize,
	}
}

func (r PageRequest) Offset() uint {
	return (r.Page - 1) * r.PageSize
}

func (r PageRequest) Limit() uint {
	return r.PageSize
}
