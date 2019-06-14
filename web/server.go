package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"republicofminer-client-go/explorer"
	"republicofminer-client-go/explorer/api"
	"republicofminer-client-go/protocol/converter/apitoprotocol"
	"strconv"

	"github.com/gorilla/mux"
)

func Run() {
	go explorer.Connect()

	router := mux.NewRouter()
	router.UseEncodedPath()
	router.StrictSlash(false)
	router.HandleFunc(`/block/{id}`, handleblock).Methods("GET")
	router.HandleFunc(`/tx/{hash}`, handletx).Methods("GET")

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 3000), router))
}

var HASHLENGTH = len("L1gwhBkBNWOAS048Dv2P+jSmLZxymCaogpvVSrfTrZY=")

func handleblock(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, ok := params["id"]
	if !ok {
		http.Error(writer, "Error getting the block id", http.StatusInternalServerError)
		return
	}

	var ledger *api.Ledger
	height, err := strconv.ParseInt(id, 10, 64)
	if err == nil {
		ledger = explorer.GetLedgerByHeight(height)
	} else if hash, err := url.QueryUnescape(id); err == nil && len(hash) == HASHLENGTH {
		ledger = explorer.GetLedgerByHash(hash)
	} else {
		http.Error(writer, "Error parsing the block id", http.StatusInternalServerError)
		return
	}

	encoded, _ := json.Marshal(ledger)
	// fmt.Println(encoded)
	writer.Write(encoded)
}

func handletx(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	hash, ok := params["hash"]
	if !ok {
		http.Error(writer, "Error getting the transaction hash", http.StatusInternalServerError)
		return
	}

	hash, err := url.QueryUnescape(hash)

	if err != nil {
		http.Error(writer, "Error parsing the transaction hash", http.StatusInternalServerError)
	}

	tx := explorer.GetTransaction(hash)

	// verify the transaction hash
	t := apitoprotocol.ToTransaction(tx)
	if t.Hash().ToBase64() != tx.Hash {
		log.Println("The hash of the transaction does not match")
	}

	encoded, _ := json.Marshal(tx)
	// fmt.Println(encoded)
	writer.Write(encoded)
}
