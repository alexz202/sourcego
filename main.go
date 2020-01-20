package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func checkVerifyParmas(c *gin.Context) (interface{}, string, bool) {
	f := c.Request.URL.Query()
	var flag = false
	var signOrignVal = ""
	var dict map[string]interface{}     //定义dict为map类型
	dict = make(map[string]interface{}) //让dict可编辑
	for k, _ := range f {
		val := c.Query(k)
		fmt.Printf("key:%s ,v:%s\n", k, val)
		dict[k] = val
	}
	signOrignVal = c.Query("sign")
	delete(dict, "sign")
	signVal := signMD5(dict, "")
	if signOrignVal == signVal {
		flag = true
	}
	fmt.Printf("signOrignVal:%s ,signVal:%s\n", signOrignVal, signVal)
	return dict, signVal, flag
}

func signMD5(data map[string]interface{}, signKey string) string {
	if signKey == "" {
		signKey = "123321!@#"
	}
	var keys []string
	var val []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var strData string
	for _, k := range keys {
		reflectA := reflect.TypeOf(data[k])
		if reflectA.Kind() == reflect.Interface {
			_json, _ := json.Marshal(data[k])
			val = append(val, string(_json))
			// strData += signKey + string(_json)
		} else if reflectA.Kind() != reflect.String {
			strData += signKey + strconv.FormatFloat(data[k].(float64), 'E', -1, 32)
			val = append(val, strconv.FormatFloat(data[k].(float64), 'E', -1, 32))
		} else {
			//strData += signKey + data[k].(string)
			val = append(val, data[k].(string))
		}
	}

	strData = strings.Join(val, signKey)
	fmt.Println(strData)
	_data := []byte(strData)
	has := md5.Sum(_data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	//router.Static("/", "./public")
	router.POST("/Upload/image", func(c *gin.Context) {
		is_random_name := c.DefaultQuery("is_random_name", "1")
		designated_path := c.Query("designated_path")
		//designated_path := c.DefaultQuery("designated_path", "img/")
		makeThumb := c.DefaultQuery("makeThumb", "0")
		thumb_w_string := c.DefaultQuery("thumb_w_string", "0")
		thumb_h_string := c.DefaultQuery("thumb_h_string", "0")
		//		compressImg := c.DefaultQuery("compressImg", "0")
		fire := c.DefaultQuery("fire", "0_0_0")
		params := map[string]string{
			"makeThumb":      makeThumb,
			"fire":           fire,
			"thumb_w_string": thumb_w_string,
			"thumb_h_string": thumb_h_string,
		}
		// Source
		svc := uploadService{}
		flag := svc.CheckImage(c)
		if flag {
			link, _ := svc.Save(c, designated_path, is_random_name, params)
			c.JSON(http.StatusOK, gin.H{
				"code": 1000,
				"msg":  "success",
				"data": gin.H{
					"link": gin.H{
						"fileUrl":  link.FileUrl,
						"fileName": link.FileName,
						"ext":      link.Ext,
						"fileUri":  link.FileUri,
						"Thumb":    link.Thumb,
					},
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "invaild type,please upload image file",
			})
		}
	})
	router.POST("/Upload/avatar", func(c *gin.Context) {
		makeThumb := c.DefaultQuery("makeThumb", "0")
		thumb_w_string := c.DefaultQuery("thumb_w_string", "0")
		thumb_h_string := c.DefaultQuery("thumb_h_string", "0")
		//		compressImg := c.DefaultQuery("compressImg", "0")
		fire := c.DefaultQuery("fire", "0_0_0")

		// Source
		var is_random_name = "1"
		svc := uploadService{}
		params := map[string]string{
			"makeThumb":      makeThumb,
			"fire":           fire,
			"thumb_w_string": thumb_w_string,
			"thumb_h_string": thumb_h_string,
		}
		link, _ := svc.base64Save(c, is_random_name, params)
		c.JSON(http.StatusOK, gin.H{
			"code": 1000,
			"msg":  "success",
			"data": gin.H{
				"link": gin.H{
					"fileUrl":  link.FileUrl,
					"fileName": link.FileName,
					"ext":      link.Ext,
					"fileUri":  link.FileUri,
				},
			},
		})
	})
	router.GET("/Tff", func(c *gin.Context) {
		dict, signV, flag := checkVerifyParmas(c)
		c.JSON(http.StatusOK, gin.H{
			"dict": dict,
			"sign": signV,
			"flag": flag,
		})
	})
	router.POST("/Upload/file", func(c *gin.Context) {
		_, _, flag := checkVerifyParmas(c)
		if flag == true {
			is_random_name := c.DefaultQuery("is_random_name", "0")
			designated_path := c.Query("designated_path")
			//designated_path := c.DefaultQuery("designated_path", "img/")
			//		compressImg := c.DefaultQuery("compressImg", "0")
			params := map[string]string{}
			// Source
			svc := uploadService{}
			link, _ := svc.Save(c, designated_path, is_random_name, params)
			c.JSON(http.StatusOK, gin.H{
				"code": 1000,
				"msg":  "success",
				"data": gin.H{
					"link": gin.H{
						"fileUrl":  link.FileUrl,
						"fileName": link.FileName,
						"ext":      link.Ext,
						"fileUri":  link.FileUri,
						// "Thumb":    link.Thumb,
					},
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 999,
				"msg":  "签名错误",
			})
		}
	})
	router.POST("/Upload/fileStream", func(c *gin.Context) {
		_, _, flag := checkVerifyParmas(c)
		if flag == true {
			// Source
			ext := c.Query("ext")
			var is_random_name = "1"
			svc := uploadService{}
			params := map[string]string{
				"ext": ext,
			}
			link, _ := svc.steamSave(c, is_random_name, params)
			c.JSON(http.StatusOK, gin.H{
				"code": 1000,
				"msg":  "success",
				"data": gin.H{
					"link": gin.H{
						"fileUrl":  link.FileUrl,
						"fileName": link.FileName,
						"ext":      link.Ext,
						"fileUri":  link.FileUri,
					},
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 999,
				"msg":  "签名错误",
			})
		}
	})

	router.Run(":8086")
}
