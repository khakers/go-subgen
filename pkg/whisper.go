package pkg

import (
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"go-subgen/internal"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
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
	context.SetThreads(Cfg.WhisperConf.Threads)
	context.SetLanguage(Cfg.TargetLang)
	context.SetSpeedup(Cfg.WhisperConf.WhisperSpeedup)

	// log.Infof("decoding wav")
	// size, _ := input.Seek(0, io.SeekEnd)
	// log.Debugf("size %v", size)
	// _, _ = input.Seek(0, io.SeekStart)

	// decoder := wav.NewDecoder(input)
	//
	// if !decoder.IsValidFile() {
	// 	err = decoder.Err()
	// 	log.Debugln(err)
	// 	return errors.New("wav file was invalid")
	// } else {
	// 	log.Debugf("wav file was valid")
	// }
	// err = decoder.Rewind()
	// if err != nil {
	// 	return err
	// }
	//
	// // TODO see if a buffer would be usable here to avoid loading all audio into memory
	// err = decoder.Err()
	// if err != nil {
	// 	return err
	// }
	//
	// pcmBuffer, err := decoder.FullPCMBuffer()
	// if err != nil {
	// 	return err
	// }
	//
	// log.Debugf("pcmBuffer frames: %v", pcmBuffer.NumFrames())
	// // todo we can probably decode directly to this from ffmpeg
	// data := pcmBuffer.AsFloat32Buffer().Data
	//
	// if len(data) == 0 {
	// 	return errors.New("empty float32 buffer")
	// }
	//
	// log.Debugf("float32len %v", len(data))
	//
	// log.Infof("finished decoding WAV")

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
