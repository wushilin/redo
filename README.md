# redo

## Install
```
go get github.com/wushilin/redo
```

## Documentation
$ godoc -http=":16666"

Browse http://localhost:16666

## Usage

```
package main

import (
	"errors"
	"fmt"
	"github.com/wushilin/redo"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Most powerful, redo up to 5 times, and with a backoff
	// function to determine how long to wait before next try
	result, ntries, err := redo.Redo(func() (interface{}, error) {
		return myEasyToFailAddFunc(1, 2)
	}, 5, redo.FixedBackoff(1))
	if err != nil {
		fmt.Println("No luck after", ntries, "tries")
	} else {
		fmt.Println("Result:", result, "after", ntries, "tries. last error is", err)
	}

	// Simply do 3 times (or less if succeeded earlier
	result, ntries, err = redo.Redo3(func() (interface{}, error) {
		return myEasyToFailAddFunc(1, 2)
	})
	if err != nil {
		fmt.Println("No luck after", ntries, "tries")
	} else {
		fmt.Println("Result:", result, "after", ntries, "tries. last error is", err)
	}

	// redo up to 3 times, sleep 1000ms after each try
	result, ntries, err = redo.RedoSleep10(func() (interface{}, error) {
		return myEasyToFailAddFunc(1, 2)
	}, 1000)
	if err != nil {
		fmt.Println("No luck after", ntries, "tries")
	} else {
		fmt.Println("Result:", result, "after", ntries, "tries. last error is", err)
	}

}

func myEasyToFailAddFunc(a, b int) (int, error) {
	// Assumes this function mostly returns a + b,
	// but in rare cases, it returns an error and an unspecified value (e.g. 0) In that case, you want to retry.
	if rand.Int()%10 != 0 {
		return 0, errors.New("Not happy")
	}
	return a + b, nil
}

```