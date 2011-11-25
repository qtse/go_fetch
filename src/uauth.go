package movo

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "os"
    )

var UAuth struct {BANNED, USER, ADMIN int}

func init() {
  UAuth.BANNED = 0
  UAuth.USER = 1
  UAuth.ADMIN = 2
}

func IsAuth(c appengine.Context) (bool, bool, os.Error) {
  u := user.Current(c)
  if u == nil {
    return false, false, nil
  }
  var auth bool
  var err os.Error
  switch auth,err = lookupUser(u, c); err {
  case datastore.ErrInvalidEntityType:
    fallthrough
  case datastore.ErrNoSuchEntity:
    q := datastore.NewQuery("AppUser").Limit(1)
    switch i,err := q.Count(c); {
    case err != nil:
      return false,true,err
    case i == 0:
      if err := newAdmin(u,c); err != nil {
        return false,true,err
      }
      auth,err = lookupUser(u,c)
      if err != nil {
        return false,true,err
      }
    default:
      if err := newUser(u,c); err != nil {
        return false,true,err
      }
      auth,err = lookupUser(u,c)
      if err != nil {
        return false,true,err
      }
    }
  default:
    return false,true,err
  }
  return auth,true,nil
}

type AppUser struct {
  Access int `datastore:",noindex"`
}

func lookupUser(u *user.User, c appengine.Context) (bool,os.Error) {
  var usr AppUser
  key := datastore.NewKey(c, "AppUser", u.Email,0,nil)
  if err := datastore.Get(c, key, &usr); err != nil {
    return false,err
  }
  //TODO
  return true,nil
}

func lookupAdmin(u *user.User, c appengine.Context) bool {
  //TODO
  return true
}

func newUser(u *user.User, c appengine.Context) os.Error {
  au := AppUser{Access:UAuth.USER}
  return saveAppUser(u.Email,au,c)
}

func newAdmin(u *user.User, c appengine.Context) os.Error {
  au := AppUser{Access:UAuth.ADMIN}
  return saveAppUser(u.Email,au,c)
}

func saveAppUser(email string, au AppUser, c appengine.Context) os.Error {
  key := datastore.NewKey(c, "AppUser", email, 0, nil)
  _,err := datastore.Put(c,key,&au)
  return err
}
