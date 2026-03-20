package quarksearch

func SupportedPlugins() map[string]Plugin {
	return map[string]Plugin{
		"quark4k":  Quark4KPlugin{},
		"quarksoo": QuarksooPlugin{},
		"qupansou": QuPanSouPlugin{},
	}
}
