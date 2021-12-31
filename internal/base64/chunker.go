package base64

import (
	"bytes"
	b64 "encoding/base64"
	"os"
)

var CHUNK_SIZE int64 = 1 * 1024 * 1024
var CHUNK_BUFF_SIZE int64 = 32 * 1024
var CHUNK_BUFF_COUNT = CHUNK_SIZE / CHUNK_BUFF_SIZE

type Chunk struct {
	Position int64
	Size     int64
	Buffer   bytes.Buffer
}

func Chunker(filename string) (chan *Chunk, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	out := make(chan *Chunk, 8)
	chunker := func() {
		b := make([]byte, CHUNK_BUFF_SIZE)
		defer f.Close()
		defer close(out)
		ok := false
		for position := int64(0); !ok; position += 1 {
			bf := &Chunk{
				Position: position,
			}
			b64 := b64.NewEncoder(b64.StdEncoding, &bf.Buffer)
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
	}

	go chunker()

	return out, nil
}
