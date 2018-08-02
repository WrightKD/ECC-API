package main

import (
  "net/http"
  "encoding/json"
  "github.com/rynobey/bn256"
)

func ECOrder(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  encoder.Encode(Response{Num: NewNumber(bn256.Order)})
}

func ECAdd(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var binaryEcOpParams BinaryEcOpParams
  err := ReadContentsIntoStruct(r, &binaryEcOpParams)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  A, err := NewECPoint(binaryEcOpParams.A.X, binaryEcOpParams.A.Y, err)
  B, err := NewECPoint(binaryEcOpParams.B.X, binaryEcOpParams.B.Y, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  ans := new(bn256.G1).Add(A, B)
  curvePoint := NewCurvePoint(ans)
  encoder.Encode(Response{P: curvePoint})
}

func ECSub(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var binaryEcOpParams BinaryEcOpParams
  err := ReadContentsIntoStruct(r, &binaryEcOpParams)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  A, err := NewECPoint(binaryEcOpParams.A.X, binaryEcOpParams.A.Y, err)
  B, err := NewECPoint(binaryEcOpParams.B.X, binaryEcOpParams.B.Y, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  ans := new(bn256.G1).Add(A, B.Neg(B))
  curvePoint := NewCurvePoint(ans)
  encoder.Encode(Response{P: curvePoint})
}

func ECMul(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var scalarEcOpParams ScalarEcOpParams
  err := ReadContentsIntoStruct(r, &scalarEcOpParams)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  s, err := NewBigInt(scalarEcOpParams.S.V, err)
  A, err := NewECPoint(scalarEcOpParams.A.X, scalarEcOpParams.A.Y, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  ans := new(bn256.G1).ScalarMult(A, s)
  curvePoint := NewCurvePoint(ans)
  encoder.Encode(Response{P: curvePoint})
}

func ECBaseMul(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var number Number
  err := ReadContentsIntoStruct(r, &number)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  s, err := NewBigInt(number.V, err)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  ans := new(bn256.G1).ScalarBaseMult(s)
  curvePoint := NewCurvePoint(ans)
  encoder.Encode(Response{P: curvePoint})
}

func ECHashToPoint(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  var text Text
  err := ReadContentsIntoStruct(r, &text)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  A := new(bn256.G1).Hash(text.T)
  if err != nil {
    encoder.Encode(Response{Err: &Error{Msg: err.Error()}})
    return
  }
  curvePoint := NewCurvePoint(A)
  encoder.Encode(Response{P: curvePoint})
}
