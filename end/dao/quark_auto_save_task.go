package dao

import (
	"ohome/global"
	"ohome/model"
	"ohome/service/dto"
	"time"
)

type QuarkAutoSaveTaskDao struct {
	BaseDao
}

func (d *QuarkAutoSaveTaskDao) GetByID(id uint, ownerUserID uint) (model.QuarkAutoSaveTask, error) {
	var task model.QuarkAutoSaveTask
	err := global.DB.Where("owner_user_id = ?", ownerUserID).First(&task, id).Error
	return task, err
}

func (d *QuarkAutoSaveTaskDao) GetOwnerUserIDByID(id uint) (uint, error) {
	type row struct {
		OwnerUserID uint `gorm:"column:owner_user_id"`
	}

	var record row
	if err := global.DB.Model(&model.QuarkAutoSaveTask{}).
		Select("owner_user_id").
		First(&record, id).Error; err != nil {
		return 0, err
	}
	return record.OwnerUserID, nil
}

func (d *QuarkAutoSaveTaskDao) GetList(listDTO *dto.QuarkAutoSaveTaskListDTO, ownerUserID uint) ([]model.QuarkAutoSaveTask, int64, error) {
	var tasks []model.QuarkAutoSaveTask
	var total int64

	filtered := global.DB.Model(&model.QuarkAutoSaveTask{}).Where("owner_user_id = ?", ownerUserID)

	if listDTO.TaskName != "" {
		filtered = filtered.Where("task_name like ?", "%"+listDTO.TaskName+"%")
	}
	if listDTO.Enabled != nil {
		filtered = filtered.Where("enabled = ?", *listDTO.Enabled)
	}

	if err := filtered.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := filtered.
		Order("id desc").
		Scopes(Paginate(listDTO.Paginate)).
		Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (d *QuarkAutoSaveTaskDao) GetEnabledTasks() ([]model.QuarkAutoSaveTask, error) {
	var tasks []model.QuarkAutoSaveTask
	err := global.DB.Model(&model.QuarkAutoSaveTask{}).
		Where("enabled = ? AND owner_user_id <> ?", true, 0).
		Find(&tasks).Error
	return tasks, err
}

func (d *QuarkAutoSaveTaskDao) AddOrUpdate(updateDTO *dto.QuarkAutoSaveTaskUpdateDTO, ownerUserID uint) error {
	var task model.QuarkAutoSaveTask
	if updateDTO.ID != 0 {
		if err := global.DB.Where("owner_user_id = ?", ownerUserID).First(&task, updateDTO.ID).Error; err != nil {
			return err
		}
	} else {
		task.OwnerUserID = ownerUserID
	}
	updateDTO.ConvertToModel(&task)
	task.OwnerUserID = ownerUserID
	return global.DB.Save(&task).Error
}

func (d *QuarkAutoSaveTaskDao) UpdateLastRunAt(id uint, at time.Time) error {
	return global.DB.Model(&model.QuarkAutoSaveTask{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"last_run_at": at,
			"updated_at":  time.Now(),
		}).Error
}

func (d *QuarkAutoSaveTaskDao) Delete(id uint, ownerUserID uint) error {
	return global.DB.Where("owner_user_id = ?", ownerUserID).Delete(&model.QuarkAutoSaveTask{}, id).Error
}

func (d *QuarkAutoSaveTaskDao) AssignMissingOwners(ownerUserID uint) error {
	if ownerUserID == 0 {
		return nil
	}

	return global.DB.Model(&model.QuarkAutoSaveTask{}).
		Where("owner_user_id = ?", 0).
		Update("owner_user_id", ownerUserID).Error
}
