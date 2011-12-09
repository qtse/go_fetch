package movo

import (
    "appengine"
    "appengine/urlfetch"
    html "fixedhtml"
    "http"
    "io"
    "os"
    "strconv"
    "strings"
    "time"
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
  Start *time.Time
  End *time.Time
  NCadets int
  NStaff int
}

type NomRoll struct {
  ActId int
  Roll []*PersDetail
}

type PersDetail struct {
  ServiceNo string
  Rank string
  RankOrder int
  LastName string
  FirstName string
  Unit string
  Sex byte
  Dob *time.Time
  HomePh string
  WorkPh string
  Mobile string
  Email string
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
///  req,err := http.NewRequest("GET", fmt.Sprintf("https://cadetone.aafc.org.au/activities/viewactivity.php?ActID=%d",actId), nil)
///  if err != nil {
///    c.Errorf("Fetch Details - actId:%d - %s", actId, err.String())
///    return err
///  }
///  session := c1SessionCookie(skey)
///  injectSession(req,session)
///  resp, err := GetClient(c).Do(req)
///  if err != nil {
///    c.Errorf("Fetch Details - actId:%d - %s", actId, err.String())
///    return err
///  }
///  defer resp.Body.Close()

///  c.Infof("Fetch Details OK - actId:%d", actId)

///  res,err := parseActDetail(c, &resp.Body)
  res,err := parseActDetail(c)
  if err != nil {
    c.Errorf("Fetch Details - actId:%d - %s", actId, err.String())
    return err
  }
  res.ActId = actId

  c.Infof(fmt.Sprint(res))

  return err
}

///func parseActDetail(c appengine.Context, r *io.ReadCloser) (*ActDetail, os.Error) {
func parseActDetail(c appengine.Context) (*ActDetail, os.Error) {
  tmp := strings.NewReader(doc)
  r := &tmp
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

  curr = nc.FindText("Start Date and Time:").Parent().NextSibling().Node.Child[0]
  res.Start,err = time.Parse("02 Jan 2006\u00a015:04", curr.Data)
  res.Start = time.SecondsToUTC(res.Start.Seconds())
  if err != nil {
    return nil,err
  }

  curr = nc.FindText("Finish Date and Time:").Parent().NextSibling().Node.Child[0]
  res.End,err = time.Parse("02 Jan 2006\u00a015:04", curr.Data)
  res.End = time.SecondsToUTC(res.End.Seconds())
  if err != nil {
    return nil,err
  }

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

func fetchNomRoll(c appengine.Context, skey string, actId int) os.Error {
  req,err := http.NewRequest("GET", fmt.Sprintf("https://cadetone.aafc.org.au/activities/nominalroll.php?ActID=%d",actId), nil)
  if err != nil {
    c.Errorf("Fetch Roll - actId:%d - %s", actId, err.String())
    return err
  }
  session := c1SessionCookie(skey)
  injectSession(req,session)
  resp, err := GetClient(c).Do(req)
  if err != nil {
    c.Errorf("Fetch Roll - actId:%d - %s", actId, err.String())
    return err
  }
  defer resp.Body.Close()

  c.Infof("Fetch Roll OK - actId:%d", actId)

  res,err := parseNomRoll(c, &resp.Body)
///  res,err := parseNomRoll(c)
  if err != nil {
    c.Errorf("Fetch Roll - actId:%d - %s", actId, err.String())
    return err
  }
  res.ActId = actId

  c.Infof(fmt.Sprint(res.ActId))
  for _,r := range res.Roll {
    c.Infof(fmt.Sprint(*r))
  }

  return err
}

func parseNomRoll(c appengine.Context, r *io.ReadCloser) (*NomRoll, os.Error) {
///func parseNomRoll(c appengine.Context) (*NomRoll, os.Error) {
///  tmp := strings.NewReader(doc2)
///  r := &tmp
  n,err := html.Parse(*r)
  if err != nil {
    c.Errorf("%s", err.String())
    return nil,err
  }

  nc := NewCursor(n)
  nc = nc.FindById("body1")
  if !nc.Valid {
    printNode(c, n, 0)
    return nil, os.NewError("Can't find body1")
  }
  nc.Prune()

  res := &NomRoll{}
  res.Roll = make([]*PersDetail, 0)

  curr := nc.FindChildElement("table")
  for curr = curr.NextSiblingElement("table"); curr.Valid;
        curr = curr.NextSiblingElement("table") {
    row := curr.FindChildElement("tr").NextSibling()
    for row = row.NextSibling(); row.Valid; row = row.NextSibling() {
      if len(row.Node.Child) == 1 {
        continue
      }

      children := row.Child()
      var p PersDetail
      p.Rank = children[0].Node.Child[0].Data
      p.RankOrder = RankOrder[p.Rank]
      p.FirstName = children[1].Node.Child[0].Data
      p.LastName = children[2].Node.Child[0].Data
      p.ServiceNo = children[3].Node.Child[0].Data
      p.Unit = children[4].Node.Child[0].Data
      res.Roll = append(res.Roll, &p)
    }
  }

  return res,nil
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

const doc string = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
</html>
`

const doc2 string =
`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
</body>
`

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
