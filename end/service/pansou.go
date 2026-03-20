package service

import (
	"context"
	"errors"
	"ohome/quarksearch"
	"strings"

	"ohome/service/dto"
)

type PansouService struct {
	BaseService
}

var embeddedQuarkSearchEngine = quarksearch.NewEngine()

func (s *PansouService) Search(ctx context.Context, req *dto.PansouSearchReq) (any, string, error) {
	if req == nil {
		return nil, "", errors.New("请求不能为空")
	}
	if strings.TrimSpace(req.Kw) == "" {
		return nil, "", errors.New("关键词不能为空")
	}

	settings, err := loadQuarkSearchSettings()
	if err != nil {
		return nil, "", err
	}

	coreReq := quarksearch.Request{
		Keyword:      strings.TrimSpace(req.Kw),
		Channels:     append([]string(nil), req.Channels...),
		Concurrency:  0,
		ForceRefresh: req.Refresh != nil && *req.Refresh,
		ResultType:   strings.TrimSpace(req.Res),
		SourceType:   strings.TrimSpace(req.Src),
		Plugins:      append([]string(nil), req.Plugins...),
		Ext:          cloneExt(req.Ext),
	}
	if req.Conc != nil {
		coreReq.Concurrency = *req.Conc
	}
	if req.Filter != nil {
		coreReq.Filter = &quarksearch.FilterConfig{
			Include: append([]string(nil), req.Filter.Include...),
			Exclude: append([]string(nil), req.Filter.Exclude...),
		}
	}

	response, err := embeddedQuarkSearchEngine.Search(ctx, settings, coreReq)
	if err != nil {
		return nil, "", err
	}
	return response, "", nil
}

func cloneExt(ext map[string]any) map[string]any {
	if len(ext) == 0 {
		return map[string]any{}
	}
	out := make(map[string]any, len(ext))
	for key, value := range ext {
		out[key] = value
	}
	return out
}
