package app

import (
	"net/http"

	"go.uber.org/zap"
)

func RequestLogger() HandlerFunc {
	return func(c *Context) {
		// 读取请求体
		body, err := c.GetRawData()
		if err != nil {
			// 处理错误
			c.JSONError(http.StatusBadRequest, "读取请求体失败")
			return
		}

		c.Logger.Info("Request", zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path), zap.String("ip", c.ClientIP()), zap.String("body", string(body)))

		c.Next()
	}
}

func ResponseLogger() HandlerFunc {
	return func(c *Context) {
		c.Next()

		// 读取响应体
		body, err := c.GetRawData()
		if err != nil {
			// 处理错误
			//c.JSONError(http.StatusBadRequest, "读取响应体失败")
			return
		}

		if len(body) > 1024 {
			body = body[:1024]
		}

		c.Logger.Info("Response", zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path), zap.String("ip", c.ClientIP()), zap.String("body", string(body)), zap.Int("status", c.Writer.Status()))
	}
}
