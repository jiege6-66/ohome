package api

import (
	"errors"
	"io"
	"ohome/global"
	"ohome/utils"
	"reflect"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BaseApi struct {
	Errors error
	Ctx    *gin.Context
	Logger *zap.SugaredLogger
}

type RequestOptions struct {
	DTO any
	Ctx *gin.Context
}

func NewBaseApi() BaseApi {
	return BaseApi{
		Logger: global.Logger,
	}

}

func (m *BaseApi) AddErrors(errNew error) {
	m.Errors = utils.AppendError(m.Errors, errNew)
}

func (m *BaseApi) GetErrors() error {
	return m.Errors
}

func (m BaseApi) Request(options RequestOptions) *BaseApi {
	var errResult error

	//绑定上下文请求
	m.Ctx = options.Ctx
	m.Errors = nil

	//绑定请求数据
	if options.DTO != nil {
		// 这里先绑定 body/query/form，再绑定 URI。
		// 原因：如果 DTO 同时包含路径参数和 body 必填字段，先 ShouldBindUri 会导致 body 必填项过早校验失败。
		// 因此，纯路径参数字段不要再额外加 binding:"required"，应在 service 层校验零值。
		err := m.Ctx.ShouldBind(options.DTO)
		if err != nil && !errors.Is(err, io.EOF) {
			errResult = utils.AppendError(errResult, err)
		}
		if hasURITag(options.DTO) {
			err := m.Ctx.ShouldBindUri(options.DTO)
			if err != nil && !errors.Is(err, io.EOF) {
				errResult = utils.AppendError(errResult, err)
			}
		}

		if errResult != nil {
			errResult = utils.ParseValidateErrors(errResult, options.DTO)
			m.AddErrors(errResult)

			utils.FailWithMessage(m.GetErrors().Error(), m.Ctx)
		}
	}
	return &m
}

func hasURITag(dto any) bool {
	return hasURITagType(reflect.TypeOf(dto))
}

func hasURITagType(t reflect.Type) bool {
	if t == nil {
		return false
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("uri") != "" {
			return true
		}

		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			if hasURITagType(fieldType) {
				return true
			}
		}
	}

	return false
}
