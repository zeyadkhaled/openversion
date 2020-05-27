package filterenc

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func Test_encryptRoundtrip(t *testing.T) {
	key, _ := hex.DecodeString("ce4f34331feab353c0a6c5f27f98097c8e81c65b1f0dac259074d0063e27eddd")

	e, err := aeadEncrypt(key, []byte(`{"test":"ok"}`))
	if err != nil {
		t.Fatal(err)
	}
	d, err := aeadDecrypt(key, e)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(d, []byte(`{"test":"ok"}`)) {
		t.Fatal("not equal")
	}

}
