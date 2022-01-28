package manager

import (
	"encoding/json"
	"io/ioutil"
)

type MstConfig struct {
	ApiPort               int      `json:"apiPort"`
	ApiReadTimeout        int      `json:"apiReadTimeout"`
	ApiWriteTimeout       int      `json:"apiWriteTimeout"`
	EtcdEndpoints         []string `json:"etcdEndpoints"`
	EtcdDialTimeout       int      `json:"etcdDialTimeout"`
	WebRoot               string   `json:"webroot"`
	MongodbUri            string   `json:"mongodbUri"`
	MongodbConnectTimeout int      `json:"mongodbConnectTimeout"`
}

var (
	// Config singleton config
	Config *MstConfig
)

// InitConfig load config from file
func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    MstConfig
	)

	// 1. read config file
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// 2. json decode
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	// 3. assign to singleton
	Config = &conf

	return
}
