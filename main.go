package main

import (
	"net/http"
	"strconv"

	"github.com/r3labs/sse/v2"
	jobqueue "go-subgen/internal/adapters"
	"go-subgen/internal/interceptors"
	"go-subgen/internal/middlewares"
	"go-subgen/internal/whisper_cpp_generator"

	log "github.com/sirupsen/logrus"
	"go-subgen/internal/api"
	"go-subgen/internal/api/webhooks"
	"go-subgen/internal/configuration"
	"go-subgen/pkg/model"
)

func main() {
	conf, err := configuration.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(conf.LogLevel)
	log.Printf("using model type %v for language %s", conf.ModelType, conf.WhisperConf.TargetLang)

	log.Debugf("%+v", configuration.Cfg)
	downloaded, err := model.IsModelPresent(conf.ModelType, conf.ModelDir)
	if err != nil {
		log.WithError(err).Errorln("Model check failed")
	}
	if downloaded {
		hash, err := model.VerifyModelHash(conf.ModelType, conf.ModelDir)
		if err != nil {
			log.WithError(err).Errorln("Model hash verification failed")
		}
		downloaded = hash
	}
	if !downloaded {
		err := model.DownloadModel(conf.ModelType, conf.ModelDir, conf.VerifyModelHash)
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

	asrArrayQueue := jobqueue.NewArrayRepository()

	sseServer := sse.New()

	asrQueue := interceptors.NewAsrEventInterceptor(asrArrayQueue)

	progressServer := api.NewSseAsrEventStreamFromSseServer(sseServer, asrQueue.EventChannels)

	handler := api.NewGenericFileHandler(asrQueue)

	jobHandler := api.NewJobHandler(asrQueue)

	subtitleGenerator := whisper_cpp_generator.NewSubtitleGenerator(conf, asrQueue)

	mux := http.NewServeMux()

	mux.HandleFunc("/webhooks/generic", handler.Serve)
	mux.HandleFunc("/webhooks/tautulli", webhooks.ServeTautulli)
	mux.HandleFunc("/webhooks/radarr", webhooks.ServeRadarr)
	mux.HandleFunc("/webhooks/sonarr", webhooks.ServeSonarr)
	mux.HandleFunc("/api/v1/jobs", jobHandler.Serve)
	mux.HandleFunc("/api/events", sseServer.ServeHTTP)

	subtitleGenerator.StartWorkers()
	progressServer.Start()

	// Sourced from https://www.jvt.me/posts/2023/09/01/golang-nethttp-global-middleware/
	wrapped := use(mux, middlewares.LoggingMiddleware, middlewares.CorsMiddleware)

	err = http.ListenAndServe(":"+strconv.Itoa(int(conf.ServerConfig.Port)), wrapped)
	if err != nil {
		log.WithError(err).Fatal("web server failure")
	}

	log.Printf("listening on %v", conf.ServerConfig.Port)

}

// Sourced from https://www.jvt.me/posts/2023/09/01/golang-nethttp-global-middleware/
func use(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = r

	for _, mw := range middlewares {
		s = mw(s)
	}

	return s
}
