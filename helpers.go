package waveshare7in5v2

import (
	"image/color"
	"time"
)

func isBlack(c color.Color, threshold uint8) bool {
	return color.GrayModel.Convert(c).(color.Gray).Y < threshold
}

func splitInChunks(data []byte) [][]byte {
	var chunks [][]byte

	for i := 0; i < len(data); i += MAX_CHUNK_SIZE {
		end := i + MAX_CHUNK_SIZE

		if end > len(data) {
			end = len(data)
		}

		chunks = append(chunks, data[i:end])
	}

	return chunks
}

func wait(d time.Duration) {
	time.Sleep(d * time.Millisecond)
}
