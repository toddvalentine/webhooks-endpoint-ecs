package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secret struct {
	WebhooksSecret string `json:"webhooks_secret"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sn := "webhooks/secret"
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("us-east-1"),
		)
		if err != nil {
			panic(err.Error())
		}

		sm := secretsmanager.NewFromConfig(cfg)
		output, err := sm.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{SecretId: aws.String(sn)})
		if err != nil {
			panic(err.Error())
		}

		var secret Secret
		err = json.Unmarshal([]byte(*output.SecretString), &secret)
		if err != nil {
			panic(err.Error())
		}

		for k, v := range r.Header {
			log.Printf("info: header %s=%s", k, v)
		}

		xhmac := r.Header.Get("X-vtypeio-Hmac-SHA256")
		if xhmac == "" {
			log.Printf("error: no xhmac header")
		} else {
			log.Print(xhmac)
		}

		mac := hmac.New(sha256.New, []byte(secret.WebhooksSecret))
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("request body: %s", string(b))

		mac.Write(b)
		expectedMac := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		if xhmac != expectedMac {
			log.Printf("error: no hmac match. expected: %s received: %s", expectedMac, xhmac)
			return
		}

		log.Printf("success: expected: %s received: %s", expectedMac, xhmac)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	port := 9000
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// log.Printf("listening on port: %d\n", port)
	log.Fatal(s.ListenAndServe())
}
