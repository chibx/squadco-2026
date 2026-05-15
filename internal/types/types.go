package types

import "math"

type Pagination struct {
	Page     int
	PageSize int
}

func (p *Pagination) Normalize() {
	if p == nil {
		return
	}
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	p.PageSize = int(math.Min(float64(p.PageSize), 50))
}
