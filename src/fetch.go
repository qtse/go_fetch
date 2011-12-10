package movo

import (
    "appengine"
    "appengine/memcache"
    "appengine/urlfetch"
    "appengine/user"
    html "fixedhtml"
    "http"
///    "io"
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

///  res := &ActDetail{}
  res,err := RetrieveActDetails(c, actId)
  if err != nil {
    c.Errorf("Fetch Detail Location - actId:%d - %s", actId, err.String())
    return err
  }

///  res,err = parseActDetail(c, &resp.Body, res)
  res,err = parseActDetail(c, res)
  if err != nil {
    c.Errorf("Fetch Details - actId:%d - %s", actId, err.String())
    return err
  }

  res.Persist(c)
  c.Infof(fmt.Sprint(res))

  return err
}

///func parseActDetail(c appengine.Context, r *io.ReadCloser, res *ActDetail) (*ActDetail, os.Error) {
func parseActDetail(c appengine.Context, res *ActDetail) (*ActDetail, os.Error) {
  tmp := strings.NewReader(sampleActHtml)
  r := &tmp
  n,err := html.Parse(*r)
  if err != nil {
    c.Errorf("%s", err.String())
    return nil,err
  }

  nc := NewCursor(n)
  nc = nc.FindById("body1")
  nc.Prune()

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
///  req,err := http.NewRequest("GET", fmt.Sprintf("https://cadetone.aafc.org.au/activities/nominalroll.php?ActID=%d",actId), nil)
///  if err != nil {
///    c.Errorf("Fetch Roll - actId:%d - %s", actId, err.String())
///    return err
///  }
///  session := c1SessionCookie(skey)
///  injectSession(req,session)
///  resp, err := GetClient(c).Do(req)
///  if err != nil {
///    c.Errorf("Fetch Roll - actId:%d - %s", actId, err.String())
///    return err
///  }
///  defer resp.Body.Close()

///  c.Infof("Fetch Roll OK - actId:%d", actId)

///  res,err := parseNomRoll(c, &resp.Body)
  res,err := parseNomRoll(c)
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

///func parseNomRoll(c appengine.Context, r *io.ReadCloser) (*NomRoll, os.Error) {
func parseNomRoll(c appengine.Context) (*NomRoll, os.Error) {
  tmp := strings.NewReader(sampleRollHtml)
  r := &tmp
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

func fetchPerson(c appengine.Context, skey string, sn string) os.Error {
  param := make(url.Values)
  param.Add("LastNametxt", "")
  param.Add("ServiceNotxt", sn)
  param.Add("Search", "Search")
  param.Add("Searchflag", "formsearch")

  req,err := http.NewRequest("POST",
      fmt.Sprintf("https://cadetone.aafc.org.au/searchmember.php?PageRef=memberdetails&amp;Members="),
      strings.NewReader(param.Encode()))
  if err != nil {
    c.Errorf("Search Person - sn:%s - %s", sn, err.String())
    return err
  }
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
///  session := c1SessionCookie(skey)
///  injectSession(req,session)
///  resp, err := GetClient(c).Do(req)
///  if err != nil {
///    c.Errorf("SearchPerson - sn:%s - %s", sn, err.String())
///    return err
///  }
///  defer resp.Body.Close()

///  c.Infof("Search Person OK - sn:%s", sn)

///  link,err := parsePSearch(c, &resp.Body)
  link,err := parsePSearch(c)
  if err != nil {
    c.Errorf("Search Person Error - sn:%s - %s", sn, err.String())
    return err
  }

  if link == "" {
    link = "http://www.cadetone.aafc.org.au/personnel/memberdetails.php?PageRef=searchmember&ID=860"
  }

///resp := &http.Response{}

///  req,err = http.NewRequest("GET", fmt.Sprintf("https://cadetone.aafc.org.au/%s", link), nil)
///  if err != nil {
///    c.Errorf("Fetch Person - link:%s - %s", link, err.String())
///    return err
///  }
///  injectSession(req,session)
///  resp, err = GetClient(c).Do(req)
///  if err != nil {
///    c.Errorf("Fetch Person - link:%s - %s", link, err.String())
///    return err
///  }
///  defer resp.Body.Close()

///  c.Infof("Fetch Person OK - link:%s", link)

///  res,err := parsePerson(c, &resp.Body)
  res,err := parsePerson(c)
  if err != nil {
    c.Errorf("Person Page Error - sn:%s - %s", sn, err.String())
    return err
  }

  c.Infof(fmt.Sprint(res))

  return err
}

///func parsePSearch(c appengine.Context, r *io.ReadCloser) (string, os.Error) {
func parsePSearch(c appengine.Context) (string, os.Error) {
  tmp := strings.NewReader(samplePSearchHtml)
  r := &tmp
  n,err := html.Parse(*r)
  if err != nil {
    c.Errorf("%s", err.String())
    return "",err
  }

  nc := NewCursor(n)
  nc = nc.FindById("main1")
  if !nc.Valid {
    printNode(c, n, 0)
    return "", os.NewError("Can't find main1")
  }
  nc.Prune()

  curr := nc.FindChildElement("table").NextSiblingElement("table")
  curr = curr.FindChildElement("tr").NextSiblingElement("tr")
  curr = curr.FindChildElement("input")
  if !curr.Valid {
    return "", os.NewError("Person not found")
  }

  for _,a := range curr.Attr() {
    if a.Key == "onclick" {
      l := strings.Split(a.Val, "'")
      return l[1], nil
    }
  }
  return "", os.NewError("No onclick attribute")
}

///func parsePerson(c appengine.Context, r *io.ReadCloser) (PersDetail, os.Error) {
func parsePerson(c appengine.Context) (PersDetail, os.Error) {
  tmp := strings.NewReader(samplePersonHtml)
  r := &tmp
  n,err := html.Parse(*r)
  if err != nil {
    c.Errorf("%s", err.String())
    return PersDetail{},err
  }

  nc := NewCursor(n)
  nc = nc.FindById("main1")
  if !nc.Valid {
    printNode(c, n, 0)
    return PersDetail{}, os.NewError("Can't find main1")
  }
  nc.Prune()

  res := PersDetail{}

  curr := nc.FindText("Sex:").FindParentElement("tr")
  res.Sex = curr.Node.Child[1].Child[0].Data[0]

  curr = nc.FindText("Date of Birth:").FindParentElement("tr")
  curr = curr.Child()[1].Child()[0]
  res.Dob,err = time.Parse("02 Jan 2006", curr.Data())
  if err != nil {
    return PersDetail{},err
  }
  res.Dob = time.SecondsToUTC(res.Dob.Seconds())

  curr = nc.FindText("Home Phone:").FindParentElement("tr")
  if !curr.FindText("Private").Valid {
    res.HomePh = curr.Node.Child[1].Child[0].Data
  }

  curr = nc.FindText("Business Phone:").FindParentElement("tr")
  if !curr.FindText("Private").Valid {
    res.WorkPh = curr.Node.Child[1].Child[0].Data
  }

  curr = nc.FindText("Mobile Phone:").FindParentElement("tr")
  if !curr.FindText("Private").Valid {
    res.Mobile = curr.Node.Child[1].Child[0].Data
  }

  curr = nc.FindText("AAFC Email Address:").FindParentElement("tr")
  res.Email = curr.Node.Child[1].Child[0].Data

  return res,nil
}

func C1Login(c appengine.Context, usr, pwd string) (string, os.Error) {
  u := user.Current(c).Email

  itm, err := memcache.Get(c, u + "__c1Sess")
  if err != nil && err != memcache.ErrCacheMiss {
    return "",err
  }
  if err == nil {
    return string(itm.Value),nil
  }

///  param := make(url.Values)
///  param.Add("ServiceNo", usr)
///  param.Add("Password", pwd)

///  client := GetClient(c)
///  resp, err := client.PostForm("https://www.cadetone.aafc.org.au/logon.php", param)
///  if err != nil {
///    return "", err
///  }
///  b := make([]byte, 1e6)
///  _,err = io.ReadFull(resp.Body, b)
///  resp.Body.Close()

///  if err != io.ErrUnexpectedEOF && err != nil {
///    return "", err
///  }

///  if strings.Contains(string(b), "Please try again") {
///    return "", C1AuthError
///  }

///  skey := c1SessionKey(resp.Cookies())
  skey := "12345"
  itm = &memcache.Item{
    Key:   u + "__c1Sess",
    Value: []byte(skey),
    Expiration: 1200,
  }
  if err := memcache.Set(c, itm); err != nil {
    return skey,err
  }

  return skey,nil
}

func C1Logout(c appengine.Context) os.Error {
  u := user.Current(c).Email

  itm, err := memcache.Get(c, u + "__c1Sess")
  if err != nil && err != memcache.ErrCacheMiss {
    return err
  }
  if err == memcache.ErrCacheMiss {
    return nil
  }
  skey := string(itm.Value)
  session := c1SessionCookie(skey)

  r,err := http.NewRequest("GET", "https://www.cadetone.aafc.org.au/logout.php", nil)
  if err != nil {
    return err
  }
  injectSession(r,session)
  client := GetClient(c)
  _,err = client.Do(r)
  if err != nil {
    return err
  }

  memcache.Delete(c, u + "__c1Sess")

  return nil
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
func c1CurrSession(c appengine.Context) (string, os.Error) {
  u := user.Current(c).Email

  itm, err := memcache.Get(c, u + "__c1Sess")
  if err != nil && err != memcache.ErrCacheMiss {
    c.Errorf("memcache Error: " + err.String())
    return "", err
  } else if err == nil {
    c.Infof("Cache hit")
    return string(itm.Value), nil
  }

  c.Infof("Cache miss")
  return "", nil
}

func c1IsLoggedIn(c appengine.Context) (bool, os.Error) {
  s,err := c1CurrSession(c)
  if err != nil {
    return false, err
  } else if s == "" {
    return false, nil
  }

  return true, nil
}

func c1LoginHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  if l,err := c1IsLoggedIn(c); err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  } else if l {
///    http.Redirect(w, r, "/", http.StatusFound)
    http.Redirect(w, r, "/index", http.StatusFound)
    return
  }

  if r.Method == "GET" {
    w.Write([]byte(
          "<html><body><form action='/c1Login' method='POST'>"+
          "<input type='text' name='usr'></input>"+
          "<input type='password' name='pwd'></input>"+
          "<input type='submit' value='Submit'></input>"+
          "</form></body></html>"))
    return
  } else if r.Method == "POST" {
    usr := r.FormValue("usr")
    pwd := r.FormValue("pwd")
    _,err := C1Login(c, usr, pwd)
    if err == C1AuthError {
      http.Error(w, err.String(), http.StatusUnauthorized)
      return
    }
    http.Redirect(w, r, "/index", http.StatusFound)
///    http.Redirect(w, r, "/", http.StatusFound)
  }
}

func c1LogoutHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  C1Logout(c)
///  http.Redirect(w, r, "/", http.StatusFound)
  http.Redirect(w, r, "/index", http.StatusFound)
}

func init() {
}
