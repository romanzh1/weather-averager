package api

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Coordinates struct {
	Lat float64
	Lon float64
}

func GetCoordLocality(locality string) (Coordinates, error) {
	addressAPI := fmt.Sprintf("https://api.geotree.ru/search.php?key=%s&level=4&term=%s", os.Getenv("GEOTREE_TOKEN"), locality)
	var coordinates Coordinates

	req, err := http.NewRequest("GET", addressAPI, nil)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // TODO write to GeoTree that the certificate is not valid
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	resp.Header.Set("Referer", os.Getenv("GEOTREE_EMAIL")) // TODO check
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	coordStr := strings.Split(strings.Split(strings.Split(string(body), `"geo_inside": {`)[1], `},`)[0], ",")
	re := regexp.MustCompile(`[0-9]+|\.`)

	coordinates.Lon, err = strconv.ParseFloat(strings.Join(re.FindAllString(coordStr[0], -1), ""), 64)
	if err != nil {
		fmt.Println(err)
	}
	coordinates.Lat, err = strconv.ParseFloat(strings.Join(re.FindAllString(coordStr[1], -1), ""), 64)
	if err != nil {
		fmt.Println(err)
	}

	return coordinates, err
}
