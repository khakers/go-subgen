package main

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"go-subgen/internal"
	"go-subgen/internal/api"
	"go-subgen/internal/api/webhooks"
	"go-subgen/internal/configuration"
	"go-subgen/pkg"
)

func main() {
	conf, err := configuration.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(conf.LogLevel)
	log.Printf("using model type %v for language %s", conf.ModelType, conf.WhisperConf.TargetLang)

	log.Debugf("%+v", configuration.Cfg)
	downloaded, err := pkg.IsModelPresent(conf.ModelType, true)
	if err != nil {
		log.WithError(err).Errorln("Model check failed (this is likely normal and can be ignored)")
	}
	if !downloaded {
		err := pkg.DownloadModel(conf.ModelType, true)
		if err != nil {
			log.WithError(err).Fatalln("failed to download model")
		}
	}

	// Handle http healthcheck endpoint
	http.Handle("/healthcheck", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Errorln(err)
		}
	}))

	http.Handle("/webhooks/generic", http.HandlerFunc(api.ServeGeneric))
	http.Handle("/webhooks/tautulli", http.HandlerFunc(webhooks.ServeTautulli))
	http.Handle("/webhooks/radarr", http.HandlerFunc(webhooks.ServeRadarr))
	http.Handle("/webhooks/sonarr", http.HandlerFunc(webhooks.ServeSonarr))

	internal.StartWorkers(conf)
	err = http.ListenAndServe(":"+strconv.Itoa(int(conf.ServerConfig.Port)), nil)
	if err != nil {
		log.WithError(err).Fatal("web server failure")
	}

}
