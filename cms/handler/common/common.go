package common

import (
	"html/template"

	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
	Env        string
	Config     *viper.Viper
	Logger     *logrus.Entry
	Assets     *hashfs.FS
	Decoder    *schema.Decoder
	Cookies    *sessions.CookieStore
	Templates  *template.Template
}

const (
	SessionCookieName          = "office-management-session"
	SessionCookieState         = "state"
	AuthCodeURL                = "somerandomstring"
	GenericErrMsg              = "Please contact the administrator."
	SessionCookieLoginRedirect = "loginRedirect"
	SessionCookieToken         = "token"

	SessionUserID       = "user-id"
	SessionEmail        = "email"
	SessionProfileImage = "profile-image"
	SessionFirstName    = "first-name"
	SessionLastName     = "last-name"
	SessionEmpFirstName = "emp-first-name"
	SessionEmpLastName  = "emp-last-name"
	SessionDesignation  = "designation"
	SessionRoleID       = "roleID"

	SUPERADMIN = "Super Admin"
)

const (
	HomePath  = "/"
	ErrorPath = "/error"

	SignInPath        = "/sign-in"
	LoginInPath       = "/log-in"
	LogoutPath        = "/logout"

	AdminDashboardPath    = "/admin/dashboard"
)
