package data

import (
	"encoding/json"
	"io"
	"net/http"
)

// ToJSON serializes the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

func Response(rw http.ResponseWriter, data interface{}, code int) {
	rw.WriteHeader(code)
	ToJSON(data, rw)
}
