package main

import (
  "fmt"
  "net/http"
  "encoding/json"
  "math/big"
  "github.com/rynobey/bn256"
  "github.com/ethereum/go-ethereum/crypto/sha3"
)

func GenerateKeccak256(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var text Text
  err := ReadContentsIntoStruct(r, &text)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  h := sha3.NewKeccak256()
  h.Reset()
  h.Write([]byte(text.T))
  out, _ := new(big.Int).SetString(fmt.Sprintf("%x", h.Sum(nil)), 16)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  encoder.Encode(Response{Num: NewNumber(out)})
}

func GenerateCommitment(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var commitmentInputs CommitmentInputs
  err := ReadContentsIntoStruct(r, &commitmentInputs)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  b, err := NewBigInt(commitmentInputs.B, err)
  v, err := NewBigInt(commitmentInputs.V, err)
  H, err := NewECPoint(commitmentInputs.H.X, commitmentInputs.H.Y, err)
  G, err := NewECPoint(commitmentInputs.G.X, commitmentInputs.G.Y, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  bbHH := new(bn256.G1).ScalarMult(H, b)
  vvGG := new(bn256.G1).ScalarMult(G, v)
  C := bbHH.Add(bbHH, vvGG)
  commitment := NewCurvePoint(C)
  encoder.Encode(Response{P: commitment})
}

func GenerateSchnorr(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var generateSchnorrInputs GenerateSchnorrInputs
  err := ReadContentsIntoStruct(r, &generateSchnorrInputs)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  X, err := NewBigInt(generateSchnorrInputs.Priv, err)
  M := generateSchnorrInputs.M
  P_out, K_out, M_out, E_out, S_out, _ := GenerateSchnorrSignature(M, X, err)
  encoder.Encode(Response{Sig: &SchnorrSignature{P: NewCurvePoint(P_out), K: NewCurvePoint(K_out), M: M_out, E: fmt.Sprintf("0x%064x", E_out), S: fmt.Sprintf("0x%064x", S_out)}})
}
