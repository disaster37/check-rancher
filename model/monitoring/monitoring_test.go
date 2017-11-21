package modelMonitoring

import (
	"github.com/disaster37/check-rancher/model/perfdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test the constructor
func TestNewMonitoring(t *testing.T) {
	md := NewMonitoring()
	assert.Equal(t, 0, md.Status())
	assert.Equal(t, make([]string, 0), md.messages)
	assert.Equal(t, make(modelPerfdata.Perfdatas, 0), md.perfdatas)
}

// Test We can get and set status
func TestGetSetStatus(t *testing.T) {

	monitoringData := NewMonitoring()

	// Normal use case
	monitoringData.SetStatus(1)
	assert.Equal(t, 1, monitoringData.Status())

	// Bad use case
	assert.Error(t, monitoringData.SetStatus(-1))
	assert.Error(t, monitoringData.SetStatus(5))
}

// Test we can get and set message
func TestGetAddMessages(t *testing.T) {
	md := NewMonitoring()
	md.AddMessage("test1")
	md.AddMessage("test2")

	assert.Equal(t, 2, len(md.Messages()))
	assert.Equal(t, "test2", md.Messages()[1])
	message, _ := md.Message(1)
	assert.Equal(t, "test2", message)
}

// Test we can get and set perfdata
func TestGetAddPerfdatas(t *testing.T) {
	md := NewMonitoring()
	md.AddPerfdata("test1", 1, "")
	md.AddPerfdata("test2", 1, "")

	assert.Equal(t, 2, len(md.Perfdatas()))
	assert.Equal(t, "test2", md.Perfdatas()[1].Label())
	perfdata, _ := md.Perfdata(1)
	assert.Equal(t, "test2", perfdata.Label())
}
