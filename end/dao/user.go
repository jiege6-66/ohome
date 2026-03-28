package dao

import (
	"errors"
	"ohome/global"
	"ohome/model"
	"ohome/service/dto"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDao struct {
	BaseDao
}

func (m *UserDao) CheckUserNameExist(stUserName string) bool {
	exists, _ := m.ExistsByName(stUserName, 0)
	return exists
}

func (m *UserDao) ExistsByName(stUserName string, excludeUserID uint) (bool, error) {
	return m.ExistsByNameWithDB(global.DB, stUserName, excludeUserID)
}

func (m *UserDao) ExistsByNameWithDB(db *gorm.DB, stUserName string, excludeUserID uint) (bool, error) {
	var nTotal int64
	query := db.Model(&model.User{}).Where("name = ?", strings.TrimSpace(stUserName))
	if excludeUserID != 0 {
		query = query.Where("id <> ?", excludeUserID)
	}
	if err := query.Count(&nTotal).Error; err != nil {
		return false, err
	}
	return nTotal > 0, nil
}

func (m *UserDao) AddUser(iUserAddDTO *dto.UserAddDTO) error {
	var iUser model.User
	iUserAddDTO.ConvertToModel(&iUser)

	err := normalizeUserMutationError(global.DB.Create(&iUser).Error)
	if err != nil {
		return err
	}

	iUserAddDTO.ID = iUser.ID
	iUserAddDTO.Password = ""
	return nil
}

func (m *UserDao) DeleteUserById(id uint) error {
	return m.DeleteUserByIdWithDB(global.DB, id)
}

func (m *UserDao) DeleteUserByIdWithDB(db *gorm.DB, id uint) error {
	return db.Delete(&model.User{}, id).Error
}

func (m *UserDao) UpdateUser(iUserUpdateDTO *dto.UserUpdateDTO) error {
	var iUser model.User

	if err := global.DB.First(&iUser, iUserUpdateDTO.ID).Error; err != nil {
		return err
	}
	iUserUpdateDTO.ConvertToModel(&iUser)

	return normalizeUserMutationError(global.DB.Omit("Role").Save(&iUser).Error)
}

func (m *UserDao) UpdateUserByModel(iModelUser *model.User) error {
	return m.UpdateUserByModelWithDB(global.DB, iModelUser)
}

func (m *UserDao) UpdateUserByModelWithDB(db *gorm.DB, iModelUser *model.User) error {
	return normalizeUserMutationError(db.Omit("Role").Save(iModelUser).Error)
}

func (m *UserDao) GetUserById(id uint) (model.User, error) {
	return m.GetUserByIdWithDB(global.DB, id)
}

func (m *UserDao) GetUserByIdWithDB(db *gorm.DB, id uint) (model.User, error) {
	var iUser model.User
	err := db.Preload("Role").First(&iUser, id).Error
	return iUser, err
}

func (m *UserDao) GetUserList(iUserListDTO *dto.UserListDTO) ([]model.User, int64, error) {
	var giUserList []model.User
	var nTotal int64

	filtered := global.DB.Model(&model.User{}).Preload("Role")

	if iUserListDTO.Name != "" {
		filtered = filtered.Where("name like ?", "%"+iUserListDTO.Name+"%")
	}

	if err := filtered.Count(&nTotal).Error; err != nil {
		return nil, 0, err
	}
	if err := filtered.
		Order("id asc").
		Scopes(Paginate(iUserListDTO.Paginate)).
		Find(&giUserList).Error; err != nil {
		return nil, 0, err
	}

	return giUserList, nTotal, nil
}

func (m *UserDao) GetUserByName(stUserName string) (model.User, error) {
	var iUser model.User
	err := global.DB.Model(&iUser).Preload("Role").Where("name = ?", stUserName).First(&iUser).Error

	return iUser, err
}

func (m *UserDao) GetLoginUserByID(id uint) (model.LoginUser, error) {
	var user model.User
	err := global.DB.Model(&model.User{}).Preload("Role").First(&user, id).Error
	if err != nil {
		return model.LoginUser{}, err
	}

	return model.LoginUser{
		ID:       user.ID,
		Name:     user.Name,
		RoleID:   user.RoleID,
		RoleCode: user.Role.Code,
	}, nil
}

func (m *UserDao) CountByRoleCode(roleCode string, excludeUserID uint) (int64, error) {
	return m.countByRoleCode(global.DB, roleCode, excludeUserID, false)
}

func (m *UserDao) CountByRoleCodeForUpdate(db *gorm.DB, roleCode string, excludeUserID uint) (int64, error) {
	return m.countByRoleCode(db, roleCode, excludeUserID, true)
}

func (m *UserDao) countByRoleCode(db *gorm.DB, roleCode string, excludeUserID uint, lock bool) (int64, error) {
	var total int64
	query := db.Model(&model.User{}).
		Joins("JOIN sys_role ON sys_role.id = sys_user.role_id").
		Where("sys_role.code = ?", roleCode)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if excludeUserID != 0 {
		query = query.Where("sys_user.id <> ?", excludeUserID)
	}
	err := query.Count(&total).Error
	return total, err
}

func normalizeUserMutationError(err error) error {
	if err == nil {
		return nil
	}
	if isDuplicateUserNameError(err) {
		return errors.New("用户名已存在")
	}
	return err
}

func isDuplicateUserNameError(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(strings.TrimSpace(err.Error()))
	if message == "" {
		return false
	}

	return ((strings.Contains(message, "duplicate entry") || strings.Contains(message, "duplicated key")) &&
		(strings.Contains(message, "uk_sys_user_name") || strings.Contains(message, "name"))) ||
		strings.Contains(message, "unique constraint failed: sys_user.name") ||
		strings.Contains(message, "unique constraint failed: user.name") ||
		strings.Contains(message, "uk_sys_user_name")
}
