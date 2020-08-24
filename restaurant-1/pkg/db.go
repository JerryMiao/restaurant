package pkg

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

var (
	//Ddb mysql connection
	Ddb *gorm.DB
	//Rds connection
	Rds *redis.Client
	//WebConfig basic config
	WebConfig *Config
)

//Config 配置内容
type Config struct {
	Driver        string //database
	ConnectString string //sql connection
	Port          string //listen port
	IsLog         bool   //whether to print sql statement
	RedisAdree    string //redis connection address
	RedisPwd      string //redis connection password
	Redislevel    int    //redis level
	RedisPort     string //redis port number
}

//readConfig read local config, load config path，return config instance
func readConfig(filename string) (*Config, error) {
	config := &Config{}
	bys, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bys, config); err != nil {
		return nil, err
	}
	return config, nil
}

//SQLOpen Return a database connection
//isLog whether to display the sql statement
func sqlOpen() *gorm.DB {
	var err error
	db, err := gorm.Open(WebConfig.Driver, WebConfig.ConnectString)
	if err != nil {
		panic(err)
	}
	db.LogMode(WebConfig.IsLog)
	return db
}

func redisInit() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: WebConfig.RedisAdree + ":" + WebConfig.RedisPort, Password: WebConfig.RedisPwd, DB: WebConfig.Redislevel})
}

//SQLInit sql init
func SQLInit() {
	cng, err := readConfig("./config.json")
	if err != nil {
		return
	}
	WebConfig = cng
	Ddb = sqlOpen()
	Rds = redisInit()
}
