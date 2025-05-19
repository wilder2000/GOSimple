package config

import (
	"errors"
	"github.com/spf13/viper"
)

type ViperConfigFile struct {
	configFile *viper.Viper
}
type ViperFileType string

const yaml ViperFileType = "yaml"
const json ViperFileType = "json"
const toml ViperFileType = "toml"

func (c *ViperConfigFile) Sub(key string) *viper.Viper {
	return c.configFile.Sub(key)
}
func ReadJSON(file string, dir string) (*ViperConfigFile, error) {
	return readFile(file, dir, json)
}
func ReadYAML(file string, dir string) (*ViperConfigFile, error) {
	return readFile(file, dir, yaml)
}
func ReadTOML(file string, dir string) (*ViperConfigFile, error) {
	return readFile(file, dir, toml)
}
func readFile(file string, dir string, fileType ViperFileType) (*ViperConfigFile, error) {

	if fileType != json && fileType != toml && fileType != yaml {
		return nil, errors.New("fileType must be JSON,YAML,TOML")
	}
	conf := &ViperConfigFile{}

	//读取文件
	conf.configFile = viper.New()
	//设置读取的配置文件
	conf.configFile.SetConfigName(file)
	//添加读取的配置文件路径
	conf.configFile.AddConfigPath(dir)
	//设置配置文件类型
	conf.configFile.SetConfigType(string(fileType))
	if err := conf.configFile.ReadInConfig(); err != nil {
		return nil, err
	}

	return conf, nil
}
func (c *ViperConfigFile) GetString(key string) string {
	return c.configFile.GetString(key)
}
func (c *ViperConfigFile) GetInt(key string) int {
	return c.configFile.GetInt(key)
}
func (c *ViperConfigFile) Unmarshal(conObject interface{}) {
	c.configFile.Unmarshal(conObject)
}
