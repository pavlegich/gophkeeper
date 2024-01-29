// Package app contains the main methods for running the client.
package app

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pavlegich/gophkeeper/internal/server/domains/user"
	_ "go.uber.org/automaxprocs"
)

// Run initialized the main app components and runs the client.
func Run() error {
	// Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	target := "http://localhost:8080/api/user"

	// =============
	// Login
	// =============

	login := "/login"
	u := user.User{
		Login:    "testUser",
		Password: "qwerty",
	}
	reqPost, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
	}

	resp, err := http.Post(target+login, "application/json", bytes.NewBuffer(reqPost))
	if err != nil {
		log.Println(err)
	}

	log.Println("login", resp.StatusCode)

	cookie := resp.Cookies()[0]

	// =============
	// Create data
	// =============

	var buf bytes.Buffer
	multipartWriter := multipart.NewWriter(&buf)
	defer multipartWriter.Close()

	// add form field

	metaPart, _ := multipartWriter.CreateFormField("metadata")
	meta, err := json.Marshal(map[string]string{"multipart": "first"})
	if err != nil {
		log.Println(err)
	}
	metaPart.Write([]byte(meta))

	textPart, _ := multipartWriter.CreateFormField("data")
	if err != nil {
		log.Println(err)
	}
	textPart.Write([]byte("hello from text mulipart client"))

	create := "/data/binary/clientNotMetadataType"

	// d := data.Data{
	// 	Name: "clientData",
	// 	Type: "text",
	// 	Data: "hello from the client",
	// 	Metadata: data.Metadata{
	// 		"first": 1,
	// 	},
	// }

	// buf, err := json.Marshal(d)
	// if err != nil {
	// 	log.Println(err)
	// }

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target+create, &buf)
	if err != nil {
		log.Println(err)
	}
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	log.Println("create data", resp.StatusCode)

	// =============
	// Update data
	// =============

	// update := "/data"

	// d := data.Data{
	// 	Name: "clientData",
	// 	Type: "text",
	// 	Data: "hello from the client",
	// 	Metadata: data.Metadata{
	// 		"first":  1,
	// 		"second": "second",
	// 	},
	// }

	// r, err := json.Marshal(d)
	// if err != nil {
	// 	log.Println(err)
	// }

	// req, err = http.NewRequestWithContext(ctx, http.MethodPut, target+update, bytes.NewBuffer(r))
	// if err != nil {
	// 	log.Println(err)
	// }
	// req.AddCookie(cookie)

	// resp, err = http.DefaultClient.Do(req)
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println("update data", resp.StatusCode)

	return nil
}
