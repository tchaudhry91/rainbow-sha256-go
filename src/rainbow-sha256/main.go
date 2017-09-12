package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

var (
	redis_host      string
	redis_connector *redis.Client
)

// Initialize Redis Connector
func init_redis() {
	redis_host = os.Getenv("REDIS_HOST")
	redis_port := "6379"
	redis_connection_string := redis_host + ":" + redis_port
	log.Infof("Attempting Redis Connection to %s", redis_connection_string)
	redis_connector = redis.NewClient(&redis.Options{
		Addr:     redis_connection_string,
		Password: "",
		DB:       0,
	})
	_, err := redis_connector.Ping().Result()
	if err != nil {
		log.Fatal("Could not create redis connection:", err)
	} else {
		log.Info("Redis Connection Established")
	}
}

// Initialize Logging
func init_logging() {
	logLevel := os.Getenv("LOG_LEVEL")
	log.SetOutput(os.Stdout)
	switch strings.ToLower(logLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

// Add a key value pair to the store
func addToStore(key string, value string) {
	err := redis_connector.Set(key, value, 0).Err()
	if err != nil {
		log.Error("Could not save to redis because:", err)
	}
}

// Lookup a value from the key value store
// Returns empty string if not found
func lookupStore(key string) string {
	val, err := redis_connector.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Infof("Key:%s not found in store")
			return ""
		}
	}
	return val
}

// Handler for /hash
// Returns hash for a given string and stores the reverse key value in a key value store
func hashHandler(writer http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	hashString := queryValues.Get("str")
	hashArray := sha256.Sum256([]byte(hashString))
	response_hex := hashArray[:]
	log.Infof("Hashed '%s':'%x'\n", hashString, response_hex)
	response := hex.EncodeToString(response_hex)
	addToStore(response, hashString)
	writer.Write([]byte(response))
}

// Handler for /reverse_hash
// Looks up key value store for already hashed values for a reverse lookup
func reverseHashHandler(writer http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	reverseHashString := queryValues.Get("str")
	response := lookupStore(reverseHashString)
	if response == "" {
		response = "String Not Found!"
	}
	writer.Write([]byte(response))
}

func startServer() {
	http.HandleFunc("/hash", hashHandler)
	http.HandleFunc("/reverse_hash", reverseHashHandler)
	http.HandleFunc("/", hashHandler)
	http.ListenAndServe(":9999", nil)
}

func main() {
	init_logging()
	init_redis()
	startServer()
}
