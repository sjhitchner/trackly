package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

type Response struct {
	Status   int         `json:"status"`
	Version  string      `json:"version"`
	Hostname string      `json:"hostname"`
	Payload  interface{} `json:"payload"`
	Error    string      `json:"error,omitempty"`
}

var hostname string

func init() {
	var err error

	hostname, err = os.Hostname()
	if err != nil {
		log.Println("Could not get hostname", err)
		hostname = "development"
	}
}

type HandlerFunc func(response http.ResponseWriter, request *http.Request)

func VersionHandler(version string) func(http.ResponseWriter, *http.Request) {

	return RecoverOnPanic(func(response http.ResponseWriter, request *http.Request) {
		msg, err := json.Marshal(Response{
			Status:   http.StatusOK,
			Version:  version,
			Hostname: hostname,
		})
		if err != nil {
			log.Println("JSON Error that should never happen", err)
			panic(err)
		}

		response.Write(msg)
		log.Printf(string(msg))
	})
}

func SeppukuHealthCheck(version string, seppukuDelaySeconds time.Duration) func(http.ResponseWriter, *http.Request) {
	healthStatus := http.StatusOK

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(
			ch,
			syscall.SIGINT,
			syscall.SIGHUP,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGUSR1,
		)

		drainTime := seppukuDelaySeconds * time.Second

		sig := <-ch
		fmt.Printf("[INFO] Got a %s, quiting in %s.\n", sig, drainTime)
		healthStatus = http.StatusServiceUnavailable
		<-time.After(drainTime)
		os.Exit(0)
	}()

	return RecoverOnPanic(func(response http.ResponseWriter, request *http.Request) {

		msg, err := json.Marshal(Response{
			Status:   http.StatusOK,
			Version:  version,
			Hostname: hostname,
		})
		if err != nil {
			panic(err)
		}

		response.Header().Add(HeaderContentType, ContentTypeJSON)
		response.WriteHeader(healthStatus)
		response.Write(msg)

		if healthStatus == http.StatusOK {
			log.Println("PING - In Service")
		} else {
			log.Println("PING - Going Out Of Service")
		}
	})
}

func RecoverOnPanic(handler HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				st := RecoverStackTrace(r)
				log.Println("PANIC", st)
				http.Error(response, "Error", http.StatusInternalServerError)
			}
		}()
		handler(response, request)
	})
}

func RecoverStackTrace(panicReason interface{}) string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Recovering from panic: - %v\r\n", panicReason))
	for i := 2; ; i += 1 {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		buffer.WriteString(fmt.Sprintf("    %s:%d\r\n", file, line))
	}
	return buffer.String()
}

func RecoverError(panicReason interface{}) error {
	switch t := panicReason.(type) {
	case string:
		return errors.New(t)
	case error:
		return t
	default:
		return errors.New("Unknown error")
	}
}
