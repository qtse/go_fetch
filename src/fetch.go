package movo

import (
    "appengine"
    "appengine/urlfetch"
    "io"
    "os"
    "url"

    "fmt"
    )

func init() {
}

type C1Person struct {
}

type C1ActDetail struct {
}

///func GetNomRoll(c appengine.Context, actID uint) ([]C1Person, os.Error) {
func GetNomRoll(c appengine.Context, actID uint) (string, os.Error) {
  client := urlfetch.Client(c)
  resp, err := client.Get("https://www.cadetone.aafc.org.au/")
  if err != nil {
///    return nil, err
    return "", err
  }
  b := make([]byte, 1e6)
  n,err := io.ReadFull(resp.Body, b)
  resp.Body.Close()

  if err != io.ErrUnexpectedEOF {
    return fmt.Sprintf("%d",n), err
  }
  return string(b), nil
}

func GetActDetail(c appengine.Context, actID uint) (C1ActDetail, os.Error)

//TODO - use correct return type
func GetPersDOB(c appengine.Context, serviceNo, session string) (int, os.Error)

func C1Login(c appengine.Context, usr, pwd string) (string, os.Error) {
  param := make(url.Values)
  param.Add("ServiceNo", usr)
  param.Add("Password", pwd)

  client := urlfetch.Client(c)
  resp, err := client.PostForm("https://www.cadetone.aafc.org.au/logon.php", param)
  if err != nil {
    return "", err
  }
  b := make([]byte, 1e6)
  n,err := io.ReadFull(resp.Body, b)
  resp.Body.Close()

  if err != io.ErrUnexpectedEOF {
    return fmt.Sprintf("%d",n), err
  }
  return string(b), nil
///  return fmt.Sprint(resp.Cookies()), nil
}
