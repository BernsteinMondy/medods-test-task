package encoder

import "encoding/base64"

type Encoder struct{}

func (e *Encoder) Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}
