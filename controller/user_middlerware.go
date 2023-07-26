package controller

import "updater-server/pkg/app"

// UserMiddleware 用户中间件
func UserMiddleware() app.HandlerFunc {
	return func(c *app.Context) {
		userId := c.GetHeader("X-User-Id")
		if userId == "" {
			c.JSONError(401, "未登录")
			return
		}

		c.Set("userId", userId)

		c.Next()
	}
}

// TeamMiddleware 团队中间件
func TeamMiddleware() app.HandlerFunc {
	return func(c *app.Context) {
		teamId := c.GetHeader("X-Team-Id")
		if teamId == "" {
			c.JSONError(401, "未登录")
			return
		}

		c.Set("teamId", teamId)

		c.Next()
	}
}
