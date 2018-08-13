package main

import (
  "fmt"
  "net/http"
  "encoding/json"
)

func VerifySchnorr(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var schnorrSignature SchnorrSignature
  err := ReadContentsIntoStruct(r, &schnorrSignature)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  P, err := NewECPoint(schnorrSignature.P.X, schnorrSignature.P.Y, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  M := schnorrSignature.M
  E, err := NewBigInt(schnorrSignature.E, err)
  S, err := NewBigInt(schnorrSignature.S, err)
  isValid, err := VerifySchnorrSignature(P, M, E, S, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  encoder.Encode(Response{Text: fmt.Sprintf("%t", isValid)})
}
