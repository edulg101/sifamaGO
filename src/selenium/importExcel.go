package selenium

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"sifamaGO/src/db"
	"sifamaGO/src/model"
	"sifamaGO/src/service"
	"sifamaGO/src/tests/geo"
	"sifamaGO/src/util"

	"github.com/360EntSecGroup-Skylar/excelize"
	"gorm.io/gorm"
)

func ImportSpreadSheet(path string, session *model.Session) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Get all the rows in the Sheet1.
	rows := f.GetRows("Planilha1")
	err = parseSpreadSheet(session, rows, db.GetDB())
	return err
}

func parseSpreadSheet(session *model.Session, rows [][]string, db *gorm.DB) error {

	var tro model.Tro
	var local model.Local
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
	var estado string
	var err error

	var previousLocal model.Local
	var listaGeo []model.Geolocation
	geo.GetDBGEO().Find(&listaGeo)

	sessionID := session.ID

	for i, row := range rows {

		if i < 1 {
			continue
		}
		// var coluna rune = 65
		for j, word := range row {
			word = strings.ToLower(word)
			word = strings.TrimSpace(word)

			if strings.Contains(word, "de ident") {
				break
			}

			if strings.TrimSpace(word) == "tro" {
				palavraChave := row[j+1]
				tro = model.Tro{PalavraChave: palavraChave}
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
				tro.SessionID = sessionID
				db.Create(&tro)

				startLocais = true
				endLocais = false
				break
			}

			if startLocais {
				if j == 0 {
					fmt.Println("idOriginal-->", word)
					nIdentidade = strings.Replace(word, ".", "", -1)
				} else if j == 1 {
					if row[j+1] == "" {
						date, hora, err = parseFullDateAsDate(word)
						if err != nil {
							date, hora, err = parseFullDateAsDecimals(word)
							if err != nil {
								return fmt.Errorf("erro. não consegui identificar a string data / hora na linha %v.... abortando", i+1)
							}
						}

					} else {
						tempDate, err := time.Parse("01-02-06", word)
						date = tempDate.Format("02/01/2006")
						if err != nil {
							date, hora, err = parseFullDateAsDecimals(word)
							if err != nil {
								return fmt.Errorf("erro. não consegui identificar a string hora na linha %v.... abortando", i+1)
							}
						}

					}
				} else if j == 2 {
					if word != "" || len(hora) != 5 {

						hora, _ = horaToString(word)
						if len(hora) != 5 {
							return fmt.Errorf("erro. não consegui identificar a string hora na linha %v.... abortando transformado: %v, word: %v", i+1, hora, word)
						}

					}
				} else if j == 4 {
					var concessionaria string
					concessionaria, rodovia, estado, err = getEstadoERodovia(word, i)
					if err != nil {
						return err
					}
					if session.Concessionaria == "" {
						session.Concessionaria = concessionaria
					}

				} else if j == 5 {
					kmInicial = word
					kmInicial = strings.Replace(kmInicial, ",", ".", -1)
					fmt.Printf("km inicial: %v\n", kmInicial)
				} else if j == 6 {
					kmFinal = word
					kmFinal = strings.Replace(kmFinal, ",", ".", -1)
					fmt.Printf("km final : %v\n", kmFinal)
				} else if j == 7 {
					if strings.HasPrefix(word, "c") {
						sentido = "Crescente"
					} else if strings.HasPrefix(word, "d") {
						sentido = "Decrescente"
					} else {
						return fmt.Errorf("erro. não consegui identificar o sentido na linha %v.... abortando", i+1)
					}
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
					local = model.Local{
						NumIdentificacao: nIdentidade,
						Data:             date,
						Hora:             hora,
						Rodovia:          rodovia,
						Estado:           estado,
						KmInicial:        kmInicial,
						KmInicialDouble:  kmInicialDouble,
						KmFinal:          kmFinal,
						KmFinalDouble:    kmFinalDouble,
						Sentido:          sentido,
						Pista:            pista,
						Tro:              tro,
					}

					fmt.Println("nIdentidade-->", local.NumIdentificacao)
					fmt.Println("DATA-->", local.Data)
					fmt.Println("HORA-->", local.Hora)

					// Check if is the same local:
					if rodovia == previousLocal.Rodovia &&
						kmInicial == previousLocal.KmInicial &&
						kmFinal == previousLocal.KmFinal &&
						sentido == previousLocal.Sentido &&
						pista == previousLocal.Pista &&
						tro.Observacao == previousLocal.Tro.Observacao {

						fmt.Println("local Repetido")
						caption, err = IsLocationValid(caption, &local)
						if err != nil {
							return err
						}

						err := service.PopulateFotosOnDB2(util.ORIGINIMAGEPATH, local.NumIdentificacao, caption, &previousLocal, listaGeo, i)
						//***
						// err := saveFotosOnLocal(local.NumIdentificacao, caption, &previousLocal, listaGeo)
						if err != nil {

							return err
						}

					} else {
						local.TroID = tro.ID
						fmt.Println("salvando local .........................")
						db.Create(&local)
						caption, err = IsLocationValid(caption, &local)
						if err != nil {
							return err
						}
						err = service.PopulateFotosOnDB2(util.ORIGINIMAGEPATH, local.NumIdentificacao, caption, &local, listaGeo, i)
						//***
						// err := saveFotosOnLocal(local.NumIdentificacao, caption, &local, listaGeo)
						if err != nil {
							fmt.Println("erro no saveFotos")
							return err

						}

						previousLocal = local
					}

				}

				if word == "tro" && !endLocais {
					endLocais = true
				}
			}
		}
	}
	err = checkForDuplicateTime()
	if err != nil {
		return err
	}

	return nil
}

func parseFullDateAsDate(word string) (string, string, error) {
	tempDate, err := time.Parse("1/2/06 15:04", word)
	if err != nil {
		return "", "", err

	}
	date := tempDate.Format("02/01/2006")
	hora := tempDate.Format("15:04")
	return date, hora, err
}
func parseFullDateAsDecimals(word string) (string, string, error) {
	windowsStartDate, err := time.Parse("01-02-2006", "01-01-1900")
	if err != nil {
		return "", "", err

	}
	splited := strings.Split(word, ".")
	days, err := strconv.ParseInt(splited[0], 10, 64)
	if err != nil {
		return "", "", err

	}
	day := windowsStartDate.Add(time.Hour * time.Duration(days-1) * 24)
	date := day.Format("02/01/06")
	hora, err := horaToString(word)
	return date, hora, err
}

func horaToString(word string) (string, error) {
	num, err := strconv.ParseFloat(word, 64)
	if err != nil {
		return word, nil
	}
	if num > 4000 {
		decimals := (num - float64(int(num)))
		fullHour := decimals * 24
		hourDecimal := (fullHour - float64(int(fullHour)))
		minute := hourDecimal * 60
		timeStr := fmt.Sprintf("%02d:%02d", int(fullHour), int(math.Round(minute)))
		return timeStr, nil
	}
	minTemp := num * 24.0
	min := int(minTemp)
	seconds := int((minTemp - float64(min)) * 60)
	timeStr := fmt.Sprintf("%02d:%02d", min, seconds)
	return timeStr, nil
}

type compareTroTime struct {
	data    string
	hora    string
	localId int
}

func checkForDuplicateTime() error {

	tros := service.FindAllTro()

	var data string
	var hora string
	var localId int

	var list []compareTroTime

	for _, tro := range tros {

		if len(tro.Locais) < 1 {
			return fmt.Errorf("ocorreu um Erro. Há TRO sem local registrado, verifique a planilha proximo ao tro com descrição: %q", tro.Observacao)
		}
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
		var local model.Local
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
	return nil
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
