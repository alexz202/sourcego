package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadService interface {
	Save(*gin.Context, string, bool) (string, error)
}

const PATH = "public/"
const MIN = 111111
const MAX = 999999
const URLPRIX = "https://dev.source.zejicert.cn/"
const REGXP_IMG = `data:image\/(.*);base64,(.*)`
const THUMB_MAX_WIDTH_STR = "60,150,315"
const THUMB_MAX_HEIGHT_STR = "80,150,420"
const THUMB_PREFIX = "small_,thumb_,big_"
const THUMB_PATH = "/thumb/"
const FIRE_W = 150
const FIRE_H = 150
const FIRE_PREFIX = "fs_"

type uploadService struct{}

type fileInfo struct {
	fileUrl  string
	fileName string
	ext      string
	fileUri  string
	Thumb    []fileInfo
}

func makeRandName(ext string) string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(MAX-MIN) + MIN
	name := fmt.Sprintf("%d%d", time.Now().Unix(), randNum)
	h := sha1.New()
	io.WriteString(h, name)
	_n := fmt.Sprintf("%x.", h.Sum(nil))
	return _n + ext
}

func parseImgStr(body string) (string, string) {
	compile := regexp.MustCompile(REGXP_IMG)
	match := compile.FindAllStringSubmatch(body, -1)
	if match == nil {
		return "", ""
	} else {
		ext := match[0][1]
		imgstr := match[0][2]
		return ext, imgstr
	}
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
	// 或者
	//return err == nil || !os.IsNotExist(err)
	// 或者
	//return !os.IsNotExist(err)
}

//form save file
func (uploadService) Save(c *gin.Context, designated_path string, is_random_name string, params map[string]string) (fileInfo, error) {
	var name string
	var xpath string
	_is_random_name, _ := strconv.Atoi(is_random_name)
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return fileInfo{}, nil
	}
	ext := path.Ext(file.Filename)
	if _is_random_name != 1 {
		name = filepath.Base(file.Filename)
	} else {
		name = makeRandName(ext)
	}
	fmt.Printf("get name:%s", name)

	if designated_path != "" {
		xpath = PATH + designated_path
		// CHECK FOLDER IF NOT EXIST MKDIR
		if !IsExist(xpath) {
			os.MkdirAll(xpath, 0777)
		}
	} else {
		xpath = PATH
	}

	fmt.Printf("get path:%s", xpath)
	filename := xpath + name
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return fileInfo{}, nil
	}
	thumb := []fileInfo{}
	//TODO make thumb
	_makeThumb, _ := strconv.Atoi(params["makeThumb"])
	if _makeThumb == 1 {
		if params["thumb_w_string"] == "0" {
			thumb_w_string := THUMB_MAX_HEIGHT_STR
			thumb_h_string := THUMB_MAX_WIDTH_STR
		} else {
			thumb_w_string := params["thumb_w_string"]
			thumb_h_string := params["thumb_h_string"]
		}

	}
	//TODO make fire
	if params["fire"] != "0_0_0" {

	}

	link := fileInfo{URLPRIX + filename, name, ext, designated_path + name, thumb}
	return link, nil
}

func (uploadService) base64Save(c *gin.Context, is_random_name string, params map[string]string) (fileInfo, error) {
	var name string
	var xpath string
	body, _ := ioutil.ReadAll(c.Request.Body)
	ext, imgstr := parseImgStr(string(body))
	if ext == "" {
		return fileInfo{}, nil
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(imgstr)
	//fmt.Printf("get imgstr:%s", decodeBytes)
	_is_random_name, _ := strconv.Atoi(is_random_name)
	if _is_random_name == 1 {
		name = makeRandName(ext)
	}
	yearMonthDay := time.Now().Format("2006/01/02")
	xpath = PATH + yearMonthDay
	// CHECK FOLDER IF NOT EXIST MKDIR
	if !IsExist(xpath) {
		os.MkdirAll(xpath, 0777)
	}
	fileName := xpath + "/" + name
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write([]byte(decodeBytes))
	}
	thumb := []fileInfo{}
	link := fileInfo{URLPRIX + fileName, name, ext, yearMonthDay + "/" + name, thumb}
	return link, nil
}

func makeThumb(thumb_w_string string, thumb_h_string string, designated_path string, filename string) []fileInfo {

	return []fileInfo{}
}
