package pixoo64

import (
	"encoding/json"
	"errors"
	"fmt"
)

type User struct {
	Nickname string
	HeadId   string // group1/M00/02/4B/L1ghblsVCSaEKftOAAAAABIDWZY45.head
	Level    int
	Score    int
	LikeCnt  int
	FansCnt  int
}

// GetUser fetches and returns information about the owner's of this device
func (dev *Pixoo64) GetUser() (*User, error) {
	// https://app.divoom-gz.com/User/GetUserData?DeviceId=300054377
	if dev.DeviceId == 0 {
		return nil, errors.New("unknown device id")
	}
	resp, err := Client.Get(fmt.Sprintf("https://app.divoom-gz.com/User/GetUserData?DeviceId=%d", dev.DeviceId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user *User

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
