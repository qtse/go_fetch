package movo

import (
    "appengine"
    "appengine/delay"
    "appengine/user"
///    "fmt"
    "http"
    "os"
    "template"
///    "time"
    )

func init() {
  http.HandleFunc("/c1Login", c1LoginHandler)
  http.HandleFunc("/c1Logout", c1LogoutHandler)
  http.HandleFunc("/", root)
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

  var url string
  var loginText string
  url,_ = user.LogoutURL(c, "/")
  loginText = "Logout"

///  skey, err := C1Login(c, usr, pwd)
///  c.Infof(skey)
///  if err != nil {
///    http.Error(w, err.String(), http.StatusInternalServerError)
///    return
///  }
///  delay.Func("key", fetchActDetail).Call(c, skey, 8525)
///  delay.Func("key", fetchNomRoll).Call(c, skey, 8525)
  delay.Func("key", fetchPerson).Call(c, skey, "8488695")
  s3 := ""

  l := []string{skey, s3}
  p := Page{Text: "", TextList: l, Url: url, LoginText: loginText}
  renderTemplate(w, "template", &p)
}
