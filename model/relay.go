package model

import (
	"errors"
	"gorm.io/gorm"
)

type Relay struct {
	gorm.Model
	ClientId        string
	HardwareVersion string
	SoftwareVersion string
	Online          bool
	On              bool
	Owner           int
}

func AddRelay(pClientId string) (*Relay, error) {
	// check exist
	relay, err := FindRelayByClientId(pClientId)
	if err == nil {
		relay.Owner += 1
		if err = db.Save(relay).Error; err != nil {
			return nil, err
		}
		return relay, nil
	}
	// internal error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// add new
	relay = &Relay{
		ClientId:        pClientId,
		HardwareVersion: "0.0",
		SoftwareVersion: "0.0",
		Online:          false,
		On:              false,
		Owner:           1,
	}
	if err = db.Create(&relay).Error; err != nil {
		return nil, err
	}
	return relay, nil
}

// FindRelayByClientId
// return gorm.ErrRecordNotFound or internal errors
func FindRelayByClientId(pClientId string) (*Relay, error) {
	var relay Relay
	err := db.Where("client_id = ?", pClientId).First(&relay).Error
	if err != nil {
		return nil, err
	}
	return &relay, nil
}

// GetAllRelayClientId
// return all relay client id and internal error
func GetAllRelayClientId() ([]string, error) {
	var relays []Relay
	var out []string
	if err := db.Find(&relays).Error; err != nil {
		return nil, err
	}
	for _, v := range relays {
		out = append(out, v.ClientId)
	}
	return out, nil
}

// Delete
// return internal error
func (r *Relay) Delete() error {
	// check owner count
	if r.Owner != 1 {
		r.Owner -= 1
		return db.Save(r).Error
	}
	return db.Delete(&r).Error
}

func (r *Relay) Save() error {
	return db.Save(r).Error
}
