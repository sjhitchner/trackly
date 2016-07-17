package redirect

import (
	"github.com/gorilla/mux"
	. "github.com/sjhitchner/trackly/domain"
	"github.com/sjhitchner/trackly/infrastructure/rest"
	"github.com/sjhitchner/trackly/interfaces/rest/common"
	"log"
	"net/http"
)

const (
	ShortUrlIdKey = "url_id"

	RedirectPath = "/{url_id:[a-zA-Z0-9]{6}}"
)

type RedirectResourceImpl struct {
	*common.BaseResourceImpl
	interactor RedirectInteractor
}

func NewRedirectResource(dte DeviceTokenEncrypter, interactor RedirectInteractor) *RedirectResourceImpl {
	return &RedirectResourceImpl{
		common.NewBaseResource(dte),
		interactor,
	}
}

func (t RedirectResourceImpl) Register(router *mux.Router) {
	router.Path(RedirectPath).
		Methods("GET").
		HandlerFunc(rest.RecoverOnPanic(t.TrackRedirect))
}

func (t RedirectResourceImpl) TrackRedirect(response http.ResponseWriter, request *http.Request) {

	cookie, deviceId, deviceInfo := t.GetDeviceData(request)
	shortUrl := getShortUrl(request)

	urlRedirect, err := t.interactor.Redirect(shortUrl, deviceId, deviceInfo)
	if err != nil {
		log.Println("Unable to find long url", err)
	}

	if err := t.SetDeviceData(response, cookie, deviceId); err != nil {
		log.Println("Unable to set cookie", err)
	}

	log.Printf("Redirecting %s => %s\n", shortUrl, urlRedirect)
	http.Redirect(response, request, string(urlRedirect), http.StatusFound)
}

func getShortUrl(request *http.Request) ShortURL {
	return ShortURL(mux.Vars(request)[ShortUrlIdKey])
}
