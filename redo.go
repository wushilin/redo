package redo

import (
  "math"
  "time"
)

type JobThatMayFail func() (interface{}, error)
type BackoffFunc func(n int) int

// Returns a backoff func that sleeps 2^i exponentially with initial sleep of initial delay seconds
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

// Returns a backoff func that sleeps fixed backoff seconds (by parameter seconds)
func FixedBackoff(seconds int) BackoffFunc {
  return func(i int) int {
    return seconds * 1000
  }
}


// retry a job that may fail, up to nRetries, with backoff policy specified by BackoffFunc. Returns the job's result, nth try (0 if first try succeeded) and an error if after nRetries it still failed. The returned error will be the last error
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

