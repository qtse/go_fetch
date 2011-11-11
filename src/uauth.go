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
  auth := lookupUser(u)
  return auth, true
}

func lookupUser(u *user.User) bool {
  //TODO
  return true
}

func lookupAdmin(u *user.User) bool {
  //TODO
  return true
}
