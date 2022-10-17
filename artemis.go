package artemis

import (
	"fmt"
	"strings"

	"github.com/artemiscloud/activemq-artemis-management/jolokia"
)

const (
	QUEUE_ALREADY_EXISTS   = "AMQ229019"
	ADDRESS_ALREADY_EXISTS = "AMQ229204"
	UNKNOWN_ERROR          = "AMQ_UNKNOWN"
)

func GetCreationError(jdata *jolokia.ResponseData) string {
	if jdata == nil {
		return UNKNOWN_ERROR
	}
	if strings.Contains(jdata.Error, QUEUE_ALREADY_EXISTS) {
		return QUEUE_ALREADY_EXISTS
	}
	if strings.Contains(jdata.Error, ADDRESS_ALREADY_EXISTS) {
		return ADDRESS_ALREADY_EXISTS
	}
	return UNKNOWN_ERROR
}

type IArtemis interface {
	NewArtemis(_ip string, _jolokiaPort string, _name string, _userName string, _password string) *Artemis
	Uptime() (*jolokia.ResponseData, error)
	CreateQueue(addressName string, queueName string) (*jolokia.ResponseData, error)
	DeleteQueue(queueName string) (*jolokia.ResponseData, error)
	ListBindingsForAddress(addressName string) (*jolokia.ResponseData, error)
	DeleteAddress(addressName string) (*jolokia.ResponseData, error)
	CreateQueueFromConfig(queueConfig string, ignoreIfExists bool) (jolokia.ResponseData, error)
	UpdateQueue(queueConfig string) (jolokia.ResponseData, error)
}

type Artemis struct {
	ip          string
	jolokiaPort string
	name        string
	userName    string
	password    string
	jolokia     jolokia.IJolokia
}

func NewArtemis(_ip string, _jolokiaPort string, _name string, _user string, _password string) *Artemis {
	return GetArtemis(_ip, _jolokiaPort, _name, _user, _password, "http")
}

func GetArtemis(_ip string, _jolokiaPort string, _name string, _user string, _password string, _protocol string) *Artemis {

	artemis := Artemis{
		ip:          _ip,
		jolokiaPort: _jolokiaPort,
		name:        _name,
		jolokia:     jolokia.GetJolokia(_ip, _jolokiaPort, "/console/jolokia", _user, _password, _protocol),
	}

	return &artemis
}

func (artemis *Artemis) Uptime() (*jolokia.ResponseData, error) {

	uptimeURL := "org.apache.activemq.artemis:broker=\"" + artemis.name + "\"/Uptime"
	data, err := artemis.jolokia.Read(uptimeURL)

	return data, err
}

func (artemis *Artemis) GetStatus() (string, error) {
	url := "org.apache.activemq.artemis:broker=\"" + artemis.name + "\"/Status"
	resp, err := artemis.jolokia.Read(url)
	if err != nil || resp == nil {
		return "", err
	}
	if resp.Status != 200 {
		return "", fmt.Errorf("unable to retrieve status %v", resp.Error)
	}
	return resp.Value, nil
}

func (artemis *Artemis) CreateQueue(addressName string, queueName string, routingType string) (*jolokia.ResponseData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	routingType = strings.ToUpper(routingType)
	parameters := `"` + addressName + `","` + queueName + `",` + `"` + routingType + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"createQueue(java.lang.String,java.lang.String,java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}

func (artemis *Artemis) UpdateQueue(queueConfig string) (*jolokia.ResponseData, error) {
	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	parameters := queueConfig
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"updateQueue(java.lang.String)","arguments":[` + parameters + `]` + ` }`

	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err

}

func (artemis *Artemis) CreateQueueFromConfig(queueConfig string, ignoreIfExists bool) (*jolokia.ResponseData, error) {
	var ignoreIfExistsValue string
	if ignoreIfExists {
		ignoreIfExistsValue = "true"
	} else {
		ignoreIfExistsValue = "false"
	}
	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	parameters := queueConfig + `,` + ignoreIfExistsValue
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"createQueue(java.lang.String,boolean)","arguments":[` + parameters + `]` + ` }`

	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}

func (artemis *Artemis) CreateAddress(addressName string, routingType string) (*jolokia.ResponseData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	routingType = strings.ToUpper(routingType)
	parameters := `"` + addressName + `","` + routingType + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"createAddress(java.lang.String,java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}

func (artemis *Artemis) DeleteQueue(queueName string) (*jolokia.ResponseData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	parameters := `"` + queueName + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"destroyQueue(java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}

func (artemis *Artemis) ListBindingsForAddress(addressName string) (*jolokia.ResponseData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	parameters := `"` + addressName + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"listBindingsForAddress(java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}

func (artemis *Artemis) DeleteAddress(addressName string) (*jolokia.ResponseData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	parameters := `"` + addressName + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"deleteAddress(java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}
