package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

var DB *shrinkflateDb
var Cache *shrinkflateCache

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic recovery", err)
		}
	}()

	err := godotenv.Load()
	if err != nil {
		log.Println("Could not load .env")
		panic(err)
	}

	// create the DB instance
	db, ctx, cancel, err := shrinkflateDb{
		host: os.Getenv("MONGO_HOST"),
		port: os.Getenv("MONGO_PORT"),
		name: os.Getenv("MONGO_DB"),
	}.New()

	DB = db

	if err != nil {
		panic(err)
	}
	defer cancel()
	defer func() {
		if err = DB.conn.Disconnect(ctx); err != nil {
			log.Println("Could not prepare database")
			panic(err)
		}
	}()

	cache, err := shrinkflateCache{
		host:     os.Getenv("REDIS_HOST"),
		port:     os.Getenv("REDIS_PORT"),
		password: os.Getenv("REDIS_PASS"),
	}.New()

	Cache = cache

	err = PrepareQueueHandler()
	if err != nil {
		log.Println("Could not prepare queue")
		panic(err)
	}

	router := configure()

	err = http.ListenAndServe(":4000", router)
	if err != nil {
		log.Println("Error starting server")
	}
}

func configure() *mux.Router {

	router := mux.NewRouter()
	controller := shrinkflateController{}

	router.HandleFunc("/", controller.Welcome).Methods("GET")
	router.HandleFunc("/compress", controller.Compress).Methods("POST")
	router.HandleFunc("/download/{id}", controller.Download).Methods("GET")
	router.PathPrefix("/frontend").Handler(http.StripPrefix("/frontend", http.FileServer(http.Dir("public/build"))))

	router.PathPrefix("/debug/").Handler(http.DefaultServeMux)

	return router
}

func init() {
	_, err := os.Stat("compressed")
	if os.IsNotExist(err) {
		err = os.Mkdir("compressed", 0777)
		if err != nil {
			log.Println("Could not create directory")
			panic(err)
		}
	}
}
