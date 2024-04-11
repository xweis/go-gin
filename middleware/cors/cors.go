package cors

import "github.com/gin-gonic/gin"

func AllowOrigin() gin.HandlerFunc {
	return func(context *gin.Context) {
		origin := context.Request.Header.Get("Origin")
		if origin != "" {
			context.Header("Access-Control-Allow-Origin", "*")
		}
		context.Next()
	}
}
