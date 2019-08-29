package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type ServerConfig struct {
	ServiceName       string `yaml:"service_name"`
	DispatchWorkerNum int    `yaml:"dispatch_worker_num"`
	SwitchWorkerNum   int    `yaml:"switch_worker_num"`
	JobWorkerNum      int    `yaml:"job_worker_num"`
	LogFile           string `yaml:"log_file"`
}

type RequestConfig struct {
	RequestUrl string `yaml:"RequestUrl"`
}

type EtcdConfig struct {
	DialTimeout    int      `yaml:"dial_timeout"`
	RequestTimeout int      `yaml:"request_timeout"`
	EndPoints      []string `yaml:"end_points"`
	Username       string   `yaml:"username"`
	Password       string   `yaml:"password"`
}

type RedisConfig struct {
	Redis [][]string `yaml:"Redis,flow"`
}

type Config struct {
	Server        ServerConfig `yaml:"Server,flow"`
	RequestConfig `yaml:",inline"`
	Etcd          EtcdConfig `yaml:"Etcd,flow"`
	RedisConfig   `yaml:",inline"`
}

var conf *Config

func LoadConfig(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("ioutil.ReadFile err:", err)
		os.Exit(0)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Println("yaml.Unmarshal err:", err)
		os.Exit(0)
	}

	//log.Println("a:", conf)
}

func GetServerConfig() ServerConfig {
	return conf.Server
}

func GetLogFile() string {
	return conf.Server.LogFile
}

func GetRequestConfig() string {
	return conf.RequestUrl
}

func GetEtcdConfig() EtcdConfig {
	return conf.Etcd
}

func GetRedisConfig() [][]string {
	return conf.Redis
}
