package modelPerfdata

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Perfdata struct
type Perfdata struct {
	unit  string
	value int
	label string
}

type Perfdatas []*Perfdata

func NewPerfdata(label string, value int, unit string) (*Perfdata, error) {

	log.Debugf("label: %s, value: %d, unit: %s", label, value, unit)

	if label == "" {
		return nil, errors.New("Label can't be empty")
	}

	perfdata := &Perfdata{
		label: label,
		value: value,
		unit:  unit,
	}

	return perfdata, nil
}

// Permit to get unit
func (p *Perfdata) Unit() string {
	return p.unit
}

// Permit to set unit
func (p *Perfdata) SetUnit(unit string) {
	log.Debug("Unit: ", unit)

	p.unit = unit
}

// Permit to get label
func (p *Perfdata) Label() string {
	return p.label
}

// Permit to set label
func (p *Perfdata) SetLabel(label string) error {
	log.Debug("Label: ", label)

	if label == "" {
		return errors.New("Label can't be empty")
	}

	p.label = label

	return nil
}

// Permit to get value
func (p *Perfdata) Value() int {
	return p.value
}

// Permit to set value
func (p *Perfdata) SetValue(value int) {
	log.Debug("Value: ", value)

	p.value = value
}
