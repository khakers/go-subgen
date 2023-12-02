package whisper_cpp_generator

import (
	"context"
	"io"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	log "github.com/sirupsen/logrus"
	"go-subgen/internal/audio"
	"go-subgen/internal/configuration"
)

func Generate(modelPath string, input []byte, subsWriter io.Writer, progressChannel chan float32, ctx context.Context) error {

	// Load the model
	model, err := whisper.New(modelPath)
	if err != nil {
		return err
	}
	defer model.Close()

	err = Process(model, input, subsWriter, configuration.Cfg.WhisperConf, progressChannel, ctx)
	if err != nil {
		return err
	}
	return nil
}

// Adapted from code at https://github.com/ggerganov/whisper.cpp/blob/2bee2650c66497b8804e3c82426373703c6d97a1/bindings/go/examples/go-whisper/process.go

func Process(model whisper.Model, input []byte, subsWriter io.Writer, whisperConfig configuration.WhisperConfig, progressChannel chan float32, ctx context.Context) error {
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

	data := audio.ConvertPCMBytes(input)

	var segmentCallback whisper.SegmentCallback

	// var progressCallback whisper.ProgressCallback

	context.ResetTimings()

	err = context.Process(data, segmentCallback, func(i int) {
		log.Tracef("progress callback %v", i)
		progressChannel <- float32(i)
	})

	if err != nil {
		return err
	}

	context.PrintTimings()
	log.Debugf("finished processing, running SRT output")
	return outputJSON(subsWriter, context)
}
