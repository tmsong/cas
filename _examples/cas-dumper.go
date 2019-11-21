package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/golang/glog"

	"gopkg.in/cas.v2"
)

type myHandler struct{}

var MyHandler = &myHandler{}

func init() {
	//flag.StringVar(&casURL, "url", "", "CAS server LoginURL")
}

func main() {
	flag.Parse()

	glog.Info("Starting up")

	m := http.NewServeMux()
	m.Handle("/", MyHandler)

	casLoginURL, _ := url.Parse("https://sso.huowanggame.top")
	casBaseURL, _ := url.Parse("https://sso.huowanggame.top/cas")
	casBackURL, _ := url.Parse("http://127.0.0.1:8080/cas/c_back")
	openURL, _ := url.Parse("https://sso.huowanggame.top")
	client := cas.NewClient(&cas.Options{
		LoginURL:       casLoginURL,
		BaseUrl:        casBaseURL,
		ValidationType: "CAS3",
		AppKey:         "sso",
		AppId:          10009,
		ClientHost:     casBackURL,
		OpenUrl:        openURL,
		Cookie: &http.Cookie{
			MaxAge:   86400,
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
		},
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: client.Handle(m),
	}

	if err := server.ListenAndServe(); err != nil {
		glog.Infof("Error from HTTP Server: %v", err)
	}

	glog.Info("Shutting down")
}

type templateBinding struct {
	Username   string
	Attributes cas.UserAttributes
}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !cas.IsAuthenticated(r) {
		cas.RedirectToLogin(w, r)
		return
	}
	if !cas.HasPermission(w, r) {
		fmt.Println("no permission")
		return
	}

	cas.UserInfo(w, r)
	cas.PermissionList(w, r,53)
	cas.RoleList(w, r)

	if r.URL.Path == "/logout" {
		cas.RedirectToLogout(w, r)
		return
	}

	w.Header().Add("Content-Type", "text/html")

	tmpl, err := template.New("index.html").Parse(index_html)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, error_500, err)
		return
	}

	binding := &templateBinding{
		Username:   cas.Username(r),
		Attributes: cas.Attributes(r),
	}

	html := new(bytes.Buffer)
	if err := tmpl.Execute(html, binding); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, error_500, err)
		return
	}

	html.WriteTo(w)
}

const index_html = `<!DOCTYPE html>
<html>
  <head>
    <title>Welcome {{.Username}}</title>
  </head>
  <body>
    <h1>Welcome {{.Username}} <a href="/logout">Logout</a></h1>
    <p>Your attributes are:</p>
    <ul>{{range $key, $values := .Attributes}}
      <li>{{$len := len $values}}{{$key}}:{{if gt $len 1}}
        <ul>{{range $values}}
          <li>{{.}}</li>{{end}}
        </ul>
      {{else}} {{index $values 0}}{{end}}</li>{{end}}
    </ul>
  </body>
</html>
`

const error_500 = `<!DOCTYPE html>
<html>
  <head>
    <title>Error 500</title>
  </head>
  <body>
    <h1>Error 500</h1>
    <p>%v</p>
  </body>
</html>
`
