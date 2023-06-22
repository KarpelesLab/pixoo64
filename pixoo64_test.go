package pixoo64_test

import (
	"image"
	"image/color"
	"log"
	"testing"

	"github.com/KarpelesLab/pixoo64"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func TestBasics(t *testing.T) {
	// this requires having an actual device
	dev, err := pixoo64.FindFirst()
	if err != nil {
		t.Fatalf("could not locate device: %s", err)
		return
	}
	log.Printf("found device: %+v", dev)

	res, err := dev.GetAllConf()
	if err != nil {
		t.Errorf("error = %s", err)
	} else {
		log.Printf("res = %+v", res)
		// map[Brightness:100 ClockTime:60 CurClockId:64 GalleryShowTimeFlag:0 GalleryTime:60 GyrateAngle:0 LightSwitch:1 MirrorFlag:0 PowerOnChannelId:1 RotationFlag:0 SingleGalleyTime:-1 TemperatureMode:0 Time24Flag:1 error_code:0]
	}

	anim := pixoo64.NewAnim()
	dr := &font.Drawer{Src: image.NewUniform(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}), Face: basicfont.Face7x13}

	for i := 0; i < 32; i++ {
		img, _ := pixoo64.NewImage(64, 100)
		dr.Dst = img
		dr.Dot = fixed.P(0, 10+i)

		dr.DrawString("Hello")

		anim.AppendFrame(img, 0)
	}

	// send it
	anim.SendTo(dev)
}
