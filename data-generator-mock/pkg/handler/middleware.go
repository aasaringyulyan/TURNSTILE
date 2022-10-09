package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	rvCtx = "rv"
)

func getRv(ctx *gin.Context) (int64, error) {
	rv := ctx.Query(rvCtx)
	
	if rv == "" {
		return 0, nil
	}

	uint64Rv, err := strconv.ParseInt(rv, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint64Rv, err
}
