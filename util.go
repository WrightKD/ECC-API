package main

import (
  "errors"
  "crypto/rand"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "math/big"
  "github.com/rynobey/bn256"
  "encoding/hex"
  "fmt"
  "github.com/ethereum/go-ethereum/crypto/sha3"
)

func Hex32ByteChunksToStr(hexChunks []string) (string) {
  hexStr := ""
  lenArr := len(hexChunks)
  for i := lenArr-1; i > 0; i-- {
    hexStr = fmt.Sprintf("%s%s", hexChunks[i][2:], hexStr)
  }
  // remove leading zeros while keeping length even
  lastChunk := hexChunks[0]
  curIndex := 2
  curDigit := lastChunk[curIndex]
  for i := 2; curDigit == '0'; i++ {
    curIndex = i
    curDigit = lastChunk[curIndex]
  }
  if curIndex % 2 == 1 { curIndex -= 1 }
  hexStr = fmt.Sprintf("0x%s%s", hexChunks[0][curIndex:], hexStr)
  return hexStr
}

func HexStrTo32ByteChunks(hexStr string) ([]string) { // 32 bytes in hex = 64 digits
  if (hexStr[:2] == "0x") { hexStr = hexStr[2:] }
  length := len(hexStr)
  lenArr := length/64
  if (length % 64 != 0) { lenArr += 1 }
  hexStrArr := make([]string, lenArr)
  start := length - 64
  end := length
  for i := lenArr-1; i > 0; i-- {
    hexStrArr[i] = fmt.Sprintf("0x%064s", hexStr[start:end])
    start -= 64
    end -= 64
  }
  hexStrArr[0] = fmt.Sprintf("0x%064s", hexStr[:end])
  return hexStrArr
}

func IsZero(num *big.Int) (bool) {
  return (num.Cmp(new(big.Int).SetInt64(0)) == 0)
}

func ReadContentsIntoStruct(r *http.Request, obj interface{}) (error) {
  contents, err := ioutil.ReadAll(r.Body)
  defer r.Body.Close()
  if err != nil { return err }
  err = json.Unmarshal(contents, &obj)
  if err != nil { return err }
  return nil
}

func NewBigInt(num string, err error) (*big.Int, error) {
  if err != nil {
    return nil, err
  } else {
    if len(num) < 3 {
      return nil, errors.New("Unable to initialize big.Int from string: too short")
    }
    bn, ok := new(big.Int).SetString(num[2:], 16)
    if !ok {
      return nil, errors.New("Failed to initialize big.Int from string")
    }
    return bn, nil
  }
}

func NewECPoint(xCoord string, yCoord string, err error) (*bn256.G1, error) {
  if err != nil {
    return nil, err
  } else {
    P := new(bn256.G1)
    marshalledPoint := fmt.Sprintf("%064s%064s", xCoord[2:], yCoord[2:])
    marshalledBytes, err := hex.DecodeString(marshalledPoint)
    if err != nil {
      return nil, err
    }
    _, err = P.Unmarshal(marshalledBytes)
    if err != nil {
      return nil, err
    }
    return P, nil
  }
}

func GenerateSchnorrSignature(P *bn256.G1, M string, X *big.Int, err error) (*bn256.G1, string, *big.Int, *big.Int, error) {
  if err != nil {
    return nil, "", nil, nil, err
  } else {
    k, _ := rand.Int(rand.Reader, bn256.Order)
    kG := new(bn256.G1).ScalarBaseMult(k)
    h := sha3.NewKeccak256()
    h.Reset()
    h.Write([]byte(fmt.Sprintf("%s%s%s", M, P, kG)))
    e, _ := new(big.Int).SetString(fmt.Sprintf("%x", h.Sum(nil)), 16)
    s := new(big.Int).Mod(new(big.Int).Add(k, new(big.Int).Mul(e, X)), bn256.Order)
    return P, M, e, s, nil
  }
}

func VerifySchnorrSignature(P *bn256.G1, M string, E, S *big.Int, err error) (bool, error) {
  if err != nil {
    return false, err
  } else {
    sG := new(bn256.G1).ScalarBaseMult(S)
    eP := new(bn256.G1).ScalarMult(P, E)
    kG := new(bn256.G1).Add(sG, eP.Neg(eP))
    h := sha3.NewKeccak256()
    h.Reset()
    h.Write([]byte(fmt.Sprintf("%s%s%s", M, P, kG)))
    e, _ := new(big.Int).SetString(fmt.Sprintf("%x", h.Sum(nil)), 16)
    return (e.Cmp(E) == 0), nil
  }
}
