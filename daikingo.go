package daikingo

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Unit struct {
	hostorip string
}

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
		return
	}

	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
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
	return u.RequestAndParse("/aircon/get_control_info", "GET", nil)
}

func (u *Unit) GetModelInfo() (data Response, err error) {
	return u.RequestAndParse("/aircon/get_model_info", "GET", nil)
}
