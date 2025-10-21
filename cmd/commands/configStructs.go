package commands

import (
	"reflect"
)

// Config 配置项结构体
type Config struct {
	Name         string
	ValueType    interface{}
	DefaultValue interface{}
	Description  string
}

// getDefaults 定义默认的配置项以及默认值
func getDefaults() map[string]interface{} {
	configs := make([]Config, 0)
	configs = append(configs,
		Config{
			Name:         "address",
			ValueType:    reflect.String,
			DefaultValue: "127.0.0.1",
			Description:  "Redis server address",
		},
		Config{
			Name:         "port",
			ValueType:    reflect.Int,
			DefaultValue: 6379,
			Description:  "Redis server port",
		},
		Config{
			Name:         "password",
			ValueType:    reflect.String,
			DefaultValue: "",
			Description:  "Redis server password",
		},
		Config{
			Name:         "pattern",
			ValueType:    reflect.String,
			DefaultValue: "*",
			Description:  "key name pattern",
		},
		Config{
			Name:         "prefer_master",
			ValueType:    reflect.Bool,
			DefaultValue: false,
			Description:  "if prefer the redis master node",
		},
		Config{
			Name:         "dial_timeout",
			ValueType:    reflect.Int,
			DefaultValue: 10,
			Description:  "connect timeout",
		})

	defaults := make(map[string]interface{})
	for _, c := range configs {
		defaults[c.Name] = c.DefaultValue
	}
	return defaults
}
