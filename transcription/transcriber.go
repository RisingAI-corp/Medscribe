package transcriber

import "context"

type Transcription interface {
	Transcribe(ctx context.Context, audio []byte) (string, error)
}
