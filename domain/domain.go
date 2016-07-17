package domain

import (
	"github.com/pborman/uuid"
)

type PropertyId string
type DeviceId string

type IP string
type UserAgent string
type AcceptLanguage string

type LongURL string
type ShortURL string

type DeviceInfo struct {
	IP        IP
	UserAgent UserAgent
	Language  AcceptLanguage
}

type PixelTrackingInteractor interface {
	Track(PropertyId, DeviceId, *DeviceInfo) error
}

type DeviceTokenEncrypter interface {
	Decrypt(deviceToken string) (DeviceId, error)
	Encrypt(deviceId DeviceId) (string, error)
	Validate() error
}

type RedirectInteractor interface {
	Redirect(ShortURL, DeviceId, *DeviceInfo) (LongURL, error)
	Register(url LongURL) (ShortURL, error)
}

func NewDeviceId() DeviceId {
	return DeviceId(uuid.New())
}
