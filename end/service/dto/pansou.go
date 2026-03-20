package dto

import (
	"encoding/json"
	"strconv"
	"strings"
)

// PansouSearchReq 搜索入参（透传到 pansou /api/search）
type PansouSearchReq struct {
	Kw         string                 `json:"kw" form:"kw" binding:"required" required_err:"关键词不能为空"`
	Channels   []string               `json:"channels,omitempty" form:"channels"`
	Conc       *int                   `json:"conc,omitempty" form:"conc"`
	Refresh    *bool                  `json:"refresh,omitempty" form:"refresh"`
	Res        string                 `json:"res,omitempty" form:"res"`
	Src        string                 `json:"src,omitempty" form:"src"`
	Plugins    []string               `json:"plugins,omitempty" form:"plugins"`
	CloudTypes []string               `json:"cloud_types,omitempty" form:"cloud_types"`
	Ext        map[string]any         `json:"ext,omitempty"`
	Filter     *PansouSearchFilterDTO `json:"filter,omitempty"`
}

type PansouSearchFilterDTO struct {
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
}

// PansouSearchQueryDTO GET 参数（逗号分隔 + JSON 字符串）
type PansouSearchQueryDTO struct {
	Kw         string `form:"kw" binding:"required" required_err:"关键词不能为空"`
	Channels   string `form:"channels"`
	Conc       string `form:"conc"`
	Refresh    string `form:"refresh"`
	Res        string `form:"res"`
	Src        string `form:"src"`
	Plugins    string `form:"plugins"`
	CloudTypes string `form:"cloud_types"`
	Ext        string `form:"ext"`
	Filter     string `form:"filter"`
}

func (q *PansouSearchQueryDTO) ToReq() (PansouSearchReq, error) {
	req := PansouSearchReq{
		Kw:  strings.TrimSpace(q.Kw),
		Res: strings.TrimSpace(q.Res),
		Src: strings.TrimSpace(q.Src),
	}

	if strings.TrimSpace(q.Channels) != "" {
		req.Channels = splitCommaList(q.Channels)
	}
	if strings.TrimSpace(q.Plugins) != "" {
		req.Plugins = splitCommaList(q.Plugins)
	}
	if strings.TrimSpace(q.CloudTypes) != "" {
		req.CloudTypes = splitCommaList(q.CloudTypes)
	}
	if strings.TrimSpace(q.Conc) != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(q.Conc)); err == nil {
			req.Conc = &n
		} else {
			return PansouSearchReq{}, err
		}
	}
	if strings.TrimSpace(q.Refresh) != "" {
		v := strings.TrimSpace(q.Refresh)
		b := strings.EqualFold(v, "true") || v == "1"
		req.Refresh = &b
	}
	if strings.TrimSpace(q.Ext) != "" {
		ext := make(map[string]any)
		if err := json.Unmarshal([]byte(q.Ext), &ext); err != nil {
			return PansouSearchReq{}, err
		}
		req.Ext = ext
	}
	if strings.TrimSpace(q.Filter) != "" {
		var filter PansouSearchFilterDTO
		if err := json.Unmarshal([]byte(q.Filter), &filter); err != nil {
			return PansouSearchReq{}, err
		}
		req.Filter = &filter
	}

	return req, nil
}

func splitCommaList(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
