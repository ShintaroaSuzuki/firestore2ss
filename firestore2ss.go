package main

import (
	"context"
  "strconv"
	"log"
  "time"
  "os"

  "github.com/joho/godotenv"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
  "google.golang.org/api/option"
  "google.golang.org/api/sheets/v4"
)

func loadEnv() string {
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatal(err)
  }

  credential := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

  return credential
}

func createClient(ctx context.Context) *firestore.Client {
  // Sets your Google Cloud Platform project ID.
  projectID := "enmoku-9e62e"
  credential := option.WithCredentialsFile(loadEnv())

  client, err := firestore.NewClient(ctx, projectID, credential)
  if err != nil {
          log.Fatalf("Failed to create client: %v", err)
  }
  // Close client when done with
  // defer client.Close()
  return client
}

func write2ss(values [][]interface{}) {
  credential := option.WithCredentialsFile(loadEnv())
  spreadsheetID := "1WPu6inVwqPCj3OAKTxDOVtOKP545vYCVt9Y1mLAq9Dk"

  srv, err := sheets.NewService(context.TODO(), credential)
  if err != nil {
    log.Fatal(err)
  }

  writeRange := "シート1!A2:E"+strconv.Itoa(len(values)+1)
  vr := &sheets.ValueRange{
    Values: values,
  }
  _, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, vr).ValueInputOption("RAW").Do()
  if err != nil {
    log.Fatalln(err)
  }
}

func main() {
  ctx := context.Background()
  client := createClient(ctx)
  iter := client.Collection("data").Documents(ctx)
  result := [][]interface{} {}
  for {
    doc, err := iter.Next()
    if err == iterator.Done {
      break
    }
    if err != nil {
      log.Fatalf("Failed to iterate: %v", err)
    }
    docData := doc.Data()
    result = append(result, []interface{} {
      docData["ts"].(time.Time).Format("2006/01/02 15:04:05"),
      docData["ip"],
      docData["action"],
      docData["category"],
      docData["label"],
    })
  }
  write2ss(result)
}
