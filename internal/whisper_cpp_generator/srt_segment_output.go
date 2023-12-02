package whisper_cpp_generator

import (
	"fmt"
	"io"
	"time"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

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
