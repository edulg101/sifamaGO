package main

import (
	"errors"
	"html/template"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"sifamaGO/db"
	"sifamaGO/util"

	"github.com/fogleman/gg"
)

func populateFotosOnDB(path string) {

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)

		}

		_, name := filepath.Split(currentPath)
		if !info.IsDir() && (strings.HasSuffix(name, "jpg") || strings.HasSuffix(name, "jpeg") || strings.HasSuffix(name, "png")) {
			name = CheckForNameSize(currentPath)

			url := filepath.Join("fotos", name)
			url = filepath.ToSlash(url)
			urlp := template.URL(url)
			image := Foto{
				Nome:    name,
				Path:    template.URL(currentPath),
				UrlPath: urlp}
			db.GetDB().Create(&image)
		}
		return err
	})
	if err != nil {
		panic(err)

	}
}

func saveFotosOnLocal(IdColuna, caption string, local *Local) {

	var fotos []Foto
	db.GetDB().Find(&fotos)
	// re := regexp2.MustCompile(IdColuna`(?!\\d+)`, 1)
	for _, foto := range fotos {
		name := foto.Nome
		// 	m, err := re.MatchString(name)

		re := regexp.MustCompile(IdColuna + `[^0-9]`)

		m := re.MatchString(name)

		if m {

			foto.LocalID = local.ID
			if foto.Legenda == "" {
				foto.Legenda = caption
			}
			db.GetDB().Save(&foto)
			local.Fotos = append(local.Fotos, foto)

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

	img, _, _ := image.DecodeConfig(file)
	defer file.Close()

	width := img.Width
	height := img.Height

	captionHeigth := 34

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

	dc := gg.NewContext(width, height+((len(lines)-1)*34))

	y := height - 34
	for i := 0; i < len(lines); i++ {
		dc.SetColor(color.White)
		dc.DrawRectangle(0, float64(y), float64(width), 34)

		dc.Fill()

		y = y + 34

		captionY := float64((height - 34) + (34 * i) + captionHeigth/2.0)

		dc.SetColor(color.Gray16{0x3030})
		if err := dc.LoadFontFace(util.FONTPATH, 18); err != nil {
			panic(err)
		}
		dc.DrawStringAnchored(lines[i], float64(width)/2, captionY, 0.5, 0.5)
	}
	im, _ := gg.LoadImage(filePath)
	dc.DrawImage(im, 0, -34)
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
