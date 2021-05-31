package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"sifamaGO/src/db"
	"sifamaGO/src/util/geo"

	"github.com/360EntSecGroup-Skylar/excelize"
	"gorm.io/gorm"
)

func ImportSpreadSheet(path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Get all the rows in the Sheet1.
	rows := f.GetRows("Planilha1")
	err = parseSpreadSheet(rows, db.GetDB())
	return err
}

func parseSpreadSheet(rows [][]string, db *gorm.DB) error {

	var tro Tro
	var local Local
	kmInicialDouble := 0.0
	kmFinalDouble := 0.0
	startLocais := false
	endLocais := false
	var kmInicial string
	var kmFinal string
	var nIdentidade string
	var date string
	var hora string
	var rodovia string
	var sentido string
	var pista string

	var previousLocal Local
	var listaGeo []Geolocation
	geo.GetDBGEO().Find(&listaGeo)

	for i, row := range rows {

		if i < 1 {
			continue
		}
		// var coluna rune = 65
		for j, word := range row {
			word = strings.ToLower(word)
			word = strings.TrimSpace(word)
			// coluna = rune(j) + coluna
			// fmt.Printf("Célula %d - %s : %s\n ", i+1, string(coluna), word)

			if strings.Contains(word, "de ident.") {
				break
			}

			if strings.TrimSpace(word) == "tro" {
				palavraChave := row[j+1]
				tro = Tro{PalavraChave: palavraChave}
				artList := GetDisposicaoLegal(palavraChave)
				if len(artList) > 1 {
					tro.DisposicaoArt = artList[0]
					tro.DisposicaoCod = artList[1]
				} else {
					return fmt.Errorf("erro. não consegui Identificar o tipo de Disposição Legal na linha %v.... abortando", i+1)

				}
				tro.Disposicao = GetDescricaoDisposicaoLegal(tro.DisposicaoCod)
				observacao := row[j+2]
				tro.Observacao = observacao
				tro.Prazo = row[j+3]
				_, err := strconv.ParseInt(tro.Prazo, 10, 32)
				if err != nil {
					return fmt.Errorf("erro. não consegui Identificar Prazo na linha %v.... abortando", i+1)

				}
				tipoPrazo := row[j+4]
				tipoPrazo = strings.ToLower(tipoPrazo)
				if strings.HasPrefix(tipoPrazo, "d") {
					tro.TipoPrazo = "2"
				} else if strings.HasPrefix(tipoPrazo, "h") {
					tro.TipoPrazo = "1"

				} else {
					return fmt.Errorf("erro. não consegui obter o tipo de prazo na linha %v.... abortando", i+1)

				}
				db.Create(&tro)

				startLocais = true
				endLocais = false
				break
			}

			if startLocais {
				if j == 0 {
					nIdentidade = strings.Replace(word, ".", "", -1)
				} else if j == 1 {
					if row[j+1] == "" {
						tempDate, err := time.Parse("1/2/06 15:04", word)
						if err != nil {
							return fmt.Errorf("erro. não consegui identificar a string data / hora na linha %v.... abortando", i+1)

						}
						date = tempDate.Format("02/01/2006")
						hora = tempDate.Format("15:04")

					} else {
						fmt.Println(word)
						tempDate, err := time.Parse("01-02-06", word)
						if err != nil {
							return fmt.Errorf("erro. não consegui identificar a string hora na linha %v.... abortando", i+1)

						}
						date = tempDate.Format("02/01/2006")
					}
				} else if j == 2 {
					if word != "" {
						hora = word
					}
				} else if j == 4 {
					if strings.Contains(word, "070") {
						rodovia = "70"
					} else if strings.Contains(word, "163") {
						rodovia = "163"
					} else if strings.Contains(word, "364") {
						rodovia = "364"
					} else {
						return fmt.Errorf("erro. não consegui identificar Rodovia na linha %v.... abortando", i+1)
					}
				} else if j == 5 {
					kmInicial = word
				} else if j == 6 {
					kmFinal = word
				} else if j == 7 {
					if strings.HasPrefix(word, "c") {
						sentido = "Crescente"
					} else if strings.HasPrefix(word, "d") {
						sentido = "Decrescente"
					} else {
						return fmt.Errorf("erro. não consegui identificar o sentido na linha %v.... abortando", i+1)
					}
					var err error
					kmInicialDouble, err = strconv.ParseFloat(kmInicial, 32)
					if err != nil {
						return fmt.Errorf("erro. não consegui identificar o km inicial na linha %v.... abortando", i+1)

					}

					kmFinalDouble, err = strconv.ParseFloat(kmFinal, 32)
					if err != nil {
						err := fmt.Errorf("erro. não consegui identificar o km final na linha %v.... abortando", i+1)
						return err
					}

					if sentido == "Crescente" {
						if kmInicialDouble > kmFinalDouble {
							kmInicialDouble, kmFinalDouble = kmFinalDouble, kmInicialDouble
						}
					} else if sentido == "Decrescente" {
						if kmInicialDouble < kmFinalDouble {
							kmFinalDouble, kmInicialDouble = kmInicialDouble, kmFinalDouble
						}
					}

					kmFinal = fmt.Sprintf("%.3f", kmFinalDouble)
					kmInicial = fmt.Sprintf("%.3f", kmInicialDouble)

					kmFinal = strings.ReplaceAll(kmFinal, ".", ",")
					kmInicial = strings.ReplaceAll(kmInicial, ".", ",")

				} else if j == 8 {
					if strings.Contains(word, "pp") {
						pista = "1"
					} else if strings.Contains(word, "pm") {
						pista = "2"
					} else {
						err := fmt.Errorf("erro. não consegui identificar a Pista na linha %v Principal ou Marginal? deve ser \"pp\" ou \"pm\".... abortando", i+1)
						return err
					}

				} else if j == 9 {
					caption := ProperTitle(word)
					local = Local{
						NumIdentificacao: nIdentidade,
						Data:             date,
						Hora:             hora,
						Rodovia:          rodovia,
						KmInicial:        kmInicial,
						KmInicialDouble:  kmInicialDouble,
						KmFinal:          kmFinal,
						KmFinalDouble:    kmFinalDouble,
						Sentido:          sentido,
						Pista:            pista,
						Tro:              tro,
					}
					// Check if is the same local:
					if rodovia == previousLocal.Rodovia &&
						kmInicial == previousLocal.KmInicial &&
						kmFinal == previousLocal.KmFinal &&
						sentido == previousLocal.Sentido &&
						pista == previousLocal.Pista {

						fmt.Println("local Repetido")
						caption = IsLocationValid(caption, &local)
						saveFotosOnLocal(local.NumIdentificacao, caption, &previousLocal, listaGeo)

					} else {
						local.TroID = tro.ID
						fmt.Println("salvando local .........................")
						db.Create(&local)
						caption = IsLocationValid(caption, &local)
						saveFotosOnLocal(local.NumIdentificacao, caption, &local, listaGeo)
						previousLocal = local
					}

				}

				if word == "tro" && !endLocais {
					endLocais = true
				}
			}
		}
	}
	checkForDuplicateTime()

	return nil
}

type compareTroTime struct {
	data    string
	hora    string
	localId int
}

func checkForDuplicateTime() {
	var t Tro

	tros := t.FindAll()

	var data string
	var hora string
	var localId int

	var list []compareTroTime

	for _, tro := range tros {

		data = tro.Locais[0].Data
		hora = tro.Locais[0].Hora
		localId = int(tro.Locais[0].ID)
		list = append(list, compareTroTime{
			data:    data,
			hora:    hora,
			localId: localId,
		})
	}

	match := getLocalIdWithDuplicated(list)

	for match != -1 {
		var local Local
		localId := list[match].localId
		db.GetDB().First(&local, localId)
		rand.Seed(time.Now().UnixNano())
		randomInt := rand.Intn(59)
		newMinutes := fmt.Sprintf("%02d", randomInt)
		oldHora := local.Hora
		horaCheia := oldHora[0:3]
		local.Hora = horaCheia + newMinutes
		db.GetDB().Save(&local)
		list[match].hora = local.Hora
		match = getLocalIdWithDuplicated(list)
	}
}

func getLocalIdWithDuplicated(list []compareTroTime) int {

	for i, v := range list {
		for j := i + 1; j < len(list); j++ {
			dataI := v.data
			horaI := v.hora
			dataJ := list[j].data
			horaJ := list[j].hora

			if dataI == dataJ && horaI == horaJ {

				return j
			}
		}
	}
	return -1
}
