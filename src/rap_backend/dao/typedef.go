package dao

type Paging struct {
	Page      int `json:"page" form:"page"`
	PageSize  int `json:"page_size" form:"page_size"`
	PageTotal int `json:"page_total"`
	Total     int `json:"total"`
}

// NewPaging 创建一个分页结构
func NewPaging(page, size int) *Paging {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 1
	}
	return &Paging{
		Page:     page,
		PageSize: size,
	}
}

// CalPageTotal 计算总页数
func (p *Paging) CalPageTotal(rowCnt int64) {
	p.Total = int(rowCnt)
	if p.Total%p.PageSize == 0 {
		p.PageTotal = p.Total / p.PageSize
	} else {
		p.PageTotal = p.Total/p.PageSize + 1
	}
}

// Offset 获取本次分页的 offset
func (p *Paging) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit 其实就是 page size，补充 mysql 的习惯
func (p *Paging) Limit() int {
	return p.PageSize
}
