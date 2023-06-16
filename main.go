package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/wspirrat/trashkv/core"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

type TemplateData struct {
	ActiveRoute ActiveRoute
	Database    map[string]interface{}
	Servers     map[string]string
}

type ActiveRoute struct {
	Dashboard bool
	Storage   bool
}

func main() {
	core.SAVE_CACHE = true
	core.REPLACE_KEY = true
	http.HandleFunc("/tkv_v1/connect", core.TkvRouteConnect)
	http.HandleFunc("/tkv_v1/save", core.TkvRouteCompareAndSave)
	http.HandleFunc("/tkv_v1/sync", core.TkvRouteSyncWithServers)
	http.HandleFunc("/tkv_v1/status", core.TkvRouteStatus)
	http.PostForm("http://localhost:80/tkv_v1/sync", nil)

	//panel

	fs := http.FileServer(http.Dir("frontend/assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "http://localhost:80/login", http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "http://localhost:80/", http.StatusMovedPermanently)
		}
	})
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout/", logout)
	http.HandleFunc("/dashboard", storage)
	http.HandleFunc("/dashboard/storage", storage)
	http.HandleFunc("/delete/", delete)
	http.HandleFunc("/addact/", addact)
	http.HandleFunc("/dashboard/servers", servers)

	http.ListenAndServe(":80", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/pages/login.html"))
	keys := r.URL.Query()

	session, _ := store.Get(r, "session")

	if r.Method == "POST" {
		_, err := core.Connect("http://localhost:80", r.FormValue("sk"))
		if err == nil {
			session.Values["authenticated"] = true
			session.Values["key"] = r.FormValue("sk")
			session.Save(r, w)

			redurl := keys.Get("redirect")
			if redurl != "" {
				http.Redirect(w, r, fmt.Sprintf("http://localhost:80%s", redurl), http.StatusMovedPermanently)
			} else {
				http.Redirect(w, r, "http://localhost:80/dashboard", http.StatusMovedPermanently)
			}
		} else {
			http.Redirect(w, r, "http://localhost:80/login?status=error", http.StatusMovedPermanently)
		}
	}
	tmpl.Execute(w, nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	session.Values["authenticated"] = false
	session.Values["key"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "http://localhost:80/login", http.StatusMovedPermanently)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "http://localhost:80/login?redirect=/dashboard", http.StatusMovedPermanently)
	} else {
		tmpl, err := template.New("").ParseFiles("frontend/pages/index.html", "frontend/TEMPLATE.html")
		if err != nil {
			fmt.Println(err)
		}
		if err = tmpl.ExecuteTemplate(w, "base", nil); err != nil {
			fmt.Println(err)
		}

		tmpl.Execute(w, nil)
	}
}

func storage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "http://localhost:80/login?redirect=/dashboard/storage", http.StatusMovedPermanently)
	} else {
		tmpl, err := template.New("").ParseFiles("frontend/pages/database.html", "frontend/TEMPLATE.html")
		if err != nil {
			fmt.Println(err)
		}

		db, err := core.Connect("http://localhost:80", session.Values["key"].(string))
		if err != nil {
			fmt.Println(err)
		}

		dataMap := make(map[string]interface{})
		dbaccess := db.Access()
		dbaccess.Range(func(k interface{}, v interface{}) bool {
			dataMap[k.(string)] = v
			return true
		})

		tmplData := TemplateData{
			ActiveRoute: ActiveRoute{
				Dashboard: false,
				Storage:   true,
			},
			Database: dataMap,
		}

		if err = tmpl.ExecuteTemplate(w, "base", tmplData); err != nil {
			fmt.Println(err)
		}
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	keys := r.URL.Query()
	key := keys.Get("key")

	if auth, ok := session.Values["authenticated"].(bool); ok || auth {
		db, err := core.Connect("http://localhost:80", session.Values["key"].(string))
		if err != nil {
			fmt.Println(err)
		}
		if key != "" {
			db.Delete(key)
			db.Save()
			http.Redirect(w, r, "http://localhost:80/dashboard/storage", http.StatusMovedPermanently)
		}
	}
}

func addact(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); ok || auth {
		db, err := core.Connect("http://localhost:80", session.Values["key"].(string))
		if err != nil {
			fmt.Println(err)
			return
		}
		if r.Method == "POST" {
			key := r.FormValue("inpkey")
			value := r.FormValue("inpval")
			if key == "" {
				http.Redirect(w, r, "http://localhost:80/dashboard/storage?status=wrongkey", http.StatusMovedPermanently)
			}

			if r.FormValue("int") != "" {
				newval, err := strconv.Atoi(value)
				if err != nil {
					fmt.Println(err)
				}
				db.Store(key, newval)
			} else if r.FormValue("array") != "" {
				newval := strings.Split(value, ",")
				db.Store(key, newval)
			} else {
				db.Store(key, value)
			}

			db.Save()
			http.Redirect(w, r, "http://localhost:80/dashboard/storage", http.StatusMovedPermanently)
		}

	}
}
