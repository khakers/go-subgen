package pkg

import (
	"fmt"
	"io"
	"time"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	log "github.com/sirupsen/logrus"
	"go-subgen/internal"
	"go-subgen/pkg/configuration"
)

func Generate(modelPath string, input []byte, subsWriter io.Writer) error {

	// Load the model
	model, err := whisper.New(modelPath)
	if err != nil {
		return err
	}
	defer model.Close()

	err = Process(model, input, subsWriter)
	if err != nil {
		return err
	}
	return nil
}

// func GenerateWithModel(model Model, reader io.Reader, subsWriter *io.Writer) error {
// 	return nil
// }

// Adapted from code at https://github.com/ggerganov/whisper.cpp/blob/2bee2650c66497b8804e3c82426373703c6d97a1/bindings/go/examples/go-whisper/process.go

func Process(model whisper.Model, input []byte, subsWriter io.Writer) error {
	context, err := model.NewContext()
	if err != nil {
		return err
	}

	log.Debugf(context.SystemInfo())

	if configuration.Cfg.WhisperConf.Threads != 0 {
		context.SetThreads(configuration.Cfg.WhisperConf.Threads)
	}
	if configuration.Cfg.WhisperConf.MaxSegmentLength != 0 {
		context.SetMaxSegmentLength(configuration.Cfg.WhisperConf.MaxSegmentLength)
	}
	if configuration.Cfg.WhisperConf.MaxTokensPerSegment != 0 {
		context.SetMaxTokensPerSegment(configuration.Cfg.WhisperConf.MaxSegmentLength)
	}
	if configuration.Cfg.WhisperConf.TokenSumThreshold != 0 {
		context.SetTokenSumThreshold(configuration.Cfg.WhisperConf.TokenSumThreshold)
	}
	if configuration.Cfg.WhisperConf.TokenThreshold != 0 {
		context.SetTokenSumThreshold(configuration.Cfg.WhisperConf.TokenThreshold)
	}

	context.SetLanguage(configuration.Cfg.TargetLang)
	context.SetSpeedup(configuration.Cfg.WhisperConf.WhisperSpeedup)

	data := internal.ConvertPCMBytes(input)

	var segmentCallback whisper.SegmentCallback

	context.ResetTimings()

	err = context.Process(data, segmentCallback)
	if err != nil {
		return err
	}

	context.PrintTimings()
	return OutputSRT(subsWriter, context)
}

func OutputSRT(writer io.Writer, context whisper.Context) (err error) {
	n := 1
	for {
		segment, err := context.NextSegment()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		_, err = fmt.Fprintln(writer, n)
		_, err = fmt.Fprintln(writer, srtTimestamp(segment.Start), " --> ", srtTimestamp(segment.End))
		_, err = fmt.Fprintln(writer, segment.Text)
		_, err = fmt.Fprintln(writer, "")
		n++
	}
	return err
}

func srtTimestamp(t time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d,%03d", t/time.Hour, (t%time.Hour)/time.Minute, (t%time.Minute)/time.Second, (t%time.Second)/time.Millisecond)
}
