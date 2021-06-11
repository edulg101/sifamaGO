package service

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestGetImage(t *testing.T) {
	t.Log("Testando")
	getImage()
	t.Fail()
}
func getImage() {

	path1 := "d:\\1.jpg"
	path2 := "d:\\2.jpg"
	print10Items(path1, path2)

}

func print10Items(path, path2 string) {

	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		panic(statsErr)
	}

	var size int64 = stats.Size()
	bytes1 := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes1)

	// 2

	file2, err := os.Open(path2)

	if err != nil {
		panic(err)
	}
	defer file2.Close()

	stats2, statsErr2 := file.Stat()
	if statsErr2 != nil {
		panic(statsErr2)
	}

	var size2 int64 = stats2.Size()
	bytes2 := make([]byte, size2)

	bufr2 := bufio.NewReader(file2)
	_, err = bufr2.Read(bytes2)

	for i := 0; i < int(size); i++ {
		if bytes1[i] != bytes2[i] {
			fmt.Println(i)
			fmt.Println(bytes1[i])
			fmt.Println(bytes2[i])
			time.Sleep(time.Second)
		}
	}

}
