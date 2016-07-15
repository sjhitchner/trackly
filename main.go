package main

import (
	"github.com/gorilla/mux"
	. "github.com/sjhitchner/trackly/infrastructure/rest"
	"github.com/sjhitchner/trackly/interfaces/rest"
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
	deviceTokenEncrypter := tracking.NewDeviceTokenEncrypter(
		"5764328pn341194qp4nq89pp99pn0qsp",
		"592baef1a02bfaf13d63b280c38775b9",
	)
	if err := deviceTokenEncrypter.Validate(); err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	//router.HandleFunc("/", http.NotFound)
	router.Handle("/", http.RedirectHandler(redirectUrl, http.StatusOK))

	trackingInteractor := tracking.NewPixelTrackingInteractor(deviceTokenEncrypter)

	trackingResource := rest.NewPixelTrackingResource(trackingInteractor)
	trackingResource.Register(router)

	router.Path("/ping").
		Methods("GET").
		HandlerFunc(SeppukuHealthCheck(VERSION, seppukuDelaySeconds))

	router.Path("/version").
		Methods("GET").
		HandlerFunc(VersionHandler(VERSION))

	http.Handle("/", router)

	// TODO: customize server
	log.Fatal(http.ListenAndServe(":8000", router))
}
