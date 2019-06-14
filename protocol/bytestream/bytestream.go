package bytestream

import "reflect"

type ByteStream struct {
	buffer []byte
	index  int
}

// Write serializes the data and returns the binary
func Write(data ByteStreamer) []byte {
	// TODO automaticaly resize the buffer
	stream := &ByteStream{make([]byte, 1024), 0}
	data.Write(stream)
	return stream.buffer[:stream.index]
}

type ByteStreamer interface {
	Write(stream *ByteStream)
}

func (stream *ByteStream) WriteByte(b byte) error {
	stream.buffer[stream.index] = b
	stream.index++
	return nil
}

func (stream *ByteStream) WriteBytes(bytes []byte) {
	for _, b := range bytes {
		stream.WriteByte(b)
	}
}

func (stream *ByteStream) WriteInt64(i int64) {
	shift := int64(256)
	for index := 0; index < 8; index++ {
		stream.WriteByte(byte(i % shift))
		i /= shift
	}
}

func (stream *ByteStream) WriteInt32(i int32) {
	shift := int32(256)
	for index := 0; index < 4; index++ {
		stream.WriteByte(byte(i % shift))
		i /= shift
	}
}

func (stream *ByteStream) WriteInt16(i int16) {
	shift := int16(256)
	for index := 0; index < 2; index++ {
		stream.WriteByte(byte(i % shift))
		i /= shift
	}
}

func (stream *ByteStream) WriteString(s string) {
	for _, c := range []byte(s) {
		stream.WriteByte(c)
	}
}

func (stream *ByteStream) WriteNullable(data interface{}) {
	switch v := data.(type) {
	case *int64:
		if v == nil {
			stream.WriteByte(0)
			return
		}
		stream.WriteByte(1)
		stream.WriteInt64(*v)
	case *int32:
		if v == nil {
			stream.WriteByte(0)
			return
		}
		stream.WriteByte(1)
		stream.WriteInt32(*v)
	case *int16:
		if v == nil {
			stream.WriteByte(0)
			return
		}
		stream.WriteByte(1)
		stream.WriteInt16(*v)
	case *byte:
		if v == nil {
			stream.WriteByte(0)
			return
		}
		stream.WriteByte(1)
		stream.WriteByte(*v)
	case ByteStreamer:
		if reflect.ValueOf(v).IsNil() {
			stream.WriteByte(0)
			return
		}
		stream.WriteByte(1)
		v.Write(stream)
	default:
		panic(data)
	}
}

func (stream *ByteStream) WriteList(length int, callback func(i int) ByteStreamer) {
	stream.WriteByte(byte(length))
	for index := 0; index < length; index++ {
		callback(index).Write(stream)
	}
}

func (stream *ByteStream) Buffer() []byte {
	return stream.buffer
}
