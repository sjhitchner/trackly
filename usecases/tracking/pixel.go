package tracking

import (
	//"github.com/pkg/errors"
	. "github.com/sjhitchner/trackly/domain"
	"log"
)

type PixelTrackingInteractorImpl struct {
}

func NewPixelTrackingInteractor() PixelTrackingInteractor {
	return &PixelTrackingInteractorImpl{}
}

func (t *PixelTrackingInteractorImpl) Track(propertyId PropertyId, deviceId DeviceId, deviceInfo *DeviceInfo) error {

	// Enqueue here
	log.Println("Fingering printing", deviceId, deviceInfo.UserAgent, deviceInfo.Language)

	return nil
}
