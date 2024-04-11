package ginLog

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go-gin/pkg/logging"
	"io"
	"strings"
	"time"
)

type bodyWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	//memory copy here!
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

func AccessLog() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 从原有Request.Body读取
		body, err := context.GetRawData()
		if err != nil {
			logging.Error("%v", err)
		}

		//body 只能读取一次， 在复制回去
		context.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		// 读取res body
		var responseBodyWriter bodyWriter
		//response写缓存
		responseBodyWriter = bodyWriter{
			bodyBuf:        bytes.NewBufferString(""),
			ResponseWriter: context.Writer}
		context.Writer = responseBodyWriter

		// 开始时间
		start := time.Now()

		// 处理请求
		context.Next()

		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		responseBody := strings.Trim(responseBodyWriter.bodyBuf.String(), "\n")

		switch {
		case context.Request.URL.Path == "/ping":
			break
		case len(context.Request.URL.Path) >= 9 && context.Request.URL.Path[:9] == "/swagger/":
			break
		default:
			bodyStr := bytes.Replace(body, []byte("\n"), []byte("\\n"), -1)
			bodyStr = bytes.Replace(bodyStr, []byte("\r\n"), []byte("\\r\\n"), -1)

			logging.Info("status:%v	ip:%v	traceId:%v	request:%v	requestBody:%v	resBody:%v	runtime:%v",
				context.Writer.Status(),
				context.ClientIP(),
				context.Keys["traceId"],
				context.Request.Method+" "+context.Request.RequestURI,
				string(bodyStr),
				responseBody,
				latency,
			)
		}
	}
}
