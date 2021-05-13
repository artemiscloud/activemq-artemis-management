package jolokia

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
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

type ResponseData interface {
	GetStatusCode() int
	GetValue() string
	GetErrorType() string
	GetError() string
}

type ExecData struct {
	Request   *ExecRequest `json:"request"`
	Value     string       `json:"value"`
	Timestamp int          `json:"timestamp"`
	ErrorType string       `json:"error_type"`
	Error     string       `json:"error"`
	Status    int          `json:"status"`
}

func (e *ExecData) GetStatusCode() int {
	return e.Status
}

func (e *ExecData) GetValue() string {
	return e.Value
}

func (e *ExecData) GetErrorType() string {
	return e.ErrorType
}

func (e *ExecData) GetError() string {
	return e.Error
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
	ErrorType string       `json:"error_type"`
	Error     string       `json:"error"`
	Status    int          `json:"status"`
}

func (r *ReadData) GetStatusCode() int {
	return r.Status
}

func (r *ReadData) GetValue() string {
	return r.Value
}

func (r *ReadData) GetErrorType() string {
	return r.ErrorType
}

func (r *ReadData) GetError() string {
	return r.Error
}

type JolokiaError struct {
	HttpCode int
	Message  string
}

func (j *JolokiaError) Error() string {
	return fmt.Sprintf("HTTP STATUS %v. Message: %v", j.HttpCode, j.Message)
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
	NewJolokia(_ip string, _port string, _path string, _user string, _password string) *Jolokia
	Read(_path string) (*ReadData, error)
	Exec(_path string) (*ExecData, error)
	Print(data *ReadData)
}

type Jolokia struct {
	ip         string
	port       string
	jolokiaURL string
	user       string
	password   string
	protocol   string
}

func NewJolokia(_ip string, _port string, _path string, _user string, _password string) *Jolokia {
	return GetJolokia(_ip, _port, _path, _user, _password, "http")
}

func GetJolokia(_ip string, _port string, _path string, _user string, _password string, _protocol string) *Jolokia {

	j := Jolokia{
		ip:         _ip,
		port:       _port,
		jolokiaURL: _ip + ":" + _port + _path,
		user:       _user,
		password:   _password,
		protocol:   _protocol,
	}
	if j.user == "" {
		j.user = "admin"
	}
	if j.password == "" {
		j.password = "admin"
	}

	return &j
}

func (j *Jolokia) getClient() *http.Client {
	if j.protocol == "https" {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: time.Second * 2, //Maximum of 2 seconds
		}
	}
	return &http.Client{
		Timeout: time.Second * 2, // Maximum of 2 seconds
	}
}

func (j *Jolokia) Read(_path string) (*ReadData, error) {

	url := j.protocol + "://" + j.user + ":" + j.password + "@" + j.jolokiaURL + "/read/" + _path

	jdata := &ReadData{
		Request: &ReadRequest{},
	}

	jolokiaClient := j.getClient()

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

		//before decoding the body, we need to check the http code
		err = CheckResponse(res, jdata)

		break
	}

	return jdata, err
}

func (j *Jolokia) Exec(_path string, _postJsonString string) (*ExecData, error) {

	url := j.protocol + "://" + j.user + ":" + j.password + "@" + j.jolokiaURL + "/exec/" + _path

	jdata := &ExecData{
		Request: &ExecRequest{},
	}

	jolokiaClient := j.getClient()

	var execErr error = nil
	for {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(_postJsonString)))
		if err != nil {
			execErr = err
			break
		}

		req.Header.Set("User-Agent", "activemq-artemis-management")
		req.Header.Set("Content-Type", "application/json")
		res, err := jolokiaClient.Do(req)

		if err != nil {
			execErr = err
			break
		}

		defer res.Body.Close()

		//before decoding the body, we need to check the http code
		err = CheckResponse(res, jdata)
		if err != nil {
			execErr = err
		}

		break
	}

	return jdata, execErr
}

func CheckResponse(resp *http.Response, jdata ResponseData) error {

	if isResponseSuccessful(resp.StatusCode) {
		//that doesn't mean it's ok, check further
		if err := json.NewDecoder(resp.Body).Decode(jdata); err != nil {
			return err
		}
		if isResponseSuccessful(jdata.GetStatusCode()) {
			return nil
		}
	}
	return &JolokiaError{
		HttpCode: jdata.GetStatusCode(),
		Message:  " Error: " + jdata.GetError(),
	}
}

func isResponseSuccessful(httpCode int) bool {
	return httpCode >= 200 && httpCode <= 299
}
