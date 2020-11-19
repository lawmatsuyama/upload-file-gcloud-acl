package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type teste struct {
	nome     string
	telefone string
}

func main() {
	var localFile string
	var destinyFile string
	var err error
	var storageClient *storage.Client

	fmt.Print("Enter local filepath: ")
	fmt.Scanf("%s\n", &localFile)
	fmt.Print("Enter destiny filepath: ")
	fmt.Scanf("%s\n", &destinyFile)

	bucket := "test_bucket_law"
	ctx := context.Background()
	data, err := ioutil.ReadFile("key.json")
	if err != nil {
		fmt.Printf("Erro1 %v\n", err.Error())
		return
	}

	conf, err := google.JWTConfigFromJSON(data, compute.DevstorageFullControlScope)
	if err != nil {
		fmt.Printf("Erro2 %v\n", err.Error())
		return
	}

	client := conf.Client(ctx)
	// storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("key.json"))

	storageClient, err = storage.NewClient(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Erro1 %v\n", err.Error())
		return
	}

	f, err := os.Open(localFile)
	if err != nil {
		fmt.Printf("Erro3 %v\n", err.Error())
		return
	}

	defer f.Close()
	obj := storageClient.Bucket(bucket).Object(destinyFile)
	sw := obj.NewWriter(ctx)
	if _, err := io.Copy(sw, f); err != nil {
		fmt.Printf("Erro4 %v\n", err.Error())
		return
	}

	if err := sw.Close(); err != nil {
		fmt.Printf("Erro5 %v\n", err.Error())
		return
	}

	obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
	obj.ACL().Set(ctx, storage.AllAuthenticatedUsers, storage.RoleWriter)

	fmt.Printf("Sucesso!\n")
}
