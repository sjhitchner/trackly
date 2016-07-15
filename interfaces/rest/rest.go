package rest

import (
	"github.com/gorilla/mux"
	"github.com/sjhitchner/trackly/domain"
	"github.com/sjhitchner/trackly/infrastructure/rest"
	"log"
	"net/http"
	"time"
)

const (
	PropertyIdKey       = "propery_id"
	PixelPath           = "/{property_id:[a-zA-Z0-9]+}/pixel.gif"
	GifMime             = "image/gif"
	GifTransparentBytes = "GIF89a\x01\x00\x01\x00\xf0\x00\x00\xff\xff\xff\x00\x00\x00!"

	HeaderUserAgent      = "User-Agent"
	HeaderAcceptLanguage = "Accept-Language"

	CookieName     = "gid-1"
	CookieDuration = 365 * 24 * time.Hour
)

type PixelTrackingResourceImpl struct {
	interactor domain.PixelTrackingInteractor
}

func NewPixelTrackingResource(interactor domain.PixelTrackingInteractor) *PixelTrackingResourceImpl {
	return &PixelTrackingResourceImpl{
		interactor,
	}
}

func (t PixelTrackingResourceImpl) Register(router *mux.Router) {
	router.NotFoundHandler = http.NotFoundHandler()

	router.Path(PixelPath).
		Methods("GET").
		HandlerFunc(rest.RecoverOnPanic(t.TrackPixel))
}

func (t PixelTrackingResourceImpl) TrackPixel(response http.ResponseWriter, request *http.Request) {
	defer ServeGif(response)

	propertyId := getPropertyId(request)
	cookie := getTrackingCookie(request)

	userAgent := domain.UserAgent(request.Header.Get(HeaderUserAgent))
	acceptLanguage := domain.AcceptLanguage(request.Header.Get(HeaderAcceptLanguage))

	deviceToken, err := t.interactor.Track(
		propertyId,
		domain.DeviceToken(cookie.Value),
		userAgent,
		acceptLanguage,
	)
	if err != nil {
		log.Println("Unable to track", err)
	}

	cookie.Value = string(deviceToken)
	cookie.Expires = time.Now().Add(CookieDuration)

	http.SetCookie(response, cookie)
}

func getPropertyId(request *http.Request) domain.PropertyId {
	return domain.PropertyId(mux.Vars(request)[PropertyIdKey])
}

func getTrackingCookie(request *http.Request) *http.Cookie {
	for _, cookie := range request.Cookies() {
		log.Println(cookie.Name)
	}

	cookie, err := request.Cookie(CookieName)
	if err != nil {
		cookie = &http.Cookie{
			Name:  CookieName,
			Value: "",
		}
	}
	return cookie
}

/*
func getPropertyID(request *http.Request) PropertyID {
	params := request.URL.Query()

	id := params.Get(cookieName)
	if id == "" {
		return ""
	}
	return PropertyID(pid)
}
*/

func ServeGif(response http.ResponseWriter) {
	response.Header().Add("Content-Type", GifMime)
	// response.WriteHeader(http.StatusOK)
	response.Write([]byte(GifTransparentBytes))
}
