package tracking

import (
	"github.com/pkg/errors"
	. "github.com/sjhitchner/trackly/domain"
	"log"
)

type PixelTrackingInteractorImpl struct {
	deviceTokenEncrypter DeviceTokenEncrypter
}

func NewPixelTrackingInteractor(deviceTokenEncrypter DeviceTokenEncrypter) PixelTrackingInteractor {
	return &PixelTrackingInteractorImpl{
		deviceTokenEncrypter: deviceTokenEncrypter,
	}
}

func (t *PixelTrackingInteractorImpl) Track(propertyId PropertyId, deviceToken DeviceToken, userAgent UserAgent, acceptLanguage AcceptLanguage) (DeviceToken, error) {

	deviceId, err := t.deviceTokenEncrypter.Decrypt(deviceToken)
	if err != nil {
		log.Println("Unable to decrypt device id. Generating new device id", err)
		deviceId = NewDeviceId()
	}

	// Enqueue here
	log.Println("Fingering printing", deviceId, userAgent, acceptLanguage)

	// Encrypt Device Id if it
	deviceToken, err = t.deviceTokenEncrypter.Encrypt(deviceId)
	if err != nil {
		return "", errors.Wrap(err, "unable to encrypt device id")
	}

	return deviceToken, nil
}
