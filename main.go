package main

import (
	"github.com/gorilla/mux"
	"github.com/sjhitchner/trackly/infrastructure/rest"
	restCommon "github.com/sjhitchner/trackly/interfaces/rest/common"
	restTracking "github.com/sjhitchner/trackly/interfaces/rest/pixel"
	restRedirect "github.com/sjhitchner/trackly/interfaces/rest/redirect"
	"github.com/sjhitchner/trackly/usecases/redirect"
	"github.com/sjhitchner/trackly/usecases/tracking"
	"log"
	"net/http"
	"time"
)

var (
	VERSION string = "development"

	redirectUrl         string
	seppukuDelaySeconds time.Duration
	pixelCookieName     string
	pixelCookieDuration time.Duration
	pixelCookieSecret   string
)

func init() {

}

func main() {
	deviceTokenEncrypter := restCommon.NewDeviceTokenEncrypter(
		"5764328pn341194qp4nq89pp99pn0qsp",
		"592baef1a02bfaf13d63b280c38775b9",
	)
	if err := deviceTokenEncrypter.Validate(); err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", http.NotFound)
	router.Handle("/", http.RedirectHandler(redirectUrl, http.StatusOK))

	trackingInteractor := tracking.NewPixelTrackingInteractor()
	redirectInteractor := redirect.NewRedirectInteractor()

	trackingResource := restTracking.NewTrackingResource(deviceTokenEncrypter, trackingInteractor)
	trackingResource.Register(router)

	redirectResource := restRedirect.NewRedirectResource(deviceTokenEncrypter, redirectInteractor)
	redirectResource.Register(router)

	router.Path("/ping").
		Methods("GET").
		HandlerFunc(rest.SeppukuHealthCheck(VERSION, seppukuDelaySeconds))

	router.Path("/version").
		Methods("GET").
		HandlerFunc(rest.VersionHandler(VERSION))

	http.Handle("/", router)

	// TODO: customize server
	log.Fatal(http.ListenAndServe(":8000", router))
}
