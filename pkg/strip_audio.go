package pkg

import (
	"bytes"
	"io"

	log "github.com/sirupsen/logrus"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func StripAudioToFile(videoFile string, output string) error {
	err := ffmpeg.
		Input(videoFile).
		Output(output,
			ffmpeg.KwArgs{
				"ar":  "16000",
				"ac":  "1",
				"c:a": "pcm_s16le",
			}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()
	if err != nil {
		return err
	}
	return err
}

func StripAudioToBuffer(videoFile string) (error, *[]byte) {
	buffer := bytes.NewBuffer(nil)

	err := StripAudio(videoFile, buffer, nil)
	if err != nil {
		return err, nil
	}

	if err != nil {
		return err, nil
	}
	log.Debugf("Read audio into buffer of size %v", len(buffer.Bytes()))
	all, err := io.ReadAll(buffer)
	if err != nil {
		return err, nil
	}
	return err, &all

}

func StripAudio(videoFile string, writer io.Writer, errOut io.Writer) error {
	err := ffmpeg.
		Input(videoFile).
		Output("pipe:",
			ffmpeg.KwArgs{
				"ar":  "16000",
				"ac":  "1",
				"c:a": "pcm_s16le",
				"f":   "wav",
			}).
		WithErrorOutput(errOut).
		WithOutput(writer).
		Run()
	return err
}

// StripAudioRaw strips the audio from a video file and writes it as raw PCM signed 16-bit little-endian audio to supplied writer
func StripAudioRaw(videoFile string, writer io.Writer, errOut io.Writer) error {
	err := ffmpeg.
		Input(videoFile).
		Output("pipe:",
			ffmpeg.KwArgs{
				"ar":       "16000",
				"ac":       "1",
				"c:a":      "pcm_s16le",
				"f":        "s16le",
				"loglevel": "error",
			}).
		WithErrorOutput(errOut).
		WithOutput(writer).
		Run()
	return err

}

func StripAudioToRawBuffer(videoFile string) (error, *[]byte) {
	buffer := bytes.NewBuffer(nil)
	err := ffmpeg.
		Input(videoFile).
		Output("pipe:",
			ffmpeg.KwArgs{
				"ar":  "16000",
				"ac":  "1",
				"c:a": "pcm_s16le",
				"f":   "raw",
			}).
		WithOutput(buffer).
		Run()
	if err != nil {
		return err, nil
	}
	log.Debugf("Read audio into buffer of size %v", len(buffer.Bytes()))

	all, err := io.ReadAll(buffer)
	if err != nil {
		return err, nil
	}
	return err, &all

}

// Requires ffprobe
func getLen(videoFile string) (error, float32) {
	probe, err := ffmpeg.Probe(videoFile, ffmpeg.KwArgs{"-select_streams": "v:0", "-show_entries": "stream=duration", "-of": "default=noprint_wrappers=1:nokey=1"})
	if err != nil {
		return err, 0
	}

	log.Printf(probe)

	return nil, 0
}
