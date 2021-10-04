package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetCoordLocality(locality string) (string, error) {
	addressAPI := fmt.Sprintf("https://api.geotree.ru/search.php?key=%s&level=4&term=%s", locality, os.Getenv("GEOTREE_TOKEN"))
	resp, err := http.Get(addressAPI)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	return string(body), err
}
