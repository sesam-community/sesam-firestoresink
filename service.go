package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	//Client Cloud Firestore client
	FirestoreClient *firestore.Client
)

func main() {
	var err error

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting service on port %s", port)

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./credentials.json")
	}

	if err != nil {
		log.Fatal(err)
	}

	credentialsFileContent := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_CONTENT")
	if credentialsFileContent != "" {
		err = ioutil.WriteFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), []byte(credentialsFileContent), 0750)
	}

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	projID := os.Getenv("GCP_PROJECT_ID")
	if projID == "" {
		log.Fatal("GCP project ID must be assigned to GCP_PROJECT_ID environment var.")
	}

	conf := &firebase.Config{ProjectID: projID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	FirestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer FirestoreClient.Close()

	router := mux.NewRouter()
	router.HandleFunc("/{collection}", PublishMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func PublishMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s\n", r.Host)
	ctx := context.Background()
	params := mux.Vars(r)
	collectionName := params["collection"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnInternalErrorResponse(err, w)
		return
	}

	var bodyJsonArray []interface{}

	err = json.Unmarshal(body, &bodyJsonArray)

	if err != nil {
		returnInternalErrorResponse(err, w)
		return
	}

	if len(bodyJsonArray) == 0 {
		log.Println("Empty batch, nothing to process, return")
		return
	}
	log.Printf("Got batch of %d entities", len(bodyJsonArray))

	var batch = FirestoreClient.Batch()
	var collection = FirestoreClient.Collection(collectionName)

	for _, item := range bodyJsonArray {
		m := item.(map[string]interface{})
		if m["_deleted"].(bool) {
			continue
		}
		var ref = collection.Doc(string(m["_id"].(string)))
		delete(m, "_id")
		batch.Set(ref, m)
	}
	//store new entities or update existing
	results, err := batch.Commit(ctx)

	if err != nil {
		returnInternalErrorResponse(err, w)
		return
	}
	log.Printf("Batch of size %d processed successfully", len(results))

}

func returnInternalErrorResponse(err error, w http.ResponseWriter) {
	log.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
