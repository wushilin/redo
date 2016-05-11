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
import "github.com/wushilin/redo"

func myEasyToFailAddFunc(a, b int) (int, error) {
  // Assumes this function mostly returns a + b, but in rare cases, it returns an error and an unspecified value (e.g. 0) In that case, you want to retry.
}

result, ntries, error := redo.Redo(func() (interface{},error) {
  return myEasyToFailAddFunc(1,2)
}, 20, redo.ExponentialBackoff(1,8))

fmt.Println(
