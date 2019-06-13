package fileload

import (
	"fmt"
	"io/ioutil"
)

// Load file load
func Load(path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	return string(dat)
}
