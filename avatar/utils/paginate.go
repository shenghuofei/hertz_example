package utils

import (
	"gorm.io/gorm"
	"strconv"
)

// PageResult 通用分页返回结构
type PageResult[T any] struct {
	Total    int64 `json:"total"`     // 总记录数
	Page     int   `json:"page"`      // 当前页
	PageSize int   `json:"page_size"` // 每页数量
	Data     []T   `json:"data"`      // 当前页数据
}

type PageReq struct {
	Page     int `default:"1" json:"page" query:"page"`            // 当前页
	PageSize int `default:"10" json:"page_size" query:"page_size"` // 每页数量
}

func GetPageInfoSafe(pageStr, pageSizeStr string) (int, int, error) {
	page, pageSize := 1, 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err != nil {
			return 0, 0, err
		} else {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err != nil {
			return 0, 0, err
		} else {
			pageSize = ps
		}
	}

	return page, pageSize, nil
}

// Paginate 通用分页查询函数
// db: gorm.DB 对象（可包含 Where / Order 等条件）
// page: 当前页，从 1 开始
// pageSize: 每页大小
// out: 查询结果切片指针
func Paginate[T any](db *gorm.DB, page, pageSize int, out *[]T) (*PageResult[T], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	if err := db.Limit(pageSize).Offset(offset).Find(out).Error; err != nil {
		return nil, err
	}

	return &PageResult[T]{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Data:     *out,
	}, nil
}
