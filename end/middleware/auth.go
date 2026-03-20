package middleware

import (
	"ohome/dao"
	"ohome/global/constants"
	"ohome/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	TOKEN_NAME   = "Authorization"
	TOKEN_PREFIX = "Bearer "
)

func Auth() func(c *gin.Context) {
	var userDao dao.UserDao

	return func(c *gin.Context) {
		token := c.GetHeader(TOKEN_NAME)
		// 允许通过查询参数 token/access_token 传递（用于流式播放等无法带头的场景）
		if token == "" {
			if q := c.Query("token"); q != "" {
				token = TOKEN_PREFIX + q
			} else if q := c.Query("access_token"); q != "" {
				token = TOKEN_PREFIX + q
			}
		}

		// Token不存在, 直接返回
		if token == "" || !strings.HasPrefix(token, TOKEN_PREFIX) {
			utils.TokenFail(c)
			return
		}

		// Token无法解析, 直接返回
		token = token[len(TOKEN_PREFIX):]
		iJwtCustClaims, err := utils.ParseToken(token)
		nUserId := iJwtCustClaims.ID
		if err != nil || nUserId == 0 {
			utils.TokenFail(c)
			return
		}

		// 将用户信息存入上下文, 方便后续处理继续使用
		loginUser, err := userDao.GetLoginUserByID(nUserId)
		if err != nil || loginUser.RoleCode == "" {
			utils.TokenFail(c)
			return
		}
		c.Set(constants.LOGIN_USER, loginUser)
		c.Next()
	}
}
