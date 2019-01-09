# Amichan package

Simple wrapper for https://github.com/ivahaev/amigo with channels for errors and events

If you need more than just get errors and events from Asterisk, please see original repository, there is much more power

## Installation

```console
go get github.com/kcasctiv/amichan
```

Maybe you will need also install amigo manually:
```console
go get github.com/ivahaev/amigo
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/kcasctiv/amichan"
)

func main() {
	port := 7080
	keepalive := true
	a := amichan.New("username", "password", "localhost", port, keepalive)
	a.Connect()

	for {
		select {
		case err := <-a.Err():
			fmt.Println(err)
		case event := <-a.Event():
			fmt.Println(event.Name())
		}
	}
}
```