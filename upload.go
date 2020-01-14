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
	"strings"
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
	FileUrl  string     `json:"fileUrl,string"`
	FileName string     `json:"fileName,string"`
	Ext      string     `json:"ext,string"`
	FileUri  string     `json:"fileUri,string"`
	Thumb    []fileInfo `json:"Thumb,omitempty"`
}

func makeRandName(ext string) string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(MAX-MIN) + MIN
	name := fmt.Sprintf("%d%d", time.Now().Unix(), randNum)
	h := sha1.New()
	io.WriteString(h, name)
	_n := fmt.Sprintf("%x", h.Sum(nil))
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
	var thumb_path string
	var thumb_w_string string
	var thumb_h_string string
	var Thumb []fileInfo
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
		xpath = PATH + designated_path + "/"
		// CHECK FOLDER IF NOT EXIST MKDIR
		if !IsExist(xpath) {
			os.MkdirAll(xpath, 0777)
		}
	} else {
		xpath = PATH
	}

	fmt.Printf("get path:%s/r/n", xpath)
	filename := xpath + name
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return fileInfo{}, nil
	}
	imgServ := imageService{}
	//make thumb
	_makeThumb, _ := strconv.Atoi(params["makeThumb"])
	if _makeThumb == 1 {
		if params["thumb_w_string"] == "0" {
			thumb_w_string = THUMB_MAX_HEIGHT_STR
			thumb_h_string = THUMB_MAX_WIDTH_STR
		} else {
			thumb_w_string = params["thumb_w_string"]
			thumb_h_string = params["thumb_h_string"]
		}
		if designated_path != "" {
			thumb_path = "/" + designated_path + THUMB_PATH
		} else {
			thumb_path = THUMB_PATH + time.Now().Format("2006/01/02")
		}
		//fmt.Printf("get thumb path:%s", thumb_path)
		Thumb = makeThumb(imgServ, thumb_w_string, thumb_h_string, thumb_path, filename, name, ext)
	}
	//make fire
	if params["fire"] != "0_0_0" {
		fire_list := strings.Split(params["fire"], "_")
		if fire_list[0] == "1" {
			fire_w := FIRE_W
			fire_h := FIRE_H
			_fire_w, _ := strconv.Atoi(fire_list[1])
			_fire_h, _ := strconv.Atoi(fire_list[2])
			if _fire_w > 0 {
				fire_w = _fire_w
			}
			if _fire_h > 0 {
				fire_h = _fire_h
			}
			go imgServ.ImageFire(xpath+FIRE_PREFIX+name, filename, fire_w, fire_h)
		}
	}
	link := fileInfo{FileUrl: URLPRIX + filename, FileName: name, Ext: ext, FileUri: designated_path + name, Thumb: Thumb}
	return link, nil
}

func (uploadService) base64Save(c *gin.Context, is_random_name string, params map[string]string) (fileInfo, error) {
	var name string
	var xpath string
	var thumb_path string
	var thumb_w_string string
	var thumb_h_string string
	var Thumb []fileInfo
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
	imgServ := imageService{}
	//make thumb
	_makeThumb, _ := strconv.Atoi(params["makeThumb"])
	if _makeThumb == 1 {
		if params["thumb_w_string"] == "0" {
			thumb_w_string = THUMB_MAX_HEIGHT_STR
			thumb_h_string = THUMB_MAX_WIDTH_STR
		} else {
			thumb_w_string = params["thumb_w_string"]
			thumb_h_string = params["thumb_h_string"]
		}
		thumb_path = THUMB_PATH + time.Now().Format("2006/01/02")
		fmt.Printf("get thumb path:%s", thumb_path)
		Thumb = makeThumb(imgServ, thumb_w_string, thumb_h_string, thumb_path, fileName, name, ext)
	}

	if params["fire"] != "0_0_0" {
		fire_list := strings.Split(params["fire"], "_")
		if fire_list[0] == "1" {
			fire_w := FIRE_W
			fire_h := FIRE_H
			_fire_w, _ := strconv.Atoi(fire_list[1])
			_fire_h, _ := strconv.Atoi(fire_list[2])
			if _fire_w > 0 {
				fire_w = _fire_w
			}
			if _fire_h > 0 {
				fire_h = _fire_h
			}
			go imgServ.ImageFire(xpath+FIRE_PREFIX+name, fileName, fire_w, fire_h)
		}
	}

	link := fileInfo{FileUrl: URLPRIX + fileName, FileName: name, Ext: ext, FileUri: yearMonthDay + "/" + name, Thumb: Thumb}
	return link, nil
}

func makeThumb(imgServ imageService, thumb_w_string string, thumb_h_string string, thumb_path string, filename string, name string, ext string) []fileInfo {
	w_list := strings.Split(thumb_w_string, ",")
	h_list := strings.Split(thumb_h_string, ",")
	prefix_list := strings.Split(THUMB_PREFIX, ",")
	var thumb []fileInfo
	var _thumb []fileInfo
	//thumb := make([]fileInfo, 3)
	w_len := len(w_list)
	_thumb_path := PATH + thumb_path
	if !IsExist(thumb_path) {
		os.MkdirAll(_thumb_path, 0777)
	}
	chs := make([]chan int, w_len)
	i := 0
	for _, w := range w_list {
		w, _ := strconv.Atoi(w)
		h, _ := strconv.Atoi(h_list[i])
		//use channel
		chs[i] = make(chan int)
		go makeOneThumb(imgServ, _thumb_path+prefix_list[i]+name, filename, w, h, chs[i])
		//_thumb[i] = fileInfo{FileUrl: URLPRIX + _thumb_path + prefix_list[i] + name, FileName: prefix_list[i] + name, Ext: ext, FileUri: thumb_path + prefix_list[i] + name, Thumb: []fileInfo{}}
		// if flag == 1 {
		// 	//fmt.Printf("get thumb w:%d,h:%d;i:%d/r/n", w, h, i)
		_thumb = append(_thumb, fileInfo{FileUrl: URLPRIX + _thumb_path + prefix_list[i] + name, FileName: prefix_list[i] + name, Ext: ext, FileUri: thumb_path + prefix_list[i] + name, Thumb: []fileInfo{}})
		// }
		i++
	}
	// check all chs has finished
	for k, ch := range chs {
		v := <-ch
		fmt.Printf("get channel k:%d,ch:%d/r/n", k, v)
		if v == 1 {
			thumb = append(thumb, _thumb[k])
		}
	}

	return thumb
}

func makeOneThumb(imgServ imageService, aimFile string, srcImage string, width int, height int, ch chan int) int {
	flag := imgServ.ImageResize(aimFile, srcImage, width, height)
	if flag == 1 {
		ch <- 1
	}
	return flag
}
