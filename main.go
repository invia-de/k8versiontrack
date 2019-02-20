package main

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"

	//"k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/api/resource"
	//"k8s.io/apimachinery/pkg/watch"
	"github.com/invia-de/K8VersionTrack/controller"
	"github.com/invia-de/K8VersionTrack/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newConfig() (*viper.Viper, error) {
	c := viper.New()
	//c.SetDefault("cookie_secret", "Az00P54fhK2SMggW")
	c.SetDefault("http_addr", ":8888")
	c.SetDefault("http_cert_file", "")
	c.SetDefault("http_key_file", "")
	c.SetDefault("http_drain_interval", "1s")
	c.SetEnvPrefix("kvt")
	c.AutomaticEnv()

	return c, nil
}

func main() {
	vcl := model.NewVersionCollector()
	prometheus.MustRegister(vcl)

	cfg, err := newConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	serverAddress := cfg.Get("http_addr").(string)

	certFile := cfg.Get("http_cert_file").(string)
	keyFile := cfg.Get("http_key_file").(string)

	drainIntervalString := cfg.Get("http_drain_interval").(string)
	drainInterval, err := time.ParseDuration(drainIntervalString)
	if err != nil {
		logrus.Fatal(err)
	}
	router := mux.NewRouter()
	router.Handle("/", http.HandlerFunc(controller.IndexAction)).Methods("GET")
	router.Handle("/metrics", promhttp.Handler())
	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: serverAddress, Handler: router},
	}

	logrus.Infoln("Running HTTP server on " + serverAddress)

	if certFile != "" && keyFile != "" {
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		logrus.Fatal(err)
	}
}
