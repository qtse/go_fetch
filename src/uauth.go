package movo

import (
    "appengine"
    "appengine/user"
    )

func IsAuth(c appengine.Context) (bool, bool) {
  u := user.Current(c)
  if u == nil {
    return false, false
  }
  return true, true
}
