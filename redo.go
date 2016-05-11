package redo

import (
  "math"
  "time"
)

type JobThatMayFail func() (interface{}, error)
type BackoffFunc func(n int) int

func ExponentialBackoff(initialDelaySeconds int, maxSeconds int) BackoffFunc {
  return func(i int) int {
    to_sleep := int(math.Pow(2.0, float64(i))) * 1000 * initialDelaySeconds
    if to_sleep > maxSeconds*1000 {
      return maxSeconds * 1000
    } else {
      return to_sleep
    }
  }
}

func FixedBackoff(seconds int) BackoffFunc {
  return func(i int) int {
    return seconds * 1000
  }
}

func Redo(f JobThatMayFail, nRetries int, backoffFunc BackoffFunc) (interface{}, int, error) {
  var err error
  var i int
  var res interface{}
  for i = 0; i < nRetries; i++ {
    res, err = f()
    if err != nil {
      if backoffFunc != nil {
        backoff := time.Duration(backoffFunc(i)) * time.Millisecond
        time.Sleep(backoff)
      } else {
        // no backoff func provided, no sleeping
      }
    } else {
      return res, i, nil
    }
  }
  return nil, i, err
}

