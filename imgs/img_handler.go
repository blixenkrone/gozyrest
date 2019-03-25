package imgs

import (
	"log"
	"net/http"
	"sync"

	"github.com/byblix/mongotest/storage"
	"go.mongodb.org/mongo-driver/bson"
)

/**
 * ? go store img in S3
 * ? go fetch img data from EXIF Lib
 * ? go insert data to mongoDB
 * ? go return img data to client
 */

var wg sync.WaitGroup

// InitImgData the endpoint
func InitImgData(w http.ResponseWriter, r *http.Request) {
	ch := make(chan []byte)
	defer r.Body.Close()
	image, err := NewImage(r.Body)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	// No wg or ch needed here because it only stores on S3
	go func() {
		if err := image.uploadImage(); err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
	}()

	wg.Add(1)
	go image.TagExif(&wg, ch)
	resJSON := <-ch
	// runtime.Breakpoint()
	close(ch)

	wg.Add(1)
	go func() {
		defer wg.Done()
		db := storage.DatabaseRef{Database: "images", Collection: "exif"}
		objectID, err := db.InsertOneItem(bson.M{"data": image.ByteVal})
		if err != nil {
			log.Panicf("Didn't store JSON: %s\n", err)
		}
		log.Println(objectID)
	}()

	wg.Wait()

	w.Write(resJSON)
}

func (img *ImageReader) uploadImage() error {
	file := storage.File{
		Body:   img.ByteVal,
		Name:   img.Name,
		Format: img.Format,
		Dir:    img.Format,
	}

	if err := file.UploadFile(); err != nil {
		return err
	}
	return nil
}

func handleError(w http.ResponseWriter, err error, code int) {
	log.Panicf("Err: %s", err)
	http.Error(w, err.Error(), code)
}
