package weather

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const WEATHER_API = "https://restapi.amap.com/v3/weather/weatherInfo"

type WeatherResponse struct {
	Status    string     `json:"status"`
	Count     string     `json:"count"`
	Info      string     `json:"info"`
	Infocode  string     `json:"infocode"`
	Forecasts []Forecast `json:"forecasts"`
}
type Forecast struct {
	City       string `json:"city"`
	Adcode     string `json:"adcode"`
	Province   string `json:"province"`
	Reporttime string `json:"reporttime"`
	Casts      []Cast `json:"casts"`
}
type Cast struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	Dayweather   string `json:"dayweather"`
	Nightweather string `json:"nightweather"`
	Daytemp      string `json:"daytemp"`
	Nighttemp    string `json:"nighttemp"`
	Daywind      string `json:"daywind"`
	Nightwind    string `json:"nightwind"`
	Daypower     string `json:"daypower"`
	Nightpower   string `json:"nightpower"`
}

func QueryWeather(url, requestBody string) (body []byte, err error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, strings.NewReader(requestBody))
	if err != nil {
		fmt.Printf("http send request url %s fails -- %v ", url, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("http send request url %s fails -- %v ", url, err)
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	//status code not in [200, 300) fail
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("response status code %d, error messge: %s", resp.StatusCode, string(body))
		return
	}
	if err != nil {
		fmt.Printf("read the result of get url %s fails, response status code %d -- %v", url, resp.StatusCode, err)
	}
	return
}
