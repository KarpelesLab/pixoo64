package pixoo64

import (
	"encoding/base64"
	"errors"
	"image"
)

type Anim struct {
	Frames []*Image
}

type GetHttpGifId struct {
	PicId int
}

// AppendFrame will add the given image to the list. If it is a pixoo64.Image then
// the time parameter will be ignored, and the image's original time will be used.
func (a *Anim) AppendFrame(img image.Image, time int) error {
	myImg, err := ConvertImage(img, time)
	if err != nil {
		return err
	}
	a.Frames = append(a.Frames, myImg)
	return nil
}

func (a *Anim) Send(p Pixoo64) error {
	if len(a.Frames) == 0 || len(a.Frames) >= 60 {
		return errors.New("animation frames count invalid")
	}

	// send the animation to the given pixoo64
	var pid *GetHttpGifId
	err := p.command("Draw/GetHttpGifId", nil, &pid)
	if err != nil {
		return err
	}
	id := pid.PicId

	for n, i := range a.Frames {
		// send frame
		err = p.command("Draw/SendHttpGif", map[string]any{"PicNum": len(a.Frames), "PicWidth": i.Bounds().Size().X, "PicOffset": n, "PicID": id, "PicSpeed": i.Time, "PicData": base64.StdEncoding.EncodeToString(i.Pix)}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
