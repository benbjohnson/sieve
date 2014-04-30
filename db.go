package sieve

import (
	"encoding/json"
	"sync"
)

// NewDB returns a new instance of a database.
func NewDB() *DB {
	return &DB{}
}

// DB represents a simple in-memory database.
// It stores data in FIFO order so that data can be retrieved by index.
type DB struct {
	sync.Mutex
	Rows        []*Row
	subscribers []chan *Row
}

// Append adds a row to the database.
func (db *DB) Append(row *Row) {
	db.Lock()
	defer db.Unlock()

	// Set the row index and update
	row.index = len(db.Rows)
	db.Rows = append(db.Rows, row)

	// Notify all subscribers.
	for _, s := range db.subscribers {
		s <- row
	}
}

// Row returns a row at a given index.
func (db *DB) Row(index int) *Row {
	return db.Rows[index]
}

// Length returns the number of rows.
func (db *DB) Length() int {
	return len(db.Rows)
}

// Subscribe returns a channel that pipes all rows since a given index.
func (db *DB) Subscribe(index int) chan *Row {
	var ch = make(chan *Row, 0)
	go func() {
		// Flush the current dataset until we reach the end. Then we
		// will simply append the channel as a subscriber.
		for {
			db.Lock()
			var row *Row
			if index < len(db.Rows) {
				row = db.Rows[index]
				index++
			} else {
				db.subscribers = append(db.subscribers, ch)
				db.Unlock()
				break
			}
			db.Unlock()

			// Send row to channel.
			ch <- row
		}
	}()
	return ch
}

// Unsubscribe removes a channel as a subscriber.
func (db *DB) Unsubscribe(ch chan *Row) {
	db.Lock()
	defer db.Unlock()

	for i, s := range db.subscribers {
		if s == ch {
			db.subscribers = append(db.subscribers[:i], db.subscribers[i+1:]...)
			return
		}
	}

	// NOTE(benbjohnson): Subscribers can be unsubscribed before they are
	// subscribed which can cause subscribers to be attached to the DB
	// indefinitely. This should probably be fixed but this is a tiny little
	// utility and is not meant to handle a bajillion requests so whatever. :-/
}

// Row represents a single element in the database.
type Row struct {
	Data  map[string]interface{}
	index int
}

// Index returns the row index in the database.
func (r *Row) Index() int {
	return r.index
}

// MarshalJSON encodes a row into a JSON structure.
func (r *Row) MarshalJSON() ([]byte, error) {
	// Construct a copy of the data map.
	var data = make(map[string]interface{})
	for k, v := range r.Data {
		data[k] = v
	}

	// Add the index to the map.
	data["index"] = r.index

	return json.Marshal(data)
}
