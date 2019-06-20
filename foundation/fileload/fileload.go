package fileload

import (
	"io/ioutil"
)

// Load file load
func Load(path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(dat)
}
