package jolokia

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
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
	Read(path string) (*ResponseData, error)
	Exec(path, postJsonString string) (*ResponseData, error)
}

type Jolokia struct {
	ip         string
	port       string
	jolokiaURL string
	auth       Auth
	protocol   string
}

type Auth interface {
	GetAuth() string
}

type BasicAuth struct {
	User     string
	Password string
}

type TokenAuth struct {
	Token string
}

func (b *BasicAuth) GetAuth() string {
	creds := b.User + ":" + url.QueryEscape(b.Password)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))
}

func (t *TokenAuth) GetAuth() string {
	return "Bearer " + url.QueryEscape(t.Token)
}

func NewJolokia(ip string, port string, path string, auth Auth) *Jolokia {
	return GetJolokia(ip, port, path, auth, "http")
}

func GetJolokia(ip string, port string, path string, auth Auth, protocol string) *Jolokia {

	j := Jolokia{
		ip:         ip,
		port:       port,
		jolokiaURL: ip + ":" + port + path,
		auth:       auth,
		protocol:   protocol,
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

func (j *Jolokia) Read(path string) (*ResponseData, error) {

	url := j.protocol + "://" + j.jolokiaURL + "/read/" + path
	jolokiaClient := j.getClient()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", j.auth.GetAuth())
	req.Header.Set("User-Agent", "activemq-artemis-management")

	res, err := jolokiaClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//decoding
	result, _, err := decodeResponseData(res)
	if err != nil {
		return result, err
	}

	//before decoding the body, we need to check the http code
	err = CheckResponse(res, result)
	return result, err
}

func (j *Jolokia) Exec(path string, postJsonString string) (*ResponseData, error) {

	url := j.protocol + "://" + j.jolokiaURL + "/exec/" + path
	jolokiaClient := j.getClient()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(postJsonString)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", j.auth.GetAuth())
	req.Header.Set("User-Agent", "activemq-artemis-management")
	req.Header.Set("Content-Type", "application/json")

	res, err := jolokiaClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//decoding
	result, _, err := decodeResponseData(res)
	if err != nil {
		return result, err
	}

	//before decoding the body, we need to check the http code
	err = CheckResponse(res, result)
	return result, err
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
