package btctx

import (
	"bytes"
	"encoding/hex"
	"io"
)

type TxOut struct {
	Value    int64   `json:"value"`
	PkScript *Script `json:"pkScript"`
}

func NewTxOut() *TxOut {
	return &TxOut{}
}

func (txOut *TxOut) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := txOut.WriteAll(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (txOut *TxOut) WriteAll(w io.Writer) error {
	if err := txOut.WriteValue(w); err != nil {
		return err
	}

	if err := txOut.WritePkScriptLength(w); err != nil {
		return err
	}

	if err := txOut.WritePkScript(w); err != nil {
		return err
	}

	return nil
}

func (txOut *TxOut) WriteValue(w io.Writer) error {
	return writeData(w, txOut.Value)
}

func (txOut *TxOut) WritePkScriptLength(w io.Writer) error {
	b, err := hex.DecodeString(txOut.PkScript.Hex)
	if err != nil {
		return err
	}

	return writeVarInt(w, len(b))
}

func (txOut *TxOut) WritePkScript(w io.Writer) error {
	return writeHex(w, txOut.PkScript.Hex)
}
