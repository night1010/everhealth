package repository

import (
	"context"
	"math"

	"gorm.io/gorm"
)

func countOffset(page int, limit int) int {
	return (page - 1) * limit
}

func countTotalPage(totalItem int64, limit *int) int {
	if limit == nil {
		return 1
	}
	div := int(math.Min(float64(totalItem), float64(*limit)))
	if div == 0 {
		div = 1
	}
	totalPage := int(totalItem) / div
	if int(totalItem)%div != 0 {
		totalPage++
	}
	return totalPage
}

func extractTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return nil
}
