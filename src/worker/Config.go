package worker

import (
	"encoding/json"
	"io/ioutil"
)

type WkConfig struct {
	EtcdEndpoints         []string `json:"etcdEndpoints"`
	EtcdDialTimeout       int      `json:"etcdDialTimeout"`
	MongodbUri            string   `json:"mongodbUri"`
	MongodbConnectTimeout int      `json:"mongodbConnectTimeout"`
	JobLogBatchSize       int      `json:"jobLogBatchSize"`
	JobLogCommitTimeout   int      `json:"jobLogCommitTimeout"`
	BashPath              string   `json:"bashPath"`
	WinBashPath           string   `json:"winBashPath"`
}

var (
	// Config singleton config
	Config *WkConfig
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    WkConfig
	)

	// 1, read config file
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// 2, json decode
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	// 3, assign to singleton
	Config = &conf

	return
}
