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
  S string
}

type C1ActDetail struct {
  S string
}

func GetNomRoll(client *http.Client, session []*http.Cookie, actID []uint) (map[uint][]C1Person, os.Error) {
  errChans := make([]chan *errorType, len(actID))
  resChans := make([]chan *c1NomRollResType, len(actID))

  for i,id := range actID {
    errChans[i] = make(chan *errorType)
    resChans[i] = make(chan *c1NomRollResType)
    go getNomRoll(client, session, id, resChans[i], errChans[i])
  }

  res := make(map[uint][]C1Person)
  for i,_ := range resChans {
    select {
      case r := <-resChans[i]:
        res[r.id] = r.roll
      case e := <-errChans[i]:
        return nil, e.err
    }
  }
  return res, nil
}

func GetActDetail(client *http.Client, session []*http.Cookie, actID []uint) (map[uint]C1ActDetail, os.Error) {
  errChans := make([]chan *errorType, len(actID))
  resChans := make([]chan *c1DetailType, len(actID))

  for i,id := range actID {
    errChans[i] = make(chan *errorType)
    resChans[i] = make(chan *c1DetailType)
    go getDetail(client, session, id, resChans[i], errChans[i])
  }

  res := make(map[uint]C1ActDetail)
  for i,_ := range resChans {
    select {
      case r := <-resChans[i]:
        res[r.id] = r.detail
      case e := <-errChans[i]:
        return nil, e.err
    }
  }
  return res, nil
}

//TODO - use correct return type
func GetPersDOB(client *http.Client, session []*http.Cookie, serviceNo []string) (int, os.Error)

func C1Login(client *http.Client, usr, pwd string) ([]*http.Cookie, os.Error) {
  param := make(url.Values)
  param.Add("ServiceNo", usr)
  param.Add("Password", pwd)

  resp, err := client.PostForm("https://www.cadetone.aafc.org.au/logon.php", param)
  if err != nil {
    return nil, err
  }
  b := make([]byte, 1e6)
  _,err = io.ReadFull(resp.Body, b)
  resp.Body.Close()

  if err != io.ErrUnexpectedEOF && err != nil {
    return nil, err
  }

  if strings.Contains(string(b), "Please try again") {
    return nil, C1AuthError
  }

  return resp.Cookies(), nil
}

func C1Logout(client *http.Client, session []*http.Cookie) ([]*http.Cookie, os.Error) {
  r,err := http.NewRequest("GET", "https://www.cadetone.aafc.org.au/logout.php", nil)
  if err != nil {
    return nil,err
  }
  injectSession(r,session)
  resp,err := client.Do(r)
  if err != nil {
    return nil,err
  }
  return resp.Cookies(),nil
}

func GetClient(c appengine.Context) *http.Client {
  return urlfetch.Client(c)
}

func injectSession(r cookieable, cs []*http.Cookie) {
  for _,c := range cs {
    r.AddCookie(c)
  }
}

func getNomRoll(client *http.Client, session []*http.Cookie, actID uint, outChan chan *c1NomRollResType, errChan chan *errorType) {
  req,err := http.NewRequest("GET", fmt.Sprintf("https://cadetone.aafc.org.au/activities/nominalroll.php?ActID=%d",actID), nil)
  if err != nil {
    errChan <- &errorType{actID, err}
    return
  }
  injectSession(req,session)
  resp, err := client.Do(req)
  if err != nil {
    errChan <- &errorType{actID, err}
  }
  b := make([]byte, 1e6)
  _,err = io.ReadFull(resp.Body, b)
  resp.Body.Close()

  if err != io.ErrUnexpectedEOF {
    errChan <- &errorType{actID, err}
  }
  res := make([]C1Person, 1)
  res[0].S = string(b)
  outChan <- &c1NomRollResType{actID, res}
}

func getDetail(client *http.Client, session []*http.Cookie, actID uint, outChan chan *c1DetailType, errChan chan *errorType) {
  req,err := http.NewRequest("GET", fmt.Sprintf("https://cadetone.aafc.org.au/activities/viewactivity.php?ActID=%d",actID), nil)
  if err != nil {
    errChan <- &errorType{actID, err}
    return
  }
  injectSession(req,session)
  resp, err := client.Do(req)
  if err != nil {
    errChan <- &errorType{actID, err}
  }
  b := make([]byte, 1e6)
  _,err = io.ReadFull(resp.Body, b)
  resp.Body.Close()

  if err != io.ErrUnexpectedEOF {
    errChan <- &errorType{actID, err}
  }
  outChan <- &c1DetailType{actID, C1ActDetail{string(b)}}
}

type cookieable interface {
  AddCookie(c *http.Cookie)
}

type errorType struct {
  id uint
  err os.Error
}

type c1NomRollResType struct {
  id uint
  roll []C1Person
}

type c1DetailType struct {
  id uint
  detail C1ActDetail
}
