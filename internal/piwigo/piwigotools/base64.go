package piwigotools

import (
	"bytes"
	"encoding/base64"
	"os"
)

var CHUNK_SIZE int64 = 1 * 1024 * 1024
var CHUNK_BUFF_SIZE int64 = 32 * 1024
var CHUNK_BUFF_COUNT = CHUNK_SIZE / CHUNK_BUFF_SIZE

type Base64Chunk struct {
	Position int64
	Size     int64
	Buffer   bytes.Buffer
}

func Base64Chunker(filename string) (out chan *Base64Chunk, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}

	out = make(chan *Base64Chunk, 8)
	go func() {
		b := make([]byte, CHUNK_BUFF_SIZE)
		defer f.Close()
		defer close(out)
		ok := false
		for position := int64(0); !ok; position += 1 {
			bf := &Base64Chunk{
				Position: position,
			}
			b64 := base64.NewEncoder(base64.StdEncoding, &bf.Buffer)
			for i := int64(0); i < CHUNK_BUFF_COUNT; i++ {
				n, _ := f.Read(b)
				if n == 0 {
					ok = true
					break
				}
				bf.Size += int64(n)
				b64.Write(b[:n])
			}
			b64.Close()
			if bf.Size > 0 {
				out <- bf
			}
		}
	}()

	return
}
