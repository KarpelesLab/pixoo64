package pixoo64

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var Client = http.DefaultClient

type Pixoo64 struct {
	DeviceName      string // "Pixoo64"
	DeviceId        int    // 300054377
	DeviceMac       string // "c8f09e3b0964"
	DevicePrivateIP string
}

func New(addr string) *Pixoo64 {
	return &Pixoo64{DevicePrivateIP: addr}
}

func (p *Pixoo64) doPost(body any) (*http.Response, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("while encoding body: %w", err)
	}

	req, err := http.NewRequest("POST", "http://"+string(p.DevicePrivateIP)+"/post", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json") // pixoo64 doesn't seem to care, but let's do things right

	return Client.Do(req)
}

func (p *Pixoo64) command(cmd string, args map[string]any, target any) error {
	if args == nil {
		args = make(map[string]any)
	}
	args["Command"] = cmd

	resp, err := p.doPost(args)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid http response %s", resp.Status)
	}

	if target == nil {
		_, err = io.Copy(io.Discard, resp.Body)
		return err
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(target)
	return err
}

// Reboot will reboot the device.
func (p *Pixoo64) Reboot() error {
	resp, err := p.doPost(map[string]any{"Command": "Device/SysReboot"})
	if err != nil {
		return err
	}
	// do not try to read response as this'll probably fail since the device will be rebooting
	resp.Body.Close()
	return nil
}

// GetAllConf retrieves all configuration parameters from the device
func (p *Pixoo64) GetAllConf() (*Config, error) {
	var obj *Config
	err := p.command("Channel/GetAllConf", nil, &obj)
	return obj, err
}

// SetBrightness sets the device's brightness between 0 and 100
func (p *Pixoo64) SetBrightness(v int) error {
	if v < 0 || v > 100 {
		return errors.New("brightness out of range")
	}
	// returns: {"error_code":0} ... should we handle error_code?
	return p.command("Channel/SetBrightness", map[string]any{"Brightness": v}, nil)
}

// SetLocation sets the location for which weather data is pulled from https://openweathermap.org/
func (p *Pixoo64) SetLocation(long, lat float64) error {
	return p.command("Sys/LogAndLat", map[string]any{"Longitude": long, "Latitude": lat}, nil)
}

func (p *Pixoo64) SetTimezone(tz string) error {
	return p.command("Sys/TimeZone", map[string]any{"TimeZoneValue": tz}, nil)
}

func (p *Pixoo64) SetTime(t time.Time) error {
	return p.command("Device/SetUTC", map[string]any{"Utc": t.Unix()}, nil)
}

// ScreenSwitch switches the screen on or off
func (p *Pixoo64) ScreenSwitch(state bool) error {
	onoff := 1
	if !state {
		onoff = 0
	}
	return p.command("Channel/OnOffScreen", map[string]any{"OnOff": onoff}, nil)
}

func (p *Pixoo64) GetDeviceTime() (*Time, error) {
	var t *Time
	err := p.command("Device/GetDeviceTime", nil, &t)
	return t, err
}

func (p *Pixoo64) GetWeatherInfo() (*WeatherInfo, error) {
	var w *WeatherInfo
	err := p.command("Device/GetWeatherInfo", nil, &w)
	return w, err
}

// Buzzer will cause the device to emit a sound. For example: Buzzer(100, 100, 500)
func (p *Pixoo64) Buzzer(activeTime, offTime, totalTime int) error {
	return p.command("Device/PlayBuzzer", map[string]any{"ActiveTimeInCycle": activeTime, "OffTimeInCycle": offTime, "PlayTotalTime": totalTime}, nil)
}

// ShortBeeps will call Buzzer in order to perform a number of beeps, this can be useful to have
// some values mean specific things, just like old time BIOS beeps.
func (p *Pixoo64) ShortBeeps(count int) error {
	return p.Buzzer(100, 100, 200*count-100)
}

func (p *Pixoo64) ResetHttpGifId() error {
	return p.command("Draw/ResetHttpGifId", nil, nil)
}
