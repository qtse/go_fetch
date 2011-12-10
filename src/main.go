package movo

import (
    "appengine"
    "appengine/delay"
///    "appengine/user"
///    "bytes"
///    "fmt"
    "http"
    "os"
    "strconv"
    "strings"
    "template"
///    "time"
    )

func init() {
  http.HandleFunc("/c1Login", c1LoginHandler)
  http.HandleFunc("/c1Logout", c1LogoutHandler)
  http.HandleFunc("/index", root)
  http.HandleFunc("/courses/", courseHandler)
}

var (
    delayedFetchActDetail = delay.Func("fetchActDetail", fetchActDetail)
    )

func courseHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  path := strings.Split(r.URL.Path, "/")[1:]
  if path[len(path)-1] == "" {
    path = path[:len(path)-1]
  }

  if len(path) == 1 {
    // All courses
    switch (r.Method) {
      case "GET":
      case "POST":
        actId,err := strconv.Atoi(r.FormValue("actId"))
        if err != nil {
          http.Error(w, "actId not found", http.StatusBadRequest)
          return
        }

        skey, err := c1CurrSession(c)
        if err != nil {
          http.Error(w, err.String(), http.StatusInternalServerError)
          return
        } else if skey == "" {
          http.Error(w, "C1 Session expired", http.StatusUnauthorized)
          return
        }
        a := ActDetail{ActId:actId}
        if err = a.Persist(c); err != nil {
          http.Error(w, err.String(), http.StatusInternalServerError)
          return
        }
///        delay.Func("key", fetchActDetail).Call(c, skey, actId)
        delayedFetchActDetail.Call(c, skey, actId)
///        fetchActDetail(c, skey, actId)
        w.WriteHeader(http.StatusAccepted)
        c.Infof(skey + ": " + strconv.Itoa(actId))
        w.Write([]byte("ActId:" + strconv.Itoa(actId) + "\n"))
      case "PUT":
        fallthrough
      case "DELETE":
        w.Header().Add("Allow", "GET")
        w.Header().Add("Allow", "POST")
        http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
        return
    }
  }
  w.Write([]byte(r.Method + " " + r.URL.String() + " " + strconv.Itoa(len(path)) + " "))
  for _,s := range path {
    w.Write([]byte("/" + s))
  }
}

type Page struct {
  LoginText string
  Url string
  Text string
  TextList []string
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  t := template.Must(template.ParseFile("template/"+tmpl+".html"))
  t.Execute(w, p)
}

func root(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var err os.Error
  if _,_,err = IsAuth(c); err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }
  if l, err := c1IsLoggedIn(c); err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  } else if !l {
    http.Redirect(w, r, "/c1Login", 302)
    return
  }

///  var url string
///  var loginText string
///  url,_ = user.LogoutURL(c, "/")
///  loginText = "Logout"

///  skey, err := C1Login(c, usr, pwd)
///  c.Infof(skey)
///  if err != nil {
///    http.Error(w, err.String(), http.StatusInternalServerError)
///    return
///  }
///  delay.Func("key", fetchActDetail).Call(c, skey, 8525)
///  delay.Func("key", fetchNomRoll).Call(c, skey, 8525)
///  delay.Func("key", fetchPerson).Call(c, skey, "8488695")
///  s3 := ""

///  l := []string{skey, s3}
///  p := Page{Text: "", TextList: l, Url: url, LoginText: loginText}
  p := Page{}
///  renderTemplate(w, "template", &p)
  renderTemplate(w, "tool", &p)
}
