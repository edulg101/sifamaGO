package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//CheckForNameSize - check if FileName is too large
// func CheckForNameSize1(fullPath string) (string, error) {
// 	var err error

// 	folder, fileName := filepath.Split(fullPath)

// 	nameLength := len(fileName)

// 	oldNameFile := fileName
// 	if nameLength > 90 {
// 		fmt.Println("reduzindo o tamanho do nome do arquivo")
// 		half := nameLength / 2
// 		extraHalf := (nameLength - 90) / 2

// 		inferiorLimit := half - extraHalf
// 		superiorLimit := half + extraHalf
// 		stringByteArray := []byte(fileName)
// 		limit1 := stringByteArray[inferiorLimit-1]
// 		limit2 := stringByteArray[superiorLimit-1]

// 		if limit1 == 194 || limit1 == 195 {
// 			inferiorLimit -= 1
// 		}
// 		if limit2 == 194 || limit2 == 195 {
// 			superiorLimit += 1
// 		}

// 		newFileName := oldNameFile[:inferiorLimit] + oldNameFile[superiorLimit:]
// 	}

// 	// newDecoder não necessário para Windows
// 	// newFileName, _ = charmap.CodePage850.NewDecoder().String(newFileName)

// 	oldPath := filepath.Join(folder, oldNameFile)
// 	newPath := filepath.Join(folder, newFileName)
// 	fileName = newFileName
// 	fmt.Println("mudando o nome do arquivo")
// 	fmt.Println("nome antigo: ", oldPath)
// 	fmt.Println("novo nome:", newPath)
// 	err = os.Rename(oldPath, newPath)

// 	return fileName, err
// }

//CheckForNameSize - check if FileName is too large
func CheckForNameSize(fullPath string) (string, error) {

	var err error

	folder, fileName := filepath.Split(fullPath)

	newFileName := removeAccentsAndRemoveWhiteSpaces(fileName)

	nameLength := len(fileName)

	if nameLength > 85 {

		// newFileName := removeAccents(fileName)

		fmt.Println("reduzindo o tamanho do nome do arquivo")
		half := nameLength / 2
		extraHalf := (nameLength - 85) / 2

		inferiorLimit := half - extraHalf
		superiorLimit := half + extraHalf

		newFileName = newFileName[:inferiorLimit] + newFileName[superiorLimit:]
	}

	// newDecoder não necessário para Windows
	// fileName, err := charmap.CodePage850.NewDecoder().String(fileName)
	oldPath := filepath.Join(folder, fileName)
	newPath := filepath.Join(folder, newFileName)

	fmt.Println("mudando o nome do arquivo")
	fmt.Println("nome antigo: ", oldPath)
	fmt.Println("novo nome:", newPath)
	err = os.Rename(oldPath, newPath)

	return newFileName, err

	// }
	// return fileName, err
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func removeAccentsAndRemoveWhiteSpaces(str string) string {
	str = removeAccents(str)
	str = strings.Replace(str, " ", "_", -1)
	str = strings.Replace(str, ",", "_", -1)
	return str
}
