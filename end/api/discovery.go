package api

import (
	"ohome/discovery"
	"ohome/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Discovery struct {
	BaseApi
}

func NewDiscoveryApi() Discovery {
	return Discovery{
		BaseApi: NewBaseApi(),
	}
}

func (d *Discovery) GetInfo(c *gin.Context) {
	manager := discovery.Default
	if manager == nil {
		utils.FailWithMessage("服务发现未初始化", c)
		return
	}

	utils.OkWithData(manager.DiscoveryInfo(requestOrigin(c, manager.Port())), c)
}

func requestOrigin(c *gin.Context, fallbackPort int) string {
	proto := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto"))
	if proto == "" {
		if c.Request.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}

	host := strings.TrimSpace(c.GetHeader("X-Forwarded-Host"))
	if host == "" {
		host = strings.TrimSpace(c.Request.Host)
	}
	if host == "" {
		host = "127.0.0.1"
		if fallbackPort > 0 {
			host = host + ":" + strconv.Itoa(fallbackPort)
		}
	}
	return proto + "://" + host
}
