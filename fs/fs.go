package fs

import (
	"fmt"
	"os"
)

func TouchFile(file string) {
	os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0777)
}

func Hello() {
	fmt.Println("this is from gokit fs package.")
}
