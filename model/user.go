package model

import (
	"errors"
	"gorm.io/gorm"
)

var ErrUserAlreadyExist = errors.New("user already exist")

var ErrUserDeviceAlreadyExist = errors.New("user device already exist")

var ErrUserDeviceNotExist = errors.New("user device not exist")

type User struct {
	gorm.Model
	Username      string
	Password      string
	Devices       string
	DeviceRemarks string
	DeviceTypes   string
}

func NewUser(pUsername, pPassword string) (*User, error) {
	var err error
	// check exist
	u, err := FindUserByUsername(pUsername)
	if err == nil {
		return nil, ErrUserAlreadyExist
	}
	// internal error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// add new
	u = &User{
		Username: pUsername,
		Password: pPassword,
	}
	if err = db.Create(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func FindUserByUsername(pUsername string) (*User, error) {
	var user User
	err := db.Where("username = ?", pUsername).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) AddDevice(pClientId string) error {
	var err error
	if !StringArrayContains(u.Devices, pClientId) {
		u.Devices, err = AppendStringArray(u.Devices, pClientId)
		if err != nil {
			return err
		}
	} else {
		return ErrUserDeviceAlreadyExist
	}
	return db.Save(u).Error
}

func (u *User) DeleteDevice(pClientId string) error {
	var err error
	var f = false
	devices := DecodeStringArray(u.Devices)
	var nds []string
	for _, v := range devices {
		if v != pClientId {
			nds = append(nds, v)
			continue
		}
		f = true
	}
	if !f {
		return ErrUserDeviceNotExist
	}
	u.Devices, err = EncodeStringArray(nds)
	if err != nil {
		return err
	}
	return db.Save(u).Error
}

func (u *User) GetDevices() []string {
	return DecodeStringArray(u.Devices)
}

func (u *User) GetDeviceRemarks() map[string]string {
	m := make(map[string]string)
	aKeyPairStr := DecodeStringArray(u.DeviceRemarks)
	for _, vs := range aKeyPairStr {
		k, v := DecodeSsKeyPair(vs)
		m[k] = v
	}
	return m
}

// DeleteDeviceRemark
// return nil or ErrUserDeviceNotExist
func (u *User) DeleteDeviceRemark(pClientId string) error {
	var f = false
	remarks := u.GetDeviceRemarks()
	var t [][]string
	for k, v := range remarks {
		if k != pClientId {
			t = append(t, []string{k, v})
			continue
		}
		f = true
	}
	if !f {
		return ErrUserDeviceNotExist
	}
	var kpe []string
	for _, v := range t {
		e, _ := EncodeSsKeyPair(v[0], v[1])
		kpe = append(kpe, e)
	}
	u.DeviceRemarks, _ = EncodeStringArray(kpe)
	return nil
}

// SetDeviceRemark
// if device not exist, create new or change it to the new one.
// return nil or ErrInvalidDataFormat
func (u *User) SetDeviceRemark(pClientId, pRemark string) error {
	aKeyPairStr := DecodeStringArray(u.DeviceRemarks)
	enc, err := EncodeSsKeyPair(pClientId, pRemark)
	if err != nil {
		return err
	}
	for i := 0; i < len(aKeyPairStr); i++ { // if exist
		k, _ := DecodeSsKeyPair(aKeyPairStr[i])
		if k == pClientId {
			aKeyPairStr[i] = enc
			u.DeviceRemarks, _ = EncodeStringArray(aKeyPairStr)
			return db.Save(u).Error
		}
	}
	// if not exist
	aKeyPairStr = append(aKeyPairStr, enc)
	u.DeviceRemarks, _ = EncodeStringArray(aKeyPairStr)
	return db.Save(u).Error
}

func (u *User) GetDeviceTypes() map[string]string {
	m := make(map[string]string)
	aKeyPairStr := DecodeStringArray(u.DeviceTypes)
	for _, vs := range aKeyPairStr {
		k, v := DecodeSsKeyPair(vs)
		m[k] = v
	}
	return m
}

// SetDeviceType
// if device type not exist, create new or change it to the new one.
// return nil or ErrInvalidDataFormat
func (u *User) SetDeviceType(pClientId, pType string) error {
	aKeyPairStr := DecodeStringArray(u.DeviceTypes)
	enc, err := EncodeSsKeyPair(pClientId, pType)
	if err != nil {
		return err
	}
	for i := 0; i < len(aKeyPairStr); i++ { // if exist
		k, _ := DecodeSsKeyPair(aKeyPairStr[i])
		if k == pClientId {
			aKeyPairStr[i] = enc
			u.DeviceTypes, _ = EncodeStringArray(aKeyPairStr)
			return db.Save(u).Error
		}
	}
	// if not exist
	aKeyPairStr = append(aKeyPairStr, enc)
	u.DeviceTypes, _ = EncodeStringArray(aKeyPairStr)
	return db.Save(u).Error
}

// DeleteDeviceType
// return nil or ErrUserDeviceNotExist
func (u *User) DeleteDeviceType(pClientId string) error {
	var f = false
	types := u.GetDeviceTypes()
	var t [][]string
	for k, v := range types {
		if k != pClientId {
			t = append(t, []string{k, v})
			continue
		}
		f = true
	}
	if !f {
		return ErrUserDeviceNotExist
	}
	var kpe []string
	for _, v := range t {
		e, _ := EncodeSsKeyPair(v[0], v[1])
		kpe = append(kpe, e)
	}
	u.DeviceTypes, _ = EncodeStringArray(kpe)
	return db.Save(u).Error
}
