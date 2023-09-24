package audio

import (
	"encoding/binary"
	"math"
)

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
