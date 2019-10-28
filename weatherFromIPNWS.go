package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

var ApiKey string = ""

func main() {
	go http.HandleFunc("/", returnWeatherData)
	http.ListenAndServe(":8086", nil)
}

func returnWeatherData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(getWeatherData(ReadUserIP(r), r.URL.Path)))
}

func getWeatherData(ip string, requestType string) string {
	fmt.Println("Accepted connection from " + ip)
	wURL := "https://national-weather-service.p.rapidapi.com/stations/" + getStation(ip) + "/observations/current"
	switch requestType {
	case "/getDaily":
		wURL = "https://national-weather-service.p.rapidapi.com/points/" + getLoc(ip) + "/forecast"
	case "/getHourly":
		wURL = "https://national-weather-service.p.rapidapi.com/stations/" + getStation(ip) + "/observations/current"
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", wURL, nil)

	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("x-rapidapi-host", "national-weather-service.p.rapidapi.com")
	request.Header.Set("x-rapidapi-key", ApiKey)
	resp, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

func getStation(ip string) string {
	wURL := "https://national-weather-service.p.rapidapi.com/points/" + getLoc(ip) + "/stations"
	client := &http.Client{}
	request, err := http.NewRequest("GET", wURL, nil)

	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("x-rapidapi-host", "national-weather-service.p.rapidapi.com")
	request.Header.Set("x-rapidapi-key", ApiKey)
	resp, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	myJson := string(body)

	value := gjson.Get(myJson, "features.0.properties.stationIdentifier")
	println(value.String())
	return value.String()
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	if idx := strings.IndexByte(IPAddress, ':'); idx >= 0 {
		IPAddress = IPAddress[:idx]
	}
	return IPAddress
}

func getLoc(ip string) string {
	url := "https://ipapi.co/" + ip + "/latlong/"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}
