package config

import (
	"io/ioutil"
	"log"
	"os"
	"gopkg.in/yaml.v2"
)

// 使用者根据自己需要修改这个结构体
var C struct {
	RunMode string  `yaml:"runMode"`
	Mysql struct {
		Bamboo   string `yaml:"bamboo_website"`
		DsnRead  string `yaml:"dsn_read"`
		DsnWrite string `yaml:"dsn_wirte"`
	} `yaml:"mysql"`
	Redis  struct{

	}
	HttpAddr  string `yaml:"http_addr"`
}

func init() {
	configFileName := "config.yml"
	configFileNameBack := "config_debug.yml"
	var findConfig, findConfigBack bool
	// 向上层查找配置文件
	// 在项目的任何地方运行(test时)都能加载到配置文件
	// 有备配置文件, 当主配置文件没找到时就会使用备配置文件
	// 优先使用最近的主备配置文件
	for i := 0; i < 10; i++ {
		_, err := os.Stat(configFileName)
		if err != nil {
			if os.IsNotExist(err) {
				configFileName = "../" + configFileName
			} else {
				panic(err)
			}
		} else {
			findConfig = true
			break
		}

		_, err = os.Stat(configFileNameBack)
		if err != nil {
			if os.IsNotExist(err) {
				configFileNameBack = "../" + configFileNameBack
			} else {
				panic(err)
			}
		} else {
			findConfigBack = true
			break
		}
	}

	var fileName string
	if findConfig {
		fileName = configFileName
	} else if findConfigBack {
		fileName = configFileNameBack
	} else {
		log.Panicf("can't find 'config.yml' or 'config_debug.yml'")
		return
	}

	log.Printf("found config file: %s", fileName)

	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicf("can't read 'config.yml'")
		return
	}

	err = yaml.Unmarshal(bs, &C)
	if err != nil {
		log.Panicf("yaml.Unmarshal err:%v; row:%s", err, bs)
		return
	}
}
