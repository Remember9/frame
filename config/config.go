package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	initOnce    sync.Once
	configName  = "config"
	configType  = "yaml"
	configPaths = []string{
		os.Getenv("ESFGO_CONFIG_PATH"),
		".",
	}
	configPathsTest = []string{
		".",
		"../",
		"../../",
		"../../../",
		"../../../../",
	}
)

type AppConfig struct {
	Name    string
	Version string
	Env     string
	MaxProc int
}

func Init() (err error) {
	initOnce.Do(func() {
		err = initViper("")
	})
	return
}

func InitTest() (err error) {
	initOnce.Do(func() {
		err = initViper("test")
	})
	return
}

func initViper(mod string) error {
	viper.SetConfigType(configType)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ESFGO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	var paths []string
	if mod == "test" {
		paths = configPathsTest
	} else {
		paths = configPaths
	}
	for _, in := range paths {
		if in != "" {
			viper.AddConfigPath(in)
			fmt.Println("config adding [", in, "] to paths to search")
		}
	}
	viper.SetConfigName(configName)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func GetConfig() *viper.Viper {
	return viper.GetViper()
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetInt32(key string) int32 {
	return viper.GetInt32(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func GetIntSlice(key string) []int {
	return viper.GetIntSlice(key)
}

func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func GetTime(key string) time.Time {
	return viper.GetTime(key)
}

func UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return viper.UnmarshalKey(key, rawVal, opts...)
}

func Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return viper.Unmarshal(rawVal, opts...)
}

func Debug() {
	viper.Debug()
}

func GetAppConfig() *AppConfig {
	return &AppConfig{
		Name:    viper.GetString("app.name"),
		Version: viper.GetString("app.version"),
		Env:     viper.GetString("app.env"),
		MaxProc: viper.GetInt("app.maxProc"),
	}
}
