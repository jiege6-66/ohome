package dao

import (
	"ohome/global"
	"ohome/model"
	"ohome/service/dto"

	"gorm.io/gorm"
)

type TodoDao struct {
	BaseDao
}

func (d *TodoDao) GetList(ownerUserID uint, listDTO *dto.TodoListDTO) ([]model.TodoItem, int64, error) {
	records := make([]model.TodoItem, 0)
	var total int64

	filtered := global.DB.Model(&model.TodoItem{}).
		Where("owner_user_id = ?", ownerUserID)

	if err := filtered.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := filtered.
		Order("completed ASC").
		Order("CASE WHEN completed = 0 AND sort_order != 0 THEN 0 ELSE 1 END ASC").
		Order("CASE WHEN completed = 0 AND sort_order != 0 THEN sort_order END ASC").
		Order("CASE WHEN completed = 0 AND sort_order = 0 THEN updated_at END DESC").
		Order("CASE WHEN completed = 1 THEN completed_at END DESC").
		Order("id DESC").
		Scopes(Paginate(listDTO.Paginate)).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (d *TodoDao) Create(item *model.TodoItem) error {
	return global.DB.Create(item).Error
}

func (d *TodoDao) GetPendingMinSortOrder(ownerUserID uint) (int64, error) {
	type row struct {
		SortOrder *int64 `gorm:"column:sort_order"`
	}

	var result row
	err := global.DB.Model(&model.TodoItem{}).
		Where("owner_user_id = ? AND completed = ?", ownerUserID, false).
		Select("MIN(sort_order) AS sort_order").
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	if result.SortOrder == nil {
		return 0, nil
	}
	return *result.SortOrder, nil
}

func (d *TodoDao) GetByOwnerAndID(id, ownerUserID uint) (model.TodoItem, error) {
	var item model.TodoItem
	err := global.DB.
		Where("id = ? AND owner_user_id = ?", id, ownerUserID).
		First(&item).Error
	return item, err
}

func (d *TodoDao) Save(item *model.TodoItem) error {
	return global.DB.Save(item).Error
}

func (d *TodoDao) CountPendingByOwner(ownerUserID uint) (int64, error) {
	var total int64
	err := global.DB.Model(&model.TodoItem{}).
		Where("owner_user_id = ? AND completed = ?", ownerUserID, false).
		Count(&total).Error
	return total, err
}

func (d *TodoDao) CountPendingByOwnerAndIDs(ownerUserID uint, ids []uint) (int64, error) {
	var total int64
	err := global.DB.Model(&model.TodoItem{}).
		Where("owner_user_id = ? AND completed = ? AND id IN ?", ownerUserID, false, ids).
		Count(&total).Error
	return total, err
}

func (d *TodoDao) ReorderPending(ownerUserID uint, updatedBy uint, ids []uint) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for index, id := range ids {
			if err := tx.Model(&model.TodoItem{}).
				Where("id = ? AND owner_user_id = ? AND completed = ?", id, ownerUserID, false).
				Updates(map[string]any{
					"sort_order": int64(index+1) * 1024,
					"updated_by": updatedBy,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *TodoDao) Delete(id uint) error {
	return global.DB.Delete(&model.TodoItem{}, id).Error
}
