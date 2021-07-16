package model

import (
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"sifamaGO/src/util"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Hash           string
	Concessionaria string
	Tros           []Tro // `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type Tro struct {
	gorm.Model
	ID            uint
	PalavraChave  string
	Observacao    string
	Prazo         string
	TipoPrazo     string
	Severidade    string
	Disposicao    string
	DisposicaoCod string
	DisposicaoArt string
	Locais        []Local // `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SessionID     uint
	Session       Session `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Foto struct {
	gorm.Model
	ID          uint
	Nome        string
	Path        template.URL
	Legenda     string
	LocalID     uint
	Local       Local `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Latitude    float64
	Longitude   float64
	GeoRodovia  string
	GeoKm       float64
	GeoMatch    bool
	UrlPath     template.URL
	OriginPath  string
	Orientation uint
}

type Local struct {
	gorm.Model
	ID               uint
	NumIdentificacao string
	Data             string
	Hora             string
	Estado           string
	Rodovia          string
	Pista            string
	KmInicial        string
	KmFinal          string
	Sentido          string
	KmInicialDouble  float64
	KmFinalDouble    float64
	TrechoDNIT       bool
	Valid            bool
	Fotos            []Foto `gorm:"ForeignKey:LocalID"`
	TroID            uint
	Tro              Tro `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Folder struct {
	FolderName string
}

type HomeModel struct {
	Folders []Folder
}

type TroModel struct {
	Title            string
	Tro              []Tro
	TotalTro         int
	Folders          []Folder
	TotalFotos       int
	LocalWithNoFotos []Local
}

type Geolocation struct {
	gorm.Model
	ID        uint `gorm:"primaryKey; autoIncrement"`
	Rodovia   string
	Km        float64
	Latitude  float64 `gorm:"precision:20"`
	Longitude float64 `gorm:"precision:20"`
}

func (foto Foto) Merge(linha int) error {

	caption := foto.Legenda
	filePath := string(foto.Path)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}

	width := img.Width
	height := img.Height
	//*
	switch foto.Orientation {
	case 6, 8:
		width, height = height, width
	}
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
			err := fmt.Errorf("extensão máxima da descrição de fotos atingida. diminua a descrição da linha %d", linha)
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
	switch foto.Orientation {
	case 6:
		im = imaging.Rotate270(im)
	case 8:
		im = imaging.Rotate90(im)
	case 3:
		im = imaging.Rotate180(im)

	}

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

	// switch foto.Orientation {
	// case 6:
	// 	image = imaging.Rotate90(image)
	// }

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
