package controller

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/golang/glog"
	"github.com/knative-sample/weather-store/pkg/tablestore"
	"github.com/knative-sample/weather-store/pkg/weather"
)

func StoreWeather() {
	// 北京市、杭州市
	cityCodes := []string{"110000", "330100"}
	client := tablestore.InitClient()
	for _, cityCode := range cityCodes {
		glog.Infof("start query city: %s", cityCode)
		key := os.Getenv("WEATHER_API_KEY")
		queryApi := fmt.Sprintf("%s?key=%s&city=%s&extensions=all", weather.WEATHER_API, key, cityCode)
		res, err := weather.QueryWeather(queryApi, "")
		if err != nil {
			glog.Errorf("QueryWeather error: %s", err.Error())
			continue
		}
		glog.Infof("weather info: %s", res)
		wr := weather.WeatherResponse{}
		err = json.Unmarshal(res, &wr)
		if err != nil {
			glog.Errorf("QueryWeather Unmarshal error: %s", err.Error())
			continue
		}
		if wr.Status == "1" && len(wr.Forecasts) > 0 {
			glog.Infof("start store city: %s", cityCode)
			err := client.Store(wr.Forecasts[0])
			if err != nil {
				glog.Errorf("Weather Store error: %s", err.Error())
				continue
			}
		}
	}
}
