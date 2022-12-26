package busi

import (
	_ "aggregate-task/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	swagHandler gin.HandlerFunc
)

func init() {
	swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
}
