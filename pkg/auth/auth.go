package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// ClientCallbackServer is a server that accepts a POST from the auth service
// containing a token.
func ClientCallbackServer() (string, error) {
	l := log.WithFields(log.Fields{
		"app": "ilogin",
		"fn":  "ClientCallbackServer",
	})
	l.Debug("start")
	var token string
	r := mux.NewRouter()
	receivedToken := make(chan string)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		l.Debug("start")
		token = r.FormValue("token")
		if token == "" {
			l.Error("token is empty")
			return
		}
		receivedToken <- token
		l.Debug("end")
	})
	port := os.Getenv("TOKEN_CALLBACK_PORT")
	if port == "" {
		port = "9889"
	}
	go func() {
		l.Debug("start")
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			Debug:            false,
		})
		h := c.Handler(r)
		http.ListenAndServe(":"+port, h)
		l.Debug("end")
	}()
	l.Debug("start")
	for token == "" {
		token = <-receivedToken
	}
	l.Debug("end")
	return token, nil
}

func openBrowser(url string) error {
	l := log.WithFields(log.Fields{
		"app": "ilogin",
		"fn":  "openBrowser",
	})
	l.Debug("start")
	fmt.Fprintf(os.Stderr, "Opening login URL in browser: %s\n", url)
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	l.Debug("end")
	return err
}

// OpenAuthWindow opens a new browser window to the auth server
// with the cookie_name parameter set to cookieName.
func OpenAuthWindow(serviceUrl string, authUrl string, cookieName string) error {
	var err error
	l := log.WithFields(log.Fields{
		"app": "ilogin",
		"fn":  "OpenAuthWindow",
	})
	l.Debug("start")
	if authUrl == "" {
		authUrl = os.Getenv("AUTH_URL")
	}
	if authUrl == "" {
		return errors.New("authUrl is empty")
	}
	if serviceUrl == "" {
		serviceUrl = os.Getenv("SERVICE_URL") + "/auth"
	}
	if serviceUrl == "" {
		return errors.New("serviceUrl is empty")
	}
	serviceUrl = serviceUrl + "/auth"
	if cookieName == "" {
		cookieName = os.Getenv("COOKIE_NAME")
	}
	if cookieName == "" {
		return errors.New("cookieName is empty")
	}
	l.Debug("end")
	// parse auth url, add redirect param and redirect to auth url
	u, err := url.Parse(serviceUrl)
	if err != nil {
		l.Error(err)
		return err
	}
	q := u.Query()
	q.Add("cookie_name", cookieName)
	q.Add("auth_url", authUrl)
	u.RawQuery = q.Encode()
	// open browser
	err = openBrowser(u.String())
	if err != nil {
		l.Error(err)
		return err
	}
	return err
}

func handleRequestAuth(w http.ResponseWriter, r *http.Request) {
	l := log.WithFields(log.Fields{
		"app": "ilogin",
		"fn":  "handleRequestAuth",
	})
	l.Debug("start")
	authUrl := r.FormValue("auth_url")
	cookieName := r.FormValue("cookie_name")
	if authUrl == "" {
		l.Error("auth_url is empty")
		http.Error(w, "auth_url is empty", http.StatusBadRequest)
		return
	}
	l.Debug("end")
	// parse auth url, add redirect param and redirect to auth url
	u, err := url.Parse(authUrl)
	if err != nil {
		l.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q := u.Query()
	ru := os.Getenv("REDIRECT_CALLBACK_URL") + "?cookie_name=" + cookieName
	q.Add("redirect", ru)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	l := log.WithFields(log.Fields{
		"app": "ilogin",
		"fn":  "handleCallback",
	})
	l.Debug("start")
	cookieName := r.FormValue("cookie_name")
	if cookieName == "" {
		l.Error("cookie_name is empty")
		http.Error(w, "cookie_name is empty", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		l.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		l.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var data struct {
		Token string
	}
	data.Token = cookie.Value
	tmpl.Execute(w, data)
}

func TokenServer() error {
	l := log.WithFields(log.Fields{
		"app": "ilogin",
		"fn":  "TokenServer",
	})
	l.Debug("start")
	r := mux.NewRouter()
	r.HandleFunc("/auth", handleRequestAuth)
	r.HandleFunc("/callback", handleCallback)
	l.Debug("start")
	port := os.Getenv("TOKEN_SERVER_PORT")
	if port == "" {
		port = "80"
	}
	return http.ListenAndServe(":"+port, r)
}
