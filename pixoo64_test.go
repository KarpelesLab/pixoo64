package pixoo64_test

import (
	"log"
	"testing"

	"github.com/KarpelesLab/pixoo64"
)

const TestDevice = "192.168.19.154"

func TestBasics(t *testing.T) {
	// this requires having an actual device
	dev := pixoo64.New(TestDevice)

	res, err := dev.GetAllConf()
	if err != nil {
		t.Errorf("error = %s", err)
	} else {
		log.Printf("res = %+v", res)
		// map[Brightness:100 ClockTime:60 CurClockId:64 GalleryShowTimeFlag:0 GalleryTime:60 GyrateAngle:0 LightSwitch:1 MirrorFlag:0 PowerOnChannelId:1 RotationFlag:0 SingleGalleyTime:-1 TemperatureMode:0 Time24Flag:1 error_code:0]
	}

	w, err := dev.GetWeatherInfo()
	log.Printf("%+v %v", w, err)

	dev.ShortBeeps(3)

}
