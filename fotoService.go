package main

import (
	"errors"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"sifamaGO/src/db"
	"sifamaGO/src/util"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

type GeoUtil struct {
	difLat  float64
	difLong float64
	index   int
}

func populateFotosOnDB(path string) {

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)

		}

		base, name := filepath.Split(currentPath)
		if !info.IsDir() && (strings.HasSuffix(name, "jpg") || strings.HasSuffix(name, "jpeg") || strings.HasSuffix(name, "png")) {
			name = CheckForNameSize(currentPath)
			currentPath = filepath.Join(base, name)

			url := filepath.Join("fotos", name)
			url = filepath.ToSlash(url)
			urlp := template.URL(url)

			imageRdr, err := os.Open(currentPath)
			if err != nil {
				fmt.Println("nao consegui ler o arquivo")
			}
			lat, long := 0.0, 0.0

			metaData, err := exif.Decode(imageRdr)
			if err != nil {
				fmt.Println("nao consegui Extrair MEtadados da imagem")
			} else {
				lat, long, err = metaData.LatLong()
				if err != nil {
					fmt.Println("nao consegui Extrair Gps da imagem")
				}
			}

			image := Foto{
				Nome:      name,
				Path:      template.URL(currentPath),
				UrlPath:   urlp,
				Latitude:  lat,
				Longitude: long,
			}

			db.GetDB().Create(&image)
		}
		return err
	})
	if err != nil {
		panic(err)

	}
}

func saveFotosOnLocal(IdColuna, caption string, local *Local, listaGeo []Geolocation) {

	var fotos []Foto

	db.GetDB().Find(&fotos)

	// re := regexp2.MustCompile(IdColuna`(?!\\d+)`, 1)
	for _, foto := range fotos {
		name := foto.Nome
		// 	m, err := re.MatchString(name)

		re := regexp.MustCompile(IdColuna + `[^0-9]`)

		m := re.MatchString(name)

		if m {
			rodovia, km, valid := GetLocation(foto.Latitude, foto.Longitude, listaGeo)
			if valid {
				foto.GeoRodovia = rodovia
				foto.GeoKm = km

				if math.Abs(local.KmInicialDouble-km) < 1 && (strings.Contains(rodovia, local.Rodovia)) {
					foto.GeoMatch = true

				}
			}

			foto.LocalID = local.ID
			if foto.Legenda == "" {
				foto.Legenda = caption
			}
			db.GetDB().Save(&foto)
			local.Fotos = append(local.Fotos, foto)
			db.GetDB().Save(&local)

			file, err := os.Open(string(foto.Path))
			errorHandle(err)
			file.Close()

			merge(string(foto.Path), foto.Legenda)
		}
		// if err != nil {
		// 	fmt.Println(err)
		// }
	}
}

func merge(filePath, caption string) {

	file, err := os.Open(filePath)
	errorHandle(err)

	img, str, err := image.DecodeConfig(file)
	width := img.Width
	height := img.Height

	fmt.Println(str) // type
	fmt.Println(err)

	file.Close()

	file, err = os.Open(filePath)
	errorHandle(err)

	imagejpg, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println(err)
		fmt.Println(file.Name())
	}

	if width > 500 {

		file.Close()

		resizeImage(imagejpg, filePath)

		file, err = os.Open(filePath)
		errorHandle(err)

		img, _, _ = image.DecodeConfig(file)

		width = img.Width
		height = img.Height

	}

	fontSize := (float64(width) * 0.036)

	captionHeigth := int(float64(height) / 19.6)

	words := strings.Split(caption, " ")

	var lines []string

	var line1 string
	var line2 string
	var line3 string

	lineFull1 := false
	lineFull2 := false
	for _, word := range words {
		if !lineFull1 && len(line1)+len(word) < 51 {
			line1 = line1 + word + " "

		} else if !lineFull2 && len(line2)+len(word) < 51 {
			lineFull1 = true
			line2 = line2 + word + " "
		} else {
			lineFull2 = true
			line3 = line3 + word + " "
		}
		if len(line3) > 50 {
			err := errors.New("extensão máxima da descrição de fotos atingida. diminua a descrição")
			errorHandle(err)
		}
	}
	lines = append(lines, line1)
	if len(line2) > 0 {
		lines = append(lines, line2)
	}
	if len(line3) > 0 {
		lines = append(lines, line3)
	}

	dc := gg.NewContext(width, height+((len(lines)-1)*captionHeigth))

	y := height - captionHeigth
	for i := 0; i < len(lines); i++ {
		dc.SetColor(color.White)
		dc.DrawRectangle(0, float64(y), float64(width), float64(captionHeigth))
		dc.Fill()
		y += captionHeigth
		captionY := float64((height - captionHeigth) + (captionHeigth * i) + captionHeigth/2.0)

		dc.SetColor(color.Gray16{0x3030})
		if err := dc.LoadFontFace(util.FONTPATH, fontSize); err != nil {
			panic(err)
		}
		dc.DrawStringAnchored(lines[i], float64(width)/2, captionY, 0.5, 0.5)
	}
	im, _ := gg.LoadImage(filePath)
	dc.DrawImage(im, 0, -captionHeigth)
	image := dc.Image()
	_, fileName := filepath.Split(filePath)

	_, err = os.Stat(util.OUTPUTIMAGEFOLDER)
	if os.IsNotExist(err) {
		os.MkdirAll(util.OUTPUTIMAGEFOLDER, os.ModePerm)
	}

	target := filepath.Join(util.OUTPUTIMAGEFOLDER, fileName)

	final, err := os.Create(target)
	if err != nil {
		panic(err)
	}
	defer final.Close()

	jpeg.Encode(final, image, nil)

}

func GetLocation(latitude, longitude float64, listaGeo []Geolocation) (string, float64, bool) {
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

func resizeImage(img image.Image, path string) *os.File {
	// decode jpeg into image.Image

	temp := filepath.Join(util.OUTPUTIMAGEFOLDER, "temp")
	_, err := os.Open(temp)
	if err != nil {
		os.MkdirAll(temp, os.ModePerm)
	}

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(500, 0, img, resize.Bicubic)

	err = os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}

	outFile, err := os.Create(path)
	errorHandle(err)

	err = jpeg.Encode(outFile, m, nil)
	if err != nil {
		fmt.Println(err)
	}

	return outFile
	// _, err = io.Copy(outFile, rc)

	// Close the file without defer to close before next iteration of loop

}
