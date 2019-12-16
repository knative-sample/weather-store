package tablestore

import (
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/knative-sample/weather-store/pkg/weather"
)

func FakeInitClient() *TableClient {
	endpoint := "https://xxx.cn-beijing.ots.aliyuncs.com" //实例访问地址
	tableName := "weather"
	instanceName := "knative-weather"
	accessKeyId := "xx"
	accessKeySecret := "xx"
	client := tablestore.NewClient(endpoint, instanceName, accessKeyId, accessKeySecret)
	return &TableClient{
		tableName: tableName,
		client:    client,
	}
}
func TestStore(t *testing.T) {
	c := FakeInitClient()
	ct := weather.Cast{
		Date:         "2019-09-24",
		Week:         "2",
		Dayweather:   "晴",
		Nightweather: "晴",
		Daytemp:      "31",
		Nighttemp:    "15",
		Daywind:      "南",
		Nightwind:    "南",
		Daypower:     "≤3",
		Nightpower:   "≤3",
	}

	f := weather.Forecast{
		City:       "北京市",
		Adcode:     "110000",
		Reporttime: "2019-09-24 20:50:56",
		Province:   "北京",
		Casts:      []weather.Cast{ct},
	}
	c.Store(f)
}
