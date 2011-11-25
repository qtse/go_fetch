package movo

import (
    "appengine"
    "appengine/user"
    "http"
    "os"
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
  t := template.Must(template.ParseFile("template/"+tmpl+".html"))
  t.Execute(w, p)
}

func root(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var loggedIn bool
  var err os.Error
  if _,loggedIn,err = IsAuth(c); err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }

  var url string
  var loginText string
  if loggedIn == false {
    url,_ = user.LoginURL(c, "/")
    loginText = "Login"
  } else {
    url,_ = user.LogoutURL(c, "/")
    loginText = "Logout"
  }
  client := GetClient(c)
  s, err := C1Login(client, "usr", "pwd")
  if err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }
  l := []string{s, "def"}
  p := Page{Text: "123", TextList: l, Url: url, LoginText: loginText}
  renderTemplate(w, "template", &p)
}
