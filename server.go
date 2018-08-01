package main

import (
  "log"
  "net/http"
  "github.com/gorilla/mux"
)

var port = "8083"

func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/generate/commitment/", GenerateCommitment).Methods("POST")
  router.HandleFunc("/generate/keccak256/", GenerateKeccak256).Methods("POST")
  router.HandleFunc("/big/add/", BigIntAdd).Methods("POST")
  router.HandleFunc("/big/submod/", BigIntSubMod).Methods("POST")
  router.HandleFunc("/big/mul/", BigIntMul).Methods("POST")
  router.HandleFunc("/big/mod/", BigIntMod).Methods("POST")
  router.HandleFunc("/ec/order", ECOrder)
  router.HandleFunc("/ec/add/", ECAdd).Methods("POST")
  router.HandleFunc("/ec/sub/", ECSub).Methods("POST")
  router.HandleFunc("/ec/mul/", ECMul).Methods("POST")
  router.HandleFunc("/ec/basemul/", ECBaseMul).Methods("POST")
  router.HandleFunc("/ec/hashtopoint/", ECHashToPoint).Methods("POST")
  log.Fatal(http.ListenAndServe(":"+port, router))
}
