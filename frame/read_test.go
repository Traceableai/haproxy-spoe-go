package frame

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

var testFrame = []byte(string(
	"\x00\x00\x00\x53" + // Size
		"\x03" + //TypeNotify
		"\x00\x00\x00\x01\xfe\x12\x01\x11\x67\x65\x74" +
		"\x2d\x69\x70\x2d\x72\x65\x70\x75\x74\x61\x74\x69\x6f\x6e\x04\x02" +
		"\x69\x70\x06\xc1\xc8\xe3\xde\x04" +
		"host" + //Host
		"\x08\x12" +
		"domain.example.com" + // authtest.ninjas.pl
		"\x0d\x61\x75\x74\x68\x6f\x72\x69\x7a\x61\x74\x69\x6f\x6e\x00\x06" +
		"\x63\x6f\x6f\x6b\x69\x65\x00",
))

func TestFrame_Read(t *testing.T) {
	r := bytes.NewBuffer(testFrame)
	f := NewFrame()
	err := f.Read(r)
	require.Nil(t, err)
	assert.Equal(t, 1, int(f.FrameID), "FrameID")
	assert.Equal(t, 542, int(f.StreamID), "StreamID")
	assert.Equal(t, TypeNotify, f.Type)
	messages := *f.Messages
	require.Len(t, messages, 1)
	host, found := messages[0].KV.Get("host")
	require.True(t, found, "key host")
	hostString, ok := host.(string)
	require.True(t, ok)
	assert.Equal(t, "domain.example.com", hostString)
}

func BenchmarkFrame_Read(b *testing.B) {
	readers := make([]io.Reader, b.N)
	// prepare readers beforehand, so we don't measure the performance of NewReader
	for idx := range readers {
		readers[idx] = bytes.NewBuffer(testFrame)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f := NewFrame()
		_ = f.Read(readers[i])
	}

}
