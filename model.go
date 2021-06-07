package main

import (
	"html/template"

	"gorm.io/gorm"
)

type Foto struct {
	gorm.Model
	ID         uint
	Nome       string
	Path       template.URL
	Legenda    string
	LocalID    uint
	Local      Local
	Latitude   float64
	Longitude  float64
	GeoRodovia string
	GeoKm      float64
	GeoMatch   bool
	UrlPath    template.URL
	OriginPath string
}

type Local struct {
	gorm.Model
	ID               uint
	NumIdentificacao string
	Data             string
	Hora             string
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
	Tro              Tro
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
	Locais        []Local // `gorm:"ForeignKey:TroID"`
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
	Longitude float64
}
