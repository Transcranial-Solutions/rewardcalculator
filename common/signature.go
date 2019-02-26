package common

import (
	"encoding/base64"
	"encoding/json"

	"github.com/icon-project/rewardcalculator/common/crypto"
)

type Signature struct {
	*crypto.Signature
}

func (sig Signature) MarshalJSON() ([]byte, error) {
	if sig.Signature == nil {
		return nil, nil
	}
	if bytes, err := sig.SerializeRSV(); err == nil {
		s := base64.StdEncoding.EncodeToString(bytes)
		return json.Marshal(s)
	} else {
		return nil, err
	}
}

func (sig *Signature) UnmarshalJSON(s []byte) error {
	var str string
	err := json.Unmarshal(s, &str)
	if err != nil {
		return err
	}
	if len(str) == 0 {
		return nil
	}
	if b, err := base64.StdEncoding.DecodeString(str); err == nil {
		if sig0, err := crypto.ParseSignature(b); err == nil {
			sig.Signature = sig0
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func (sig *Signature) MarshalBinary() ([]byte, error) {
	return sig.Signature.SerializeRSV()
}

func (sig *Signature) UnmarshalBinary(s []byte) error {
	sig0, err := crypto.ParseSignature(s)
	if err == nil {
		sig.Signature = sig0
	}
	return err
}
