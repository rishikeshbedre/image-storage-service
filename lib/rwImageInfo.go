package lib

import (
	"io/ioutil"
	"log"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ImageInfoChan - channel to send signal to write current state of imageInfo map into file
var ImageInfoChan = make(chan bool, 1000)

// WriteImageInfoToFile function writes current state of imageInfo map into file
func WriteImageInfoToFile(imageInfoChan chan bool) {
	for {
		_, ok := <-imageInfoChan
		if !ok {
			log.Println("ImageInfoChan closed:: writing to file storage is corrupted")
			return
		}

		ImageInfo.Lock()
		tempJSONByte, marshalErr := json.Marshal(ImageInfo.AlbumMap)
		ImageInfo.Unlock()
		if marshalErr != nil {
			log.Println("Error in marshalling ImageInfo: ", marshalErr)
		}

		fileWriteErr := ioutil.WriteFile("image-db/image-info.json", tempJSONByte, 0755)
		if fileWriteErr != nil {
			log.Println("Cannot write into image-info.json: ", fileWriteErr)
		}
	}
}
