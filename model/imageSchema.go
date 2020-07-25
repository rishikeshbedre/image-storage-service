package model

import (
	"sync"
)

// ImageField stores info of individual image
type ImageField struct {
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
}

// ImageSchema stores info of all albums in a map
type ImageSchema struct {
	sync.RWMutex
	AlbumMap map[string]map[string]ImageField
}
