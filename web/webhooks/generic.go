package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"go-subgen/pkg"
)

type GenericWebhookData struct {
	Files []string `json:"files"`
}

func ServeGeneric(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Received generic webhook")

	var data GenericWebhookData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.WithError(err).Errorln("Failed to decode webhook JSON data")

		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	log.WithField("data", fmt.Sprintf("%+v", data)).Debugln("Decoded webhook json data")

	for _, file := range data.Files {
		log.Debugf("Queued %v from radarr", file)
		pkg.EnqueueSub(file)
	}
	w.WriteHeader(200)
}
