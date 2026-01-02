package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AddUserRequest struct {
	CodeforcesHandle string `json:"codeforces_handle" binding:"required"`
}


func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run add_user.go <handle>")
		return
	}

	reqBody := AddUserRequest{
		CodeforcesHandle: os.Args[1],
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	fmt.Println("Sending:", string(data))

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/api/v1/users",
		bytes.NewBuffer(data),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(body))
}