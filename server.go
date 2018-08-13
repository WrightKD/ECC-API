package main

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
)

var port = "8083"

func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/isalive", IsAlive).Methods("GET")
  router.HandleFunc("/generate/keccak256/", GenerateKeccak256).Methods("POST")
  router.HandleFunc("/generate/commitment/", GenerateCommitment).Methods("POST")
  router.HandleFunc("/generate/schnorr/", GenerateSchnorr).Methods("POST")
  router.HandleFunc("/verify/schnorr/", VerifySchnorr).Methods("POST")
  router.HandleFunc("/big/add/", BigIntAdd).Methods("POST")
  router.HandleFunc("/big/submod/", BigIntSubMod).Methods("POST")
  router.HandleFunc("/big/invmod/", BigIntInvMod).Methods("POST")
  router.HandleFunc("/big/mul/", BigIntMul).Methods("POST")
  router.HandleFunc("/big/mod/", BigIntMod).Methods("POST")
  router.HandleFunc("/big/rand", CryptoRandBigInt).Methods("GET")
  router.HandleFunc("/ec/order", ECOrder)
  router.HandleFunc("/ec/add/", ECAdd).Methods("POST")
  router.HandleFunc("/ec/sub/", ECSub).Methods("POST")
  router.HandleFunc("/ec/mul/", ECMul).Methods("POST")
  router.HandleFunc("/ec/basemul/", ECBaseMul).Methods("POST")
  router.HandleFunc("/ec/hashtopoint/", ECHashToPoint).Methods("POST")
  fmt.Printf("Listening on port %s\n", port)
  log.Fatal(http.ListenAndServe(":"+port, router))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  encoder.Encode(Response{Text: "It's alive!"})
}
