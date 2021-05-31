package util

import (
	"os"
	"path/filepath"
)

func ToAscII(path string) {
	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)

		}

		base, name := filepath.Split(currentPath)

		bytes := []byte(name)

		change := false

		for i, ch := range bytes {
			if ch > 123 {
				bytes[i] = 97
				change = true
			}
		}

		if change {
			os.Rename(currentPath, filepath.Join(base, string(bytes)))
		}
		return err
	})
	if err != nil {
		panic(err)

	}
}
