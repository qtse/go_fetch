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
///  client := GetClient(c)
///  session, err := C1Login(client, usr, pwd)
///  if err != nil {
///    http.Error(w, err.String(), http.StatusInternalServerError)
///    return
///  }
///  s := fmt.Sprint(session)
///  skey := c1SessionKey(session)
///  delay.Func("key", fetchActDetail).Call(c, skey, 8525)
///  delay.Func("key", fetchNomRoll).Call(c, skey, 8525)
  delay.Func("key", fetchPerson).Call(c, skey, "8488695")
  s3 := ""

///  s2,err := GetActDetail(client, session, []uint{8524})
///  if err != nil {
///    http.Error(w, err.String(), http.StatusInternalServerError)
///    return
///  }
///  session,err = C1Logout(client, session)
///  if err != nil {
///    http.Error(w, err.String(), http.StatusInternalServerError)
///    return
///  }
///  l := []string{s, s3, fmt.Sprint(s2)}
///  l := []string{s, s3}
  l := []string{skey, s3}
  p := Page{Text: "", TextList: l, Url: url, LoginText: loginText}
  renderTemplate(w, "template", &p)
}
