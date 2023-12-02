package whisper_cpp_generator

import (
	"encoding/json"
	"io"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

func outputJSON(writer io.Writer, context whisper.Context) (err error) {
	n := 1
	var segments []whisper.Segment
	for {
		segment, err := context.NextSegment()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		segments = append(segments, segment)
		n++
	}

	marshal, err := json.Marshal(segments)
	if err != nil {
		return err
	}
	_, err = writer.Write(marshal)
	if err != nil {
		return err
	}
	return nil
}
