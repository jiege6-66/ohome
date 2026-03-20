package dto

// QuarkApplicationDTO Quark 路由中用于定位配置的 application 参数
//
// 注意：BaseApi.Request 同时调用 ShouldBind 和 ShouldBindUri，
// 如果这里加 binding:"required" 会导致在无 body 的 GET/DELETE 上先校验失败。
type QuarkApplicationDTO struct {
	Application string `json:"application" form:"application" uri:"application"`
}
