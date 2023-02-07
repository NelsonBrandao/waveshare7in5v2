package waveshare7in5v2

import (
	"image/color"
	"time"
)

func isWhite(c color.Color) bool {
	r, g, b, _ := c.RGBA()

	return r == 0 && g == 0 && b == 0
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
