package dao

import (
	"ohome/global"
	"ohome/model"
	"ohome/service/dto"
)

type FileDao struct {
	BaseDao
}

func (m *FileDao) AddFile(iFileAddDTO *dto.FileAddDTO) (model.File, error) {
	var iFile model.File
	iFile.Status = "1"
	iFileAddDTO.ConvertToModel(&iFile)

	err := global.DB.Save(&iFile).Error

	return iFile, err
}

func (m *FileDao) GetFileById(id uint) (model.File, error) {
	var iFile model.File
	err := global.DB.First(&iFile, id).Error
	return iFile, err
}

func (m *FileDao) GetFileList(listDTO *dto.FileListDTO) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	query := global.DB.Model(&model.File{}).
		Scopes(Paginate(listDTO.Paginate))

	if listDTO.Name != "" {
		query = query.Where("name like ?", "%"+listDTO.Name+"%")
	}
	if listDTO.Type != "" {
		query = query.Where("`type` = ?", listDTO.Type)
	}
	if listDTO.Status != "" {
		query = query.Where("`status` = ?", listDTO.Status)
	}
	if listDTO.UploaderId != 0 {
		query = query.Where("uploader_id = ?", listDTO.UploaderId)
	}

	query.Find(&files).Offset(-1).Limit(-1).Count(&total)

	return files, total, query.Error
}

func (m *FileDao) UpdateFile(updateDTO *dto.FileUpdateDTO) error {
	var iFile model.File

	if err := global.DB.First(&iFile, updateDTO.ID).Error; err != nil {
		return err
	}

	updateDTO.ApplyToModel(&iFile)

	return global.DB.Save(&iFile).Error
}

func (m *FileDao) DeleteFile(id uint) error {
	return global.DB.Delete(&model.File{}, id).Error
}
