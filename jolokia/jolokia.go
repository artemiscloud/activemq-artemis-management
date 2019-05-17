package jolokia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type IData interface {
	Print()
}

type ExecRequest struct {
	MBean     string   `json:"mbean"`
	Arguments []string `json:"arguments"`
	Type      string   `json:"type"`
	Operation string   `json:"operation"`
}

type ExecData struct {
	Request   *ExecRequest `json:"request"`
	Value     string       `json:"value"`
	Timestamp int          `json:"timestamp"`
	Status    int          `json:"status"`
}

type ReadRequest struct {
	MBean     string `json:"mbean"`
	Attribute string `json:"attribute"`
	Type      string `json:"type"`
}

type ReadData struct {
	Request   *ReadRequest `json:"request"`
	Value     string       `json:"value"`
	Timestamp int          `json:"timestamp"`
	Status    int          `json:"status"`
}

func (data *ReadData) Print() {
	fmt.Println(data.Request)
	fmt.Println(data.Value)
	fmt.Println(data.Timestamp)
	fmt.Println(data.Status)
}

func (data *ExecData) Print() {
	fmt.Println(data.Request)
	fmt.Println(data.Value)
	fmt.Println(data.Timestamp)
	fmt.Println(data.Status)
}

type IJolokia interface {
	NewJolokia(_ip string, _port string, _path string) *Jolokia
	Read(_path string) (*ReadData, error)
	Exec(_path string) (*ExecData, error)
	Print(data *ReadData)
}

type Jolokia struct {
	ip         string
	port       string
	jolokiaURL string
}

func NewJolokia(_ip string, _port string, _path string) *Jolokia {

	j := Jolokia{
		ip:         _ip,
		port:       _port,
		jolokiaURL: "http://admin:admin@" + _ip + ":" + _port + _path,
	}

	return &j
}

func (j *Jolokia) Read(_path string) (*ReadData, error) {

	url := j.jolokiaURL + "/read/" + _path

	jdata := &ReadData{
		Request: &ReadRequest{},
	}

	jolokiaClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 seconds
	}

	var err error = nil
	for {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			break
		}
		req.Header.Set("User-Agent", "activemq-artemis-management")

		res, err := jolokiaClient.Do(req)

		if err != nil {
			break
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			break
		}

		bodyString := string(body)
		err = json.Unmarshal([]byte(bodyString), jdata)
		if err != nil {
			break
		}

		break
	}

	return jdata, err
}

func (j *Jolokia) Exec(_path string, _postJsonString string) (*ExecData, error) {

	url := j.jolokiaURL + "/exec/" + _path

	jdata := &ExecData{
		Request: &ExecRequest{},
	}

	jolokiaClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 seconds
	}

	var err error = nil
	for {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(_postJsonString)))
		if err != nil {
			break
		}

		req.Header.Set("User-Agent", "activemq-artemis-management")
		req.Header.Set("Content-Type", "application/json")
		res, err := jolokiaClient.Do(req)

		if err != nil {
			break
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			break
		}

		bodyString := string(body)
		err = json.Unmarshal([]byte(bodyString), jdata)
		if err != nil {
			break
		}

		break
	}

	return jdata, err
}
