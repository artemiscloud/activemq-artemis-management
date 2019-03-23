package artemis

import (
	"github.com/roddiekieley/activemq-artemis-management/jolokia"
)

type IArtemis interface {
	NewArtemis(_ip string, _jolokiaPort string, _name string) (*Artemis)
	Uptime()
}

type Artemis struct {
	ip 			string
	jolokiaPort string
	name 		string
	jolokia 	*jolokia.Jolokia
}

func NewArtemis(_ip string, _jolokiaPort string, _name string) (*Artemis) {
	artemis := Artemis {
		ip: _ip,
		jolokiaPort: _jolokiaPort,
		name: _name,
		jolokia: jolokia.NewJolokia(_ip, _jolokiaPort, "/console/jolokia"),
	}

	return &artemis
}

func (artemis *Artemis) Uptime() (*jolokia.Data) {

	uptimeURL := "org.apache.activemq.artemis:broker=\"" + artemis.name + "\"/Uptime"
	data := artemis.jolokia.Read(uptimeURL)

	return data
}
