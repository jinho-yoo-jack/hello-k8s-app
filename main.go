package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// version 은 빌드할 때 -ldflags "-X main.version=..." 로 바뀐다 (Jenkins 가 빌드 번호와 커밋 해시를 넣는다).
var version = "dev"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname() // 쿠버네티스에서는 Pod 이름이 hostname 이다
	fmt.Fprintf(w, "Hello, Kubernetes! version=%s pod=%s\n", version, hostname)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/healthz", healthHandler)
	return mux
}

func main() {
	log.Printf("hello-k8s-app version=%s listening on :8080", version)
	log.Fatal(http.ListenAndServe(":8080", newMux()))
}
