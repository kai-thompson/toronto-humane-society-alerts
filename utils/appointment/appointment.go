package appointment

import (
	"fmt"
	"os"
	"strings"

	"github.com/kai-thompson/toronto-humane-society-alerts/utils/helpers"
)

type ID int8

const (
	WellnessServiceID ID = iota

	CanineNeuterID

	CanineSpayID

	FelineNeuterID

	FelineSpayID

	DentalServiceID
)

type Appointment struct {
	id  ID
	name string
}

func (a *Appointment) ID() ID {
	return a.id
}

func (a *Appointment) Name() string {
	return a.name
}

func (a *Appointment) TopicARN() string {
	arnEnvPrefix := strings.ToUpper(helpers.ToSnakeCase(string(a.name)))
	return os.Getenv(fmt.Sprintf("%s_TOPIC_ARN", arnEnvPrefix))
}

func New(id ID) (*Appointment, error) {
	appointmentName, ok := appointmentIDToName[id]
	if !ok {
		return nil, fmt.Errorf("appointment with ID %d does not exist", id)
	}

	return &Appointment{id: id, name: appointmentName}, nil
}

func IDs() []ID {
	ids := make([]ID, 0, len(appointmentIDToName))
	for id := range appointmentIDToName {
		ids = append(ids, id)
	}

	return ids
}

var appointmentIDToName = map[ID]string{
	WellnessServiceID: "Wellness Service",
	CanineNeuterID:    "Canine Neuter",
	CanineSpayID:      "Canine Spay",
	FelineNeuterID:    "Feline Neuter",
	FelineSpayID:      "Feline Spay",
	DentalServiceID:   "Dental Service",
}
