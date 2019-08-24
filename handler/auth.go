package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netdisk/util"
)

func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		//验证登录token是否有效
		if len(username) < 3 || !IsTokenValid(token) {
			c.Abort()
			res := util.NewRespMsg(
				10004,
				"token无效",
				nil,
			)
			c.JSON(http.StatusOK, res)
			return
		}
		c.Next()
	}
}
