package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"os"
	"path"
	"path/filepath"
)

// Get 读取环境变量ENV，读取参数 --env
// 读取配置的目录：程序所在目录，程序运行目录
// 配置文件读取顺序：config/default.properties，config/${env}.properties，后面的覆盖前面的
func Get(envParams ...string) (*properties.Properties, error) {
	var env = ""
	if len(envParams) > 0 {
		env = envParams[0]
	}

	env = GetEnv(env)

	var filenames = make([]string, 0)

	executableDir, err := getExecutableDir()
	if err == nil {
		filenames = append(filenames,
			path.Join(executableDir, configPath, defaultConfigFile),
			path.Join(executableDir, configPath, fmt.Sprintf("%s%s", env, fileType)),
		)
	}

	wordDir, err := os.Getwd()
	if err == nil {
		filenames = append(filenames,
			path.Join(wordDir, configPath, defaultConfigFile),
			path.Join(wordDir, configPath, fmt.Sprintf("%s%s", env, fileType)),
		)
	}

	confDir := getConfDir()
	if confDir != "" {
		filenames = append(filenames,
			path.Join(confDir, defaultConfigFile),
			path.Join(confDir, fmt.Sprintf("%s%s", env, fileType)),
		)
	}

	if len(filenames) == 0 {
		return nil, errors.New("cannot read config path")
	}

	props, err := properties.LoadFiles(filenames, properties.UTF8, true)
	if err != nil {
		return nil, err
	}

	err = fixVariableConfig(props)
	if err != nil {
		return nil, err
	}
	return props, nil
}

func fixVariableConfig(props *properties.Properties) error {
	keys := props.Keys()
	for _, k := range keys {
		v, ok := props.Get(k)
		if ok {
			_, _, err := props.Set(k, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

const configPath = "config"
const fileType = ".properties"
const defaultConfigFile = "default.properties"

const defaultEnv = "local"

var envFlag = flag.String("env", "", "环境变量，默认为local")
var confFlag = flag.String("conf", "", "配置目录")

// GetEnv 获取环境变量
func GetEnv(env string) string {
	if env != "" {
		return env
	}

	flag.Parse()
	if *envFlag != "" {
		return *envFlag
	}

	env = os.Getenv("ENV")
	if env != "" {
		return env
	}
	return defaultEnv
}

func getConfDir() string {
	flag.Parse()
	if *confFlag != "" {
		return *confFlag
	}
	return os.Getenv("CONF")
}

func getExecutableDir() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(dir), nil
}
