package pixoo64_test

import (
	"image"
	"image/draw"
	"image/gif"
	"log"
	"net/http"
	"testing"

	"github.com/KarpelesLab/pixoo64"
)

func TestGif(t *testing.T) {
	// this requires having an actual device
	dev, err := pixoo64.FindFirst()
	if err != nil {
		t.Fatalf("could not locate device: %s", err)
		return
	}
	log.Printf("found device: %+v", dev)

	resp, err := http.Get("https://pixeljoint.com/files/icons/avatar4_8.gif")
	if err != nil {
		t.Skipf("download failed: %s", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Skipf("download failed: %s", resp.Status)
		return
	}
	src, err := gif.DecodeAll(resp.Body)

	if len(src.Image) >= 60 {
		t.Skipf("too many frames")
		return
	}

	// apply gif images
	prevImg := src.Image[0]
	for n, img := range src.Image {
		// build image
		newImg := &image.Paletted{}
		*newImg = *prevImg
		newImg.Pix = dup(prevImg.Pix)

		draw.Draw(newImg, newImg.Bounds(), img, image.Point{}, draw.Over)
		src.Image[n] = newImg
		prevImg = newImg
	}

	anim := pixoo64.NewAnim()
	for n, i := range src.Image {
		delay := src.Delay[n] * 10
		log.Printf("delay = %d", delay)
		imc, err := pixoo64.ConvertImage(i, delay)
		if err != nil {
			t.Errorf("convert image failed: %s", err)
			continue
		}
		anim.AppendFrame(imc, 0)
	}

	// send it
	anim.SendTo(dev)
}

func dup(v []byte) []byte {
	r := make([]byte, len(v))
	for n, b := range v {
		r[n] = b
	}
	return r
}
