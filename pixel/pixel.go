package pixel

import (
	"github.com/sjhitchner/trackly/domain"
	"time"
)

const (
	SignatureSecret = "test"
	EncryptionKey   = "qwerty"
)

type PixelTrackingInteractor struct {
}

func (t PixelTrackingInteractor) Track(propertyId PropertyId, deviceToken DeviceToken, headers http.Header) (DeviceToken, error) {

	if err := deviceToken.Validate(); err != nil {
	}

	deviceId, lastSeen, err := deviceToken.Decrypt()
	if err != nil {
		deviceId = NewDeviceId()
		lastSeen = time.Now().UTC()
	}

	t.RecordImpression(propertyId, deviceId, headers)

	deviceToken := deviceId.Encrypt()
	return deviceToken, nil
}

func (t PixelTrackingInteractor) RecordImpression(propertyId PropertyId, deviceId DeviceId, headers http.Header) {

}
