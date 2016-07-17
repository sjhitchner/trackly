package common

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	. "github.com/sjhitchner/trackly/domain"
	"log"
	"net/http"
	"time"
)

const (
	HeaderUserAgent      = "User-Agent"
	HeaderAcceptLanguage = "Accept-Language"

	CookieName     = "gid-1"
	CookieDuration = 365 * 24 * time.Hour
)

type Resource interface {
	Register(*mux.Router)
}

type BaseResourceImpl struct {
	deviceTokenEncrypter DeviceTokenEncrypter
}

func NewBaseResource(dte DeviceTokenEncrypter) *BaseResourceImpl {
	return &BaseResourceImpl{dte}
}

func (t BaseResourceImpl) GetDeviceData(request *http.Request) (*http.Cookie, DeviceId, *DeviceInfo) {
	//for _, cookie := range request.Cookies() {
	//	log.Println(cookie.Name)
	//}

	cookie, err := request.Cookie(CookieName)
	if err != nil {
		cookie = &http.Cookie{
			Name:  CookieName,
			Value: "",
		}
	}

	deviceId, err := t.deviceTokenEncrypter.Decrypt(cookie.Value)
	if err != nil {
		log.Println("Unable to decrypt device id. Generating new device id", err)
		deviceId = NewDeviceId()
	}

	deviceInfo := &DeviceInfo{
		IP:        "test",
		UserAgent: UserAgent(request.Header.Get(HeaderUserAgent)),
		Language:  AcceptLanguage(request.Header.Get(HeaderAcceptLanguage)),
	}

	return cookie, deviceId, deviceInfo
}

func (t BaseResourceImpl) SetDeviceData(response http.ResponseWriter, cookie *http.Cookie, deviceId DeviceId) error {
	deviceToken, err := t.deviceTokenEncrypter.Encrypt(deviceId)
	if err != nil {
		return errors.Wrap(err, "unable to encrypt device id")
	}

	cookie.Value = deviceToken
	cookie.Expires = time.Now().Add(CookieDuration)

	http.SetCookie(response, cookie)

	return nil
}
