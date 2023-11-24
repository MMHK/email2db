package pkg

type Pagination struct {
	Current  uint `json:"current"`
	LastPage uint `json:"lastPage"`
	Total    uint `json:"total"`
	PageSize uint `json:"pageSize"`
}

func NewPagination(total uint, current uint, pageSize uint) *Pagination {
	lastPage := total / pageSize + 1

	return &Pagination{
		Total: total,
		Current: current,
		LastPage: lastPage,
		PageSize: pageSize,
	}
}

func (this *Pagination) GetLastPage() uint {
	return this.LastPage
}

func (this *Pagination) GetTotal() uint {
	return this.Total
}

func (this *Pagination) CurrentPage() uint {
	if this.Current > this.LastPage {
		this.Current = this.LastPage
	}
	if this.Current < 0 {
		this.Current = 1
	}
	return this.Current
}

