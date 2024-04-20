package boot

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	g "main/app/global"
	"os"
)

// 定义配置环境和配置文件的常量
const (
	configEnv  = "CONFIG_PATH"
	configFile = "manifest/config/config.yaml"
)

// ViperSetup 函数设置配置，...接收一个可变数量的字符串参数
func ViperSetup(path ...string) {
	var configPath string

	// 获取配置文件路径，优先级：参数 > 命令行 > 环境变量 > 默认值
	// get config file path
	// priority: param > command line > environment > default
	if len(path) != 0 {
		// param
		//如果提供了路径参数，那么就使用第一个参数作为配置文件的路径
		configPath = path[0]
	} else {
		// command line
		//如果没有提供路径参数，那么就尝试从命令行参数中获取配置文件的路径。这里，"c"是命令行参数的名称
		flag.StringVar(&configPath, "c", "", "set config path")
		flag.Parse()

		if configPath == "" {
			//如果命令行参数中没有提供配置文件的路径，那么就尝试从环境变量中获取配置文件的路径
			if configPath = os.Getenv(configEnv); configPath != "" {
				// environment
			} else {
				// default
				//如果环境变量中没有提供配置文件的路径，那么就使用默认的配置文件路径
				configPath = configFile
			}
		}
	}
	fmt.Printf("get config path: %s\n", configPath)

	v := viper.New()            // 创建一个新的viper实例
	v.SetConfigFile(configPath) // 设置配置文件路径
	v.SetConfigType("yaml")     // 设置配置文件类型为yaml
	err := v.ReadInConfig()     // 读取配置文件
	if err != nil {
		// 如果读取配置文件失败，抛出panic
		panic(fmt.Errorf("get config file failed, err: %v", err))
	}

	if err := v.Unmarshal(&g.Config); err != nil {
		// 如果解析配置失败，抛出panic
		panic(fmt.Errorf("unmarshal config failed, err: %v", err))
	}
}
