package main

// import (
// 	"encoding/csv"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// )

// const COEFF = 1.1855

// var (
// 	OsrmURLTruck    string
// 	OsrmURLGazelle  string
// 	OsrmURLFiveTonn string
// )

// type LatLng struct {
// 	Lat float64 `json:"latitude"`
// 	Lng float64 `json:"longitude"`
// }

// func main() {
// 	OsrmURLGazelle = "http://wbdispatch-ingress-controller.wbdispatch.k8s.dev-el/osrm-gazelle"
// 	OsrmURLFiveTonn = "http://wbdispatch-ingress-controller.wbdispatch.k8s.dev-el/osrm-fivetonn"
// 	OsrmURLTruck = "http://wbdispatch-ingress-controller.wbdispatch.k8s.dev-el/osrm-longtruck"

// 	f, err := os.Open("Routes_MC_cartype.csv")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	resultFile, _ := os.Create("result_source_cartype.csv")
// 	writer := csv.NewWriter(resultFile)
// 	writer.Comma = ';'

// 	defer f.Close()

// 	// read csv values using csv.Reader
// 	csvReader := csv.NewReader(f)
// 	csvReader.Comma = ','
// 	var (
// 		OsrmURL string
// 		record  []string
// 		srcLL   LatLng
// 		dstLL   LatLng
// 	)
// 	_, err = csvReader.Read() // read header
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for {
// 		rec, err := csvReader.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// do something with read line
// 		fmt.Printf("%+v\n", rec)

// 		err = json.Unmarshal([]byte(rec[3]), &srcLL)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		err = json.Unmarshal([]byte(rec[5]), &dstLL)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if rec[1] == "1" { // gazelle
// 			OsrmURL = OsrmURLGazelle
// 		} else if rec[1] == "2" { // fivetonn
// 			OsrmURL = OsrmURLFiveTonn
// 		} else {
// 			OsrmURL = OsrmURLTruck
// 		}

// 		time_truck, dist_truck, err := gettime(OsrmURL, [][]float64{{srcLL.Lng, srcLL.Lat}, {dstLL.Lng, dstLL.Lat}})
// 		if err != nil {
// 			panic(err)
// 		}

// 		time_truck *= COEFF
// 		hours_in_route := time_truck / (60.0 * 60.0)
// 		solid := int(hours_in_route / 8)
// 		yy := solid / 2
// 		yyy := solid%2 + yy
// 		rest_time := float64(yy)*8.0 + float64(yyy)*2.0

// 		record = append(rec, fmt.Sprintf("%f", time_truck), fmt.Sprintf("%f hours in route", hours_in_route), fmt.Sprintf("%f", dist_truck),
// 			fmt.Sprintf("%f km/h", (dist_truck/1000.0)/(time_truck/(60*60))),
// 			fmt.Sprintf("rest: %f", rest_time), fmt.Sprintf("total time: %f ", rest_time+hours_in_route), OsrmURL)

// 		writer.Write(record)
// 	}

// 	writer.Flush()

// }

// type OSRMApiRouteResp struct {
// 	Code   string `json:"code"`
// 	Routes []struct {
// 		Distance   float64 `json:"distance"`
// 		Duration   float64 `json:"duration"`
// 		WeightName string  `json:"weight_name"`
// 		Weight     float64 `json:"weight"`
// 	} `json:"routes"`
// }

// func gettime(ourl string, points [][]float64) (float64, float64, error) {
// 	qsParts := make([]string, 0)
// 	for _, p := range points {
// 		qsParts = append(qsParts, fmt.Sprintf("%f,%f", p[0], p[1]))
// 	}

// 	resp, err := http.Get(fmt.Sprintf("%s/route/v1/driving/%s?geometries=geojson", ourl, strings.Join(qsParts, ";")))
// 	if err != nil {
// 		fmt.Println(err)
// 		panic("!!!")
// 	}
// 	apiResp := OSRMApiRouteResp{}
// 	if json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
// 		return 0, 0, fmt.Errorf("can't unmarshal osrm api resp: %s", err)
// 	}
// 	if len(apiResp.Routes) == 0 {
// 		fmt.Println(ourl, qsParts)
// 		return 0, 0, fmt.Errorf("not route")
// 	}

// 	return apiResp.Routes[0].Duration, apiResp.Routes[0].Distance, nil
// }
