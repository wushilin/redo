package redo

import (
  "math"
  "time"
)

type JobThatMayFail func() (interface{}, error)
type BackoffFunc func(n int) int

var BACKOFF_IMMEDIATE = func(i int) int { return 0 }
var BACKOFF_1_SECOND = func(i int) int { return 1000 }

// Returns a backoff func that sleeps 2^i exponentially with initial sleep of initial delay seconds
func ExponentialBackoff(initialMS int, maxMS int) BackoffFunc {
  return func(i int) int {
    to_sleep := int(math.Pow(2.0, float64(i))) * initialMS
    if to_sleep > maxMS {
      return maxMS
    } else {
      return to_sleep
    }
  }
}

// Returns a backoff func that sleeps fixed backoff seconds (by parameter seconds)
func FixedBackoff(milliseconds int) BackoffFunc {
  return func(i int) int {
    return milliseconds
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
        if backoff > 0 {
          time.Sleep(backoff)
        }
      } else {
        // no backoff func provided, no sleeping
      }
    } else {
      return res, i, nil
    }
  }
  return nil, i, err
}

// retry a job that may fail, up to nRetries, and sleep sleepms milliseconds between tries
func RedoSleep(f JobThatMayFail, nRetries int, sleepms int) (interface{}, int, error) {
  return Redo(f, nRetries, FixedBackoff(sleepms))
}

// retry a job that may fail, up to 3 times, without sleep
func Redo3(f JobThatMayFail) (interface{}, int, error) {
  return Redo(f, 3, nil)
}

// retry a job that may fail, up to 5 times, without sleep
func Redo5(f JobThatMayFail) (interface{}, int, error) {
  return Redo(f, 5, nil)
}

// retry a job that may fail, up to 10 times, without sleep
func Redo10(f JobThatMayFail) (interface{}, int, error) {
  return Redo(f, 10, nil)
}

// retry a job that may fail, up to 3 time, sleep sleepms milliseconds betweeen tries
func RedoSleep3(f JobThatMayFail, sleepms int) (interface{}, int, error) {
  return RedoSleep(f, 3, sleepms)
}

// retry a job that may fail, up to 5 time, sleep sleepms milliseconds betweeen tries
func RedoSleep5(f JobThatMayFail, sleepms int) (interface{}, int, error) {
  return RedoSleep(f, 5, sleepms)
}

// retry a job that may fail, up to 10 time, sleep sleepms milliseconds betweeen tries
func RedoSleep10(f JobThatMayFail, sleepms int) (interface{}, int, error) {
  return RedoSleep(f, 10, sleepms)

}

