package main

import (
	"apitools"
	"fmt"
	"net/http"
	"strings"

	"github.com/lolerwaffles/GoAPITools"
)

var headersmap = map[string]string{
	"x-rapidapi-host": "national-weather-service.p.rapidapi.com",
	"HeaderType":      "",
}

func main() {
	go http.HandleFunc("/", returnWeatherData)
	http.ListenAndServe(":8086", nil)
}

func returnWeatherData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ip := ReadUserIP(r)
	fmt.Println("Accepted connection from " + ip)
	requestType := r.URL.Path
	wURL := "https://national-weather-service.p.rapidapi.com/stations/" + getStation(ip) + "/observations/current"
	switch requestType {
	case "/getDaily":
		wURL = "https://national-weather-service.p.rapidapi.com/points/" + getLoc(ip) + "/forecast"
	case "/getHourly":
		wURL = "https://national-weather-service.p.rapidapi.com/stations/" + getStation(ip) + "/observations/current"
	}
	w.Write([]byte(GoAPITools.CallAPIReturnString(wURL, headersmap)))
}

func getStation(ip string) string {
	wURL := "https://national-weather-service.p.rapidapi.com/points/" + getLoc(ip) + "/stations"
	return apitools.CallAPIReturnString(wURL, headersmap)
}

func getLoc(ip string) string {
	wURL := "https://ipapi.co/" + ip + "/latlong/"
	return apitools.CallAPIReturnString(wURL, map[string]string{})

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
