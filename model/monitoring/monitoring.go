package modelMonitoring

import (
	"bytes"
	"github.com/disaster37/check-rancher/model/perfdata"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
)

// Nagios status
const (
	STATUS_UNKNOWN  = 3
	STATUS_CRITICAL = 2
	STATUS_WARNING  = 1
	STATUS_OK       = 0
)

// Monitoring data
type Monitoring struct {
	messages  []string
	status    int
	perfdatas modelPerfdata.Perfdatas
}

// Get Monitoring object
// The status is inizialized to 0
func NewMonitoring() *Monitoring {
	monitoringData := &Monitoring{
		status:    STATUS_OK,
		perfdatas: make(modelPerfdata.Perfdatas, 0, 3),
		messages:  make([]string, 0, 3),
	}

	return monitoringData
}

// Set monitoring status
func (m *Monitoring) SetStatus(status int) error {
	log.Debugf("Status: %d", status)
	if status > 3 || status < 0 {
		return errors.New("Status can't be greater than 3")
	}

	if status > m.status {
		log.Debugf("New monitoring status is %d", status)
		m.status = status
	}

	return nil
}

// Get monitoring status
func (m *Monitoring) Status() int {
	return m.status
}

// Add message to display in monitoring tools
func (m *Monitoring) AddMessage(message string) {
	log.Debugf("Message: %s", message)

	m.messages = append(m.messages, message)

}

// Get all messages
func (m *Monitoring) Messages() []string {
	return m.messages
}

// Get message on given index
func (m *Monitoring) Message(index int) (string, error) {
	if index >= len(m.messages) {
		return "", errors.New("Index is out of list messages")
	}

	return m.messages[index], nil
}

// Add perfdata to display in monitoring tools
func (m *Monitoring) AddPerfdata(label string, value int, unit string) error {
	log.Debugf("Label: %s, Value: %d, Unit: %s", label, value, unit)

	perfdata, err := modelPerfdata.NewPerfdata(label, value, unit)
	if err != nil {
		return errors.Wrap(err, "Error appear when tryp to create new perfdata")
	}

	m.perfdatas = append(m.perfdatas, perfdata)

	return nil
}

func (m *Monitoring) Perfdatas() modelPerfdata.Perfdatas {
	return m.perfdatas
}

// Get perfdata on given index
func (m *Monitoring) Perfdata(index int) (*modelPerfdata.Perfdata, error) {
	if index >= len(m.perfdatas) {
		return nil, errors.New("Index is out of list messages")
	}

	return m.perfdatas[index], nil
}

// Get string from monitoring data
func (m *Monitoring) ToString() string {

	var buffer bytes.Buffer

	for idx, message := range m.messages {
		if idx == 0 {
			buffer.WriteString(fmt.Sprintf("%s", message))
		} else {
			buffer.WriteString(fmt.Sprintf("\n%s", message))
		}
	}

	if len(m.Perfdatas()) > 0 {
		buffer.WriteString("|")
		for _, perfdata := range m.Perfdatas() {
			buffer.WriteString(fmt.Sprintf("%s=%d%s;;;; ", perfdata.Label(), perfdata.Value(), perfdata.Unit()))
		}
	}

	return buffer.String()

}

// Print in stdout the monitoring data
func (m *Monitoring) ToSdtOut() {
	var status string
	switch m.Status() {
	case STATUS_UNKNOWN:
		status = "UNKNOWN"
	case STATUS_CRITICAL:
		status = "CRITICAL"
	case STATUS_WARNING:
		status = "WARNING"
	case STATUS_OK:
		status = "OK"
	}

	fmt.Printf("%s - %s\n", status, m.ToString())
	os.Exit(int(m.Status()))
}
