package movo

import (
    "appengine"
    "appengine/urlfetch"
    "http"
    "io"
    "os"
    "strings"
    "url"

    "fmt"
    )

var (
    C1AuthError = os.NewError("C1 Authentication Error")
    )

func init() {
}

type C1Person struct {
}

type C1ActDetail struct {
}

///func GetNomRoll(c appengine.Context, actID uint) ([]C1Person, os.Error) {
func GetNomRoll(client *http.Client, actID uint) (string, os.Error) {
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

func GetActDetail(client *http.Client, actID uint) (C1ActDetail, os.Error)

//TODO - use correct return type
func GetPersDOB(client *http.Client, serviceNo, session string) (int, os.Error)

func C1Login(client *http.Client, usr, pwd string) (string, os.Error) {
  param := make(url.Values)
  param.Add("ServiceNo", usr)
  param.Add("Password", pwd)

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

  if strings.Contains(string(b), "Please try again") {
    return "", C1AuthError
  }

  return string(b), nil
///  return fmt.Sprint(resp.Cookies()), nil
}

func GetClient(c appengine.Context) *http.Client {
  return urlfetch.Client(c)
}
