package api

import (
	"ohome/service/dto"
	"ohome/utils"

	"github.com/gin-gonic/gin"
)

type File struct {
	BaseApi
}

func NewFileApi() File {
	return File{
		BaseApi: NewBaseApi(),
	}
}

// Upload 上传文件
func (f *File) Upload(context *gin.Context) {
	var fileUploadDTO dto.FileUploadDTO

	if err := f.Request(RequestOptions{Ctx: context, DTO: &fileUploadDTO}).GetErrors(); err != nil {
		return
	}

	addFile, err := fileService.AddFile(context, fileUploadDTO)

	if err != nil {
		utils.FailWithMessage(err.Error(), context)
		return
	}

	utils.OkWithData(addFile, context)
}

func (f *File) GetFileById(context *gin.Context) {
	var iCommonIDDTO dto.CommonIDDTO
	if err := f.Request(RequestOptions{Ctx: context, DTO: &iCommonIDDTO}).GetErrors(); err != nil {
		return
	}

	fileInfo, err := fileService.GetFileById(&iCommonIDDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), context)
		return
	}

	utils.OkWithData(fileInfo, context)
}

func (f *File) GetFileList(context *gin.Context) {
	var fileListDTO dto.FileListDTO
	if err := f.Request(RequestOptions{Ctx: context, DTO: &fileListDTO}).GetErrors(); err != nil {
		return
	}

	fileList, total, err := fileService.GetFileList(&fileListDTO)

	if err != nil {
		utils.FailWithMessage(err.Error(), context)
		return
	}

	utils.OkWithData(gin.H{
		"records": fileList,
		"total":   total,
	}, context)
}

func (f *File) UpdateFile(context *gin.Context) {
	var fileUpdateDTO dto.FileUpdateDTO

	if err := f.Request(RequestOptions{Ctx: context, DTO: &fileUpdateDTO}).GetErrors(); err != nil {
		return
	}

	err := fileService.UpdateFile(&fileUpdateDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), context)
		return
	}

	utils.Ok(context)
}

func (f *File) DeleteFile(context *gin.Context) {
	var iCommonIDDTO dto.CommonIDDTO

	if err := f.Request(RequestOptions{Ctx: context, DTO: &iCommonIDDTO}).GetErrors(); err != nil {
		return
	}

	err := fileService.DeleteFile(&iCommonIDDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), context)
		return
	}

	utils.OkWithMessage("删除成功", context)
}
