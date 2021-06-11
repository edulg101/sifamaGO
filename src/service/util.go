package service

import (
	"fmt"
	"os"
	"path/filepath"
)

//CheckForNameSize - check if FileName is too large
func CheckForNameSize(fullPath string) (string, error) {
	var err error

	folder, fileName := filepath.Split(fullPath)

	nameLength := len(fileName)
	if nameLength > 90 {
		oldNameFile := fileName

		fmt.Println("reduzindo o tamanho do nome do arquivo")
		half := nameLength / 2
		extraHalf := (nameLength - 90) / 2
		newFileName := oldNameFile[:half-extraHalf] + oldNameFile[half+extraHalf:]

		// newDecoder não necessário para Windows
		// fileName, err := charmap.CodePage850.NewDecoder().String(fileName)
		oldPath := filepath.Join(folder, oldNameFile)
		newPath := filepath.Join(folder, newFileName)
		fileName = newFileName
		fmt.Println("mudando o nome do arquivo")
		fmt.Println("nome antigo: ", oldPath)
		fmt.Println("novo nome:", newPath)
		err = os.Rename(oldPath, newPath)

	}
	return fileName, err
}
