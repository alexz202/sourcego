package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/", "./public")
	router.POST("/Upload/image", func(c *gin.Context) {
		is_random_name := c.DefaultQuery("is_random_name", "1")
		designated_path := c.Query("designated_path")
		//designated_path := c.DefaultQuery("designated_path", "img/")
		//makeThumb := c.DefaultQuery("makeThumb", "0")
		//		thumb_w_string := c.Query("thumb_w_string")
		//		thumb_h_string := c.Query("thumb_h_string")
		//		compressImg := c.DefaultQuery("compressImg", "0")
		//		fire := c.DefaultQuery("fire", "0_0_0")

		// Source
		svc := uploadService{}
		link, _ := svc.Save(c, designated_path, is_random_name)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "success",
			"data": gin.H{
				"link": gin.H{
					"fileUrl":  link.fileUrl,
					"fileName": link.fileName,
					"ext":      link.ext,
					"fileUri":  link.fileUri,
					"Thumb":    link.Thumb,
				},
			},
		})
	})
	router.POST("/Upload/avatar", func(c *gin.Context) {
		//makeThumb := c.DefaultQuery("makeThumb", "0")
		//		thumb_w_string := c.Query("thumb_w_string")
		//		thumb_h_string := c.Query("thumb_h_string")
		//		compressImg := c.DefaultQuery("compressImg", "0")
		//		fire := c.DefaultQuery("fire", "0_0_0")

		// Source
		var is_random_name = "1"
		svc := uploadService{}
		link, _ := svc.base64Save(c, is_random_name)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "success",
			"data": gin.H{
				"link": gin.H{
					"fileUrl":  link.fileUrl,
					"fileName": link.fileName,
					"ext":      link.ext,
					"fileUri":  link.fileUri,
					"Thumb":    link.Thumb,
				},
			},
		})
	})

	router.Run(":8080")
}
