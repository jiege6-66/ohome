package service

import (
	"testing"
	"time"

	puresqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"ohome/global"
	"ohome/model"
	"ohome/service/dto"
)

func TestEnsureTaskOwnershipAndTaskIsolation(t *testing.T) {
	restore := setupPhase2TestDB(t)
	defer restore()

	superRole := seedRole(t, "超级管理员", model.RoleCodeSuperAdmin)
	userRole := seedRole(t, "普通用户", model.RoleCodeUser)

	superAdmin := seedUser(t, "admin", superRole.ID)
	member := seedUser(t, "member", userRole.ID)

	legacyTask := model.QuarkAutoSaveTask{
		TaskName:     "legacy",
		ShareURL:     "https://example.com/legacy",
		SavePath:     "/legacy",
		ScheduleType: "daily",
		RunTime:      "08:00",
		Enabled:      true,
	}
	if err := global.DB.Create(&legacyTask).Error; err != nil {
		t.Fatalf("Create(legacyTask) error = %v", err)
	}

	memberTask := model.QuarkAutoSaveTask{
		OwnerUserID:  member.ID,
		TaskName:     "member-task",
		ShareURL:     "https://example.com/member",
		SavePath:     "/member",
		ScheduleType: "daily",
		RunTime:      "09:00",
		Enabled:      true,
	}
	if err := global.DB.Create(&memberTask).Error; err != nil {
		t.Fatalf("Create(memberTask) error = %v", err)
	}

	legacySourceTaskID := legacyTask.ID
	legacyTransfer := model.QuarkTransferTask{
		DisplayName:  "legacy-transfer",
		ShareURL:     "https://example.com/legacy-transfer",
		SavePath:     "/legacy",
		SourceType:   model.QuarkTransferTaskSourceSyncSchedule,
		SourceTaskID: &legacySourceTaskID,
		Status:       model.QuarkTransferTaskStatusSuccess,
	}
	if err := global.DB.Create(&legacyTransfer).Error; err != nil {
		t.Fatalf("Create(legacyTransfer) error = %v", err)
	}

	memberTransfer := model.QuarkTransferTask{
		OwnerUserID: member.ID,
		DisplayName: "member-transfer",
		ShareURL:    "https://example.com/member-transfer",
		SavePath:    "/member",
		SourceType:  model.QuarkTransferTaskSourceSearchManual,
		Status:      model.QuarkTransferTaskStatusSuccess,
	}
	if err := global.DB.Create(&memberTransfer).Error; err != nil {
		t.Fatalf("Create(memberTransfer) error = %v", err)
	}

	if err := EnsureTaskOwnership(); err != nil {
		t.Fatalf("EnsureTaskOwnership() error = %v", err)
	}

	var reloadedLegacyTask model.QuarkAutoSaveTask
	if err := global.DB.First(&reloadedLegacyTask, legacyTask.ID).Error; err != nil {
		t.Fatalf("First(legacyTask) error = %v", err)
	}
	if reloadedLegacyTask.OwnerUserID != superAdmin.ID {
		t.Fatalf("legacy task owner = %d, want %d", reloadedLegacyTask.OwnerUserID, superAdmin.ID)
	}

	var reloadedLegacyTransfer model.QuarkTransferTask
	if err := global.DB.First(&reloadedLegacyTransfer, legacyTransfer.ID).Error; err != nil {
		t.Fatalf("First(legacyTransfer) error = %v", err)
	}
	if reloadedLegacyTransfer.OwnerUserID != superAdmin.ID {
		t.Fatalf("legacy transfer owner = %d, want %d", reloadedLegacyTransfer.OwnerUserID, superAdmin.ID)
	}

	autoService := &QuarkAutoSaveTaskService{}
	autoList, autoTotal, err := autoService.GetList(&dto.QuarkAutoSaveTaskListDTO{
		Paginate: dto.Paginate{Page: 1, Limit: 20},
	}, member.ID)
	if err != nil {
		t.Fatalf("GetList(auto) error = %v", err)
	}
	if autoTotal != 1 || len(autoList) != 1 || autoList[0].ID != memberTask.ID {
		t.Fatalf("member auto task list mismatch: total=%d len=%d first=%d", autoTotal, len(autoList), firstAutoTaskID(autoList))
	}
	if _, err := autoService.GetByID(&dto.CommonIDDTO{ID: legacyTask.ID}, member.ID); err == nil {
		t.Fatal("GetByID() should reject cross-owner auto task access")
	}

	transferService := &QuarkTransferTaskService{}
	transferList, transferTotal, err := transferService.GetList(&dto.QuarkTransferTaskListDTO{
		Paginate: dto.Paginate{Page: 1, Limit: 20},
	}, member.ID)
	if err != nil {
		t.Fatalf("GetList(transfer) error = %v", err)
	}
	if transferTotal != 1 || len(transferList) != 1 || transferList[0].ID != memberTransfer.ID {
		t.Fatalf("member transfer task list mismatch: total=%d len=%d first=%d", transferTotal, len(transferList), firstTransferTaskID(transferList))
	}
	if err := transferService.DeleteByID(&dto.CommonIDDTO{ID: legacyTransfer.ID}, member.ID); err == nil {
		t.Fatal("DeleteByID() should reject cross-owner transfer task access")
	}
}

func TestUserMutationGuards(t *testing.T) {
	restore := setupPhase2TestDB(t)
	defer restore()

	superRole := seedRole(t, "超级管理员", model.RoleCodeSuperAdmin)
	userRole := seedRole(t, "普通用户", model.RoleCodeUser)

	admin := seedUser(t, "admin", superRole.ID)

	first := &dto.UserAddDTO{
		Name:     "alice",
		Password: "123456",
		RoleID:   userRole.ID,
		RoleCode: userRole.Code,
	}
	if err := userDao.AddUser(first); err != nil {
		t.Fatalf("AddUser(first) error = %v", err)
	}

	second := &dto.UserAddDTO{
		Name:     "alice",
		Password: "654321",
		RoleID:   userRole.ID,
		RoleCode: userRole.Code,
	}
	if err := userDao.AddUser(second); err == nil || err.Error() != "用户名已存在" {
		t.Fatalf("AddUser(second) error = %v, want 用户名已存在", err)
	}

	userService := &UserService{}
	err := userService.UpdateUser(&dto.UserUpdateDTO{
		ID:       admin.ID,
		Name:     admin.Name,
		RealName: admin.RealName,
		Avatar:   admin.Avatar,
		RoleCode: model.RoleCodeUser,
	}, model.LoginUser{
		ID:       admin.ID,
		Name:     admin.Name,
		RoleCode: model.RoleCodeSuperAdmin,
	})
	if err == nil || err.Error() != "至少需要保留一个超级管理员" {
		t.Fatalf("UpdateUser(downgrade last super admin) error = %v, want 至少需要保留一个超级管理员", err)
	}
}

func setupPhase2TestDB(t *testing.T) func() {
	t.Helper()

	previousDB := global.DB
	dsn := "file:phase2_test_" + time.Now().Format("20060102150405.000000000") + "?mode=memory&cache=shared"
	db, err := gorm.Open(puresqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "sys_",
			SingularTable: true,
		},
	})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}
	if err := db.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.QuarkAutoSaveTask{},
		&model.QuarkTransferTask{},
	); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}

	global.DB = db
	return func() {
		global.DB = previousDB
	}
}

func seedRole(t *testing.T, name string, code string) model.Role {
	t.Helper()

	role := model.Role{
		Name: name,
		Code: code,
	}
	if err := global.DB.Create(&role).Error; err != nil {
		t.Fatalf("Create(role=%s) error = %v", code, err)
	}
	return role
}

func seedUser(t *testing.T, name string, roleID uint) model.User {
	t.Helper()

	user := model.User{
		Name:     name,
		RealName: name,
		Password: "123456",
		RoleID:   roleID,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		t.Fatalf("Create(user=%s) error = %v", name, err)
	}
	if err := global.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		t.Fatalf("Reload(user=%s) error = %v", name, err)
	}
	return user
}

func firstAutoTaskID(tasks []model.QuarkAutoSaveTask) uint {
	if len(tasks) == 0 {
		return 0
	}
	return tasks[0].ID
}

func firstTransferTaskID(tasks []model.QuarkTransferTask) uint {
	if len(tasks) == 0 {
		return 0
	}
	return tasks[0].ID
}
