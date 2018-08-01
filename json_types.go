package main

import (
  "github.com/rynobey/bn256"
  "math/big"
  "fmt"
)

type Response struct {
  Num   *Number         `json:"number,omitempty"`
  P     *CurvePoint     `json:"curvepoint,omitempty"`
  Err   *Error          `json:"error,omitempty"`
}

type Text struct {
  T   string        `json:"t"`
}

type BinaryEcOpParams struct {
  A   *CurvePoint   `json:"a"`
  B   *CurvePoint   `json:"b"`
}

type ScalarEcOpParams struct {
  S   *Number       `json:"s"`
  A   *CurvePoint   `json:"a"`
}

type CurvePoint struct {
  X   string      `json:"x"`
  Y   string      `json:"y"`
}

func NewCurvePoint(P *bn256.G1) (*CurvePoint) {
  marshalledPoint := P.Marshal()
  x := fmt.Sprintf("0x%064x", marshalledPoint[0:32])
  y := fmt.Sprintf("0x%064x", marshalledPoint[32:64])
  return &CurvePoint{X: x, Y: y}
}

type BinaryOpParams struct {
  A   string      `json:"a"`
  B   string      `json:"b"`
}

type TernaryOpParams struct {
  A   string      `json:"a"`
  B   string      `json:"b"`
  C   string      `json:"c"`
}

type CommitmentInputs struct {
  B   string        `json:"b"`
  V   string        `json:"v"`
  H   *CurvePoint   `json:"h"`
  G   *CurvePoint   `json:"g"`
}

type Number struct {
  V   string    `json:"v"`
}

func NewNumber(num *big.Int) (*Number) {
  return &Number{V: fmt.Sprintf("0x%x", num)}
}

type Error struct {
  Msg string    `json:"msg"`
}
