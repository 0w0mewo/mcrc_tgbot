package config

var defaultConfig *configManager

func init() {
	defaultConfig = newConfigManager("./config.yaml")
}

func ConfigChanged() chan bool {
	return defaultConfig.ConfigChanged()
}

func GetConfigFile() ConfigType {
	return defaultConfig.GetConfigFile()
}

func RegisterModuleConfig(module string, c IConfig) {
	defaultConfig.RegisterModuleConfig(module, c)
}

func GetModuleConfig(module string) IConfig {
	return defaultConfig.GetModuleConfig(module)
}

func Reload() {
	defaultConfig.Reload()
}

func Table() map[string]IConfig {
	return defaultConfig.Table()
}
