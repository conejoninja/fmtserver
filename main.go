// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/julienschmidt/httprouter"
)

type Response struct {
	Msg   string `json:"msg,omitempty"`
	Code  string `json:"code"`
	Error int    `json:"error"`
}

func main() {
	router := httprouter.New()
	router.GET("/", handlerInfo)
	router.POST("/", handler)
	fmt.Println("Server Up")
	log.Fatal(http.ListenAndServe(":8737", router))
	fmt.Println("Server Down")
}

func handlerInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("INFO", r)
	fmt.Fprintf(w, jsonResponse("Server is up", 0, ""))
}

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("E1")
		http.Error(w, jsonResponse("invalid request", 1, ""), http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		fmt.Println("E2")
		http.Error(w, jsonResponse("no body", 2, ""), http.StatusBadRequest)
		return
	}

	cmd := exec.Command("gofmt")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("E3", err)
		http.Error(w, jsonResponse(err.Error(), 3, ""), http.StatusBadRequest)
		return
	}

	/*bodyStr := string(body)
	fmt.Println("BODY", bodyStr)
	//bodyStr = bodyStr[1:len(bodyStr)-1]
	bodyDecoded, err := base64.StdEncoding.DecodeString(bodyStr)
	fmt.Println("E4", err)
	http.Error(w, jsonResponse(err.Error(), 4, ""), http.StatusBadRequest) */


	io.WriteString(stdin, string(body))
	stdin.Close()

	if err := cmd.Run(); err != nil {
		fmt.Println("E5", err)
		http.Error(w, jsonResponse(err.Error(), 5, ""), http.StatusBadRequest)
		return
	}

	fmt.Println("O5", stdout.String())
	fmt.Fprintf(w, jsonResponse("ok", 0, stdout.String()))
}

func jsonResponse(msg string, err int, code string) string {
	var r Response
	r.Msg = msg
	r.Error = err
	r.Code = code
	j, _ := json.Marshal(r)
	return string(j)
}
