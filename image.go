package main

import (
	"image"
	"log"

	"github.com/disintegration/imaging"
)

type ImageService interface {
	ImageResize(string, string, int, int) string
	ImageFire(string, string, int, int) string
}

type imageService struct{}

func (imageService) ImageResize(aimFile string, srcImage string, width int, height int) int {
	src, err := imaging.Open(srcImage)
	var flag = 0
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	dstImageResize := imaging.Resize(src, width, height, imaging.Lanczos)
	flag = save(dstImageResize, aimFile)
	return flag
}

func (imageService) ImageFire(aimFile string, srcImage string, width int, height int) int {
	src, err := imaging.Open(srcImage)
	var flag = 0
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	dstImageFill := imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)
	flag = save(dstImageFill, aimFile)
	return flag
}

func save(dst image.Image, aimFile string) int {
	err := imaging.Save(dst, aimFile)
	var flag = 0
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	} else {
		flag = 1
	}
	return flag
}
