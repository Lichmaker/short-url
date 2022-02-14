package config

import (
	"os"
	"shorturl/pkg/helpers"

	"github.com/spf13/cast"
	viperlib "github.com/spf13/viper"
)

// viper 实例
var viper *viperlib.Viper

// ConfigFunc 用于动态加载配置信息
type ConfigFunc func() map[string]interface{}

// ConfigFunc 先加载到这里，然后 loadConfig 动态生成配置信息
var ConfigFuncs map[string]ConfigFunc

func init() {
	// 初始化viper
	viper = viperlib.New()
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("appenv") // 设置前缀，同一环境用以区分
	viper.AutomaticEnv()         // 读取环境变量

	ConfigFuncs = make(map[string]ConfigFunc)
}

// 初始化配置信息，完成环境变量和config的加载
func InitConfig(env string) {
	loadEnv(env)

	loadConfig()
}

func loadConfig() {
	for name, fn := range ConfigFuncs {
		viper.Set(name, fn())
	}
}

func loadEnv(envSuffix string) {
	envPath := ".env"
	if len(envSuffix) > 0 {
		filePath := ".env" + envSuffix
		if _, err := os.Stat(filePath); err == nil {
			envPath = filePath
		}
	}

	viper.SetConfigName(envPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.WatchConfig()
}

// 读取环境变量
func Env(envName string, defaultValues ...interface{}) interface{} {
	if len(defaultValues) > 0 {
		return internalGet(envName, defaultValues[0])
	}
	return internalGet(envName)
}

func Get(path string, defaultValues ...interface{}) string {
	return GetString(path, defaultValues...)
}

// 新增配置项
func Add(name string, configFn ConfigFunc) {
	ConfigFuncs[name] = configFn
}

func internalGet(path string, defaultValues ...interface{}) interface{} {
	// config 或者环境变量都不存在，使用default
	if !viper.IsSet(path) || helpers.Empty(viper.Get(path)) {
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		return nil
	}
	return viper.Get(path)
}

// get 的 string 版本
func GetString(path string, defaultValues ...interface{}) string {
	return cast.ToString(internalGet(path, defaultValues...))
}

// get 的 int 版本
func GetInt(path string, defaultValues ...interface{}) int {
	return cast.ToInt(internalGet(path, defaultValues...))
}

func GetFloat64(path string, defaultValues ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValues...))
}

func GetInt64(path string, defaultValues ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValues...))
}

func GetUint(path string, defaultValues ...interface{}) uint {
	return cast.ToUint(internalGet(path, defaultValues...))
}

func GetBool(path string, defaultValues ...interface{}) bool {
	return cast.ToBool(internalGet(path, defaultValues...))
}

func GetStringMapString(path string, defaultValues ...interface{}) map[string]string {
	return viper.GetStringMapString(path)
}
