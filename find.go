package pixoo64

import (
	"encoding/json"
	"errors"
	"io/fs"
)

type pixoo64LanResponse struct {
	ReturnCode    int    // 0
	ReturnMessage string // ""
	DeviceList    []*Pixoo64
}

func SameLANDevices() ([]*Pixoo64, error) {
	resp, err := Client.Get("https://app.divoom-gz.com/Device/ReturnSameLANDevice")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r *pixoo64LanResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&r)
	if err != nil {
		return nil, err
	}

	// TODO check ReturnCode?
	return r.DeviceList, nil
}

func FindFirst() (*Pixoo64, error) {
	lst, err := SameLANDevices()
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, errors.New("no Pixoo64 found on local network")
	}

	return lst[0], nil
}

func FindId(id int) (*Pixoo64, error) {
	lst, err := SameLANDevices()
	if err != nil {
		return nil, err
	}

	for _, dev := range lst {
		if dev.DeviceId == id {
			return dev, nil
		}
	}
	return nil, fs.ErrNotExist
}

func FindMac(mac string) (*Pixoo64, error) {
	lst, err := SameLANDevices()
	if err != nil {
		return nil, err
	}

	for _, dev := range lst {
		if dev.DeviceMac == mac {
			return dev, nil
		}
	}
	return nil, fs.ErrNotExist
}

func FindName(name string) (*Pixoo64, error) {
	lst, err := SameLANDevices()
	if err != nil {
		return nil, err
	}

	for _, dev := range lst {
		if dev.DeviceName == name {
			return dev, nil
		}
	}
	return nil, fs.ErrNotExist
}
