package dao

import (
	"ohome/global"
	"ohome/model"
	"ohome/service/dto"
)

type UserDao struct {
	BaseDao
}

func (m *UserDao) CheckUserNameExist(stUserName string) bool {
	var nTotal int64
	var iUser model.User
	global.DB.Model(&iUser).Where("name = ? ", stUserName).
		Count(&nTotal)

	return nTotal > 0
}

func (m *UserDao) AddUser(iUserAddDTO *dto.UserAddDTO) error {
	var iUser model.User
	iUserAddDTO.ConvertToModel(&iUser)

	err := global.DB.Save(&iUser).Error
	if err != nil {
		return err
	} else {
		iUserAddDTO.ID = iUser.ID
		iUserAddDTO.Password = ""
		return nil
	}

}

func (m *UserDao) DeleteUserById(id uint) error {
	return global.DB.Delete(&model.User{}, id).Error
}

func (m *UserDao) UpdateUser(iUserUpdateDTO *dto.UserUpdateDTO) error {
	var iUser model.User

	if err := global.DB.First(&iUser, iUserUpdateDTO.ID).Error; err != nil {
		return err
	}
	iUserUpdateDTO.ConvertToModel(&iUser)

	return global.DB.Omit("Role").Save(&iUser).Error
}

func (m *UserDao) UpdateUserByModel(iModelUser *model.User) error {
	return global.DB.Omit("Role").Save(iModelUser).Error
}

func (m *UserDao) GetUserById(id uint) (model.User, error) {
	var iUser model.User
	err := global.DB.Preload("Role").First(&iUser, id).Error
	return iUser, err
}

func (m *UserDao) GetUserList(iUserListDTO *dto.UserListDTO) ([]model.User, int64, error) {
	var giUserList []model.User
	var nTotal int64

	query := global.DB.Model(&model.User{}).
		Preload("Role").
		Scopes(Paginate(iUserListDTO.Paginate))

	if iUserListDTO.Name != "" {
		query = query.Where("name like ?", "%"+iUserListDTO.Name+"%")
	}

	err := query.Find(&giUserList).Error
	if err != nil {
		return giUserList, 0, err
	}
	err = query.Offset(-1).Limit(-1).Count(&nTotal).Error

	return giUserList, nTotal, err
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
	var total int64
	query := global.DB.Model(&model.User{}).
		Joins("JOIN sys_role ON sys_role.id = sys_user.role_id").
		Where("sys_role.code = ?", roleCode)
	if excludeUserID != 0 {
		query = query.Where("sys_user.id <> ?", excludeUserID)
	}
	err := query.Count(&total).Error
	return total, err
}
