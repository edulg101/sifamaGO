package selenium

import (
	"fmt"
	"regexp"
	"sifamaGO/src/util"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func reg(id string) string {

	r, _ := regexp.Compile(`TRO0*(\d+)202(\d)`)

	subString := r.FindStringSubmatch(id)

	return subString[1]

}

func checkForMissingTros(tros [][]string, listaTrosEmOrdem []int) error {
	flag := false
	var err error
	var troFaltantes []string
	for i, listaTroExcel := range tros {
		if i < 1 {
			continue
		}
		for _, listaDisponivelNoSistema := range listaTrosEmOrdem {
			v := strconv.Itoa(listaDisponivelNoSistema)
			if strings.Contains(v, listaTroExcel[0]) {

				flag = true
				break
			} else {
				flag = false
			}
		}
		if !flag {
			troFaltantes = append(troFaltantes, listaTroExcel[0])
		}
	}

	if len(troFaltantes) > 0 {
		fmt.Printf("faltando os Seguintes Tros - (%d) : \n", len(troFaltantes))
		message := fmt.Sprintln("Os Seguintes TROs não estão disponívels para verificação:")
		for _, v := range troFaltantes {
			fmt.Println(v)
			message += fmt.Sprintf("%s\n", v)
		}
		fmt.Println(message)
		return fmt.Errorf(message)
	}
	return err
}

func getInfoFromExcel() [][]string {
	f, err := excelize.OpenFile(util.CHECKSPREADSHEETPATH)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Get all the rows in the Sheet1.
	rows := f.GetRows("Planilha1")
	return rows
}
