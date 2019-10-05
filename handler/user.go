package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netdisk/model"
	"netdisk/util"
)

const (
	//自定义加盐的值
	pwdSalt = "*#717"
)

func SignUpHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "./static/view/signup.html")

}
//注册接口
func DoSignUpHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	if len(username) < 3 || len(password) < 6 {

		c.JSON(http.StatusOK, gin.H{
			"msg":  "Invalid Parameter",
			"code": 10000,
		})
		return
	}

	//对密码进行加盐及取sha1值加密
	encPwd := util.Sha1(util.QuickStringToBytes(password + pwdSalt))
	user := &model.User{Username:username,Password:encPwd}
	suc := user.SignUp()
	if suc == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "SignUp succeeded",
			"code": 10001,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "SignUp failed",
			"code": 10000,
		})
	}
}

//处理登录请求
func DoSignInHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	encPwd := util.Sha1(util.QuickStringToBytes(password + pwdSalt))

	flag := model.SignIn(username, encPwd)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "LogIn failed",
			"code": 10000,
		})
		return
	}

	token := util.GenToken(username)
	upRes := model.UpdateToken(username, token)
	if upRes != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "LogIn failed(fail to write token)",
			"code": 10000,
		})
	}

	res := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Username string
			Token string
		}{
			Username: username,
			Token: token,
		},
	}
	c.Data(http.StatusOK, "application/json", res.JSONBytes())
}



// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}

