package storage

import (
	"log"
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

var StorageClient *storage_go.Client

func InitStorage() {
	StorageClient = storage_go.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("ANON_KEY"), nil)
	log.Println("Storage client initialized")
}
