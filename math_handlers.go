package main

import (
  "net/http"
  "encoding/json"
  "math/big"
  "github.com/rynobey/bn256"
  "crypto/rand"
)

func CryptoRandBigInt(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  num, _ := rand.Int(rand.Reader, bn256.Order)
  encoder.Encode(Response{Num: NewNumber(num)})
}

func BigIntAdd(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var binaryOpParams BinaryOpParams
  err := ReadContentsIntoStruct(r, &binaryOpParams)
  a, err := NewBigInt(binaryOpParams.A, err)
  b, err := NewBigInt(binaryOpParams.B, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  ans := new(big.Int).Add(a, b)
  encoder.Encode(Response{Num: NewNumber(ans)})
}

func BigIntSubMod(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var ternaryOpParams TernaryOpParams
  err := ReadContentsIntoStruct(r, &ternaryOpParams)
  a, err := NewBigInt(ternaryOpParams.A, err)
  b, err := NewBigInt(ternaryOpParams.B, err)
  c, err := NewBigInt(ternaryOpParams.C, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  ans := new(big.Int).Sub(a, b)
  ans.Mod(ans, c)
  encoder.Encode(Response{Num: NewNumber(ans)})
}

func BigIntInvMod(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var binaryOpParams BinaryOpParams
  err := ReadContentsIntoStruct(r, &binaryOpParams)
  a, err := NewBigInt(binaryOpParams.A, err)
  b, err := NewBigInt(binaryOpParams.B, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  ans := new(big.Int).ModInverse(a, b)
  encoder.Encode(Response{Num: NewNumber(ans)})
}

func BigIntMul(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var binaryOpParams BinaryOpParams
  err := ReadContentsIntoStruct(r, &binaryOpParams)
  a, err := NewBigInt(binaryOpParams.A, err)
  b, err := NewBigInt(binaryOpParams.B, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  ans := new(big.Int).Mul(a, b)
  encoder.Encode(Response{Num: NewNumber(ans)})
}

func BigIntMod(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var binaryOpParams BinaryOpParams
  err := ReadContentsIntoStruct(r, &binaryOpParams)
  a, err := NewBigInt(binaryOpParams.A, err)
  b, err := NewBigInt(binaryOpParams.B, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  ans := new(big.Int).Mod(a, b)
  encoder.Encode(Response{Num: NewNumber(ans)})
}
