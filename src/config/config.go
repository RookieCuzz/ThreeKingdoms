package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

const configFileUrl = "conf/config.yaml" // 相对路径

// 导出的结构体，首字母大写
type ConfigStruct struct {
	LoginServer  loginServerStruct  `yaml:"loginServer"`
	MySqlSection mysqlSectionStruct `yaml:"mysql"`
	WebServer    webServerStruct    `yaml:"webServer"`
}

type webServerStruct struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}
type mysqlSectionStruct struct {
	Dsn             string        `yaml:"dsn"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	MaxConnLifetime time.Duration `yaml:"maxConnLifetime"`
}

type loginServerStruct struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

// 导出的配置变量
var Config ConfigStruct

func Init() {

	var configPath string
	lenth := len(os.Args)
	if lenth > 1 {
		tempDir := os.Args[1]
		if tempDir != "" {
			configPath = filepath.Join(tempDir, configFileUrl)
		}
	} else {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		// 使用 filepath.Join 来构建路径，以确保跨平台兼容
		configPath = filepath.Join(dir, configFileUrl)
	}

	file, err := os.Open(configPath)
	fmt.Println(configPath)
	if err != nil {
		log.Fatal("无法打开配置文件:", err)

	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("关闭文件时出错:", err)
		}
	}(file)

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("解码错误:", err)
		return
	}
	// 输出配置项
	log.Printf("Server Port: %d\n", Config.LoginServer.Port)
	log.Printf("Server Host: %s\n", Config.LoginServer.Host)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil && !os.IsExist(err)
}
