package main

import (
	"fmt"
	"log"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login_with_idp.html") // 簡易的なログインフォーム
	return
}

func main() {
	// ログインページ
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", loginHandler)

	fmt.Println("IdP running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
