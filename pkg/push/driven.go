package push

// DataDestinationFactory exposes methods to create new datadestinations.
type DataDestinationFactory interface {
	New(url string, schema string) DataDestination
}

// DataDestination to write in the push process.
type DataDestination interface {
	Open(plan Plan, mode Mode, disableConstraints bool) *Error
	Commit() *Error
	RowWriter(table Table) (RowWriter, *Error)
	Close() *Error
}

// RowWriter write row to destination table
type RowWriter interface {
	Write(row Row) *Error
}

type NoErrorCaptureRowWriter struct{}

func (necrw NoErrorCaptureRowWriter) Write(row Row) *Error {
	return &Error{"No error capture configured"}
}

// RowIterator iter over a collection of rows
type RowIterator interface {
	Next() bool
	Value() *Row
	Error() *Error
	Close() *Error
}

// Logger for events.
type Logger interface {
	Trace(msg string)
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

// Nologger default implementation do nothing.
type Nologger struct{}

// Trace event.
func (l Nologger) Trace(msg string) {}

// Debug event.
func (l Nologger) Debug(msg string) {}

// Info event.
func (l Nologger) Info(msg string) {}

// Warn event.
func (l Nologger) Warn(msg string) {}

// Error event.
func (l Nologger) Error(msg string) {}
