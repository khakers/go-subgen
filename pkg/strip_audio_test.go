package pkg

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"go-subgen/internal/projectpath"
)

func TestStripAudioToFile(t *testing.T) {
	type args struct {
		videoFile string
		output    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "a", args: args{
			videoFile: filepath.Join(projectpath.TestSamples, "bbb_sunflower_1080p_60fps_normal.mp4"),
			output:    filepath.Join(projectpath.TestOut, "bbb_sample.wav"),
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StripAudioToFile(tt.args.videoFile, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("StripAudio() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStripAudio(t *testing.T) {
	open, err := os.Create(filepath.Join(projectpath.TestOut, "bbb_sample_buffered.wav"))
	if err != nil {
		t.Fatal(err)
	}
	err = StripAudio(filepath.Join(projectpath.TestSamples, "bbb_sunflower_1080p_60fps_normal.mp4"), open, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	err = open.Close()
	if err != nil {
		t.Fatal(err)
	}
}

// func TestStripAudio(t *testing.T) {
// 	type args struct {
// 		videoFile string
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		wantWriter string
// 		wantErrOut string
// 		wantErr    bool
// 	}{
// 		{
// 			name: "test bbb",
// 			args: args{
// 				videoFile: filepath.Join(projectpath.TestSamples, "bbb_sunflower_1080p_60fps_normal.mp4"),
//
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			writer := &bytes.Buffer{}
// 			errOut := &bytes.Buffer{}
// 			err := StripAudio(tt.args.videoFile, writer, errOut)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("StripAudio() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
// 				t.Errorf("StripAudio() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
// 			}
// 			if gotErrOut := errOut.String(); gotErrOut != tt.wantErrOut {
// 				t.Errorf("StripAudio() gotErrOut = %v, want %v", gotErrOut, tt.wantErrOut)
// 			}
// 		})
// 	}
// }

// func Test_getLen(t *testing.T) {
// 	type args struct {
// 		videoFile string
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		want  error
// 		want1 float32
// 	}{
// 		{
// 			name: "bbb_sunflower_mp4",
// 			args: args{videoFile: filepath.Join(projectpath.TestSamples, "bbb_sunflower_1080p_60fps_normal.mp4")},
// 			want: nil,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1 := getLen(tt.args.videoFile)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("getLen() got = %v, want %v", got, tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("getLen() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }
