package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetWeatherNow(locality string) (string, error) {
	addressAPI := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=ru", locality, os.Getenv("OWM_TOKEN"))
	resp, err := http.Get(addressAPI)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// log.Println(string(body))
	return string(body), err
}
