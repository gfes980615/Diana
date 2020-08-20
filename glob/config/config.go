package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config ...
var (
	Config *ConfigSetup
)

func init() {
	InitConfig()
}

// LoadConfig ...
func LoadConfig(file string) {
	if Config != nil {
		return
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return
	}

	Config = new(ConfigSetup)

	viper.SetConfigType("yaml")
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}

	viper.Unmarshal(&Config)
}

// InitConfig ...
func InitConfig() {
	LoadConfig("app.yaml")
}

// ConfigSetup
type ConfigSetup struct {
	LogConfig LogConfig `yaml:"LogConfig"`
	// NotifyConfig   NotifyConfig   `yaml:"NotifyConfig"`
	GINConfig      GINConfig      `yaml:"GINConfig"`
	DatabaseConfig DatabaseConfig `yaml:"DatabaseConfig"`
	// ConsumerConfig ConsumerConfig `yaml:"ConsumerConfig"`
}

// // ConsumerConfig ...
// type ConsumerConfig struct {
// 	Brokers   []string `yaml:"Brokers"`
// 	Increment string   `yaml:"Increment"`
// }

// // NotifyConfig
// type NotifyConfig struct {
// 	Slack Slack `yaml:"Slack"`
// }

// // Slack
// type Slack struct {
// 	Channel string `yaml:"Channel"`
// 	Hook    bool   `yaml:"Hook"`
// 	API     string `yaml:"API"`
// }

// DatabaseConfig
type DatabaseConfig struct {
	Mysql Mysql `yaml:"Mysql"`
	// Mongo Mongo `yaml:"Mongo"`
}

// DataBases
type DataBases struct {
	Database string `yaml:"Database"`
	Enable   bool   `yaml:"Enable"`
	Name     string `yaml:"Name"`
	Address  string `yaml:"Address"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

// LogConfig
type LogConfig struct {
	HistoryPath   string `yaml:"HistoryPath"`
	FullColor     bool   `yaml:"FullColor"`
	FullTimestamp bool   `yaml:"FullTimestamp"`
	Name          string `yaml:"Name"`
	Env           string `yaml:"Env"`
	Level         string `yaml:"Level"`
	Duration      string `yaml:"Duration"`
}

// GINConfig
type GINConfig struct {
	Address string `yaml:"Address"`
}

// Mysql
type Mysql struct {
	MaxIdle        int       `yaml:"MaxIdle"`
	MaxOpen        int       `yaml:"MaxOpen"`
	ConnMaxLifeMin int       `yaml:"ConnMaxLifeMin"`
	DataBases      DataBases `yaml:"DataBases"`
	LogMode        bool      `yaml:"LogMode"`
}

// // Mongo
// type Mongo struct {
// 	DataBases []DataBases `yaml:"DataBases"`
// 	MaxIdle   int         `yaml:"MaxIdle"`
// 	MaxOpen   int         `yaml:"MaxOpen"`
// }
