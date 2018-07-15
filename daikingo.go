package daikingo

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Unit struct {
	hostorip string
	Control  ControlInfo
}

type ControlInfo struct {
	Power       PowerMode
	Mode        Mode
	Temperature Temperature
	FanRate     FanRate
	FanDir      FanDirection
	Humidity    Humidity
}

type PowerMode string
type Mode string
type Temperature int
type FanRate string
type FanDirection string
type Humidity string

const (
	POWER_ON  PowerMode = "1"
	POWER_OFF PowerMode = "0"
)

const (
	MODE_AUTO  Mode = "0"
	MODE_DEHUM Mode = "2"
	MODE_COOL  Mode = "3"
	MODE_HOT   Mode = "4"
	MODE_FAN   Mode = "6"
)

const (
	FAN_AUTO   FanRate = "A"
	FAN_SILENT FanRate = "B"
	FAN_LVL1   FanRate = "3"
	FAN_LVL2   FanRate = "4"
	FAN_LVL3   FanRate = "5"
	FAN_LVL4   FanRate = "6"
	FAN_LVL5   FanRate = "7"
)

const (
	FAN_STOP  FanDirection = "0"
	FAN_VERT  FanDirection = "1"
	FAN_HORIZ FanDirection = "2"
	FAN_3D    FanDirection = "3"
)

type Response map[string]string

func NewUnit(ip string) *Unit {
	var unit = new(Unit)
	unit.hostorip = ip
	return unit
}

func (u *Unit) RequestAndParse(path string, method string, params url.Values) (data Response, err error) {
	urlstring := "http://" + u.hostorip + path
	var resp *http.Response
	if method == "GET" {
		if params != nil {
			urlstring += "?" + params.Encode()
		}
		resp, err = http.Get(urlstring)
	} else if method == "POST" {
		resp, err = http.PostForm(urlstring, params)
	} else {
		return nil, errors.New("Not a valid method: " + method)
	}

	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	data = make(Response)
	for _, val := range strings.Split(string(content), ",") {
		v := strings.Split(val, "=")
		data[v[0]] = v[1]
		if v[0] == "name" {
			data[v[0]], err = url.QueryUnescape(v[1])
		}
	}

	if data["ret"] != "" && data["ret"] != "OK" {
		return nil, errors.New(data["ret"] + ", " + data["msg"])
	}
	return
}

func (u *Unit) SetControlInfo(params url.Values) (data Response, err error) {
	return u.RequestAndParse("/aircon/set_control_info", "GET", params)
}

func (u *Unit) GetBasicInfo() (data Response, err error) {
	return u.RequestAndParse("/common/basic_info", "GET", nil)
}

func (u *Unit) GetSensorInfo() (data Response, err error) {
	return u.RequestAndParse("/aircon/get_sensor_info", "GET", nil)
}

func (u *Unit) GetControlInfo() (data Response, err error) {
	data, err = u.RequestAndParse("/aircon/get_control_info", "GET", nil)
	u.Control.Power = PowerMode(data["pow"])
	u.Control.Mode = Mode(data["mode"])
	u.Control.FanRate = FanRate(data["f_rate"])
	u.Control.FanDir = FanDirection(data["f_dir"])
	temp, converr := strconv.Atoi(data["stemp"])
	if converr != nil {
		temp = 21
	}
	u.Control.Temperature = Temperature(temp)
	u.Control.Humidity = Humidity(data["shum"])
	return
}

func (u *Unit) GetModelInfo() (data Response, err error) {
	return u.RequestAndParse("/aircon/get_model_info", "GET", nil)
}
