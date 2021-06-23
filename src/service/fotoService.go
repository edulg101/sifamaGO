package service

import (
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"sifamaGO/src/db"
	"sifamaGO/src/model"
	"sifamaGO/src/util"

	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

type GeoUtil struct {
	difLat  float64
	difLong float64
	index   int
}

func PopulateFotosOnDB2(path, localId, caption string, local *model.Local, listaGeo []model.Geolocation, linha int) error {

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

			name, err = CheckForNameSize(currentPath)
			if err != nil {
				return err
			}
			url := filepath.Join("fotos", name)
			url = filepath.ToSlash(url)
			urlp := template.URL(url)

			l := *local

			foto := model.Foto{
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
			err = foto.Merge(linha)
			if err != nil {
				fmt.Println("erro no foto merge")
				return err
			}

		}
		return err
	})
	return err
}

func PopulateFotosOnDB(path string, sessionId string) {

	var lat float64
	var long float64

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)

		}

		base, name := filepath.Split(currentPath)
		if !info.IsDir() && (strings.HasSuffix(name, "jpg") || strings.HasSuffix(name, "jpeg") || strings.HasSuffix(name, "png")) {
			name, err = CheckForNameSize(currentPath)
			if err != nil {
				return err
			}
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

			image := model.Foto{
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

// func saveFotosOnLocal(IdColuna, caption string, local *model.Local, listaGeo []model.Geolocation) error {

// 	var fotos []model.Foto

// 	db.GetDB().Find(&fotos)

// 	for _, foto := range fotos {
// 		name := foto.Nome

// 		re := regexp.MustCompile(IdColuna + `[^0-9]`)
// 		m := re.MatchString(name)

// 		if m {
// 			rodovia, km, valid := GetLocation(foto.Latitude, foto.Longitude, listaGeo)
// 			if valid {
// 				foto.GeoRodovia = rodovia
// 				foto.GeoKm = km

// 				if math.Abs(local.KmInicialDouble-km) < 1 && (strings.Contains(rodovia, local.Rodovia)) {
// 					foto.GeoMatch = true

// 				}
// 			}

// 			foto.LocalID = local.ID
// 			if foto.Legenda == "" {
// 				foto.Legenda = caption
// 			}
// 			db.GetDB().Save(&foto)
// 			local.Fotos = append(local.Fotos, foto)
// 			db.GetDB().Save(&local)

// 			file, err := os.Open(string(foto.Path))
// 			if err != nil {
// 				fmt.Println("erro no open")
// 				return err
// 			}

// 			file.Close()

// 			err = foto.Merge()
// 			if err != nil {
// 				fmt.Println("erro no foto merge")
// 				return err
// 			}

// 		}
// 	}
// 	return nil
// }

func resizeImage(img image.Image, path string, size uint) error {

	m := resize.Resize(size, 0, img, resize.Bicubic)

	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}

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
	var command string
	if strings.Contains(runtime.GOOS, "window") {
		command = util.EXIFTOOL
	} else {
		command = "exiftool"
	}
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

	var command string
	if strings.Contains(runtime.GOOS, "window") {
		command = util.EXIFTOOL
	} else {
		command = "exiftool"
	}

	fmt.Println("tentando copiar metadados de:")
	fmt.Printf("%s para %s\n", originPath, targetPath)

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
	fmt.Println(oldFilePath)
	fmt.Println(imagePath)

	err = copyAllMetadata(oldFilePath, imagePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = os.Remove(oldFilePath)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

//ResizeAllImageInFolder will resize all image in a given Folder and Save with metadata.
func ResizeAllImagesInFolder(path string, width uint) (string, error) {

	totalImagesDone := 0

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			if !strings.Contains(err.Error(), "temp.jpg") {
				panic(err)
			}

		}
		base, name := filepath.Split(currentPath)
		if len(name) > 90 {
			name, err = CheckForNameSize(currentPath)
			if err != nil {
				return err
			}
			currentPath = filepath.Join(base, name)

		}

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
					return fmt.Errorf("não foi possível reduzir o arquivo: %s\n", currentPath+err.Error())
				} else {
					totalImagesDone++
				}

			}

		}
		return err
	})

	var returnMessage string
	if totalImagesDone > 0 {
		returnMessage = fmt.Sprintf("Sucesso ! %d imagens Compactadas", totalImagesDone)
	} else {
		returnMessage = "Não há nenhuma imagem para ser compactada na pasta selecionada"
	}
	fmt.Println(returnMessage)
	return returnMessage, err
}

func ResizeImageAndCopyMetadataFromOriginal(imagePath, originPath string, size uint) error {

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
