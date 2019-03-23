package jolokia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type IData interface {
	Print()
}

type Request struct {
	MBean string `json:"mbean"`
	Attribute string `json:"attribute"`
	Type string `json:"type"`
}

type Data struct {
	Request *Request `json:"request"`
	Value string `json:"value"`
	Timestamp int `json:"timestamp"`
	Status int `json:"status"`
}

func (data *Data) Print() {
	fmt.Println(data.Request)
	fmt.Println(data.Value)
	fmt.Println(data.Timestamp)
	fmt.Println(data.Status)
}

type IJolokia interface {
	NewJolokia(_ip string, _port string, _path string) (*Jolokia)
	Read(_path string) (*Data)
	Print(data *Data)
}

type Jolokia struct {
	ip			string
	port		string
	jolokiaURL  string
}

func NewJolokia(_ip string, _port string, _path string) (*Jolokia) {
	j := Jolokia {
		ip: _ip,
		port: _port,
		jolokiaURL: "http://admin:admin@" + _ip + ":" + _port + _path,
	}

	return &j
}

func (j *Jolokia) Read(_path string) (*Data) {

	url := j.jolokiaURL + "/read/" + _path

	jdata := &Data{
		Request: &Request{},
	}

	jolokiaClient := http.Client{
		Timeout: time.Second *2, // Maximum of 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("err")
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "activemq-artemis-management")
	res, getErr := jolokiaClient.Do(req)
	if getErr != nil {
		fmt.Println("getErr")
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("readErr")
		log.Fatal(readErr)
	}

	bodyString := string(body)
	//fmt.Println("bodyString:" + bodyString)
	unmarshalErr := json.Unmarshal([]byte(bodyString), jdata)
	if unmarshalErr != nil {
		fmt.Println("unmarshalErr")
		log.Fatal(unmarshalErr)
	}

	return jdata
}

