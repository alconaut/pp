// Package models provides functions and methods for accessing
// data extracted from the remote API and types for its representation.
package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	// remoteAPI is the address of a service that
	// returns information about purchases, users, and products.
	remoteAPI string
)

type data map[string]interface{}

// get does a request to the specified URI, makes sure the status
// code is 200 OK, and returns the response body or a non-nil error.
func get(uri string) ([]byte, error) {
	// Do a request and make sure it is successfull.
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Make sure the status code is 200 OK.
	if sc := res.StatusCode; sc != http.StatusOK {
		return nil, fmt.Errorf(`%s: unexpected status code "%d"`, uri, sc)
	}

	// Try to read the body of the response.
	// NB: We assume here that the purchases API microservice is trusted
	// and thus wouldn't intentionally return an insanely large response.
	// Otherwise, it would be a good idea to limit the number of bytes
	// we read.
	return ioutil.ReadAll(res.Body)
}

// objectFromURN gets a URN of an API, makes a GET request to get the response,
// parses it, and unmarshals the result.
// If something goes wrong, an error is returned.
func objectFromURN(urn string, obj interface{}) error {
	// Do a GET request to the remote server.
	uri := remoteAPI + urn
	res, err := get(uri)
	if err != nil {
		return err
	}

	// Make sure the result is not "{}". This is to fix behaviour of
	// daw-purchases service that doesn't indicate non-existent user/product
	// with a status code. But instead returns an empty json.
	// Ideally, that should be fixed.
	// NB: Theoretically it could also return the "{}" with spaces inside
	// the brackets or outside. So, this check can be not enough.
	if string(res) == "{}" {
		return fmt.Errorf("%s: empty response, requested object was not found", uri)
	}

	// Try to unmarshal the body of the response.
	return json.Unmarshal(res, &obj)
}

// Init is a functions that initializes the models.
func Init(api string) {
	remoteAPI = api
}
