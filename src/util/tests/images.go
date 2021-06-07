package tests

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/barasher/go-exiftool"
)

func getImage() {
	fmt.Println("test1")
	log.Println("test2")

	fileName := "Afundamento de-BR163MT-C-PP-km 613,050-9302302.jpg"

	baseWithGPS := "D:\\Documentos\\Users\\Eduardo\\Documentos\\ANTT\\OneDrive - ANTT- Agencia Nacional de Transportes Terrestres\\CRO\\Relatórios RTA\\2021_05_28 ACOMPANHAMENTO LOTE 07 CRO\\Anexos\\PAV - 24.05 a 28.05.2021"
	baseNoGPS := "D:\\sifamadocs\\imagens"

	pathWithGPS := filepath.Join(baseWithGPS, fileName)
	pathNoGPS := filepath.Join(baseNoGPS, fileName)
	fmt.Println(pathNoGPS)

	// fmt.Println(pathNoGPS)
	InsertGPSIntoImage(pathWithGPS)

}

func readBinaryFileDUMP(file *os.File) {

	reader := bufio.NewReader(file)
	buf := make([]byte, 256)

	for {
		_, err := reader.Read(buf)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		fmt.Printf("%s", hex.Dump(buf))
	}

}

func toSliceBytes(file *os.File) {

	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	str := fmt.Sprint(img)

	// str := fmt.Sprintln(image)

	// index1 := strings.Index(str, "[")
	// index2 :=  strings.Index(str, "]")
	// str1 := str[index1 + 1: index2]

	writeToFile("Gps.txt", str)

	// fmt.Printf("%x", 10)

	// arr := strings.Split(str, " ")

	// var arrByte []byte
	// for _, v := range arr {
	// 	fmt.Println("t", v)
	// 	i, err := strconv.ParseUint(v, 8, 16)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(i)
	// }

	// for _, v := range arr {
	// 	i := fmt.Sprintf("%X", v)
	// 	fmt.Println(v, i)
	// 	time.Sleep(time.Second / 5)
	// }

}

func writeToFile(name, str string) {
	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)

	writer.WriteString(str)

	writer.Flush()

}

func exifTool(filepath string) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		return
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(filepath)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			fmt.Printf("[%v] %v\n", k, v)
		}
	}
}

func InsertGPSIntoImage(filepath string) error {

	command := "exiftool"
	lat := -14.345924722222223
	long := -56.148826666666665
	longRef := "West"
	latRef := "South"

	newLat := fmt.Sprintf("-XMP:GPSLatitude='%f'", lat)
	newLong := fmt.Sprintf("-XMP:GPSLongitude='%f'", long)
	newLatRef := fmt.Sprintf("-GPSLongitudeRef='%s'", longRef)
	newLongRef := fmt.Sprintf("-GPSLatitudeRef='%s'", latRef)

	cmd := exec.Command(command, newLat, newLong, newLatRef, newLongRef, filepath)
	err := cmd.Run()

	time.Sleep(time.Second / 2)

	if err != nil {
		fmt.Println("deu erro: ", err)
		return err
	}
	fmt.Println("nao deu erro")
	return nil

}
