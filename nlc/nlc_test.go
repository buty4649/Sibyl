package nlc

import (
    "bytes"
    "testing"
)

func TestbuildPacket(t *testing.T) {
    buf := new(bytes.Buffer)

    err := buildPacket(buf)
    if err != nil {
        t.Fatalf("failed to call buildPacket(): %s", err)
    }
}
