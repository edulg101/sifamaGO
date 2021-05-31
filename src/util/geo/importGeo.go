package geo

import (
	"fmt"
	"log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var dbsql *gorm.DB

type Geolocation struct {
	gorm.Model
	ID        uint `gorm:"primaryKey; autoIncrement"`
	Rodovia   string
	Km        float64
	Latitude  float64 `gorm:"precision:20"`
	Longitude float64
}

func GetDBGEO() *gorm.DB {
	return createDBGEO()
}

func createDBGEO() *gorm.DB {
	var err error
	// databaseName := ":memory:"
	databaseName := "geoDatabase.db"

	db, err = gorm.Open(sqlite.Open(databaseName), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Geolocation{})

	// ImportSpreadSheet()
	// fmt.Println("print")
	// printAllRows()
	// populateMockDB()

	return db

}

func ImportSpreadSheet() error {
	path := "D:\\Documentos\\Users\\Eduardo\\Documentos\\ANTT\\OneDrive - ANTT- Agencia Nacional de Transportes Terrestres\\CRO\\geo\\geoarrumada.xlsx"
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Get all the rows in the Sheet1.
	rows := f.GetRows("Planilha1")
	err = parseSpreadSheet(rows, db)
	return err
}

func parseSpreadSheet(rows [][]string, db *gorm.DB) error {

	for i, row := range rows {
		if i > -1 && i < 10000 {
			log.Println(i)
			rodovia := row[0]
			fmt.Println(row[1], row[1], row[2], row[3])
			km, err := strconv.ParseFloat(row[1], 64)

			if err != nil {
				fmt.Println("erro na linha", i)
				panic(err)
			}
			latitude, err := strconv.ParseFloat(row[2], 64)
			if err != nil {
				fmt.Println("erro na linha", i)
				panic(err)
			}
			longitude, err := strconv.ParseFloat(row[3], 64)
			if err != nil {
				fmt.Println("erro na linha", i)
				panic(err)
			}

			// fmt.Printf("imprimindo depois do parse. lat: %.20f, long: %.20f\n", latitude, longitude)
			location := Geolocation{
				Rodovia:   rodovia,
				Km:        km,
				Latitude:  latitude,
				Longitude: longitude,
			}
			db.Create(&location)
		}

	}
	return nil
}

func printAllRows() {
	var locations []Geolocation
	db.Find(&locations)

	for _, location := range locations {
		fmt.Printf("Rodovia %s. Km: %f, lat: %.15f, long: %.15f\n",
			location.Rodovia,
			location.Km,
			location.Latitude,
			location.Longitude)
	}
}

// Conectar funcao com banco de dados
func ConectarSql() {
	stringConexao := "eduardo:123456@/crogeo?charset=utf8&parseTime=True&loc=Local"

	var err error
	dbsql, err = gorm.Open(mysql.Open(stringConexao), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Geolocation{})

	ImportSpreadSheet()
	fmt.Println("print")
	printAllRows()
	// populateMockDB()
}
