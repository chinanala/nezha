package mygin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/p14yground/nezha/model"
	"github.com/p14yground/nezha/service/dao"
)

// AuthorizeOption ..
type AuthorizeOption struct {
	Guest    bool
	Member   bool
	IsPage   bool
	Msg      string
	Redirect string
	Btn      string
}

// Authorize ..
func Authorize(opt AuthorizeOption) func(*gin.Context) {
	return func(c *gin.Context) {
		token, err := c.Cookie(dao.Conf.Site.CookieName)
		var code uint64 = http.StatusForbidden
		if opt.Guest {
			code = http.StatusBadRequest
		}
		commonErr := ErrInfo{
			Title: "访问受限",
			Code:  code,
			Msg:   opt.Msg,
			Link:  opt.Redirect,
			Btn:   opt.Btn,
		}
		var isLogin bool
		if err == nil {
			isLogin = token == dao.Admin.Token && dao.Admin.Token != "" &&
				dao.Admin.TokenExpired.After(time.Now())
		}
		c.Set(model.CtxKeyIsUserLogin, isLogin)
		// 已登录且只能游客访问
		if isLogin && opt.Guest {
			ShowErrorPage(c, commonErr, opt.IsPage)
			return
		}
		// 未登录且需要登录
		if !isLogin && opt.Member {
			ShowErrorPage(c, commonErr, opt.IsPage)
			return
		}
	}
}