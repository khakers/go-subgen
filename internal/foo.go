package internal

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-audio/wav"
	log "github.com/sirupsen/logrus"
	"go-subgen/internal/projectpath"
)

// cursed test method
func foo() {

	file, err := os.ReadFile(filepath.Join(projectpath.TestSamples, "micro_machines_sample.wav"))
	if err != nil {
		log.Fatal(err)
		return
	}

	pcmData, err := os.ReadFile(filepath.Join(projectpath.TestSamples, "micro_machines_sample.pcm"))
	if err != nil {
		log.Fatal(err)
		return
	}

	transformedPCM := toInt16(pcmData)
	pcmFloat32 := floatify(transformedPCM)

	decoder := wav.NewDecoder(bytes.NewReader(file))

	if !decoder.IsValidFile() {
		err := decoder.Err()
		log.Fatal(err)
	} else {
		log.Debugf("wav file was valid")
	}
	err = decoder.Rewind()
	if err != nil {
		log.Fatal(err)
	}

	err = decoder.Err()
	if err != nil {
		log.Fatal(err)
	}

	pcmBuffer, err := decoder.FullPCMBuffer()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("pcmBuffer frames: %v", pcmBuffer.NumFrames())
	float32s := pcmBuffer.AsFloat32Buffer().Data

	log.Printf("uint8          len %v", len(pcmBuffer.Data))
	log.Printf("pcmdata        len %v", len(pcmData))
	log.Printf("transformedPCM len %v", len(transformedPCM))
	log.Printf("float32        len %v", len(float32s))
	log.Printf("pcmFloat32     len %v", len(pcmFloat32))

	// for i := 0; i < len(float32s); i++ {
	// 	log.Printf("val: %v", float32s[i])
	// }
	log.Printf("uint8          val: %v", pcmBuffer.Data[0])
	log.Printf("pcmdata        val: %v", pcmData[0])
	log.Printf("transformedPCM val: %v", transformedPCM[0])
	log.Printf("float32        val: %v", float32s[0])
	log.Printf("pcmfloat32     val: %v", pcmFloat32[0])

	log.Printf("last val: %v %v", pcmBuffer.Data[len(pcmBuffer.Data)-2], pcmBuffer.Data[len(pcmBuffer.Data)-1])
	log.Printf("last val: %v %v", pcmData[len(pcmData)-2], pcmData[len(pcmData)-1])
	log.Printf("last val: %v %v", transformedPCM[len(transformedPCM)-2], transformedPCM[len(transformedPCM)-1])
	log.Printf("last val: %v %v", float32s[len(float32s)-2], float32s[len(float32s)-1])
	log.Printf("last val: %v %v", pcmFloat32[len(pcmFloat32)-2], pcmFloat32[len(pcmFloat32)-1])

	if reflect.DeepEqual(pcmBuffer.Data, transformedPCM) {
		log.Printf("buffers are equal")
	} else {
		log.Errorln("buffers aren't equal")
	}

	if reflect.DeepEqual(float32s, pcmFloat32) {
		log.Printf("floats are equal")
	} else {
		log.Errorln("floats aren't equal")
	}

	// log.Printf("last hex: %v", hex.EncodeToString([]byte(pcmBuffer.Data[(len(pcmBuffer.Data)-2):])))

}

func ConvertPCMBytes(data []byte) []float32 {
	return floatify(toInt16(data))

}

func floatify(buf []int16) []float32 {
	// newB := &Float32Buffer{}

	newData := make([]float32, len(buf))
	// buf.SourceBitDepth = 16
	// newB.SourceBitDepth = buf.SourceBitDepth
	factor := math.Pow(2, float64(16)-1)
	for i := 0; i < len(buf); i++ {
		newData[i] = float32(float64(buf[i]) / factor)
	}
	return newData
}

func toInt16(raw []byte) []int16 {
	const SIZEOF_INT16 = 2 // bytes

	data := make([]int16, len(raw)/SIZEOF_INT16)
	for i := range data {
		// assuming little endian
		data[i] = int16(
			binary.LittleEndian.Uint16(raw[i*SIZEOF_INT16 : (i+1)*SIZEOF_INT16]))
	}
	return data
}
