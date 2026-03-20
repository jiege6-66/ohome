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

func (d *QuarkAutoSaveTaskDao) GetByID(id uint) (model.QuarkAutoSaveTask, error) {
	var task model.QuarkAutoSaveTask
	err := global.DB.First(&task, id).Error
	return task, err
}

func (d *QuarkAutoSaveTaskDao) GetList(listDTO *dto.QuarkAutoSaveTaskListDTO) ([]model.QuarkAutoSaveTask, int64, error) {
	var tasks []model.QuarkAutoSaveTask
	var total int64

	query := global.DB.Model(&model.QuarkAutoSaveTask{}).Scopes(Paginate(listDTO.Paginate))

	if listDTO.TaskName != "" {
		query = query.Where("task_name like ?", "%"+listDTO.TaskName+"%")
	}
	if listDTO.Enabled != nil {
		query = query.Where("enabled = ?", *listDTO.Enabled)
	}

	query.Find(&tasks).Offset(-1).Limit(-1).Count(&total)
	return tasks, total, query.Error
}

func (d *QuarkAutoSaveTaskDao) GetEnabledTasks() ([]model.QuarkAutoSaveTask, error) {
	var tasks []model.QuarkAutoSaveTask
	err := global.DB.Model(&model.QuarkAutoSaveTask{}).
		Where("enabled = ?", true).
		Find(&tasks).Error
	return tasks, err
}

func (d *QuarkAutoSaveTaskDao) AddOrUpdate(updateDTO *dto.QuarkAutoSaveTaskUpdateDTO) error {
	var task model.QuarkAutoSaveTask
	if updateDTO.ID != 0 {
		if err := global.DB.First(&task, updateDTO.ID).Error; err != nil {
			return err
		}
	}
	updateDTO.ConvertToModel(&task)
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

func (d *QuarkAutoSaveTaskDao) Delete(id uint) error {
	return global.DB.Delete(&model.QuarkAutoSaveTask{}, id).Error
}
