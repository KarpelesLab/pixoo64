package pixoo64

import (
	"encoding/json"
	"errors"
	"net/http"
)

type pixoo64LanResponse struct {
	ReturnCode    int    // 0
	ReturnMessage string // ""
	DeviceList    []*Pixoo64
}

func SameLANDevices() ([]*Pixoo64, error) {
	req, err := http.NewRequest("GET", "https://app.divoom-gz.com/Device/ReturnSameLANDevice", nil)
	if err != nil {
		return nil, err
	}

	resp, err := Client.Do(req)
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
