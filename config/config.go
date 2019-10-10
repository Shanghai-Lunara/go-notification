package config

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type DispatchConfig struct {
	ServiceName       string `yaml:"service_name"`
	LogFile           string `yaml:"log_file"`
	IP                string `yaml:"ip"`
	InternalIP        string `yaml:"internal_ip"`
	Port              int    `yaml:"port"`
	HeartBeatInternal int    `yaml:"heart_beat_internal"`
}

type LogicConfig struct {
	ServiceName string `yaml:"service_name"`
	LogFile     string `yaml:"log_file"`
	WorkerNum   int    `yaml:"worker_num"`
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
	Dispatch       DispatchConfig `yaml:"Dispatch,flow"`
	Logic          LogicConfig    `yaml:"Logic,flow"`
	HttpRequestApi string         `yaml:"HttpRequestAPI"`
	Etcd           EtcdConfig     `yaml:"Etcd,flow"`
	RedisConfig    `yaml:",inline"`
	ConfigPath     string
}

var (
	conf       *Config
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configPath", "push.yml", "configuration file path")
}

func Init() (err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile(configPath); err != nil {
		return errors.New(fmt.Sprintf("ioutil.ReadFile err:%v", err))
	}
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return errors.New(fmt.Sprintf("yaml.Unmarshal err:%v", err))
	}
	conf.ConfigPath = configPath
	//log.Println("a:", conf)
	return nil
}

func GetConfigPath() string {
	return conf.ConfigPath
}

func GetConfig() *Config {
	return conf
}

func GetEtcdConfig() EtcdConfig {
	return conf.Etcd
}

func GetDispatchLogFile() string {
	return conf.Dispatch.LogFile
}

func GetLogicLogFile() string {
	return conf.Logic.LogFile
}

func (c *Config) GetRedisConfig() [][]string {
	return c.Redis
}

func GetDispatchListenAddr() string {
	return fmt.Sprintf("%s:%d", conf.Dispatch.IP, conf.Dispatch.Port)
}

func GetDispatcHeartBeatInternal() int {
	return conf.Dispatch.HeartBeatInternal
}
