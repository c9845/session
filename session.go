/*
Package session provides tooling for storing of user session data in cookies
for handling user sessions. This package provides some boilerplate around the
gorilla/sessions package to provide some common functionality that is typically
reused in web apps.

Data stored in a sessions is stored in a cookie. The cookie data is encrypted
and hashed to prevent tampering and viewing of the data client side. This data
can be read, altered, and added to as needed on the server side using this
package. While gorilla/sessions allows for alternative "stores", ex.: storing
sessions on disk, we only support cookies since this is typically how sessions
are handled.

To use, you will need to initialize your session store using NewConfig() or
DefaultConfig() and then call Init(). Once this has been done, you can get
a session for a request using GetSession(), add data to the session using
AddValue(), and read data from the session in subsequent requests using GetValue().

We allow storing the session config globally using the package-level config
variable, or you can store your Config elsewhere and pass it to requests. Most
use cases will store the config in the package-level config variable just for
ease of use.
*/
package session

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

//Config is the set of configuration settings for working with templates.
type Config struct {
	//Domain is the domain to serve the cookie under. The default is ".".
	Domain string

	//Path is the path off the domain to serve the cookie under. The default is "/"
	//so that the cookie is served on any path for the domain.
	Path string

	//MaxAge is the time until the session cookie will expire.
	MaxAge time.Duration

	//HTTPOnly stops client side scripts form having access to the cookie. The default
	//value is true. There really is not need for client side scripts to access the cookie
	//since it will be encrypted anyway.
	HTTPOnly bool

	//Secure sets the cookie that stores the session data to only be served over HTTPS.
	//The default value is false since we want to support HTTP requests as well.
	Secure bool

	//SameSite sets the SameSite value for the cookie to reduce leaking information during
	//requests. This is a privacy setting. The default is http.SameSiteStrictMode.
	SameSite http.SameSite

	//CookieName is the name of the cookie used for storing session data. The default is
	//"session_cookie".
	CookieName string

	//AuthKey is a 64 character long string used for authenticating the cookie stored value.
	//If this is not provided, a random value is assigned upon app start up.
	AuthKey string

	//EncryptKey is a 32 character long string used for encrypting the cookie stored value. If
	//this is not provided, a random value is assigned upon app start up. This makes the cookie
	//stored value unusable by anthing (i.e: client side scripts) other than your app.
	EncryptKey string

	//store stores the session data
	store *sessions.CookieStore
}

//defaults
const (
	defaultDomain     = "."
	defaultPath       = "/"
	defaultMaxAge     = 1 * time.Hour
	defaultHTTPOnly   = true
	defaultSecure     = false
	defaultSameSite   = http.SameSiteStrictMode
	defaultCookieName = "session"

	authKeyLength    = 64
	encryptKeyLength = 32
)

//errors
var (
	//ErrAuthKeyWrongSize is returned when user provided an AuthKey value that isn't 64 characters.
	ErrAuthKeyWrongSize = errors.New("session: auth key is invalid, must be exactly 64 characters")

	//ErrEncyptKeyWrongSize is returned when user provided an EncryptKey value that isn't 32 characters.
	ErrEncyptKeyWrongSize = errors.New("session: encrypt key is invalid, must be exactly 32 characters")

	//ErrLifetimeTooShort is returned when user provided a MaxAge value less that 1 second.
	ErrMaxAgeTooShort = errors.New("session: max age is invalid, must be greater than 1 second")

	//ErrKeyNotFound is returned when a desired key is not found in the session.
	ErrKeyNotFound = errors.New("session: key not found in session data")
)

//config is the package level saved config. This stores your config when you want to store
//it for global use. It is populated when you use one of the Default...Config() funcs.
var config Config

//NewConfig returns a config for managing your session setup with some defaults set.
func NewConfig() *Config {
	return &Config{
		Domain:     defaultDomain,
		Path:       defaultPath,
		MaxAge:     defaultMaxAge,
		HTTPOnly:   defaultHTTPOnly,
		Secure:     defaultSecure,
		SameSite:   defaultSameSite,
		CookieName: defaultCookieName,
	}
}

//DefaultConfig initializes the package level config with some defaults set. This wraps
//NewConfig() and saves the config to the package.
func DefaultConfig() {
	cfg := NewConfig()
	config = *cfg
}

//validate handles validation of a provided config.
func (c *Config) validate() (err error) {
	if strings.TrimSpace(c.Domain) == "" {
		c.Domain = defaultDomain
	}

	if strings.TrimSpace(c.Path) == "" {
		c.Path = defaultPath
	}

	if c.MaxAge < 1*time.Second {
		return ErrMaxAgeTooShort
	}

	//min and max taken from http\cookie from standard lib.
	if c.SameSite < 1 || c.SameSite > 4 {
		c.SameSite = defaultSameSite
	}

	//if auth and encrypt keys were not provided, generate values
	//switch is just cleaner than if/elseif/else in.
	switch len(c.AuthKey) {
	case 0:
		c.AuthKey = string(securecookie.GenerateRandomKey(authKeyLength))
	case authKeyLength:
	default:
		return ErrAuthKeyWrongSize
	}

	switch len(c.EncryptKey) {
	case 0:
		c.EncryptKey = string(securecookie.GenerateRandomKey(encryptKeyLength))
	case encryptKeyLength:
	default:
		return ErrEncyptKeyWrongSize
	}

	return
}

//getOptions returns the options for setting up the session store. This is a helper func
//to clean up code in Init() and Extend().
func (c *Config) getOptions() *sessions.Options {
	return &sessions.Options{
		Domain:   c.Domain,
		Path:     c.Path,
		MaxAge:   int(c.MaxAge.Seconds()),
		HttpOnly: c.HTTPOnly,
		Secure:   c.Secure,
		SameSite: c.SameSite,
	}
}

//Init initializes the session store for the given config.
func (c *Config) Init() (err error) {
	//validate the config
	err = c.validate()
	if err != nil {
		return
	}

	//initialize the session
	c.store = sessions.NewCookieStore(
		[]byte(c.AuthKey),
		[]byte(c.EncryptKey),
	)
	c.store.Options = c.getOptions()
	return
}

//Init initializes the session using the defaul package level config.
func Init() (err error) {
	return config.Init()
}

//GetConfig returns the current state of the package level config.
func GetConfig() (c *Config) {
	return &config
}

//GetSession returns an existing session for a request or a new session if none existed. The
//field IsNew of the returned sessions.Session will be true if session was just created.
func (c *Config) GetSession(r *http.Request) (*sessions.Session, error) {
	return c.store.Get(r, c.CookieName)
}

//GetSession returns the session using the default package level config.
func GetSession(r *http.Request) (*sessions.Session, error) {
	return config.GetSession(r)
}

//Destroy delete a session for a request. This is typically used when you log a user out.
func (c *Config) Destroy(w http.ResponseWriter, r *http.Request) (err error) {
	s, err := c.GetSession(r)
	if err != nil {
		return
	}

	s.Options = c.getOptions()
	s.Options.MaxAge = -1 //setting MaxAge to a negative value marks it as expired immediately

	err = s.Save(r, w)
	return
}

//Destroy deletes a session using the default package level config.
func Destroy(w http.ResponseWriter, r *http.Request) (err error) {
	return config.Destroy(w, r)
}

//Extend extends the expiration of a session and cookie. This is typically used for keeping
//a used logged in by reseting the expiration each time a user visits a page.
func (c *Config) Extend(w http.ResponseWriter, r *http.Request) (err error) {
	s, err := c.GetSession(r)
	if err != nil {
		return
	}

	//each time we get the options, the new expiration date of the cookie is calculated
	//from the MaxAge.
	s.Options = c.getOptions()

	err = s.Save(r, w)
	return
}

//Extend handles expiration for sessions using the package level config.
func Extend(w http.ResponseWriter, r *http.Request) (err error) {
	return config.Extend(w, r)
}

//AddValue adds a key-value pair to a session.
func (c *Config) AddValue(w http.ResponseWriter, r *http.Request, key, value string) (err error) {
	s, err := c.GetSession(r)
	if err != nil {
		return
	}

	s.Values[key] = value

	err = s.Save(r, w)
	return
}

//AddValue adds a key-value pair to a session using the default package level config.
func AddValue(w http.ResponseWriter, r *http.Request, key, value string) (err error) {
	return config.AddValue(w, r, key, value)
}

//GetValue retrieves the value stored for a key in the session.
func (c *Config) GetValue(r *http.Request, key string) (value string, err error) {
	s, err := c.GetSession(r)
	if err != nil {
		return
	}

	value, exists := s.Values[key].(string)
	if !exists {
		return "", ErrKeyNotFound
	}

	return
}

//GetValue retrieves a value for a key in the session using the default package level config.
func GetValue(r *http.Request, key string) (value string, err error) {
	return config.GetValue(r, key)
}

//GetAllValues retrieves all key value pairs stored in the session.
func (c *Config) GetAllValues(r *http.Request) (kv map[string]string, err error) {
	s, err := c.GetSession(r)
	if err != nil {
		return
	}

	//convert the keys and values to strings since that is the type we use when adding values
	//to the session and the type we use when returning the value for a specific key. just for
	//consistency.
	kv = make(map[string]string)
	for k, v := range s.Values {
		ks := k.(string)
		vs := v.(string)
		kv[ks] = vs
	}

	return
}

//Secure sets the Secure field on the package level config.
func Secure(yes bool) {
	config.Secure = yes
}

//HTTPOnly sets the HTTPOnly field on the package level config.
func HTTPOnly(yes bool) {
	config.HTTPOnly = yes
}

//Domain sets the Domain field on the package level config.
func Domain(domain string) {
	config.Domain = domain
}

//Path sets the Path field on the package level config.
func Path(path string) {
	config.Path = path
}

//MaxAge sets the MaxAge field on the package level config.
func MaxAge(maxAge time.Duration) {
	config.MaxAge = maxAge
}

//Keys sets the AuthKey and EncryptKey fields on the package level config.
func Keys(authKey, encryptkey string) {
	config.AuthKey = authKey
	config.EncryptKey = encryptkey
}

//CookieName sets the CookieName field on the package level config.
func CookieName(cookieName string) {
	config.CookieName = cookieName
}

//SameSite sets the SameSite field on the package level config.
func SameSite(sameSite http.SameSite) {
	config.SameSite = sameSite
}
