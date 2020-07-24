package model

// SuccessMessage - success message to rest client
type SuccessMessage struct {
	Message string `json:"message"`
}

// ErrorMessage - error message to rest client
type ErrorMessage struct {
	Error string `json:"error"`
}

// ListAlbumJSON - JSON struct for response of get list of albums
type ListAlbumJSON struct {
	AlbumList []string `json:"albumList"`
}

// AlbumJSON - JSON struct for add album route
type AlbumJSON struct {
	AlbumName string `json:"albumName" binding:"required"`
}

// ListImageJSON - JSON struct for response of get list of images in an album
type ListImageJSON struct {
	ImageList []string `json:"imageList"`
}

// ImageJSON - JSON struct for updating image name
type ImageJSON struct {
	ImageName string `json:"imageName" binding:"required"`
}
