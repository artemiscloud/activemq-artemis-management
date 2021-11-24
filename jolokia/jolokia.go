package jolokia

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type IData interface {
	Print()
}

type ResponseData struct {
	Status    int
	Value     string
	ErrorType string
	Error     string
}

type ReadRequest struct {
	MBean     string `json:"mbean"`
	Attribute string `json:"attribute"`
	Type      string `json:"type"`
}

type JolokiaError struct {
	HttpCode int
	Message  string
}

func (j *JolokiaError) Error() string {
	return fmt.Sprintf("HTTP STATUS %v. Message: %v", j.HttpCode, j.Message)
}

type IJolokia interface {
	NewJolokia(_ip string, _port string, _path string, _user string, _password string) *Jolokia
	Read(_path string) (*ResponseData, error)
	Exec(_path string) (*ResponseData, error)
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
	} else {
		//encode password in case it has special chars
		j.password = url.QueryEscape(j.password)
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

func (j *Jolokia) Read(_path string) (*ResponseData, error) {

	url := j.protocol + "://" + j.user + ":" + j.password + "@" + j.jolokiaURL + "/read/" + _path

	jolokiaClient := j.getClient()

	var err error = nil
	var jdata *ResponseData = nil

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

		//decoding
		result, _, err := decodeResponseData(res)
		if err != nil {
			return result, err
		}

		jdata = result

		//before decoding the body, we need to check the http code
		err = CheckResponse(res, result)

		break
	}

	return jdata, err
}

func (j *Jolokia) Exec(_path string, _postJsonString string) (*ResponseData, error) {

	url := j.protocol + "://" + j.user + ":" + j.password + "@" + j.jolokiaURL + "/exec/" + _path

	jolokiaClient := j.getClient()

	var jdata *ResponseData
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

		//decoding
		result, _, err := decodeResponseData(res)
		if err != nil {
			return result, err
		}

		jdata = result

		err = CheckResponse(res, result)
		if err != nil {
			execErr = err
		}

		break
	}

	return jdata, execErr
}

func CheckResponse(resp *http.Response, jdata *ResponseData) error {

	if isResponseSuccessful(resp.StatusCode) {
		//that doesn't mean it's ok, check further
		if isResponseSuccessful(jdata.Status) {
			return nil
		}
		errCode := jdata.Status
		errType := jdata.ErrorType
		errMsg := jdata.Error
		errData := jdata.Value
		internalErr := fmt.Errorf("Error response code %v, type %v, message %v and data %v", errCode, errType, errMsg, errData)
		return internalErr
	}
	return &JolokiaError{
		HttpCode: resp.StatusCode,
		Message:  " Error: " + resp.Status,
	}
}

func isResponseSuccessful(httpCode int) bool {
	return httpCode >= 200 && httpCode <= 299
}

func decodeResponseData(resp *http.Response) (*ResponseData, map[string]interface{}, error) {
	result := &ResponseData{}
	rawData := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		return nil, rawData, err
	}

	//fill in response data
	if v, ok := rawData["error"]; ok {
		if v != nil {
			result.Error = v.(string)
		}
	}
	if v, ok := rawData["error_type"]; ok {
		if v != nil {
			result.ErrorType = v.(string)
		}
	}
	if v, ok := rawData["status"]; ok {
		if v != nil {
			result.Status = int(v.(float64))
		}
	}
	if v, ok := rawData["value"]; ok {
		if v != nil {
			result.Value = v.(string)
		}
	}

	return result, rawData, nil

}
