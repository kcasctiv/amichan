// Package amichan implements simple wrapper for github.com/ivahaev/amigo
// with channels for errors and events
package amichan

import (
	"errors"
	"strconv"
	"time"

	"github.com/ivahaev/amigo"
)

// Amichan presents struct for interact with Asterisk
type Amichan struct {
	err   chan error
	event chan Event
	amigo *amigo.Amigo
}

// New returns new Amichan example
func New(username, password, host string, port int, keepalive bool) *Amichan {
	s := &amigo.Settings{
		Username:  username,
		Password:  password,
		Host:      host,
		Port:      strconv.Itoa(port),
		Keepalive: keepalive,
	}
	a := amigo.New(s)

	ac := &Amichan{
		amigo: a,
		err:   make(chan error),
		event: make(chan Event, 100),
	}

	a.RegisterDefaultHandler(ac.eventHandler)
	a.On("error", ac.errorHandler)

	return ac
}

// Err returns errors channel
func (ac *Amichan) Err() <-chan error {
	return ac.err
}

// Event returns events channel
func (ac *Amichan) Event() <-chan Event {
	return ac.event
}

// Connect connects with Asterisk
func (ac *Amichan) Connect() {
	ac.amigo.Connect()
}

// Event presents AMI event interface
type Event interface {
	// Field returns event field value by name
	Field(name string) (val string, ok bool)
	// Fields returns all field names in event
	Fields() []string
	// Name returns event name
	Name() string
	// Time returns event time
	Time() time.Time
}

type event map[string]string

func (e event) Field(name string) (string, bool) {
	val, ok := e[name]
	return val, ok
}

func (e event) Fields() []string {
	fields := make([]string, len(e))
	var idx int
	for k := range e {
		fields[idx] = k
	}
	return fields
}

func (e event) Name() string {
	return e["Event"]
}

func (e event) Time() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, e["Time"])
	return t
}

func (ac *Amichan) eventHandler(m map[string]string) {
	ac.event <- event(m)
}

func (ac *Amichan) errorHandler(message string) {
	ac.err <- errors.New(message)
}
