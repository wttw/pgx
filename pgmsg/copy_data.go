package pgmsg

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
)

type CopyData struct {
	Data []byte
}

func (*CopyData) Backend()  {}
func (*CopyData) Frontend() {}

func (dst *CopyData) UnmarshalBinary(src []byte) error {
	buf := bytes.NewBuffer(src)

	if buf.Len() < 4 {
		return &invalidMessageFormatErr{messageType: "CopyData"}
	}

	dataSize := int(int32(binary.BigEndian.Uint32(buf.Next(4))))

	if buf.Len() != dataSize {
		return &invalidMessageFormatErr{messageType: "CopyData"}
	}

	dst.Data = make([]byte, buf.Len())
	copy(dst.Data, buf.Bytes())

	return nil
}

func (src *CopyData) MarshalBinary() ([]byte, error) {
	var bigEndian BigEndianBuf
	buf := &bytes.Buffer{}

	buf.WriteByte('d')
	buf.Write(bigEndian.Uint32(uint32(4 + len(src.Data))))
	buf.Write(src.Data)

	return buf.Bytes(), nil
}

func (src *CopyData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
		Data string
	}{
		Type: "CopyData",
		Data: hex.EncodeToString(src.Data),
	})
}
