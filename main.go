package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
)

// サンプルの認証情報
var users = map[string]string{
	"testuser": "password",
}

// セッションストア
var store = sessions.NewCookieStore([]byte("secret-key"))

// 認可コードを保存するための簡易データ構造
var authCodes = map[string]string{}

const (
	clientID     = "sample-client"
	clientSecret = "sample-secret"
	redirectURI  = "http://localhost:8081/callback"
)

// ログインページのハンドラー
func loginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// OPTIONSリクエストの処理
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "login.html") // 簡易的なログインフォーム
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// 認証
		if pass, ok := users[username]; ok && pass == password {
			// 認可コードを生成
			authCode := fmt.Sprintf("%d", time.Now().UnixNano())
			authCodes[authCode] = username

			// リダイレクト
			http.Redirect(w, r, fmt.Sprintf("%s?code=%s", redirectURI, authCode), http.StatusFound)
			return
		}

		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
}

// トークンエンドポイントのハンドラー
func tokenHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	clientIDParam := r.FormValue("client_id")
	code := r.FormValue("code")

	// クライアント認証
	if clientIDParam != clientID {
		http.Error(w, "Invalid client credentials", http.StatusUnauthorized)
		return
	}

	// 認可コードの検証
	username, ok := authCodes[code]
	if !ok {
		http.Error(w, "Invalid authorization code", http.StatusBadRequest)
		return
	}

	// IDトークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// 応答
	resp := map[string]string{
		"id_token": tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	// ログインページ
	http.HandleFunc("/login", loginHandler)
	// トークンエンドポイント
	http.HandleFunc("/token", tokenHandler)

	fmt.Println("IdP running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
