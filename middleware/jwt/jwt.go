package jwt

import (
	"github.com/gin-gonic/gin"
	"go-gin/pkg/app"
	"go-gin/pkg/e"
	"go-gin/pkg/logging"
	utilJwt "go-gin/pkg/util/jwt"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取请求头中的Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			app.FailedSetHttpCode(c, http.StatusUnauthorized, e.ERROR_AUTH, nil)
			c.Abort()
			return
		}

		// 检查Authorization头是否以"Bearer "开头
		if !strings.HasPrefix(authHeader, "Bearer ") {
			app.FailedSetHttpCode(c, http.StatusUnauthorized, e.ERROR_AUTH, nil)
			c.Abort()
			return
		}

		// 提取JWT
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utilJwt.ParseToken(tokenString)
		if err != nil {
			logging.WithCtx(c).Error("%v", err)
			app.FailedSetHttpCode(c, http.StatusForbidden, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
			c.Abort()
			return
		}
		c.Set("email", claims.Username)
		c.Next()
	}
}
