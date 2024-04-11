package app

import (
	"github.com/gin-gonic/gin"
	"go-gin/pkg/convertor"
	"go-gin/pkg/e"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  e.GetMsg(http.StatusOK),
		Data: data,
	})
	return
}

func SuccessCodeStr(ctx *gin.Context, Code int, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": convertor.ToString(Code),
		"msg":  e.GetMsg(Code),
		"data": data,
	})
	return
}

func Failed(ctx *gin.Context, errCode int, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
	return
}

func FailedSetHttpCode(ctx *gin.Context, httpCode int, code int, data interface{}) {
	ctx.AbortWithStatusJSON(httpCode, gin.H{
		"code": convertor.ToString(code),
		"msg":  e.GetMsg(code),
		"data": data,
	})
	return
}
