package movo

import (
    "appengine"
    "appengine/datastore"
///    "appengine/memcache"
    "appengine/urlfetch"
///    "appengine/user"
    "html"
    "http"
    "io"
    "os"
    "strconv"
    "strings"
    "url"

    "fmt"
    )

var (
    C1AuthError = os.NewError("C1 Authentication Error")
    )

type ActDetail struct {
  ActId int
  Type string
  Name string
  Desc string
  Location string
  Start datastore.Time
  End datastore.Time
  NCadets int
  NStaff int
}

func c1SessionCookie(skey string) []*http.Cookie {
  c := &http.Cookie{Name:"PHPSESSID", Value:skey}
  return []*http.Cookie{c}
}

func c1SessionKey(cs []*http.Cookie) string {
  for _,c := range cs {
    if c.Name == "PHPSESSID" {
      return c.Value
    }
  }
  return ""
}

func fetchActDetail(c appengine.Context, skey string, actId int) os.Error {
  req,err := http.NewRequest("GET", fmt.Sprintf("https://cadetone.aafc.org.au/activities/viewactivity.php?ActID=%d",actId), nil)
  if err != nil {
    c.Errorf("Fetch Details - actId:%d - %s", actId, err.String())
    return err
  }
  session := c1SessionCookie(skey)
  injectSession(req,session)
  resp, err := GetClient(c).Do(req)
  if err != nil {
    c.Errorf("Fetch Details - actId:%d - %s", actId, err.String())
    return err
  }
  defer resp.Body.Close()

  c.Infof("Fetch Details OK - actId:%d", actId)

  res,err := parseActDetail(c, &resp.Body)
  res.ActId = actId

  c.Infof(fmt.Sprint(res))

  return err
}

func parseActDetail(c appengine.Context, r *io.ReadCloser) (*ActDetail, os.Error) {
  n,err := html.Parse(*r)
  if err != nil {
    c.Errorf("%s", err.String())
    return nil,err
  }

  nc := NewCursor(n)
  nc = nc.FindById("body1")
  nc.Prune()

  res := &ActDetail{}

  curr := nc.FindText("Activity Type:").Parent().NextSibling().Node.Child[0]
  res.Type = curr.Data

  curr = nc.FindText("Name:").Parent().NextSibling().Node.Child[0]
  res.Name = curr.Data

  curr = nc.FindText("Description:").Parent().NextSibling().Node.Child[0]
  res.Desc = curr.Data

  curr = nc.FindText("Location:").Parent().NextSibling().Node.Child[0]
  res.Location = curr.Data

///  res.Start datastore.Time
///  res.End datastore.Time

  curr = nc.FindText("No of Cadets:").Parent().NextSibling().Node.Child[0]
  res.NCadets,err = strconv.Atoi(curr.Data)
  if err != nil {
    return nil,err
  }

  curr = nc.FindText("No of Staff:").Parent().NextSibling().Node.Child[0]
  res.NStaff,err = strconv.Atoi(curr.Data)
  if err != nil {
    return nil,err
  }

  return res,nil
}

func fetchNomRoll(c appengine.Context, skey string, actId int) {
  req,err := http.NewRequest("GET", fmt.Sprintf("https://cadetone.aafc.org.au/activities/nominalroll.php?ActID=%d",actId), nil)
  if err != nil {
    c.Errorf("Fetch Roll - actId:%d - %s", actId, err.String())
    return
  }
  session := c1SessionCookie(skey)
  injectSession(req,session)
  resp, err := GetClient(c).Do(req)
  if err != nil {
    c.Errorf("Fetch Roll - actId:%d - %s", actId, err.String())
    return
  }
  b := make([]byte, 1e6)
  _,err = io.ReadFull(resp.Body, b)
  resp.Body.Close()

  if err != io.ErrUnexpectedEOF {
    c.Errorf("Fetch Roll - actId:%d - %s", actId, err.String())
    return
  }
  c.Infof("Fetch Roll OK - actId:%d", actId)
}

type C1Person struct {
  S string
}

type C1ActDetail struct {
  S string
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

type cookieable interface {
  AddCookie(c *http.Cookie)
}


func toString(n *html.Node) (res string, skip bool) {
  skip = false
  switch (n.Type) {
    case 0:
      res = "Error"
    case 1:
      res = "Text"
    case 2:
      res = "Doc"
    case 3:
      res = "Elm"
    case 4:
      res = "Cmt"
    case 5:
      res = "DocType"
  }

  res += "(" + strings.TrimSpace(n.Data) + ")"
  return
}

func printNode(c appengine.Context, n *html.Node, ind int) {
  s := ""
  for i := 0; i < ind; i++ {
    s += "  "
  }
  ss, skip := toString(n)
  if !skip {
    s += ss
    if len(strings.TrimSpace(s)) > 0 {
      c.Debugf(s)
    }
    for _,nn := range n.Child {
      printNode(c, nn, ind+1)
    }
  }
}

func init() {
}
