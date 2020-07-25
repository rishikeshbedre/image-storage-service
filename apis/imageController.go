package apis

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	//"log"

	"github.com/rishikeshbedre/image-storage-service/lib"
	"github.com/rishikeshbedre/image-storage-service/model"

	"github.com/gin-gonic/gin"
)

// ListImages handler returns list of images of the specified album
// @Summary list of images
// @Description returns list of images of the specified album
// @Tags images
// @Produce json
// @Param albumName path string true "Album name"
// @Success 200 {object} model.ListImageJSON
// @Failure 400 {object} model.ErrorMessage
// @Failure 500 {object} model.ErrorMessage
// @Router /albums/{albumName}/images [get]
func ListImages(c *gin.Context) {
	albumname := c.Param("albumName")
	tempList, listImagesErr := lib.OpListImages(albumname)
	if listImagesErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": listImagesErr.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"imageList": tempList})
	}
}

// GetImage handler sends the specified image to rest client
// @Summary get a image
// @Description get the specified image from the server
// @Tags images
// @Produce octet-stream
// @Param albumName path string true "Album name"
// @Param imageName path string true "Image name"
// @Success 200 {file} octet-stream
// @Failure 400 {object} model.ErrorMessage
// @Failure 500 {object} model.ErrorMessage
// @Router /albums/{albumName}/images/{imageName} [get]
func GetImage(c *gin.Context) {
	albumname := c.Param("albumName")
	imagename := c.Param("imageName")

	lib.ImageInfo.RLock()
	_, isAlbumPresent := lib.ImageInfo.AlbumMap[albumname]
	lib.ImageInfo.RUnlock()
	if !isAlbumPresent {
		c.JSON(http.StatusBadRequest, gin.H{"error": albumname + " album is not present"})
		return
	}

	lib.ImageInfo.RLock()
	_, isImagePresent := lib.ImageInfo.AlbumMap[albumname][imagename]
	lib.ImageInfo.RUnlock()
	if !isImagePresent {
		c.JSON(http.StatusBadRequest, gin.H{"error": imagename + " image is not present"})
		return
	}

	_, statErr := os.Stat(lib.StoragePath + albumname + "/" + imagename)
	if os.IsNotExist(statErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": imagename + " image is present database but not in storage"})
		return
	}

	header := c.Writer.Header()
	header["Content-type"] = []string{"application/octet-stream"}
	header["Content-Disposition"] = []string{"attachment; filename=" + imagename}

	file, err := os.Open(lib.StoragePath + albumname + "/" + imagename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not sent image " + imagename + " for download"})
		return
	}
	defer file.Close()

	io.Copy(c.Writer, file)
}

// AddImage handler adds an image to the specified album
// @Summary add a image
// @Description adds an image to the specified album
// @Tags images
// @Accept multipart/form-data
// @Produce json
// @Param albumName path string true "Album name"
// @Param file formData file true "Add Image"
// @Success 200 {object} model.SuccessMessage
// @Failure 400 {object} model.ErrorMessage
// @Failure 500 {object} model.ErrorMessage
// @Router /albums/{albumName}/images [post]
func AddImage(c *gin.Context) {
	albumname := c.Param("albumName")
	file, formErr := c.FormFile("file")
	if formErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formErr.Error()})
		return
	}

	filename := filepath.Base(file.Filename)

	lib.ImageInfo.RLock()
	_, isAlbumPresent := lib.ImageInfo.AlbumMap[albumname]
	lib.ImageInfo.RUnlock()
	if !isAlbumPresent {
		c.JSON(http.StatusBadRequest, gin.H{"error": albumname + " album is not present"})
		return
	}

	lib.ImageInfo.RLock()
	_, isImagePresent := lib.ImageInfo.AlbumMap[albumname][filename]
	lib.ImageInfo.RUnlock()
	if isImagePresent {
		c.JSON(http.StatusBadRequest, gin.H{"error": filename + " image already present"})
		return
	}

	_, statErr := os.Stat(lib.StoragePath + albumname + "/" + filename)
	if !os.IsNotExist(statErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": filename + " image already present in storage but not present in database"})
		return
	}

	if fileSaveErr := c.SaveUploadedFile(file, lib.StoragePath+albumname+"/"+filename); fileSaveErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fileSaveErr.Error()})
		return
	}

	lib.ImageInfo.Lock()
	lib.ImageInfo.AlbumMap[albumname][filename] = model.ImageField{
		FileName: filename,
		FilePath: lib.StoragePath + albumname + "/" + filename,
	}
	lib.ImageInfo.Unlock()
	lib.NotifierChan <- filename + " image added to album " + albumname
	lib.ImageInfoChan <- true
	c.JSON(http.StatusOK, gin.H{"message": filename + " image added to album " + albumname})
}

// DeleteImage handler deletes a image from the specified album
// @Summary delete a image
// @Description deletes a image from the specified album
// @Tags images
// @Produce json
// @Param albumName path string true "Album name"
// @Param imageName path string true "Image name"
// @Success 200 {object} model.SuccessMessage
// @Failure 400 {object} model.ErrorMessage
// @Failure 500 {object} model.ErrorMessage
// @Router /albums/{albumName}/images/{imageName} [delete]
func DeleteImage(c *gin.Context) {
	albumname := c.Param("albumName")
	imagename := c.Param("imageName")
	if respCode, deleteImageErr := lib.OpDeleteImage(albumname, imagename); deleteImageErr != nil {
		c.JSON(respCode, gin.H{"error": deleteImageErr.Error()})
	} else {
		c.JSON(respCode, gin.H{"message": imagename + " image deleted from album " + albumname})
	}
}

// UpdateImage handler updates the image in specified album
// @Summary update a image
// @Description updates the image in specified album
// @Tags images
// @Accept json
// @Produce json
// @Param albumName path string true "Album name"
// @Param imageName path string true "Update Image Old Name"
// @Param newImageName body model.ImageJSON true "Update Image New Name"
// @Success 200 {object} model.SuccessMessage
// @Failure 400 {object} model.ErrorMessage
// @Failure 500 {object} model.ErrorMessage
// @Router /albums/{albumName}/images/{imageName} [patch]
func UpdateImage(c *gin.Context) {
	albumname := c.Param("albumName")
	oldimagename := c.Param("imageName")

	var updateImage model.ImageJSON
	if jsonbinderr := c.ShouldBindJSON(&updateImage); jsonbinderr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jsonbinderr.Error()})
		return
	}

	if respCode, updateImageErr := lib.OpUpdateImage(albumname, oldimagename, updateImage.ImageName); updateImageErr != nil {
		c.JSON(respCode, gin.H{"error": updateImageErr.Error()})
	} else {
		c.JSON(respCode, gin.H{"message": oldimagename + " image modified to " + updateImage.ImageName})
	}
}
