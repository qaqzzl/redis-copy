package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	_redis "github.com/qaqzzl/redis-copy/library/cache/redis"
	"github.com/spf13/viper"
)

var SourceRds *redis.Pool
var TargetRds *redis.Pool

func main() {

	SourceRds = _redis.NewPool(_redis.Config{
		Dial: viper.GetString("source.address"),
		Auth: viper.GetString("source.password"),
	})
	connSource := SourceRds.Get()
	defer SourceRds.Close()
	connSource.Do("SELECT", 6)

	TargetRds = _redis.NewPool(_redis.Config{
		Dial: viper.GetString("target.address"),
		Auth: viper.GetString("target.password"),
	})
	connTarget := TargetRds.Get()
	defer TargetRds.Close()
	connTarget.Do("SELECT", 6)

	result, _ := connSource.Do("KEYS", "*")
	resultarr := result.([]interface{})

	for _, value := range resultarr {
		key := string(value.([]uint8))
		_type, _ := connSource.Do("TYPE", key)
		fmt.Println(key)
		fmt.Println(_type)
		switch _type {
		case "hash":
			hvall, _ := connSource.Do("HGETALL", key)
			hvallarr := hvall.([]interface{})
			for hk, hv := range hvallarr {
				if (hk % 2) == 0 {
					connTarget.Do("HMSET", key, string(hv.([]uint8)), string(hvallarr[hk+1].([]uint8)))
				}
			}
		case "string":
			fmt.Println("string")
		default:
			fmt.Println("不支持的数据类型", _type)

		}
	}

}

func init() {

	viper.SetConfigName("config")            // 配置文件名
	viper.SetConfigType("yaml")              // 配置文件类型，可以是yaml、json、xml。。。
	viper.AddConfigPath("G:\\redis-copy\\.") // 配置文件路径

	err := viper.ReadInConfig() // 读取配置文件信息
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
