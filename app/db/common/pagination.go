package common

import "math"

type PaginationHelper struct{}

func NewPaginationHelper() *PaginationHelper {
	return &PaginationHelper{}
}

func (ph *PaginationHelper) CalculateTotalPages(total, pageSize int64) int {
	if pageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

func (ph *PaginationHelper) CalculateOffset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * pageSize
}

func (ph *PaginationHelper) ValidatePage(page, pageSize int, total int64) (validPage int, totalPage int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	totalPage = ph.CalculateTotalPages(total, int64(pageSize))
	if page > totalPage && totalPage > 0 {
		page = totalPage
	}

	return page, totalPage
}
