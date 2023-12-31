package env

import (
	"strings"

	"github.com/spf13/viper"
)

// Config ...
type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetStringSlice(key string) []string
	Init()
}

type viperConfig struct{}

// Init ...
func (v *viperConfig) Init() {
	viper.SetEnvPrefix(`airpax-gateway`)
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(`.`, `_`)
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType(`json`)
	viper.SetConfigFile(`config.json`)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// GetString ...
func (v *viperConfig) GetString(key string) string {
	return viper.GetString(key)
}

// GetInt ...
func (v *viperConfig) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool ...
func (v *viperConfig) GetBool(key string) bool {
	return viper.GetBool(key)
}

//GetStringSlice ...
func (v *viperConfig) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// NewViperConfig ...
func NewViperConfig() Config {
	v := &viperConfig{}
	v.Init()

	return v
}
