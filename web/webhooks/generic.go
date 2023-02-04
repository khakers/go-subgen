package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"go-subgen/pkg"
)

type GenericWebhookData struct {
	Files []string `json:"file"`
}

func ServeGeneric(w http.ResponseWriter, r *http.Request) {
	var data GenericWebhookData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(400)
		fmt.Fprintf(w, "error: %v", err)
		return
	}
	for _, s := range data.Files {
		pkg.EnqueueSub(s)
	}
	w.WriteHeader(200)
}
