package main

import(
	"os"
	"log"

	"github.com/rishikeshbedre/image-storage-service/apis"
	docs "github.com/rishikeshbedre/image-storage-service/docs"
	"github.com/rishikeshbedre/image-storage-service/model"
	"github.com/rishikeshbedre/image-storage-service/lib"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @contact.name Rishikesh Bedre
// @contact.email rishikeshbedre@gmail.com

func main() {

	hostIP := os.Getenv("HOSTIP")
	port := ":3333"

	lib.ImageInfo.AlbumMap = make(map[string]map[string]model.ImageField)

	go lib.WriteImageInfoToFile(lib.ImageInfoChan)
	go lib.NotifierService(lib.NotifierChan, hostIP)

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Image Storage Service"
	docs.SwaggerInfo.Description = "A microservice based on REST APIs to store and retrieve images."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = hostIP+port
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		albums := v1.Group("/albums")
		{
			// CRUD on albums
			albums.GET("", apis.ListAlbums)
			albums.POST("", apis.AddAlbum)
			albums.DELETE(":albumName", apis.DeleteAlbum)
			albums.PATCH(":albumName", apis.UpdateAlbum)

			// CRUD on images
			albums.GET(":albumName/images", apis.ListImages)
			albums.GET(":albumName/images/:imageName", apis.GetImage)
			albums.POST(":albumName/images", apis.AddImage)
			albums.DELETE(":albumName/images/:imageName", apis.DeleteImage)
			albums.PATCH(":albumName/images/:imageName", apis.UpdateImage)
		}
	}

	log.Println("Image Storage Service started at port", port)
	_, envPresent := os.LookupEnv("PRODMODE")
	if !envPresent {
		log.Println("Swagger Documentation available at http://"+hostIP+port+"/swagger/index.html")
		router.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "PRODMODE"))
	}

	// listen and serve on 3333
	router.Run(port) 
}
