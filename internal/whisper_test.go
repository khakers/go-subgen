package internal

// func TestGenerate(t *testing.T) {
// 	type args struct {
// 		modelPath string
// 		input     string
// 	}
// 	tests := []struct {
// 		name           string
// 		args           args
// 		wantSubsWriter string
// 		wantErr        bool
// 	}{
// 		{
// 			name: "micro machines sample",
// 			args: struct {
// 				modelPath string
// 				input     string
// 			}{modelPath: filePath.Join(projectpath.Test, "ggml-base.bin"), input: filePath.Join(projectpath.TestSamples, "micro_machines_sample.wav")},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			subsWriter := &bytes.Buffer{}
// 			err := Generate(tt.args.modelPath, tt.args.input, subsWriter)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotSubsWriter := subsWriter.String(); gotSubsWriter != tt.wantSubsWriter {
// 				t.Errorf("Generate() gotSubsWriter = %v, want %v", gotSubsWriter, tt.wantSubsWriter)
// 			}
// 		})
// 	}
// }

// func TestGenerate(t *testing.T) {
// 	microMachines, err := os.ReadFile(filePath.Join(projectpath.TestSamples, "micro_machines_sample.wav"))
// 	if err != nil {
// 		return
// 	}
// 	bush, err := os.ReadFile(filePath.Join(projectpath.TestSamples, "WBUSH_RADIO_SAMPLE.wav"))
// 	if err != nil {
// 		return
// 	}
//
// 	type args struct {
// 		modelPath string
// 		input     io.ReadSeeker
// 	}
// 	tests := []struct {
// 		name           string
// 		args           args
// 		wantSubsWriter string
// 		wantErr        bool
// 	}{
// 		{
// 			name: "micro machines sample",
// 			args: args{
// 				modelPath: filePath.Join(projectpath.TestModels, "ggml-base.bin"),
// 				input:     bytes.NewReader(microMachines),
// 			},
// 			wantSubsWriter: "",
// 			wantErr:        false,
// 		},
// 		{
// 			name: "bush radio sample",
// 			args: args{
// 				modelPath: filePath.Join(projectpath.TestModels, "ggml-base.bin"),
// 				input:     bytes.NewReader(bush),
// 			},
// 			wantSubsWriter: "",
// 			wantErr:        false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			subsWriter := &bytes.Buffer{}
// 			err := Generate(tt.args.modelPath, tt.args.input, subsWriter)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotSubsWriter := subsWriter.String(); gotSubsWriter != tt.wantSubsWriter {
// 				t.Errorf("Generate() gotSubsWriter = %v, want %v", gotSubsWriter, tt.wantSubsWriter)
// 			}
// 		})
// 	}
// }
