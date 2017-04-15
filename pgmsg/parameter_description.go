package pgmsg

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

type ParameterDescription struct {
	ParameterOIDs []uint32
}

func (*ParameterDescription) Backend() {}

func (dst *ParameterDescription) UnmarshalBinary(src []byte) error {
	buf := bytes.NewBuffer(src)

	if buf.Len() < 2 {
		return &invalidMessageFormatErr{messageType: "ParameterDescription"}
	}
	parameterCount := int(binary.BigEndian.Uint16(buf.Next(2)))
	if buf.Len() != parameterCount*4 {
		return &invalidMessageFormatErr{messageType: "ParameterDescription"}
	}

	*dst = ParameterDescription{ParameterOIDs: make([]uint32, parameterCount)}

	for i := 0; i < parameterCount; i++ {
		dst.ParameterOIDs[i] = binary.BigEndian.Uint32(buf.Next(4))
	}

	return nil
}

func (src *ParameterDescription) MarshalBinary() ([]byte, error) {
	var bigEndian BigEndianBuf
	buf := &bytes.Buffer{}

	buf.WriteByte('t')
	buf.Write(bigEndian.Uint32(uint32(4 + 2 + 4*len(src.ParameterOIDs))))

	buf.Write(bigEndian.Uint16(uint16(len(src.ParameterOIDs))))

	for _, oid := range src.ParameterOIDs {
		buf.Write(bigEndian.Uint32(oid))
	}

	return buf.Bytes(), nil
}

func (src *ParameterDescription) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type          string
		ParameterOIDs []uint32
	}{
		Type:          "ParameterDescription",
		ParameterOIDs: src.ParameterOIDs,
	})
}
