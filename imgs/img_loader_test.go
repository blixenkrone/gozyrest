package imgs_test

import (
	"bytes"
	"io/ioutil"
)

// NewTestImage ..
func NewTestImage() error {
	b, err := ioutil.ReadFile("image.jpeg")
	if err != nil {
		return err
	}
	reader := bytes.NewReader(b)
	_ = reader
	return nil
}
