package service

import (
	"math"
	"sifamaGO/src/model"
)

func GetLocation(latitude, longitude float64, listaGeo []model.Geolocation) (string, float64, bool) {
	if latitude < -17.5197411486680 || latitude == -1.0 {
		return "", -1, false
	}

	difLat := 0.0
	difLong := 0.0
	precisaoEmMetros := 50.0 // pode mudar
	precisaoEmGraus := precisaoEmMetros / 111139
	var filteredGeoList []GeoUtil

	for i, loc := range listaGeo {
		difLat = math.Abs(math.Abs(loc.Latitude) - math.Abs(latitude))
		difLong = math.Abs(math.Abs(loc.Longitude) - math.Abs(longitude))

		if difLat <= precisaoEmGraus && difLong <= precisaoEmGraus {
			geoUtil := GeoUtil{
				difLat:  difLat,
				difLong: difLong,
				index:   i,
			}
			filteredGeoList = append(filteredGeoList, geoUtil)
		}
	}
	if len(filteredGeoList) < 1 {
		return "", -1, false
	}

	var listClosests []closestsLocations

	for _, x := range filteredGeoList {
		avgDif := (x.difLong + x.difLat) / 2
		listClosests = append(listClosests, closestsLocations{avgDif, x.index})
	}

	minorIndex := listClosests[len(listClosests)-1].index
	for z := 0; z < len(listClosests); z++ {
		for h := z + 1; h < len(listClosests); h++ {
			if listClosests[z].avgDif < listClosests[h].avgDif {
				minorIndex = listClosests[z].index
			}
		}
	}

	return listaGeo[minorIndex].Rodovia, listaGeo[minorIndex].Km, true

}

type closestsLocations struct {
	avgDif float64
	index  int
}
