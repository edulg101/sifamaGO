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
	"os/exec"
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

func populateFotosOnDB2(path, localId, caption string, local *Local, listaGeo []Geolocation) error {

	var lat float64
	var long float64
	var geoMatch bool

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)

		}

		dir, name := filepath.Split(currentPath)

		re := regexp.MustCompile(localId + `[^0-9]`)
		m := re.MatchString(name)

		if m {
			imageRdr, err := os.Open(currentPath)
			if err != nil {
				fmt.Println("nao consegui ler o arquivo")
			}

			metaData, err := exif.Decode(imageRdr)

			// config, _, _ := image.DecodeConfig(imageRdr)
			// width := config.Width
			// fmt.Println(width)

			imageRdr.Close()

			if err != nil {
				fmt.Println("nao consegui Extrair MEtadados da imagem")
			} else {

				lat, long, err = metaData.LatLong()

				if err != nil {
					fmt.Println("nao consegui Extrair Gps da imagem")
				}
			}

			rodovia, km, valid := GetLocation(lat, long, listaGeo)
			if valid {

				if math.Abs(local.KmInicialDouble-km) < 1 && (strings.Contains(rodovia, local.Rodovia)) {
					geoMatch = true

				}
			}

			url := filepath.Join("fotos", name)
			url = filepath.ToSlash(url)
			urlp := template.URL(url)

			name = CheckForNameSize(currentPath)

			l := *local

			foto := Foto{
				Nome:       name,
				Path:       template.URL(filepath.Join(dir, name)),
				Legenda:    caption,
				LocalID:    local.ID,
				Local:      l,
				Latitude:   lat,
				Longitude:  long,
				GeoRodovia: rodovia,
				GeoKm:      km,
				GeoMatch:   geoMatch,
				UrlPath:    urlp,
				OriginPath: currentPath,
			}

			db.GetDB().Create(&foto)
			// local.Fotos = append(local.Fotos, foto)
			// db.GetDB().Save(&local)
			err = foto.merge()
			if err != nil {
				fmt.Println("erro no foto merge")
				return err
			}

		}
		return err
	})
	return err
}

func populateFotosOnDB(path string) {

	var lat float64
	var long float64

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

			metaData, err := exif.Decode(imageRdr)

			// config, _, _ := image.DecodeConfig(imageRdr)
			// width := config.Width
			// fmt.Println(width)

			imageRdr.Close()

			if err != nil {
				fmt.Println("nao consegui Extrair MEtadados da imagem")
			} else {

				lat, long, err = metaData.LatLong()

				if err != nil {
					fmt.Println("nao consegui Extrair Gps da imagem")
				}
			}

			// ***

			// if width > int(util.MAXIMAGEWIDTH) {

			// 	err = resizeImageAndCopyMetadata(currentPath, util.MAXIMAGEWIDTH)

			// 	if err != nil {
			// 		fmt.Println(err)
			// 	}
			// }

			image := Foto{
				Nome:       name,
				Path:       template.URL(currentPath),
				UrlPath:    urlp,
				Latitude:   lat,
				Longitude:  long,
				OriginPath: currentPath,
			}

			db.GetDB().Create(&image)
		}
		return err
	})
	if err != nil {
		panic(err)

	}
}

func saveFotosOnLocal(IdColuna, caption string, local *Local, listaGeo []Geolocation) error {

	var fotos []Foto

	db.GetDB().Find(&fotos)

	for _, foto := range fotos {
		name := foto.Nome

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
			if err != nil {
				fmt.Println("erro no open")
				return err
			}

			file.Close()

			err = foto.merge()
			if err != nil {
				fmt.Println("erro no foto merge")
				return err
			}

		}
	}
	return nil
}

func (foto Foto) merge() error {

	caption := foto.Legenda
	filePath := string(foto.Path)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}

	width := img.Width
	height := img.Height

	defer file.Close()

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
			return err
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
			if err != nil {
				return err
			}
		}
		dc.DrawStringAnchored(lines[i], float64(width)/2, captionY, 0.5, 0.5)
	}
	im, _ := gg.LoadImage(filePath)
	dc.DrawImage(im, 0, -captionHeigth)
	image := dc.Image()

	currentWidth := dc.Width()

	if currentWidth > int(util.MAXIMAGEWIDTH) {
		image = resize.Resize(util.MAXIMAGEWIDTH, 0, image, resize.Bicubic)
	}

	_, fileName := filepath.Split(filePath)

	_, err = os.Stat(util.OUTPUTIMAGEFOLDER)
	if os.IsNotExist(err) {
		os.MkdirAll(util.OUTPUTIMAGEFOLDER, os.ModePerm)
	}

	target := filepath.Join(util.OUTPUTIMAGEFOLDER, fileName)

	final, err := os.Create(target)
	if err != nil {
		return err
	}

	err = jpeg.Encode(final, image, nil)
	if err != nil {
		return fmt.Errorf("não foi possivel codificar para jpeg o arquivo %s", foto.Path)
	}

	final.Close()

	//  ***
	// fmt.Printf("inserindo GPS no arquivo %s\n", target)
	// err = copyAllMetadata(string(foto.Path), target)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	return nil

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

func resizeImage(img image.Image, path string, size uint) error {

	m := resize.Resize(size, 0, img, resize.Bicubic)

	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}

	outFile, err := os.Create(path)
	errorHandle(err)

	err = jpeg.Encode(outFile, m, nil)
	if err != nil {
		fmt.Println(err)
	}

	err = outFile.Close()

	if err != nil {
		return err
	}

	return nil
}

func insertGPSIntoImage(filepath string, lat, long float64) error {

	command := "exiftool"
	longRef := "West"
	latRef := "South"

	newLat := fmt.Sprintf("-XMP:GPSLatitude='%f'", lat)
	newLong := fmt.Sprintf("-XMP:GPSLongitude='%f'", long)
	newLatRef := fmt.Sprintf("-GPSLongitudeRef='%s'", longRef)
	newLongRef := fmt.Sprintf("-GPSLatitudeRef='%s'", latRef)
	newLatExif := fmt.Sprintf("-exif:gpslatitude='%f'", lat)
	newLatExifRef := "-exif:gpslongituderef=W"
	newLongExif := fmt.Sprintf("-exif:gpslongitude='%f'", long)
	newLongExifRef := "-exif:gpslatituderef=S"

	override := "-overwrite_original"

	cmd := exec.Command(command, newLat, newLong, newLongExif, newLatExif, newLatExifRef, newLongExifRef, newLatRef, newLongRef, filepath, override)
	err := cmd.Run()

	if err != nil {
		fmt.Println(filepath)
		fmt.Println("lat ", lat)
		fmt.Println("long ", long)
		return fmt.Errorf("não foi possível inserir as informações de GPS no arquivo %s", filepath)
	}
	return nil

}

func copyAllMetadata(originPath, targetPath string) error {

	fmt.Println("tentando copiar metadados de:")
	fmt.Printf("%s para %s\n", originPath, targetPath)

	command := "exiftool"
	override := "-overwrite_original"
	cmd := exec.Command(command, "-TagsFromFile", originPath, targetPath, override)
	returnMessage, err := cmd.Output()
	fmt.Println(string(returnMessage))
	if err != nil {
		return fmt.Errorf(string(returnMessage))
	}
	return nil

}
func resizeImageAndCopyMetadata(imagePath string, size uint) error {

	var img image.Image

	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}

	img, err = jpeg.Decode(file)

	if err != nil {
		return err
	}

	file.Close()

	dir := filepath.Dir(imagePath)
	fmt.Println(dir)

	oldFilePath := filepath.Join(dir, "temp.jpg")

	err = os.Rename(imagePath, oldFilePath)

	if err != nil {
		fmt.Println(err)
		return err
	}

	img = resize.Resize(size, 0, img, resize.Bicubic)

	file, err = os.Create(imagePath)
	if err != nil {
		os.MkdirAll(imagePath, os.ModePerm)
		file, err = os.Create(imagePath)
		if err != nil {
			return err
		}
	}

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return err
	} else {
		fmt.Println("Reduzindo arquivo:", imagePath)
	}

	file.Close()

	err = copyAllMetadata(oldFilePath, imagePath)
	if err != nil {
		return err
	}
	err = os.Remove(oldFilePath)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

//ResizeAllImageInFolder will resize all image in a given Folder and Save with metadata.
func ResizeAllImagesInFolder(path string, width uint) error {

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)

		}
		_, name := filepath.Split(currentPath)

		if !info.IsDir() && (strings.HasSuffix(name, "jpg") || strings.HasSuffix(name, "jpeg") || strings.HasSuffix(name, "png")) {
			imgFile, err := os.Open(currentPath)
			if err != nil {
				return fmt.Errorf("não foi possível abrir imagem: %s", currentPath)
			}
			conf, _, err := image.DecodeConfig(imgFile)
			imgFile.Close()
			if err != nil {
				return fmt.Errorf("não consegui abrir o arquivo: %s", currentPath)
			}
			if conf.Width > int(width) {
				err = resizeImageAndCopyMetadata(currentPath, width)
				if err != nil {
					return fmt.Errorf("não foi possível reduzir o arquivo: %s", currentPath)
				}
			}

		}
		return err
	})
	return err
}

func resizeImageAndCopyMetadataFromOriginal(imagePath, originPath string, size uint) error {

	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}

	config, _, _ := image.DecodeConfig(file)

	width := config.Width

	file.Close()

	if uint(width) > size {
		file, _ := os.Open(imagePath)
		img, _ := jpeg.Decode(file)
		img = resize.Resize(size, 0, img, resize.Bicubic)
		file.Close()
		err = os.Remove(imagePath)
		fmt.Println(err)
		file, err = os.Create(imagePath)
		fmt.Println(err)
		err = jpeg.Encode(file, img, nil)
		fmt.Println(err)
		file.Close()

	}

	err = copyAllMetadata(originPath, imagePath)
	if err != nil {
		return err
	}
	return err

}
