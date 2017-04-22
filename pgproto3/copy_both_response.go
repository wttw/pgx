package pgproto3

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

type CopyBothResponse struct {
	OverallFormat     byte
	ColumnFormatCodes []uint16
	ParameterOIDs     []uint32
}

func (*CopyBothResponse) Backend() {}

func (dst *CopyBothResponse) UnmarshalBinary(src []byte) error {
	buf := bytes.NewBuffer(src)

	if buf.Len() < 2 {
		return &invalidMessageFormatErr{messageType: "CopyBothResponse"}
	}
	parameterCount := int(binary.BigEndian.Uint16(buf.Next(2)))
	if buf.Len() != parameterCount*4 {
		return &invalidMessageFormatErr{messageType: "CopyBothResponse"}
	}

	*dst = CopyBothResponse{ParameterOIDs: make([]uint32, parameterCount)}

	for i := 0; i < parameterCount; i++ {
		dst.ParameterOIDs[i] = binary.BigEndian.Uint32(buf.Next(4))
	}

	return nil
}

func (src *CopyBothResponse) MarshalBinary() ([]byte, error) {
	var bigEndian BigEndianBuf
	buf := &bytes.Buffer{}

	buf.WriteByte('W')
	buf.Write(bigEndian.Uint32(uint32(4 + 2 + 4*len(src.ParameterOIDs))))

	buf.Write(bigEndian.Uint16(uint16(len(src.ParameterOIDs))))

	for _, oid := range src.ParameterOIDs {
		buf.Write(bigEndian.Uint32(oid))
	}

	return buf.Bytes(), nil
}

func (src *CopyBothResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type          string
		ParameterOIDs []uint32
	}{
		Type:          "CopyBothResponse",
		ParameterOIDs: src.ParameterOIDs,
	})
}
