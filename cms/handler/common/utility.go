package common

import (
	"html/template"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/gorilla/csrf"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

type Authenticator struct {
	BaseURL   string
	LogoutURL string
}

func (s *Server) ParseTemplates() error {
	templates := template.New("cms-templates").Funcs(template.FuncMap{
		"assetHash": func(n string) string {
			return path.Join("/", s.Assets.HashName(strings.TrimPrefix(path.Clean(n), "/")))
		},
		"activeStatus": func(status int32) string {
			if status == 1 {
				return "Active"
			}
			return "Inactive"
		},
		"incrementKey": func(status int) int {
			return status + 1
		},
		"formatDate": func(ts *tspb.Timestamp, layout string) string {
			if !ts.IsValid() {
				return ""
			}
			return ts.AsTime().Format(layout)
		},

		"countPaginate": func(a, b int32) int32 {
			if a > 0 {
				c := a / b
				if a%b != 0 {
					c = c + 1
				}
				return c
			}
			return 0
		},
		"noScape": func(str string) template.HTML {
			if str == "" {
				return template.HTML("<h1>Content not found</h1>")
			}
			return template.HTML(str)
		},
		"nowTime": func() string {
			return time.Now().Format("02 Jan 2006")
		},
		"permissionChecked": func(res string, act string, allPerm map[string][]string) string {
			if val, ok := allPerm[res]; ok {
				for _, v := range val {
					if v == act {
						return "checked"
					}
				}
				return ""
			}
			return ""
		},
		"permission": func(res string) bool {
			return true
		}, "urls": func(url string, params ...string) string {
			for _, v := range params {
				a := strings.Split(v, "_")
				if len(a) == 2 {
					url = strings.Replace(url, "{"+a[0]+"}", a[1], 1)
				}
			}
			return url
		},
	}).Funcs(sprig.FuncMap())

	tmpl, err := templates.ParseFS(s.Assets, "templates/*/*.html")
	if err != nil {
		return err
	}
	s.Templates = tmpl
	return nil
}

func (s *Server) GetErrorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.DoTemplate(w, r, "error.html", http.StatusTemporaryRedirect); err != nil {
			s.Logger.WithError(err).Error("unable to load error template")
		}
	})
}

func IsPartialTemplate(name string) bool {
	return strings.HasSuffix(name, ".part.html")
}

func (s *Server) DoTemplate(w http.ResponseWriter, r *http.Request, name string, status int) error {
	template := s.LookupTemplate(name)
	if template == nil || IsPartialTemplate(name) {
		template, status = s.Templates.Lookup("error.html"), http.StatusNotFound
	}

	w.WriteHeader(status)
	return template.Execute(w, s.TemplateData(r))
}

func (s *Server) LookupTemplate(name string) *template.Template {
	if s.Env == "development" {
		if err := s.ParseTemplates(); err != nil {
			s.Logger.WithError(err).Error("template reload")
			return nil
		}
	}
	return s.Templates.Lookup(name)
}

type TemplateData struct {
	Env       string
	CSRFField template.HTML
	Form      TemplateForm
}

type TemplateForm struct {
	ErrorCode    string
	ErrorDetails string
}

func (s *Server) TemplateData(r *http.Request) TemplateData {
	return TemplateData{
		Env:       s.Env,
		CSRFField: csrf.TemplateField(r),
		Form: TemplateForm{
			ErrorCode:    "500",
			ErrorDetails: "Internal error",
		},
	}
}

func CacheStaticFiles(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if asset is hashed extend cache to 180 days
		e := `"4FROTHS24N"`
		w.Header().Set("Etag", e)
		w.Header().Set("Cache-Control", "max-age=15552000")
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, e) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
