package test

import (
	"GoRedis/interface/redis"
	"GoRedis/lib/utils"
	"GoRedis/redis/parser"
	"GoRedis/redis/protocol"
	"bytes"
	"io"
	"testing"
)

func TestParseStream(t *testing.T) {
	replies := []redis.Reply{
		protocol.NewIntReply(1),
		protocol.NewStatusReply("OK"),
		protocol.NewErrReply("ERR unknown"),
		protocol.NewBulkReply([]byte("a\r\nb")), // test binary safe
		protocol.NewNullBulkReply(),
		protocol.NewMultiBulkReply([][]byte{
			[]byte("a"),
			[]byte("\r\n"),
		}),
		protocol.NewEmptyMultiBulkReply(),
	}
	reqs := bytes.Buffer{}
	for _, re := range replies {
		reqs.Write(re.ToBytes())
	}
	reqs.Write([]byte("set a a" + protocol.CRLF)) // test text protocol
	expected := make([]redis.Reply, len(replies))
	copy(expected, replies)
	expected = append(expected, protocol.NewMultiBulkReply([][]byte{
		[]byte("set"), []byte("a"), []byte("a"),
	}))

	ch := parser.ParseStream(bytes.NewReader(reqs.Bytes()))
	i := 0
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF {
				return
			}
			t.Error(payload.Err)
			return
		}
		if payload.Data == nil {
			t.Error("empty data")
			return
		}
		exp := expected[i]
		i++
		if !utils.BytesEquals(exp.ToBytes(), payload.Data.ToBytes()) {
			t.Error("parse failed: " + string(exp.ToBytes()))
		}
	}
}

func TestParseOne(t *testing.T) {
	replies := []redis.Reply{
		protocol.NewIntReply(1),
		protocol.NewStatusReply("OK"),
		protocol.NewErrReply("ERR unknown"),
		protocol.NewBulkReply([]byte("a\r\nb")), // test binary safe
		protocol.NewNullBulkReply(),
		protocol.NewMultiBulkReply([][]byte{
			[]byte("a"),
			[]byte("\r\n"),
		}),
		protocol.NewEmptyMultiBulkReply(),
	}
	for _, re := range replies {
		result, err := parser.ParseOne(re.ToBytes())
		if err != nil {
			t.Error(err)
			continue
		}
		if !utils.BytesEquals(result.ToBytes(), re.ToBytes()) {
			t.Error("parse failed: " + string(re.ToBytes()))
		}
	}
}
