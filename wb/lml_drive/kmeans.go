package main

import (
	"fmt"
	"lml-drive/rand"
	"lml-drive/types"
)

type Kmeans struct {
	deltaThreshold     float64
	iterationThreshold int
}

func NewKmeans() Kmeans {
	return Kmeans{
		deltaThreshold:     0.001,
		iterationThreshold: 1000,
	}
}

type Cluster struct {
	pts    []int
	center types.Point
}

func (m Kmeans) Do(pts []types.Point, k int) ([][]int, error) {
	n := len(pts)

	if k > n || n == 0 || k == 0 {
		return [][]int{}, fmt.Errorf("the size of the data set must at least equal k")
	}

	clusters := make([]Cluster, k)
	rndIndex := 0
	for i := 0; i < k; i++ { // init k clusters
		rndIndex = (int(float64((n-rndIndex)/k)*rand.Float64()) + rndIndex) % n
		clusters[i] = Cluster{
			center: types.Point{Lat: pts[rndIndex].Lat, Lng: pts[rndIndex].Lng},
			pts:    []int{},
		}
	}

	mappoints := map[int]int{}
	changes := 1
	for i := 0; changes > 0; i++ {
		changes = 0
		for j := 0; j < k; j++ { // reset clusters
			clusters[j].pts = []int{}
		}

		for p, point := range pts { // распределили types.Points по кластерам. если какое-то изменение произошло, то обновляем changes
			dist := -1.0
			ci := 0
			for i, cl := range clusters {
				d := CalcSphereDistance(cl.center, point)
				if dist < 0 || d < dist {
					dist = d
					ci = i
				}
			}

			clusters[ci].pts = append(clusters[ci].pts, p)
			if mappoints[p] != ci {
				mappoints[p] = ci
				changes++
			}
		}

		for ci := 0; ci < len(clusters); ci++ {
			if len(clusters[ci].pts) == 0 {
				// During the iterations, if any of the cluster centers has no
				// data types.Points associated with it, assign a random data types.Point
				// to it.
				// Also see: http://user.ceng.metu.edu.tr/~tcan/ceng465_f1314/Schedule/KMeansEmpty.html
				var ri int
				for {
					// find a cluster with at least two data types.Points, otherwise
					// we're just emptying one cluster to fill another
					ri = int(rand.Float64() * float64(len(pts))) //nolint:gosec // rand.Intn is good enough for this
					if len(clusters[mappoints[ri]].pts) > 1 {
						break
					}
				}
				clusters[ci].pts = append(clusters[ci].pts, ri)
				mappoints[ri] = ci

				// Ensure that we always see at least one more iteration after
				// randomly assigning a data types.Point to a cluster
				changes = len(pts)
			}
		}

		l := 0
		if changes > 0 {
			for i := 0; i < k; i++ {
				l = len(clusters[i].pts)
				if l == 0 {
					continue
				}

				cc := types.Point{}
				for _, p := range clusters[i].pts {
					cc.Lat += pts[p].Lat
					cc.Lng += pts[p].Lng
				}

				clusters[i].center = types.Point{Lat: cc.Lat / float64(l), Lng: cc.Lng / float64(l)}
			}
		}
		if i == m.iterationThreshold ||
			changes < int(float64(len(pts))*m.deltaThreshold) {

			fmt.Println("kmins finished", i, changes)
			break
		}
	}

	res := make([][]int, len(clusters))

	for i, cl := range clusters {
		res[i] = cl.pts
	}

	return res, nil
}
