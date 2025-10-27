package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    // レスポンスのContent-TypeをJSONに設定
    w.Header().Set("Content-Type", "application/json")
    
    // 1. クエリパラメータのマップを取得
	queryValues := r.URL.Query()

    // 2. 特定のパラメータを取得 (Get() はキーが存在しない場合、空文字列 "" を返します)
	name := queryValues.Get("name") 

    // 3. 取得した値を利用してJSONレスポンスを作成
    var response map[string]string
    if name == "" {
        response = map[string]string{"message": "Hello, stdlib API Server!"}
    }else {
        response = map[string]string{"message": fmt.Sprintf("Hello, %s", name)}
    }

    // JSONにエンコードして書き込む
    json.NewEncoder(w).Encode(response)

    // 200 OK を返します
	w.WriteHeader(http.StatusOK) 
}

func helloPathHandler(w http.ResponseWriter, r *http.Request) {
    // レスポンスのContent-TypeをJSONに設定
    w.Header().Set("Content-Type", "application/json")

    // 1. リクエストからパス全体を取得する
    path := r.URL.Path

    // 2. パスをスラッシュで分割する
    // 例: "/hello/john" -> ["", "hello", "john"]
	segments := strings.Split(path, "/")

    // 3. パラメータ（ID）を抽出する
	name := segments[2]

    // 抽出したIDが空でないか、または追加でバリデーションを行う
    if name == "" {
        http.Error(w, `{"error": "name is missing"}`, http.StatusBadRequest)
        return
    }

    // JSONレスポンスを作成
    response := map[string]string{"message": fmt.Sprintf("Hello, %s", name)}
    
    // JSONにエンコードして書き込む
    json.NewEncoder(w).Encode(response)

    // 200 OK を返します
	w.WriteHeader(http.StatusOK) 

}

// User構造体: クライアントから受け取るJSONデータに対応するGoの構造体
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// userCreateHandler は新しいユーザーを作成するためのPOSTリクエストを処理します。
func userCreateHandler(w http.ResponseWriter, r *http.Request) {
	// 1. HTTPメソッドの確認
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// 2. リクエストボディの読み取り
	var newUser User
	
	// json.NewDecoderを使ってストリームとして読み取る
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		// JSONの形式が不正な場合など
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// 3. データのバリデーションと処理
	if newUser.Name == "" || newUser.Email == "" {
		http.Error(w, `{"error": "Name and Email are required"}`, http.StatusBadRequest)
		return
	}
	
	// ここでデータベースにデータを保存するなどの実際の処理を行います
	log.Printf("Received new user: Name=%s, Email=%s", newUser.Name, newUser.Email)

	// 4. 成功レスポンスの返却
	w.Header().Set("Content-Type", "application/json")
	
	// 一般的にPOST成功時は 201 Created を返します
	w.WriteHeader(http.StatusCreated) 
	
	response := map[string]string{
		"message": "User created successfully",
		"id":      "auto-generated-id-123",
	}
	json.NewEncoder(w).Encode(response)
}

// userModifyHandler は /users/{id} の形式のリクエストを処理します。
func userModifyHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

	// 1. HTTPメソッドの確認
	if r.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// 2. リクエストボディの読み取り
	var newUser User
	
	// json.NewDecoderを使ってストリームとして読み取る
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		// JSONの形式が不正な場合など
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// 3. データのバリデーションと処理
	if newUser.Name == "" || newUser.Email == "" {
		http.Error(w, `{"error": "Name and Email are required"}`, http.StatusBadRequest)
		return
	}
	
	// ここでデータベースにデータを保存するなどの実際の処理を行います
	log.Printf("Received modify user: Name=%s, Email=%s", newUser.Name, newUser.Email)

	// 4. 成功レスポンスの返却
	w.Header().Set("Content-Type", "application/json")
	
	// 204 No Content を返す
	w.WriteHeader(http.StatusNoContent) 
}

func main() {
    // ルーター（ServeMux）を作成
    mux := http.NewServeMux()
    
    // hello apiのエンドポイントを追加する
    mux.HandleFunc("/hello", helloHandler)
    
    // hello name apiのエンドポイントを追加する
    mux.HandleFunc("/hello/", helloPathHandler)

    // /users エンドポイントに登録ハンドラーを登録
	mux.HandleFunc("/users", userCreateHandler)

    // /users エンドポイントに更新ハンドラーを登録
    // STLをだけで実装する場合はHTTPメソッドによって処理を分岐するということはできない。
    // パスを分けるかcaseなどを使用する必要がある。
	mux.HandleFunc("/users/modify", userModifyHandler)

    log.Println("Server listening on :8080...")
    // サーバーを起動
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}