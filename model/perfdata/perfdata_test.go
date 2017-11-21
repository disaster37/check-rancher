package modelPerfdata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test the constructor
func TestNewPerfdata(t *testing.T) {

	// Normal use case
	p, _ := NewPerfdata("test", 1, "Kb")
	assert.Equal(t, "test", p.label)
	assert.Equal(t, 1, p.value)
	assert.Equal(t, "Kb", p.unit)

	// Bad use case
	p, err := NewPerfdata("", 1, "Kb")
	assert.Error(t, err)

}

// Test get and set label
func TestGetSetLabel(t *testing.T) {

	// Normal use case
	p, _ := NewPerfdata("test", 1, "kb")
	p.SetLabel("test2")
	assert.Equal(t, "test2", p.Label())

	// Bad use case
	p, err := NewPerfdata("", 1, "kb")
	assert.Error(t, err)
}

// Test get and set value
func TestGetSetValue(t *testing.T) {
	p, _ := NewPerfdata("test", 1, "kb")
	p.SetValue(2)
	assert.Equal(t, 2, p.Value())
}

// Test get and set unit
func TestGetSetUnit(t *testing.T) {
	p, _ := NewPerfdata("test", 1, "kb")

	p.SetUnit("gb")
	assert.Equal(t, "gb", p.Unit())
}
