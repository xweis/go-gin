package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"go-gin/docs"
	"go-gin/middleware/cors"
	"go-gin/middleware/ginLog"
	"go-gin/middleware/jwt"
	"go-gin/middleware/recovery"
	"go-gin/middleware/traceId"
	"go-gin/pkg/setting"
	v1 "go-gin/routers/api/v1"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	//拦截报错
	//恢复错误并抛出500
	r.Use(recovery.Recovery())
	//添加trace id
	r.Use(traceId.TraceId())
	//添加跨域信息
	r.Use(cors.AllowOrigin())
	//添加access log
	r.Use(ginLog.AccessLog())

	/*
		r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
		r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
		r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	*/

	//debug 模式开启 swagger
	if setting.ServerSetting.RunMode == "debug" {
		docs.SwaggerInfo.BasePath = "/api/v1"
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	//api v1
	apiV1 := r.Group("/api/v1")
	{
		//获取标签列表
		apiV1.POST("/login", v1.Login)
		apiV1.POST("/add_user", v1.AddUser)
	}

	//jwt 认证
	jwtApiV1 := r.Group("/api/v1")
	jwtApiV1.Use(jwt.AuthMiddleware())
	{
		jwtApiV1.GET("/get_user", v1.GetUser)
	}

	return r
}
