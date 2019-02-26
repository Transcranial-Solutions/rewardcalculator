package common

import (
	"bytes"
	"testing"

	"github.com/icon-project/rewardcalculator/common/codec"
	"github.com/icon-project/rewardcalculator/common/crypto"
)

func TestSignatureCoding(t *testing.T) {
	obs := []byte("01234567890123456789012345678901234567890123456789012345678901234")
	var sig Signature
	var err error
	sig.Signature, err = crypto.ParseSignature(obs)
	if err != nil {
		t.Fail()
	}
	sigBS, err := codec.MP.MarshalToBytes(sig)
	if err != nil {
		t.Fail()
	}
	var sig2 Signature
	_, err = codec.MP.UnmarshalFromBytes(sigBS, &sig2)
	if err != nil {
		t.Fail()
	}
	rsv, err := sig2.SerializeRSV()
	if err != nil {
		t.Fail()
	}
	if !bytes.Equal(obs, rsv) {
		t.Fail()
	}
}
