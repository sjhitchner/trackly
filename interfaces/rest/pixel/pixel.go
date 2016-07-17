package pixel

import (
	"github.com/gorilla/mux"
	. "github.com/sjhitchner/trackly/domain"
	"github.com/sjhitchner/trackly/infrastructure/rest"
	"github.com/sjhitchner/trackly/interfaces/rest/common"
	"log"
	"net/http"
)

const (
	PropertyIdKey       = "propery_id"
	PixelPath           = "/{property_id:[a-zA-Z0-9]+}/pixel.gif"
	GifMime             = "image/gif"
	GifTransparentBytes = "GIF89a\x01\x00\x01\x00\xf0\x00\x00\xff\xff\xff\x00\x00\x00!"
)

type TrackingResourceImpl struct {
	*common.BaseResourceImpl
	interactor PixelTrackingInteractor
}

func NewTrackingResource(dte DeviceTokenEncrypter, interactor PixelTrackingInteractor) *TrackingResourceImpl {
	return &TrackingResourceImpl{
		common.NewBaseResource(dte),
		interactor,
	}
}

func (t TrackingResourceImpl) Register(router *mux.Router) {
	router.Path(PixelPath).
		Methods("GET").
		HandlerFunc(rest.RecoverOnPanic(t.TrackPixel))
}

func (t TrackingResourceImpl) TrackPixel(response http.ResponseWriter, request *http.Request) {
	defer ServeGif(response)

	cookie, deviceId, deviceInfo := t.GetDeviceData(request)

	propertyId := getPropertyId(request)
	if err := t.interactor.Track(propertyId, deviceId, deviceInfo); err != nil {
		log.Println("Unable to track", err)
	}

	if err := t.SetDeviceData(response, cookie, deviceId); err != nil {
		log.Println("Unable to set cookie", err)
	}
}

func getPropertyId(request *http.Request) PropertyId {
	return PropertyId(mux.Vars(request)[PropertyIdKey])
}

func ServeGif(response http.ResponseWriter) {
	response.Header().Add("Content-Type", GifMime)
	response.Write([]byte(GifTransparentBytes))
}
