package dto

type TodoListDTO struct {
	Paginate
}

type TodoItemCreateDTO struct {
	Title string `json:"title" form:"title" binding:"required"`
}

type TodoItemUpdateDTO struct {
	// 注意：BaseApi.Request 会先执行 ShouldBind，再执行 ShouldBindUri。
	// 如果这里加 binding:"required"，在仅通过路径传 id 的 PUT 请求里会先报校验错误。
	ID    uint   `json:"id" form:"id" uri:"id"`
	Title string `json:"title" form:"title" binding:"required"`
}

type TodoItemStatusDTO struct {
	// 注意：BaseApi.Request 会先执行 ShouldBind，再执行 ShouldBindUri。
	// 如果这里加 binding:"required"，在仅通过路径传 id 的 PUT 请求里会先报校验错误。
	ID        uint  `json:"id" form:"id" uri:"id"`
	Completed *bool `json:"completed" form:"completed" binding:"required"`
}

type TodoReorderDTO struct {
	IDs []uint `json:"ids" form:"ids"`
}
