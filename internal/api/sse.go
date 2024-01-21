package api

import (
	"encoding/json"

	"github.com/r3labs/sse/v2"
	log "github.com/sirupsen/logrus"
	"go-subgen/internal/interceptors"
)

const channelID = "progress"

// todo may make more sense to convert this into an asr repository middleware instead of the weird event channel thing

type SseProgressStream struct {
	ProgressStream *sse.Stream
	sseServer      *sse.Server
	interceptors.EventChannels
}

// // NewSseProgressStream Creates a new SseProgressStream
// func NewSseProgressStream() *SseProgressStream {
// 	return NewSseAsrEventStreamFromSseServer(sse.New())
// }

// NewSseAsrEventStreamFromSseServer Creates a new SseProgressStream from an existing sse server
func NewSseAsrEventStreamFromSseServer(server *sse.Server, channels interceptors.EventChannels) *SseProgressStream {

	stream := server.CreateStream("progress")
	stream.AutoReplay = false

	stream.OnSubscribe = func(streamID string, sub *sse.Subscriber) {
		log.WithFields(log.Fields{
			"streamID": streamID,
			"subUrl":   sub.URL,
		}).Debugf("Client subscribed to SSE")
	}

	stream.OnUnsubscribe = func(streamID string, sub *sse.Subscriber) {
		log.WithFields(log.Fields{
			"streamID": streamID,
			"subUrl":   sub.URL,
		}).Debugf("Client unsubscribed from SSE")
	}

	return &SseProgressStream{
		ProgressStream: stream,
		sseServer:      server,
		EventChannels:  channels,
	}
}

// Start Starts the progress channel worker
func (s *SseProgressStream) Start() {
	go progressChannelWorker(s.EventChannels.ProgressChannel, s.sseServer, channelID)
	go jobChannelWorker(s.EventChannels.JobChangeEventChannel, s.sseServer, channelID)
	go statusChannelWorker(s.EventChannels.StatusChannel, s.sseServer, channelID)
}

// Consumes events from the progress channel and publishes them to the sse server as json
func progressChannelWorker(channel chan interceptors.AsrProgressEvent, sseServer *sse.Server, SseChannelId string) {
	for i := range channel {
		log.Debugf("Dispatching progress event to SSE: %+v", i)

		data, err := json.Marshal(i)
		if err != nil {
			log.WithError(err).
				WithFields(log.Fields{"id": i.JobId, "progress": i.Progress}).
				Errorln("failed to marshal progress event")
		}

		sseServer.Publish(SseChannelId, &sse.Event{
			Event: []byte("progressUpdate"),
			Data:  data,
		})
		if err != nil {
			log.WithError(err).
				WithFields(log.Fields{"id": i.JobId, "progress": i.Progress}).
				Errorf("failed to publish event")
		}
	}
}

// Consumes events from the progress channel and publishes them to the sse server as json
func jobChannelWorker(channel chan interceptors.AsrJobEvent, sseServer *sse.Server, SseChannelId string) {
	for i := range channel {
		log.Debugf("Dispatching job event to SSE: %+v", i)

		data, err := json.Marshal(i)
		if err != nil {
			log.WithError(err).
				WithFields(log.Fields{"id": i.JobId, "event": i}).
				Errorln("failed to marshal job event")
		}

		sseServer.Publish(SseChannelId, &sse.Event{
			Event: []byte("jobUpdate"),
			Data:  data,
		})
		if err != nil {
			log.WithError(err).
				WithFields(log.Fields{"id": i.JobId, "event": i}).
				Errorf("failed to publish event")
		}
	}
}

// Consumes events from the progress channel and publishes them to the sse server as json
func statusChannelWorker(channel chan interceptors.AsrStatusEvent, sseServer *sse.Server, SseChannelId string) {
	for i := range channel {
		log.Debugf("Dispatching status event to SSE: %+v", i)
		data, err := json.Marshal(i)
		if err != nil {
			log.WithError(err).
				WithFields(log.Fields{"id": i.JobId, "event": i}).
				Errorln("failed to marshal status event")
		}

		sseServer.Publish(SseChannelId, &sse.Event{
			Event: []byte("statusChange"),
			Data:  data,
		})
		if err != nil {
			log.WithError(err).
				WithFields(log.Fields{"id": i.JobId, "event": i}).
				Errorf("failed to publish event")
		}
	}
}
