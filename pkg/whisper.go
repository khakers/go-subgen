package pkg

import (
	"fmt"
	"io"
	"time"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	log "github.com/sirupsen/logrus"
	"go-subgen/internal"
	"go-subgen/internal/configuration"
)

func Generate(modelPath string, input []byte, subsWriter io.Writer) error {

	// Load the model
	model, err := whisper.New(modelPath)
	if err != nil {
		return err
	}
	defer model.Close()

	err = Process(model, input, subsWriter, configuration.Cfg.WhisperConf)
	if err != nil {
		return err
	}
	return nil
}

// func GenerateWithModel(model Model, reader io.Reader, subsWriter *io.Writer) error {
// 	return nil
// }

// Adapted from code at https://github.com/ggerganov/whisper.cpp/blob/2bee2650c66497b8804e3c82426373703c6d97a1/bindings/go/examples/go-whisper/process.go

func Process(model whisper.Model, input []byte, subsWriter io.Writer, whisperConfig configuration.WhisperConfig) error {
	context, err := model.NewContext()
	if err != nil {
		return err
	}

	log.Debugf(context.SystemInfo())

	if whisperConfig.Threads != 0 {
		context.SetThreads(whisperConfig.Threads)
	}
	if whisperConfig.MaxSegmentLength != 0 {
		context.SetMaxSegmentLength(whisperConfig.MaxSegmentLength)
	}
	if whisperConfig.MaxTokensPerSegment != 0 {
		context.SetMaxTokensPerSegment(whisperConfig.MaxSegmentLength)
	}
	if whisperConfig.TokenSumThreshold != 0 {
		context.SetTokenSumThreshold(whisperConfig.TokenSumThreshold)
	}
	if whisperConfig.TokenThreshold != 0 {
		context.SetTokenSumThreshold(whisperConfig.TokenThreshold)
	}

	err = context.SetLanguage(whisperConfig.TargetLang)
	if err != nil {
		return err
	}
	context.SetSpeedup(whisperConfig.WhisperSpeedup)

	data := internal.ConvertPCMBytes(input)

	var segmentCallback whisper.SegmentCallback

	context.ResetTimings()

	err = context.Process(data, segmentCallback)
	if err != nil {
		return err
	}

	context.PrintTimings()
	return outputSRT(subsWriter, context)
}

func outputSRT(writer io.Writer, context whisper.Context) (err error) {
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
