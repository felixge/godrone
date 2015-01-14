package main

import (
	"image"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/nfnt/resize"
)

const framesPerSec = 1

// imgMu protects access to the images
var imgMu sync.Mutex
var imgForward image.Image
var imgDown image.Image

func getImages() (fwd, down image.Image) {
	imgMu.Lock()
	defer imgMu.Unlock()

	return imgForward, imgDown
}

func setImages(fwd, down image.Image) (image.Image, image.Image) {
	imgMu.Lock()
	defer imgMu.Unlock()

	oldFwd, oldDown := imgForward, imgDown
	imgForward = fwd
	imgDown = down
	return oldFwd, oldDown
}

// fetchVideo runs in it's own goroutine
func fetchVideo() {
	if *dummy {
		return
	}

	for {
		fwd := fetchForward()
		setImages(fwd, nil)

		time.Sleep(1 / framesPerSec * time.Second)
	}
}

func fetchForward() image.Image {
	cmd := exec.Command("yavta", "-c1", "-F/tmp/frame", "-f", "UYVY", "-s", "1280x720", "/dev/video1")
	err := cmd.Run()
	if err != nil {
		log.Print("front image capture error: ", err)
		return nil
	}
	frame, err := ioutil.ReadFile("/tmp/frame")
	if err != nil {
		log.Print("front image read error: ", err)
		return nil
	}
	os.Remove("/tmp/frame")

	im := image.NewYCbCr(image.Rect(0, 0, 1280, 720), image.YCbCrSubsampleRatio422)
	frameToImage(frame, im)
	small := resize.Resize(300, 0, im, resize.NearestNeighbor)

	return small
}

// frameToImage copies frame into image i.
func frameToImage(frame []byte, i image.Image) {
	// Format UVUY into Y, Cb and Cr planes
	// http://linuxtv.org/downloads/v4l-dvb-apis/V4L2-PIX-FMT-UYVY.html
	// U = Cb, V = Cr

	im := i.(*image.YCbCr)
	y, br := 0, 0
	for i := 0; i < len(frame); i += 4 {
		im.Cb[br] = frame[i+0]
		im.Y[y] = frame[i+1]
		im.Cr[br] = frame[i+2]
		im.Y[y+1] = frame[i+3]
		br += 1
		y += 2
	}
}
