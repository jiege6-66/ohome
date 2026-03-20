package api

import (
	"ohome/service"
)

var (
	userService              service.UserService
	configService            service.ConfigService
	todoService              service.TodoService
	dictService              service.DictService
	fileService              service.FileService
	quarkConfigService       service.QuarkConfigService
	quarkFsService           service.QuarkFsService
	userMediaHistoryService  service.UserMediaHistoryService
	doubanService            service.DoubanService
	pansouService            service.PansouService
	quarkAutoSaveTaskService service.QuarkAutoSaveTaskService
	quarkTransferTaskService service.QuarkTransferTaskService
	dropsItemService         service.DropsItemService
	dropsEventService        service.DropsEventService
	appMessageService        service.AppMessageService
)
