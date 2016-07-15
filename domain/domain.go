package domain

import (
	"github.com/pborman/uuid"
)

type PixelTrackingInteractor interface {
	Track(PropertyId, DeviceToken, UserAgent, AcceptLanguage) (DeviceToken, error)
}

type DeviceTokenEncrypter interface {
	Decrypt(deviceToken DeviceToken) (DeviceId, error)
	Encrypt(deviceId DeviceId) (DeviceToken, error)
	Validate() error
}

type PropertyId string

type DeviceToken string

type DeviceId string

func NewDeviceId() DeviceId {
	return DeviceId(uuid.New())
}

type UserAgent string
type AcceptLanguage string
