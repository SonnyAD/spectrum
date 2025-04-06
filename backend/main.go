package main

import (
	"encoding/xml"
	"net/http"
	"os"

	"SonnyAD/spectrum/api"
	_ "SonnyAD/spectrum/docs"
	"SonnyAD/spectrum/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func EmptyResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// @Summary		Healthcheck
// @Description	Get the status of the API
// @Tags			health
// @Produce		json,xml,application/yaml,plain
// @Success		200	{object}	Health
// @Router			/status [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	var health Health
	health.Status = "up"

	version, present := os.LookupEnv("API_VERSION")
	if present {
		health.Version = version
	}

	utils.Output(w, r.Header["Accept"], health, health.Status)
}

type Health struct {
	XMLName xml.Name `json:"-" xml:"health" yaml:"-"`
	Version string   `json:"version,omitempty" xml:"version,omitempty" yaml:"version,omitempty"`
	Status  string   `json:"status" xml:"status" yaml:"status"`
}

func initLogging() {
	log.SetLevel(log.DebugLevel)
	//log.SetFormatter(&log.JSONFormatter{})
}

// @title			SonnyAD Spectrum API
// @version		1.0
// @description	The backend powering SonnyAD online spectrum platform.
//
// @contact.name	API Support
// @contact.email	api@utile.space
//
// @license.name	utile.space API License
// @license.url	https://utile.space/api/
//
// @BasePath		/api
func main() {
	initLogging()

	router := mux.NewRouter()

	router.Use(utils.EnableCors)

	apiRouter := router.PathPrefix("/api").Subrouter()

	router.HandleFunc("/", EmptyResponse).Methods(http.MethodGet)
	apiRouter.HandleFunc("/", EmptyResponse).Methods(http.MethodGet)

	apiRouter.HandleFunc("/spectrum/ws", api.SpectrumWebsocket).Methods(http.MethodGet)

	apiRouter.HandleFunc("/status", HealthCheck).Methods(http.MethodGet)

	apiRouter.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	port, present := os.LookupEnv("PORT")
	if !present {
		port = "3000"
	}

	log.Info("Starting server on port ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
