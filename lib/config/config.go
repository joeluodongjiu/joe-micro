package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// 使用者根据自己需要修改这个结构体
var C struct {
	Debug bool `yaml:"debug"`

	Service struct {
		Name    string `yaml:"name"`
		Port    string `yaml:"port"`
		Version string `yaml:"version"`
	}

	Jwt  struct{
		SignKey string `yaml:"signKey"`
		TimeOut  time.Duration   `yaml:"timeOut"`
	}   `yaml:"jwt"`


	Mysql struct {
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
		DbName   string `yaml:"db_name"`
	} `yaml:"mysql"`

	Redis struct {
		Key  string `yaml:"key"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Auth string `yaml:"auth"`
		Db   int    `yaml:"db"`
	} `yaml:"redis"`

	Consul string `yaml:"consul"`

	Jaeger string `yaml:"jaeger"`

	Nsq struct {
		Address     string `yaml:"address"`
		Lookup      string `yaml:"lookup"`
		MaxInFlight int    `yaml:"maxInFlight"`
	} `yaml:"nsq"`
}

func init() {
	configFileName := "config.yaml"
	var findConfig bool

	for i := 0; i < 10; i++ {
		_, err := os.Stat(configFileName)
		if err != nil {
			if os.IsNotExist(err) {
				configFileName = "./" + configFileName
			} else {
				panic(err)
			}
		} else {
			findConfig = true
			break
		}
	}

	var fileName string
	if findConfig {
		fileName = configFileName
	} else {
		log.Panicf("can't find 'config.yml' ")
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
