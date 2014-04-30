package sieve_test

import (
	"encoding/json"
	"testing"

	. "github.com/benbjohnson/sieve"
)

// Ensure that a rows can be appended to the database.
func TestDB_Append(t *testing.T) {
	var db DB
	db.Append(&Row{Data: map[string]interface{}{"foo": "bar", "num": 100}})
	db.Append(&Row{Data: map[string]interface{}{"baz": "bat"}})
	equals(t, 2, db.Length())
	equals(t, "bar", db.Row(0).Data["foo"])
	equals(t, 100, db.Row(0).Data["num"])
	equals(t, "bat", db.Row(1).Data["baz"])
}

// Ensure that a row can be marshalled to JSON.
func TestRow_MarshalJSON(t *testing.T) {
	var db DB
	db.Append(&Row{})
	db.Append(&Row{})
	db.Append(&Row{Data: map[string]interface{}{"foo": 100, "bar": "baz"}})
	b, err := json.Marshal(db.Row(2))
	assert(t, err == nil, "unexpected error: %v", err)
	equals(t, `{"bar":"baz","foo":100,"index":2}`, string(b))
}

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		tb.Fatalf(msg, v...)
	}
}

func equals(tb testing.TB, exp, act interface{}) {
	assert(tb, exp == act, "exp: %#v, got: %#v", exp, act)
}
