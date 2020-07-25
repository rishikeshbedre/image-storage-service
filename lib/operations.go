package lib

import (
	"errors"
	"net/http"
	"os"

	"github.com/rishikeshbedre/image-storage-service/model"
)

// ImageInfo stores lock and map of image storage
var ImageInfo = model.ImageSchema{}

// StoragePath stores relative to image storage
var StoragePath = "image-db/image-store/"

// OpListAlbums function returns an array of albums present in imageinfo
func OpListAlbums() []string {
	ImageInfo.RLock()
	tempAlbumsList := make([]string, 0, len(ImageInfo.AlbumMap))
	for tempAlbumName := range ImageInfo.AlbumMap {
		tempAlbumsList = append(tempAlbumsList, tempAlbumName)
	}
	ImageInfo.RUnlock()
	return tempAlbumsList
}

// OpAddAlbum function creates folder in image-store and adds entry in image-info json
func OpAddAlbum(albumname string) (int, error) {
	// check if album already present
	ImageInfo.RLock()
	_, isPresent := ImageInfo.AlbumMap[albumname]
	ImageInfo.RUnlock()
	if isPresent {
		return http.StatusBadRequest, errors.New(albumname + " album already present")
	}

	// create album folder if not present
	_, statErr := os.Stat(StoragePath + albumname)
	if os.IsNotExist(statErr) {
		errDir := os.MkdirAll(StoragePath+albumname, 0755)
		if errDir != nil {
			return http.StatusInternalServerError, errors.New("Cannot create " + albumname + " album")
		}
	} else {
		return http.StatusInternalServerError, errors.New(albumname + " album present in storage but not in database")
	}

	// add album entry into the imageinfo map
	ImageInfo.Lock()
	ImageInfo.AlbumMap[albumname] = map[string]model.ImageField{}
	ImageInfo.Unlock()
	NotifierChan <- albumname + " album created"
	ImageInfoChan <- true
	return http.StatusOK, nil
}

// OpDeleteAlbum function deletes the folder in image-store and deletes from image-info if present
func OpDeleteAlbum(albumname string) (int, error) {
	// check if album is present
	ImageInfo.RLock()
	_, isPresent := ImageInfo.AlbumMap[albumname]
	ImageInfo.RUnlock()
	if !isPresent {
		return http.StatusBadRequest, errors.New(albumname + " album is not present")
	}

	//delete album if present
	_, statErr := os.Stat(StoragePath + albumname)
	if os.IsNotExist(statErr) {
		return http.StatusInternalServerError, errors.New(albumname + " album present in database but not in storage")
	}
	errDir := os.RemoveAll(StoragePath + albumname)
	if errDir != nil {
		return http.StatusInternalServerError, errors.New(albumname + " album present but cannot delete")
	}

	//delete album entry from imageinfo map
	ImageInfo.Lock()
	delete(ImageInfo.AlbumMap, albumname)
	ImageInfo.Unlock()
	NotifierChan <- albumname + " album deleted"
	ImageInfoChan <- true
	return http.StatusOK, nil
}

// OpUpdateAlbum function modifies the folder name of the album
func OpUpdateAlbum(oldname, newname string) (int, error) {
	// check if album is present
	ImageInfo.RLock()
	_, isPresent := ImageInfo.AlbumMap[oldname]
	ImageInfo.RUnlock()
	if !isPresent {
		return http.StatusBadRequest, errors.New(oldname + " album is not present")
	}

	//update album if present
	_, statErr := os.Stat(StoragePath + oldname)
	if os.IsNotExist(statErr) {
		return http.StatusInternalServerError, errors.New(oldname + " album present in database but not in storage")
	}
	errDir := os.Rename(StoragePath+oldname, StoragePath+newname)
	if errDir != nil {
		return http.StatusInternalServerError, errors.New("Cannot modify the album")
	}

	//update album entry in imageinfo
	ImageInfo.Lock()
	ImageInfo.AlbumMap[newname] = ImageInfo.AlbumMap[oldname]
	delete(ImageInfo.AlbumMap, oldname)
	ImageInfo.Unlock()
	NotifierChan <- oldname + " album modified to " + newname
	ImageInfoChan <- true
	return http.StatusOK, nil
}

// OpListImages function returns array of image names present in the specified album
func OpListImages(albumname string) ([]string, error) {
	ImageInfo.RLock()
	_, isAlbumPresent := ImageInfo.AlbumMap[albumname]
	ImageInfo.RUnlock()
	if !isAlbumPresent {
		return nil, errors.New(albumname + " album is not present")
	}

	ImageInfo.RLock()
	tempImagesList := make([]string, 0, len(ImageInfo.AlbumMap[albumname]))
	for tempImageName := range ImageInfo.AlbumMap[albumname] {
		tempImagesList = append(tempImagesList, tempImageName)
	}
	ImageInfo.RUnlock()
	return tempImagesList, nil
}

// OpDeleteImage function deletes the image from the specified album
func OpDeleteImage(albumname, imagename string) (int, error) {
	ImageInfo.RLock()
	_, isAlbumPresent := ImageInfo.AlbumMap[albumname]
	ImageInfo.RUnlock()
	if !isAlbumPresent {
		return http.StatusBadRequest, errors.New(albumname + " album is not present")
	}

	ImageInfo.RLock()
	_, isImagePresent := ImageInfo.AlbumMap[albumname][imagename]
	ImageInfo.RUnlock()
	if !isImagePresent {
		return http.StatusBadRequest, errors.New(imagename + " image is not present")
	}

	_, statErr := os.Stat(StoragePath + albumname + "/" + imagename)
	if os.IsNotExist(statErr) {
		return http.StatusInternalServerError, errors.New(imagename + " image is present database but not in storage")
	}

	errFile := os.RemoveAll(StoragePath + albumname + "/" + imagename)
	if errFile != nil {
		return http.StatusInternalServerError, errors.New(imagename + " image present but cannot delete")
	}

	ImageInfo.Lock()
	delete(ImageInfo.AlbumMap[albumname], imagename)
	ImageInfo.Unlock()
	NotifierChan <- imagename + " image deleted from album " + albumname
	ImageInfoChan <- true
	return http.StatusOK, nil
}

// OpUpdateImage function modifies the image name in the storage and database
func OpUpdateImage(albumname, oldimagename, newimagename string) (int, error) {
	ImageInfo.RLock()
	_, isAlbumPresent := ImageInfo.AlbumMap[albumname]
	ImageInfo.RUnlock()
	if !isAlbumPresent {
		return http.StatusBadRequest, errors.New(albumname + " album is not present")
	}

	ImageInfo.RLock()
	_, isImagePresent := ImageInfo.AlbumMap[albumname][oldimagename]
	ImageInfo.RUnlock()
	if !isImagePresent {
		return http.StatusBadRequest, errors.New(oldimagename + " image is not present")
	}

	_, statErr := os.Stat(StoragePath + albumname + "/" + oldimagename)
	if os.IsNotExist(statErr) {
		return http.StatusInternalServerError, errors.New(oldimagename + " image is present database but not in storage")
	}

	errFile := os.Rename(StoragePath+albumname+"/"+oldimagename, StoragePath+albumname+"/"+newimagename)
	if errFile != nil {
		return http.StatusInternalServerError, errors.New("Cannot modify the image")
	}

	ImageInfo.Lock()
	ImageInfo.AlbumMap[albumname][newimagename] = model.ImageField{
		FileName: newimagename,
		FilePath: StoragePath + albumname + "/" + newimagename,
	}
	delete(ImageInfo.AlbumMap[albumname], oldimagename)
	ImageInfo.Unlock()
	NotifierChan <- oldimagename + " image modified to " + newimagename
	ImageInfoChan <- true
	return http.StatusOK, nil
}
