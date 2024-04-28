package main

import (
	"log"
	"net/http"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/server"
)

func main() {
	// TODO:

	// test my Serialise and Deserialise more? (esp error paths!)

	// read array of BulkStrings: The first (and sometimes also the second) bulk string in the array is the command's name. Subsequent elements of the array are the arguments for the command.
	// route this command through a router!

	// need some form of interface Data, implements Execute, which takes in the commands (if it is not a valid data type: throw an error)
	// find some way to propagate errors

	// protocol description: https://redis.io/docs/latest/develop/reference/protocol-spec/#resp-protocol-description

	router := server.New("6379") // TODO: take port as flag
	err := router.Serve()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("server serve error: %v", err)
	}
}
