package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/byblix/mongotest/imgs"
	"github.com/byblix/mongotest/tips"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No env file, getting from server?")
	}
	if err := initEndpoints(); err != nil {
		log.Fatalf("Server failed: %s\n", err)
	}
}

// TODO No CORS yet
// TODO No auth() (JWT/Auth0?) yet
func initEndpoints() error {
	r := mux.NewRouter()
	port := os.Getenv("PORT")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Dont go there, mate")
	}).Methods("GET")

	r.HandleFunc("/create-tip", tips.CreateTip).Methods("POST")
	r.HandleFunc("/get-tip", tips.GetAllTips).Methods("GET")
	r.HandleFunc("/exif-img", imgs.InitImgData).Methods("POST")

	fmt.Printf("Listening on port: %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		return err
	}
	return nil
}
