package rql

import (
	"fmt"
	"io"
	"time"
)

// SecureDB interface defines the set of functions that RQL
// driver expects.
type SecureDB interface {
	Set(key string, data interface{}, expireIn time.Duration)
	Get(key string) (interface{}, bool)
	Authenticate(username string, password string) bool
}

// Driver is the RQL driver which acts as an interface between a database client and
// the querying language parser
//
// Driver takes in a database which conforms to the RQL DB interaface.
// This ensures that the RQL driver isn't tied to a single implementation
// of the database. Any database API that conforms this interface will work
type Driver struct {
	db SecureDB
}

// New function returns a pointer to an instance of RQL driver
func New(db SecureDB) *Driver {
	return &Driver{db}
}

// Operate method can take in any RQL query and perform action
// based on query.
//
// This method doesn't returns anything even if the query is invalid
// instead it will use the io.Writer to write the response to the
// specified stream
func (d *Driver) Operate(src string, w io.Writer) {
	// Parse the src
	ast, err := Parse(src)
	if err != nil {
		errRes(err.Error(), w)
		return
	}
	if ast == nil {
		return
	}

	for _, stmt := range ast.Statements {
		switch stmt.Typ {
		case SetType:
			res(d.set(stmt.SetStatement), w)
		case GetType:
			res(d.get(stmt.GetStatement), w)
		case AuthType:
			res(d.auth(stmt.AuthStatement), w)
		}
	}
}

// set method calls the set method on the database by providing
// appropriate parameters
func (d *Driver) set(stmt *SetStatement) string {
	d.db.Set(stmt.key, stmt.val, convertToDuration(stmt.exp))

	return "Success"
}

// get method calls the get method on the database by providing
// appropriate parameters
// it ignores the "keys" which do not exists in the database and places
// nil in the slice for them
//
// It returns the stringified slice
func (d *Driver) get(stmt *GetStatement) string {
	var res []interface{}

	for _, key := range stmt.keys {
		val, ok := d.db.Get(key)
		if ok {
			res = append(res, val)
		} else {
			res = append(res, nil)
		}
	}

	return stringify(res)
}

func (d *Driver) auth(stmt *AuthStatement) string {
	if d.db.Authenticate(stmt.username, stmt.password) {
		return "Successfully authenticated"
	}

	return "Invalid Credentials"
}

// errRes function is supposed to write error messages to the
// specified stream using io.writer
//
// TODO: Customise the error messages to be more helpful
func errRes(msg string, w io.Writer) {
	w.Write([]byte("ERROR: " + msg + "\n"))
}

// res function is similar to the err method except it is used
// to simply write the resonse generated by driver directly to
// the specified stream
func res(msg string, w io.Writer) {
	w.Write([]byte(msg + "\n"))
}

// ============================ HELPER FUNCTIONS ===================================

// convertToDuration converts uint to time.Duration object.
// This uint is supposed to be in MILLISECONDS.
// It's internally converted into nanoseconds and is then casted into
// time.Duration object
func convertToDuration(t uint) time.Duration {
	return time.Duration(t * 1000)
}

// stringify function can be used to stringify any data type
// It internally uses fmt.Sprintf("%v", ...) to perform the conversion
// which internally uses the String() method on the objects to perform the conversion
func stringify(any interface{}) string {
	return fmt.Sprintf("%v", any)
}
