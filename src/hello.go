package movo

import (
    "appengine"
    "appengine/user"
    "http"
///    "io/ioutil"
    "template"
///    "time"
    )

func init() {
  http.HandleFunc("/", root)
}

type Page struct {
  LoginText string
  Url string
  Text string
  TextList []string
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  t := template.Must(template.ParseFile(tmpl+".html"))
  t.Execute(w, p)
}

func root(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  _, ok := IsAuth(c)

  var url string
  var loginText string
  if ok == false {
    url,_ = user.LoginURL(c, "/")
    loginText = "Login"
  } else {
    url,_ = user.LogoutURL(c, "/")
    loginText = "Logout"
  }
  l := []string{"abc", "def"}
  p := Page{Text: "123", TextList: l, Url: url, LoginText: loginText}
  renderTemplate(w, "template/template", &p)
}
