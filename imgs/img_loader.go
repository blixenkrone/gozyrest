package imgs

import (
	"bytes"
	"fmt"
	"image"

	// Keep this import so the compiler knows the format
	_ "image/jpeg"
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/rwcarlsen/goexif/exif"
)

// ImageReader contains image info
type ImageReader struct {
	Image   image.Image
	Name    string
	Format  string
	ByteVal []byte
	Reader  io.Reader
}

// NewImage returns structure of new Image
func NewImage(r io.Reader) (*ImageReader, error) {
	var buf bytes.Buffer
	teeRead := io.TeeReader(r, &buf)
	src, format, err := image.Decode(teeRead)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Image format is: %s\n", format)

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		handleWriteErr(err)
	}
	return &ImageReader{
		Image:   src,
		Format:  format,
		Name:    string(uuid)[:13] + "." + format,
		ByteVal: buf.Bytes(),
		Reader:  &buf,
	}, nil
}

// TagExif returns the bytes of the tiff in ch
func (img *ImageReader) TagExif(wg *sync.WaitGroup, ch chan<- []byte) {
	defer wg.Done()
	out, err := exif.Decode(img.Reader)
	if err != nil {
		if exif.IsCriticalError(err) {
			log.Fatalf("exif.Decode, critical error: %v", err)
		}
		log.Printf("exif.Decode, warning: %v", err)
	}
	val, err := out.MarshalJSON()
	if err != nil {
		log.Fatalf("Error marshalling JSON: %s", err)
	}
	ch <- val
}

// func imgFormatToPNG(w io.Writer, r io.Reader, format string) error {
// 	switch format {
// 	case "jpeg":
// 		fmt.Println("Its JPG!")
// 		img, err := jpeg.Decode(r)
// 		if err != nil {
// 			return err
// 		}
// 		return png.Encode(w, img)
// 	case "png":
// 		fmt.Println("Its PNG!")
// 	}
// 	return nil
// }

func handleWriteErr(err error) {
	fmt.Printf("This error happened %s\n", err)
	panic("Err:\n" + err.Error())
}
