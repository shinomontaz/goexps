package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	apiName    = "erofeev_api"
	apiVersion = "erofeev_api.1"
)

func main() {
	prepare("addresses_fin.csv", "geocoded_final.csv")
}

type LatLong struct {
	Lat  float64
	Long float64
}

type Order struct {
	Coords  *LatLong
	Address string
	Guid    string
}

type Task struct {
	Hub    *Order
	Points []*Order
}

type YandexJson struct {
	Response struct {
		GeoObjectCollection struct {
			FeatureMember []struct {
				GeoObject struct {
					MetaDataProperty struct {
						GeocoderMetaData struct {
							Precision string `json:"precision"`
						}
					}
					Point struct {
						Pos string `json:"pos"`
					}
				}
			} `json:"featureMember"`
		}
	} `json:"response"`
}

func Geocode(address string, useProxy bool) (x, y float64, err error) {
	query := fmt.Sprintf("https://geocode-maps.yandex.ru/1.x/?format=json&geocode=%s", address)

	timestart := time.Now()

	var resp *http.Response
	if rand.Float64() < 0.4 {
		time.Sleep(1 * time.Second)
	}
	if useProxy {
		resp, err = http.PostForm("http://localhost:6789", url.Values{"destination": {query}})
	} else {
		resp, err = http.Get(query)
	}

	if err != nil && !useProxy {
		return x, y, fmt.Errorf("Bad request")
	} else if err != nil && useProxy {
		fmt.Println("failed to geocode throw proxy, use plain access")
		return Geocode(address, false)
	}

	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		if useProxy {
			fmt.Println("failed to geocode throw proxy, use plain access")
			return Geocode(address, false)
		}
		return x, y, err
	}

	var yJson YandexJson
	err = json.Unmarshal(result, &yJson)
	if err != nil {
		if useProxy {
			fmt.Println("failed to geocode throw proxy, use plain access")
			return Geocode(address, false)
		}
		return x, y, err
	}

	if len(yJson.Response.GeoObjectCollection.FeatureMember) < 1 {
		if useProxy {
			fmt.Println("failed to geocode throw proxy, use plain access")
			return Geocode(address, false)
		}
		return x, y, fmt.Errorf("empty geocode")
	}

	res := strings.Split(yJson.Response.GeoObjectCollection.FeatureMember[0].GeoObject.Point.Pos, " ")
	x, _ = strconv.ParseFloat(res[0], 64)
	y, _ = strconv.ParseFloat(res[1], 64)

	timeend := time.Now()
	fmt.Println("geocode takes:", timeend.Unix()-timestart.Unix(), "seconds")

	return x, y, nil
}
