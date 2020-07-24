package apis

import(
	"net/http"

	"github.com/rishikeshbedre/image-storage-service/model"
	"github.com/rishikeshbedre/image-storage-service/lib"

	"github.com/gin-gonic/gin"
)

// ListAlbums handler returns list of albums present
// @Summary list of albums
// @Description get list of albums present in the storage
// @Tags albums
// @Produce json
// @Success 200 {object} model.ListAlbumJSON
// @Router /albums [get]
func ListAlbums(c *gin.Context){
	tempList := lib.OpListAlbums()
	c.JSON(http.StatusOK, gin.H{"albumList": tempList})
}

// AddAlbum handler creates a album
// @Summary create a album
// @Description creates new album if not present
// @Tags albums
// @Accept json
// @Produce json
// @Param album body model.AlbumJSON true "Add Album"
// @Success 200 {object} model.SuccessMessage
// @Failure 400 {object} model.ErrorMessage
// @Failure 500 {object} model.ErrorMessage
// @Router /albums [post]
func AddAlbum(c *gin.Context){
	var addAlbum model.AlbumJSON
	if jsonbinderr := c.ShouldBindJSON(&addAlbum); jsonbinderr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jsonbinderr.Error()})
		return
	}
	if respCode, addAlbumErr := lib.OpAddAlbum(addAlbum.AlbumName); addAlbumErr != nil {
		c.JSON(respCode, gin.H{"error": addAlbumErr.Error()})
	} else {
		c.JSON(respCode, gin.H{"message": addAlbum.AlbumName+" album added"})
	}
}

// DeleteAlbum handler deletes a album
// @Summary delete a album
// @Description deletes the specified album if present
// @Tags albums
// @Produce json
// @Param albumName path string true "Delete Album"
// @Success 200 {object} model.SuccessMessage
// @Failure 400 {object} model.ErrorMessage
// @Failure 500 {object} model.ErrorMessage
// @Router /albums/{albumName} [delete]
func DeleteAlbum(c *gin.Context){
	albumname := c.Param("albumName")
	if respCode, deleteAlbumErr := lib.OpDeleteAlbum(albumname); deleteAlbumErr != nil {
		c.JSON(respCode, gin.H{"error": deleteAlbumErr.Error()})
	} else {
		c.JSON(respCode, gin.H{"message": albumname+" album deleted"})
	}
}

// UpdateAlbum handler updates a album
// @Summary update a album
// @Description update the specified album if present
// @Tags albums
// @Accept json
// @Produce json
// @Param albumName path string true "Update Album Old Name"
// @Param newAlbumName body model.AlbumJSON true "Update Album New Name"
// @Success 200 {object} model.SuccessMessage
// @Failure 400 {object} model.ErrorMessage
// @Failure 500 {object} model.ErrorMessage
// @Router /albums/{albumName} [patch]
func UpdateAlbum(c *gin.Context){
	oldalbumname := c.Param("albumName")
	var updateAlbum model.AlbumJSON
	if jsonbinderr := c.ShouldBindJSON(&updateAlbum); jsonbinderr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jsonbinderr.Error()})
		return
	}
	if respCode, updateAlbumErr := lib.OpUpdateAlbum(oldalbumname, updateAlbum.AlbumName); updateAlbumErr != nil {
		c.JSON(respCode, gin.H{"error": updateAlbumErr.Error()})
	} else {
		c.JSON(respCode, gin.H{"message": "album "+oldalbumname+" modified to "+updateAlbum.AlbumName})
	}
}