package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const COEFF = 1.1855

var (
	OsrmURLTruck    string
	OsrmURLGazelle  string
	OsrmURLFiveTonn string
)

type LatLng struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

func main() {
	OsrmURLGazelle = "http://osrm-gazelle.wbdispatch.k8s.prod-xc"
	OsrmURLFiveTonn = "http://osrm-fivetonn.wbdispatch.k8s.prod-xc"
	OsrmURLTruck = "http://osrm-longtruck.wbdispatch.k8s.prod-xc"

	f, err := os.Open("KBT_cost_transports_all.csv")
	if err != nil {
		log.Fatal(err)
	}

	resultFile, _ := os.Create("KBT_cost_transports_res_all.csv")
	writer := csv.NewWriter(resultFile)
	writer.Comma = ','

	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	var (
		record []string
	)
	_, err = csvReader.Read() // read header
	if err != nil {
		log.Fatal(err)
	}

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", rec)

		srcLat, err := strconv.ParseFloat(rec[6], 64)
		if err != nil {
			log.Fatal(err)
		}
		srcLng, err := strconv.ParseFloat(rec[7], 64)
		if err != nil {
			log.Fatal(err)
		}
		dstLat, err := strconv.ParseFloat(rec[9], 64)
		if err != nil {
			log.Fatal(err)
		}
		dstLng, err := strconv.ParseFloat(rec[10], 64)
		if err != nil {
			log.Fatal(err)
		}

		time_truck1, dist_truck1, err := gettime(OsrmURLGazelle, [][]float64{{srcLng, srcLat}, {dstLng, dstLat}})
		//		time.Sleep(10 * time.Millisecond)
		// if err != nil {
		// 	//			panic(err)
		// }

		// time_truck2, dist_truck2, err := gettime(OsrmURLFiveTonn, [][]float64{{srcLng, srcLat}, {dstLng, dstLat}})
		// if err != nil {
		// 	panic(err)
		// }

		// time_truck3, dist_truck3, err := gettime(OsrmURLTruck, [][]float64{{srcLng, srcLat}, {dstLng, dstLat}})
		// if err != nil {
		// 	panic(err)
		// }

		time_truck1 *= COEFF
		hours_in_route1 := time_truck1 / (60.0 * 60.0)
		solid1 := int(hours_in_route1 / 8)
		yy1 := solid1 / 2
		yyy1 := solid1%2 + yy1
		rest_time1 := float64(yy1)*8.0 + float64(yyy1)*2.0

		// time_truck2 *= COEFF
		// hours_in_route2 := time_truck2 / (60.0 * 60.0)
		// solid2 := int(hours_in_route2 / 8)
		// yy2 := solid2 / 2
		// yyy2 := solid2%2 + yy2
		// rest_time2 := float64(yy2)*8.0 + float64(yyy2)*2.0

		// time_truck3 *= COEFF
		// hours_in_route3 := time_truck3 / (60.0 * 60.0)
		// solid3 := int(hours_in_route3 / 8)
		// yy3 := solid3 / 2
		// yyy3 := solid1%2 + yy3
		// rest_time3 := float64(yy3)*8.0 + float64(yyy3)*2.0

		// record = append(rec, fmt.Sprintf("%f", time_truck1), fmt.Sprintf("%f hours in route", hours_in_route1), fmt.Sprintf("%f", dist_truck1),
		// 	fmt.Sprintf("%f km/h", (dist_truck1/1000.0)/(time_truck1/(60*60))),
		// 	fmt.Sprintf("rest: %f", rest_time1), fmt.Sprintf("total time: %f ", rest_time1+hours_in_route1),
		// 	fmt.Sprintf("%f", time_truck2), fmt.Sprintf("%f hours in route", hours_in_route2), fmt.Sprintf("%f", dist_truck2),
		// 	fmt.Sprintf("%f km/h", (dist_truck2/1000.0)/(time_truck2/(60*60))),
		// 	fmt.Sprintf("rest: %f", rest_time2), fmt.Sprintf("total time: %f ", rest_time2+hours_in_route2),
		// 	fmt.Sprintf("%f", time_truck3), fmt.Sprintf("%f hours in route", hours_in_route3), fmt.Sprintf("%f", dist_truck3),
		// 	fmt.Sprintf("%f km/h", (dist_truck3/1000.0)/(time_truck3/(60*60))),
		// 	fmt.Sprintf("rest: %f", rest_time3), fmt.Sprintf("total time: %f ", rest_time3+hours_in_route3))

		record = append(rec,
			fmt.Sprintf("%f", time_truck1),
			fmt.Sprintf("%f hours in route", hours_in_route1), fmt.Sprintf("%f", dist_truck1),
			fmt.Sprintf("%f km/h", (dist_truck1/1000.0)/(time_truck1/(60*60))),
			fmt.Sprintf("rest: %f", rest_time1),
			fmt.Sprintf("total time: %f ", rest_time1+hours_in_route1))
		writer.Write(record)
	}

	writer.Flush()

}

type OSRMApiRouteResp struct {
	Code   string `json:"code"`
	Routes []struct {
		Distance   float64 `json:"distance"`
		Duration   float64 `json:"duration"`
		WeightName string  `json:"weight_name"`
		Weight     float64 `json:"weight"`
	} `json:"routes"`
}

func gettime(ourl string, points [][]float64) (float64, float64, error) {
	qsParts := make([]string, 0)
	for _, p := range points {
		qsParts = append(qsParts, fmt.Sprintf("%f,%f", p[0], p[1]))
	}
	str := fmt.Sprintf("%s/route/v1/driving/%s?geometries=geojson", ourl, strings.Join(qsParts, ";"))

	resp, err := http.Get(str)
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}
	apiResp := OSRMApiRouteResp{}
	if json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, 0, fmt.Errorf("can't unmarshal osrm api resp: %s", err)
	}
	if len(apiResp.Routes) == 0 {
		fmt.Println(ourl, qsParts)
		return 0, 0, fmt.Errorf("not route")
	}

	return apiResp.Routes[0].Duration, apiResp.Routes[0].Distance, nil
}
