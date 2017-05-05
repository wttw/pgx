package pgproto3

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/pgio"
)

const (
	protocolVersionNumber = 196608 // 3.0
	sslRequestNumber      = 80877103
)

type StartupMessage struct {
	ProtocolVersion uint32
	Parameters      map[string]string
}

func (*StartupMessage) Frontend() {}

func (dst *StartupMessage) Decode(src []byte) error {
	src, dst.ProtocolVersion = pgio.NextUint32(src)
	if dst.ProtocolVersion == sslRequestNumber {
		return fmt.Errorf("can't handle ssl connection request")
	}

	if dst.ProtocolVersion != protocolVersionNumber {
		return fmt.Errorf("Bad startup message version number. Expected %d, got %d", protocolVersionNumber, msg.ProtocolVersion)
	}

	msg.Parameters = make(map[string]string)
	for {
		var ok bool
		var key, value string
		src, key, ok = pgio.NextCString(src)
		if !ok {
			return &invalidMessageFormatErr{messageType: "StartupMessage"}
		}
		src, value, ok = pgio.NextCString(src)
		if !ok {
			return &invalidMessageFormatErr{messageType: "StartupMessage"}
		}

		msg.Parameters[key] = value

		if len(src) == 1 {
			if src[0] != 0 {
				return fmt.Errorf("Bad startup message last byte. Expected 0, got %d", b)
			}
			break
		}
	}

	return nil
}

func (src *StartupMessage) MarshalBinary() ([]byte, error) {
	var bigEndian BigEndianBuf
	buf := &bytes.Buffer{}
	buf.Write(bigEndian.Uint32(0))
	buf.Write(bigEndian.Uint32(src.ProtocolVersion))
	for k, v := range src.Parameters {
		buf.WriteString(k)
		buf.WriteByte(0)
		buf.WriteString(v)
		buf.WriteByte(0)
	}
	buf.WriteByte(0)

	binary.BigEndian.PutUint32(buf.Bytes()[0:4], uint32(buf.Len()))

	return buf.Bytes(), nil
}

func (src *StartupMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type            string
		ProtocolVersion uint32
		Parameters      map[string]string
	}{
		Type:            "StartupMessage",
		ProtocolVersion: src.ProtocolVersion,
		Parameters:      src.Parameters,
	})
}
