package main

import (
  "testing"
  "crypto/rand"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "encoding/hex"
  "github.com/rynobey/bn256"
  "github.com/ethereum/go-ethereum/crypto/sha3"
  "math/big"
  "bytes"
  "fmt"
)

func TestECOrder(t *testing.T) {
  q, _ := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
  response, err := http.Get("http://localhost:" + port + "/ec/order")
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  qAPI, ok := new(big.Int).SetString(res.Num.V[2:], 16)
  if !ok {
    t.Errorf("An error occurred while initializing big.Int from string")
    return
  }
  if (q.Cmp(qAPI) != 0) {
    t.Errorf("Invalid value returned\n")
  }
}

func TestECAdd(t *testing.T) {
  A := new(bn256.G1).ScalarBaseMult(new(big.Int).SetInt64(10))
  marshalledBytesA := A.Marshal()
  B := new(bn256.G1).ScalarBaseMult(new(big.Int).SetInt64(100))
  marshalledBytesB := B.Marshal()
  a_pt := &CurvePoint{X: fmt.Sprintf("0x%x", marshalledBytesA[:32]), Y:fmt.Sprintf("0x%x", marshalledBytesA[32:64])}
  b_pt := &CurvePoint{X: fmt.Sprintf("0x%x", marshalledBytesB[:32]), Y:fmt.Sprintf("0x%x", marshalledBytesB[32:64])}
  binaryEcOpParams := BinaryEcOpParams{A: a_pt, B: b_pt}
  marshalledJSON, _ := json.Marshal(binaryEcOpParams)
  response, err := http.Post("http://localhost:" + port + "/ec/add/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  pt := res.P
  P := new(bn256.G1)
  marshalledPoint := pt.X[2:] + pt.Y[2:]
  marshalledBytes, err := hex.DecodeString(marshalledPoint)
  if err != nil {
    t.Errorf("An error occurred while decoding hex string: %s\n", err)
    return
  }
  _, err = P.Unmarshal(marshalledBytes)
  if err != nil {
    t.Errorf("An error occurred while unmarshalling BN256 curve point: %s\n", err)
    return
  }
  Ptest := new(bn256.G1).Add(A, B)
  if (P.String() != Ptest.String()) {
    t.Errorf("Incorrect answer returned\n")
  }
}

func TestECSub(t *testing.T) {
  A := new(bn256.G1).ScalarBaseMult(new(big.Int).SetInt64(10))
  marshalledBytesA := A.Marshal()
  B := new(bn256.G1).ScalarBaseMult(new(big.Int).SetInt64(100))
  marshalledBytesB := B.Marshal()
  a_pt := &CurvePoint{X: fmt.Sprintf("0x%064x", marshalledBytesA[:32]), Y:fmt.Sprintf("0x%064x", marshalledBytesA[32:64])}
  b_pt := &CurvePoint{X: fmt.Sprintf("0x%064x", marshalledBytesB[:32]), Y:fmt.Sprintf("0x%064x", marshalledBytesB[32:64])}
  binaryEcOpParams := BinaryEcOpParams{A: a_pt, B: b_pt}
  marshalledJSON, _ := json.Marshal(binaryEcOpParams)
  response, err := http.Post("http://localhost:" + port + "/ec/sub/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  pt := res.P
  P := new(bn256.G1)
  marshalledPoint := pt.X[2:] + pt.Y[2:]
  marshalledBytes, err := hex.DecodeString(marshalledPoint)
  if err != nil {
    t.Errorf("An error occurred while decoding hex string: %s\n", err)
    return
  }
  _, err = P.Unmarshal(marshalledBytes)
  if err != nil {
    t.Errorf("An error occurred while unmarshalling BN256 curve point: %s\n", err)
    return
  }
  Ptest := new(bn256.G1).Add(A, B.Neg(B))
  if (P.String() != Ptest.String()) {
    t.Errorf("Incorrect answer returned\n")
  }
}

func TestECMul(t *testing.T) {
  s := new(big.Int).SetInt64(100)
  A := new(bn256.G1).ScalarBaseMult(new(big.Int).SetInt64(10))
  marshalledBytesA := A.Marshal()
  a_pt := &CurvePoint{X: fmt.Sprintf("0x%064x", marshalledBytesA[:32]), Y:fmt.Sprintf("0x%064x", marshalledBytesA[32:64])}
  s_val := NewNumber(s)
  scalarEcOpParams := ScalarEcOpParams{S: s_val, A: a_pt}
  marshalledJSON, _ := json.Marshal(scalarEcOpParams)
  response, err := http.Post("http://localhost:" + port + "/ec/mul/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  pt := res.P
  P := new(bn256.G1)
  marshalledPoint := pt.X[2:] + pt.Y[2:]
  marshalledBytes, err := hex.DecodeString(marshalledPoint)
  if err != nil {
    t.Errorf("An error occurred while decoding hex string: %s\n", err)
    return
  }
  _, err = P.Unmarshal(marshalledBytes)
  if err != nil {
    t.Errorf("An error occurred while unmarshalling BN256 curve point: %s\n", err)
    return
  }
  Ptest := new(bn256.G1).ScalarMult(A, s)
  if (P.String() != Ptest.String()) {
    t.Errorf("Incorrect answer returned\n")
  }
}

func TestECBaseMul(t *testing.T) {
  s := new(big.Int).SetInt64(100)
  s_val := NewNumber(s)
  marshalledJSON, _ := json.Marshal(s_val)
  response, err := http.Post("http://localhost:" + port + "/ec/basemul/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  pt := res.P
  P := new(bn256.G1)
  marshalledPoint := pt.X[2:] + pt.Y[2:]
  marshalledBytes, err := hex.DecodeString(marshalledPoint)
  if err != nil {
    t.Errorf("An error occurred while decoding hex string: %s\n", err)
    return
  }
  _, err = P.Unmarshal(marshalledBytes)
  if err != nil {
    t.Errorf("An error occurred while unmarshalling BN256 curve point: %s\n", err)
    return
  }
  Ptest := new(bn256.G1).ScalarBaseMult(s)
  if (P.String() != Ptest.String()) {
    t.Errorf("Incorrect answer returned\n")
  }
}

func TestECHashToPoint(t *testing.T) {
  str := "input to hash function"
  t_val := Text{T: str}
  marshalledJSON, _ := json.Marshal(t_val)
  response, err := http.Post("http://localhost:" + port + "/ec/hashtopoint/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  pt := res.P
  P := new(bn256.G1)
  marshalledPoint := pt.X[2:] + pt.Y[2:]
  marshalledBytes, err := hex.DecodeString(marshalledPoint)
  if err != nil {
    t.Errorf("An error occurred while decoding hex string: %s\n", err)
    return
  }
  _, err = P.Unmarshal(marshalledBytes)
  if err != nil {
    t.Errorf("An error occurred while unmarshalling BN256 curve point: %s\n", err)
    return
  }
  Ptest := new(bn256.G1)
  marshalledPointTest := "05167cc8fa5cc80098f29d6df8c65a1f1c956dd407073f33a5a75e2dee028a7e08cf95e70b3c8a91341b9fea15e985c7660136dcff651789925d79e98eed2d15"
  marshalledBytesTest, err := hex.DecodeString(marshalledPointTest)
  _, err = Ptest.Unmarshal(marshalledBytesTest)
  if (P.String() != Ptest.String()) {
    t.Errorf("Incorrect answer returned\n")
  }
}

func TestGenerateCommitment(t *testing.T) {
  var testBlind = int64(4563452349857)
  b := new(big.Int).SetInt64(testBlind)
  var testValue = int64(123445)
  v := new(big.Int).SetInt64(testValue)
  H := new(bn256.G1).ScalarBaseMult(new(big.Int).SetInt64(10))
  marshalledBytesA := H.Marshal()
  G := new(bn256.G1).ScalarBaseMult(new(big.Int).SetInt64(100))
  marshalledBytesB := G.Marshal()
  h_pt := &CurvePoint{X: fmt.Sprintf("0x%064x", marshalledBytesA[:32]), Y:fmt.Sprintf("0x%064x", marshalledBytesA[32:64])}
  g_pt := &CurvePoint{X: fmt.Sprintf("0x%064x", marshalledBytesB[:32]), Y:fmt.Sprintf("0x%064x", marshalledBytesB[32:64])}
  commitmentInputs := CommitmentInputs{B: fmt.Sprintf("0x%x", testBlind), V: fmt.Sprintf("0x%x", testValue), H: h_pt, G: g_pt}
  marshalledJSON, _ := json.Marshal(commitmentInputs)
  response, err := http.Post("http://localhost:" + port + "/generate/commitment/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  commitment := res.P
  C := new(bn256.G1)
  marshalledPoint := commitment.X[2:] + commitment.Y[2:]
  marshalledBytes, err := hex.DecodeString(marshalledPoint)
  if err != nil {
    t.Errorf("An error occurred while decoding hex string: %s\n", err)
    return
  }
  _, err = C.Unmarshal(marshalledBytes)
  if err != nil {
    t.Errorf("An error occurred while unmarshalling BN256 curve point: %s\n", err)
    return
  }
  bbHH := new(bn256.G1).ScalarMult(H, b)
  vvGG := new(bn256.G1).ScalarMult(G, v)
  Ct := bbHH.Add(bbHH, vvGG)
  if (C.String() != Ct.String()) {
    t.Errorf("Invalid commitment returned\n")
  }
}

func TestGenerateSchnorr(t *testing.T) {
  x, _ := rand.Int(rand.Reader, bn256.Order)
  P := new(bn256.G1).ScalarBaseMult(x)
  m := "This is the message to be signed"
  generateSchnorrInputs := GenerateSchnorrInputs{P: NewCurvePoint(P), X: fmt.Sprintf("0x%064x", x), M: m}
  marshalledJSON, _ := json.Marshal(generateSchnorrInputs)
  response, err := http.Post("http://localhost:" + port + "/generate/schnorr/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  sig := res.Sig
  P_out := new(bn256.G1)
  marshalledPoint := sig.P.X[2:] + sig.P.Y[2:]
  marshalledBytes, err := hex.DecodeString(marshalledPoint)
  if err != nil {
    t.Errorf("An error occurred while decoding hex string: %s\n", err)
    return
  }
  _, err = P_out.Unmarshal(marshalledBytes)
  if err != nil {
    t.Errorf("An error occurred while unmarshalling BN256 curve point: %s\n", err)
    return
  }
  E_out, _ := new(big.Int).SetString(sig.E[2:], 16)
  S_out, _ := new(big.Int).SetString(sig.S[2:], 16)
  M_out := sig.M
  isValid, err := VerifySchnorrSignature(P_out, M_out, E_out, S_out, err)
  if (!isValid) {
    t.Errorf("Invalid Schnorr signature generated")
  }
}

func TestVerifySchnorr(t *testing.T) {
  x, _ := rand.Int(rand.Reader, bn256.Order)
  P := new(bn256.G1).ScalarBaseMult(x)
  m := "This is the message to be signed"
  P_out, _, M_out, E_out, S_out, err := GenerateSchnorrSignature(P, m, x, nil)
  schnorrSignature := SchnorrSignature{P: NewCurvePoint(P_out), M: M_out, E: fmt.Sprintf("0x%064x", E_out), S: fmt.Sprintf("0x%064x", S_out)}
  marshalledJSON, _ := json.Marshal(schnorrSignature)
  response, err := http.Post("http://localhost:" + port + "/verify/schnorr/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  if (res.Text != "true") {
    t.Errorf("Invalid Schnorr signature generated")
  }
}

func TestBigAdd(t *testing.T) {
  a, _ := new(big.Int).SetString("20222222222222222222222222222222222222222222222222222222222222222222222222222", 10)
  b, _ := new(big.Int).SetString("11111111111111111111111111111111111111111111111111111111111111111111111111111", 10)
  binaryOpParams := BinaryOpParams{A: fmt.Sprintf("0x%x", a), B: fmt.Sprintf("0x%x", b)}
  marshalledJSON, _ := json.Marshal(binaryOpParams)
  response, err := http.Post("http://localhost:" + port + "/big/add/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  number := res.Num
  ansAPI, _ := new(big.Int).SetString(number.V[2:], 16)
  ans := new(big.Int).Add(a, b)
  if(ans.Cmp(ansAPI) != 0) {
    t.Errorf("Wrong answer returned")
    return
  }
}

func TestBigSubMod(t *testing.T) {
  a, _ := new(big.Int).SetString("20222222222222222222222222222222222222222222222222222222222222222222222222222", 10)
  b, _ := new(big.Int).SetString("11111111111111111111111111111111111111111111111111111111111111111111111111111", 10)
  c, _ := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
  ternaryOpParams := TernaryOpParams{A: fmt.Sprintf("0x%x", a), B: fmt.Sprintf("0x%x", b), C: fmt.Sprintf("0x%x", c)}
  marshalledJSON, _ := json.Marshal(ternaryOpParams)
  response, err := http.Post("http://localhost:" + port + "/big/submod/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  number := res.Num
  ansAPI, _ := new(big.Int).SetString(number.V[2:], 16)
  ans := new(big.Int).Sub(a, b)
  ans.Mod(ans, c)
  if(ans.Cmp(ansAPI) != 0) {
    t.Errorf("Wrong answer returned")
    return
  }
}

func TestBigMul(t *testing.T) {
  a, _ := new(big.Int).SetString("20222222222222222222222222222222222222222222222222222222222222222222222222222", 10)
  b, _ := new(big.Int).SetString("11111111111111111111111111111111111111111111111111111111111111111111111111111", 10)
  binaryOpParams := BinaryOpParams{A: fmt.Sprintf("0x%x", a), B: fmt.Sprintf("0x%x", b)}
  marshalledJSON, _ := json.Marshal(binaryOpParams)
  response, err := http.Post("http://localhost:" + port + "/big/mul/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  number := res.Num
  ansAPI, _ := new(big.Int).SetString(number.V[2:], 16)
  ans := new(big.Int).Mul(a, b)
  if(ans.Cmp(ansAPI) != 0) {
    t.Errorf("Wrong answer returned")
    return
  }
}

func TestBigMod(t *testing.T) {
  a, _ := new(big.Int).SetString("50222222222222222222222222222222222222222222222222222222222222222222222222222", 10)
  b, _ := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
  binaryOpParams := BinaryOpParams{A: fmt.Sprintf("0x%x", a), B: fmt.Sprintf("0x%x", b)}
  marshalledJSON, _ := json.Marshal(binaryOpParams)
  response, err := http.Post("http://localhost:" + port + "/big/mod/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  number := res.Num
  ansAPI, _ := new(big.Int).SetString(number.V[2:], 16)
  ans := new(big.Int).Mod(a, b)
  if(ans.Cmp(ansAPI) != 0) {
    t.Errorf("Wrong answer returned")
    return
  }
}

func TestInvMod(t *testing.T) {
  a, _ := new(big.Int).SetString("20222222222222222222222222222222222222222222222222222222222222222222222222222", 10)
  b, _ := new(big.Int).SetString("11111111111111111111111111111111111111111111111111111111111111111111111111111", 10)
  binaryOpParams := BinaryOpParams{A: fmt.Sprintf("0x%x", a), B: fmt.Sprintf("0x%x", b)}
  marshalledJSON, _ := json.Marshal(binaryOpParams)
  response, err := http.Post("http://localhost:" + port + "/big/invmod/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  number := res.Num
  ansAPI, _ := new(big.Int).SetString(number.V[2:], 16)
  ans := new(big.Int).ModInverse(a, b)
  if(ans.Cmp(ansAPI) != 0) {
    t.Errorf("Wrong answer returned")
    return
  }
}

func TestGenerateKeccak256(t *testing.T) {
  str := "input to hash function"
  t_val := Text{T: str}
  marshalledJSON, _ := json.Marshal(t_val)
  response, err := http.Post("http://localhost:" + port + "/generate/keccak256/", "application/json", bytes.NewBuffer(marshalledJSON))
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  if res.Err != nil && res.Err.Msg != "" {
    t.Errorf(fmt.Sprintf("An error occurred: %s\n", res.Err.Msg))
    return
  }
  number := res.Num
  ansAPI, _ := new(big.Int).SetString(number.V[2:], 16)
  h := sha3.NewKeccak256()
  h.Reset()
  h.Write([]byte(str))
  ans, _ := new(big.Int).SetString(fmt.Sprintf("%x", h.Sum(nil)), 16)
  if (ans.String() != ansAPI.String()) {
    t.Errorf("Incorrect answer returned\n")
  }
}

func TestIsAlive(t *testing.T) {
  response, err := http.Get("http://localhost:" + port + "/isalive")
  if err != nil {
    t.Errorf("An error occurred while making request to API: %s\n", err)
    return
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    t.Errorf("An error occurred while reading response body: %s\n", err)
    return
  }
  var res Response
  err = json.Unmarshal(contents, &res)
  if err != nil {
    t.Errorf("An error occurred while reading into JSON object: %s\n", err)
    return
  }
  s := res.Text
  if (s != "It's alive!") {
    t.Errorf("Invalid value returned\n")
  }
}
