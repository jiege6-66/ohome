package service

import (
	"ohome/model"
	"ohome/quarksearch"
	"ohome/service/dto"
	"strings"

	"gorm.io/gorm"
)

const (
	quarkSearchHTTPProxyKey      = "quark_search_http_proxy"
	quarkSearchHTTPSProxyKey     = "quark_search_https_proxy"
	quarkSearchChannelsKey       = "quark_search_channels"
	quarkSearchEnabledPluginsKey = "quark_search_enabled_plugins"
	quarkSearchDefaultChannels   = "tgsearchers4,Aliyun_4K_Movies,bdbdndn11,yunpanx,bsbdbfjfjff,yp123pan,sbsbsnsqq,yunpanxunlei,tianyifc,BaiduCloudDisk,txtyzy,peccxinpd,gotopan,PanjClub,kkxlzy,baicaoZY,MCPH01,MCPH02,MCPH03,bdwpzhpd,ysxb48,jdjdn1111,yggpan,MCPH086,zaihuayun,Q66Share,ucwpzy,shareAliyun,alyp_1,dianyingshare,Quark_Movies,XiangxiuNBB,ydypzyfx,ucquark,xx123pan,yingshifenxiang123,zyfb123,tyypzhpd,tianyirigeng,cloudtianyi,hdhhd21,Lsp115,oneonefivewpfx,qixingzhenren,taoxgzy,Channel_Shares_115,tyysypzypd,vip115hot,wp123zy,yunpan139,yunpan189,yunpanuc,yydf_hzl,leoziyuan,Q_dongman,yoyokuakeduanju,TG654TG,WFYSFX02,QukanMovie,yeqingjie_GJG666,movielover8888_film3,Baidu_netdisk,D_wusun,FLMdongtianfudi,KaiPanshare,QQZYDAPP,rjyxfx,PikPak_Share_Channel,btzhi,newproductsourcing,cctv1211,duan_ju,QuarkFree,yunpanNB,kkdj001,xxzlzn,pxyunpanxunlei,jxwpzy,kuakedongman,liangxingzhinan,xiangnikanj,solidsexydoll,guoman4K,zdqxm,kduanju,cilidianying,CBduanju,SharePanFilms,dzsgx,BooksRealm,Oscar_4Kmovies,douerpan,baidu_yppan,Q_jilupian,Netdisk_Movies,yunpanquark,ammmziyuan,ciliziyuanku,cili8888,jzmm_123pan,Q_dianying,domgmingapk,dianying4k,q_dianshiju,tgbokee,ucshare,godupan,gokuapan"
	quarkSearchDefaultPlugins    = "ddys,erxiao,hdr4k,jutoushe,labi,libvio,lou1,panta,susu,wanou,xuexizhinan,zhizhen,ahhhhfs,alupan,ash,clxiong,discourse,djgou,duoduo,dyyj,hdmoli,huban,jsnoteclub,kkmao,leijing,meitizy,mikuclub,muou,nsgame,ouge,panyq,shandian,xinjuc,ypfxw,yunsou,aikanzy,bixin,cldi,clmao,cyg,daishudj,feikuai,fox4k,haisou,hunhepan,jikepan,kkv,miaoso,mizixing,nyaa,pan666,pansearch,panwiki,pianku,qingying,quark4k,quarksoo,qupanshe,qupansou,sdso,sousou,wuji,xb6v,xdpan,xdyh,xiaoji,xiaozhang,xys,yiove,zxzj"
)

func loadQuarkSearchSettings() (quarksearch.Settings, error) {
	configs, err := configDao.GetByKeys(quarkSearchConfigKeys())
	if err != nil {
		return quarksearch.Settings{}, err
	}
	return buildQuarkSearchSettings(configs), nil
}

func buildQuarkSearchSettings(configs []model.Config) quarksearch.Settings {
	values := make(map[string]string, len(configs))
	present := make(map[string]struct{}, len(configs))
	for _, cfg := range configs {
		values[strings.TrimSpace(cfg.Key)] = cfg.Value
		present[strings.TrimSpace(cfg.Key)] = struct{}{}
	}

	channelsValue := values[quarkSearchChannelsKey]
	channels := splitConfigList(channelsValue)
	if _, exists := present[quarkSearchChannelsKey]; !exists {
		channels = splitConfigList(quarkSearchDefaultChannels)
	}

	pluginsValue := values[quarkSearchEnabledPluginsKey]
	if _, exists := present[quarkSearchEnabledPluginsKey]; !exists {
		pluginsValue = quarkSearchDefaultPlugins
	}

	return quarksearch.Settings{
		HTTPProxy:      strings.TrimSpace(values[quarkSearchHTTPProxyKey]),
		HTTPSProxy:     strings.TrimSpace(values[quarkSearchHTTPSProxyKey]),
		Channels:       channels,
		EnabledPlugins: splitConfigList(pluginsValue),
	}
}

func quarkSearchConfigKeys() []string {
	return []string{
		quarkSearchHTTPProxyKey,
		quarkSearchHTTPSProxyKey,
		quarkSearchChannelsKey,
		quarkSearchEnabledPluginsKey,
	}
}

func splitConfigList(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(strings.ToLower(part))
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func EnsureDefaultQuarkSearchConfigs() error {
	defaults := []struct {
		key         string
		name        string
		value       string
		remark      string
		legacyValue string
	}{
		{
			key:         quarkSearchHTTPProxyKey,
			name:        "夸克搜索 HTTP 代理",
			value:       "",
			remark:      "HTTP 请求代理地址，留空则直连",
			legacyValue: "",
		},
		{
			key:         quarkSearchHTTPSProxyKey,
			name:        "夸克搜索 HTTPS 代理",
			value:       "",
			remark:      "HTTPS 请求代理地址，留空则直连",
			legacyValue: "",
		},
		{
			key:         quarkSearchChannelsKey,
			name:        "夸克搜索 TG 频道",
			value:       quarkSearchDefaultChannels,
			remark:      "默认搜索 TG 频道，多个频道用逗号分隔",
			legacyValue: "tgsearchers3",
		},
		{
			key:         quarkSearchEnabledPluginsKey,
			name:        "夸克搜索 启用插件",
			value:       quarkSearchDefaultPlugins,
			remark:      "指定启用插件，多个插件用逗号分隔",
			legacyValue: "",
		},
	}

	for _, item := range defaults {
		if err := ensureConfigDefault(item.key, item.name, item.value, item.remark, item.legacyValue); err != nil {
			return err
		}
	}
	return nil
}

func ensureConfigDefault(key, name, value, remark, legacyValue string) error {
	config, err := configDao.GetByKey(key)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return configDao.AddOrUpdateConfig(&dto.ConfigUpdateDTO{
				Name:   name,
				Key:    key,
				Value:  value,
				IsLock: "1",
				Remark: remark,
			})
		}
		return err
	}

	current := strings.TrimSpace(config.Value)
	if current != "" && current != strings.TrimSpace(legacyValue) {
		return nil
	}

	return configDao.AddOrUpdateConfig(&dto.ConfigUpdateDTO{
		ID:     config.ID,
		Name:   name,
		Key:    key,
		Value:  value,
		IsLock: "1",
		Remark: remark,
	})
}
