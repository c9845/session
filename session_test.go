package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	if cfg == nil {
		t.Fatal("No config was returned as expected")
		return
	}

	if cfg.Domain != defaultDomain {
		t.Fatal("Default Domain not set as expected")
		return
	}
	if cfg.Path != defaultPath {
		t.Fatal("Default Path not set as expected")
		return
	}
	if cfg.MaxAge != defaultMaxAge {
		t.Fatal("Default MaxAge not set as expected")
		return
	}
	if cfg.HTTPOnly != defaultHTTPOnly {
		t.Fatal("Default HTTPOnly not set as expected")
		return
	}
	if cfg.Secure != defaultSecure {
		t.Fatal("Default Secure not set as expected")
		return
	}
	if cfg.SameSite != defaultSameSite {
		t.Fatal("Default SameSite not set as expected")
		return
	}
	if cfg.CookieName != defaultCookieName {
		t.Fatal("Default CookieName not set as expected")
		return
	}
}

func TestValidate(t *testing.T) {
	//Test Start>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//Provide a known good config.
	cfg := NewConfig()
	err := cfg.validate()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
	//Test End<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	//Test Start>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//Set a blank domain and make sure default was used.
	cfg = NewConfig()
	cfg.Domain = ""
	err = cfg.validate()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
	if cfg.Domain != defaultDomain {
		t.Fatal("Default Domain should have been set but wasnt")
		return
	}
	//Test End<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	//Test Start>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//Set a blank path and make sure default was used.
	cfg = NewConfig()
	cfg.Path = ""
	err = cfg.validate()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
	if cfg.Path != defaultPath {
		t.Fatal("Default Path should have been set but wasnt")
		return
	}
	//Test End<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	//Test Start>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//Make sure a max age must be set.
	cfg = NewConfig()
	cfg.MaxAge = 0
	err = cfg.validate()
	if err == nil {
		t.Fatal("ErrMaxAgeTooShort should have occured but didn't")
		return
	}
	//Test End<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	//Test Start>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//Make sure a default same site is set.
	cfg = NewConfig()
	cfg.SameSite = 0
	err = cfg.validate()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
	if cfg.SameSite != defaultSameSite {
		t.Fatal("Default SameSite should have been set but wasnt")
		return
	}
	//Test End<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	//Test Start>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//Make sure an auth and encrypt key is set if neither is provided.
	cfg = NewConfig()
	err = cfg.validate()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
	if cfg.AuthKey == "" {
		t.Fatal("AuthKey not set as expected")
		return
	}
	if cfg.EncryptKey == "" {
		t.Fatal("EncryptKey not set as expected")
		return
	}
	//Test End<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	//Test Start>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//Check if auth key is incorrect length
	cfg = NewConfig()
	cfg.AuthKey = "too short"
	err = cfg.validate()
	if err != ErrAuthKeyWrongSize {
		t.Fatal("ErrAuthKeyWrongSize should have occured but didnt")
		return
	}

	cfg.AuthKey = "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdf"
	err = cfg.validate()
	if err != nil {
		t.Fatal("Error occued but should not have", err)
		return
	}
	//Test End<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	//Test Start>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//Check if encrypt key is incorrect length
	cfg = NewConfig()
	cfg.EncryptKey = "too short"
	err = cfg.validate()
	if err != ErrEncyptKeyWrongSize {
		t.Fatal("ErrEncyptKeyWrongSize should have occured but didnt")
		return
	}

	cfg.EncryptKey = "asdfasdfasdfasdfasdfasdfasdfasdf"
	err = cfg.validate()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
	//Test End<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
}

func TestGetOptions(t *testing.T) {
	cfg := NewConfig()
	ops := cfg.getOptions()
	if ops.Domain != cfg.Domain {
		t.Fatal("Domain not set in options correctly")
		return
	}

	if ops.Path != cfg.Path {
		t.Fatal("Path not set in options correctly")
		return
	}

	if ops.MaxAge != int(cfg.MaxAge.Seconds()) {
		t.Fatal("MaxAge not set in options correctly")
		return
	}

	if ops.HttpOnly != cfg.HTTPOnly {
		t.Fatal("HTTPOnly not set in options correctly")
		return
	}

	if ops.Secure != cfg.Secure {
		t.Fatal("Secure not set in options correctly")
		return
	}

	if ops.SameSite != cfg.SameSite {
		t.Fatal("SameSite not set in options correctly")
		return
	}
}

func TestInit(t *testing.T) {
	//Test Start>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//Test with something that will fail validation.
	cfg := NewConfig()
	cfg.EncryptKey = "too short"
	err := cfg.Init()
	if err != ErrEncyptKeyWrongSize {
		t.Fatal("ErrEncyptKeyWrongSize should have occured but didnt")
		return
	}
	//Test End<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	//Test Start>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//Test with good config and make sure store was initialized.
	cfg = NewConfig()
	err = cfg.Init()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}

	if cfg.store == nil {
		t.Fatal("store not initialized as expected")
		return
	}
	//Test End<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
}

func TestGetSession(t *testing.T) {
	cfg := NewConfig()
	err := cfg.Init()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}

	req := httptest.NewRequest("GET", "/", nil)
	s, err := cfg.GetSession(req)
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
	if s == nil {
		t.Fatal("No session was returned")
		return
	}
}

func TestDestroy(t *testing.T) {
	cfg := NewConfig()
	err := cfg.Init()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	err = cfg.Destroy(w, req)
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
}

func TestExtend(t *testing.T) {
	cfg := NewConfig()
	err := cfg.Init()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	err = cfg.Extend(w, req)
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
}

func TestAddAndGetValue(t *testing.T) {
	cfg := NewConfig()
	err := cfg.Init()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}

	//add value
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	key := "key"
	value := "value"
	err = cfg.AddValue(w, req, key, value)
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}

	//get value
	getValue, err := cfg.GetValue(req, key)
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
	if getValue != value {
		t.Fatal("value not retrieved")
		return
	}

	//get value for key that doesn't exist
	_, err = cfg.GetValue(req, "wrong key")
	if err != ErrKeyNotFound {
		t.Fatal("ErrKeyNotFound should have occued but didn't", err)
		return
	}
}

func TestGetAllValues(t *testing.T) {
	cfg := NewConfig()
	err := cfg.Init()
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}

	//add value
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	key := "key"
	value := "value"
	err = cfg.AddValue(w, req, key, value)
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}

	//get values
	values, err := cfg.GetAllValues(req)
	if err != nil {
		t.Fatal("Error occured but should not have", err)
		return
	}
	if len(values) != 1 {
		t.Fatal("incorrect list of values returned")
		return
	}
}

func TestDefaultConfig(t *testing.T) {
	DefaultConfig()

	//getting the default config
	c := GetConfig()
	if c == nil {
		t.Fatal("no config returned")
		return
	}

	//modifying the default config
	Secure(true)
	c = GetConfig()
	if !c.Secure {
		t.Fatal("Secure field not set correctly")
		return
	}

	HTTPOnly(true)
	c = GetConfig()
	if !c.HTTPOnly {
		t.Fatal("HTTPOnly field not set correctly")
		return
	}

	Domain("example.com")
	c = GetConfig()
	if c.Domain != "example.com" {
		t.Fatal("Domain field not set correctly")
		return
	}

	Path("/example/path/")
	c = GetConfig()
	if c.Path != "/example/path/" {
		t.Fatal("Path field not set correctly")
		return
	}

	MaxAge(2 * time.Hour)
	c = GetConfig()
	if c.MaxAge != 2*time.Hour {
		t.Fatal("MaxAge field not set correctly")
		return
	}

	auth := "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdf"
	encypt := "asdfasdfasdfasdfasdfasdfasdfasdf"
	Keys(auth, encypt)
	c = GetConfig()
	if c.AuthKey != auth || c.EncryptKey != encypt {
		t.Fatal("Keys field not set correctly")
		return
	}

	CookieName("test")
	c = GetConfig()
	if c.CookieName != "test" {
		t.Fatal("CookieName field not set correctly")
		return
	}

	SameSite(http.SameSiteLaxMode)
	c = GetConfig()
	if c.SameSite != http.SameSiteLaxMode {
		t.Fatal("SameSite field not set correctly")
		return
	}

	//setting the default config
	err := c.Init()
	if err != nil {
		t.Fatal("Could not Init default config.", err)
		return
	}

	//get session, set data, get data
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	_, err = GetSession(req)
	if err != nil {
		t.Fatal("Could not get session", err)
		return
	}

	k := "defaultKey"
	v := "defaultVal"
	err = AddValue(w, req, k, v)
	if err != nil {
		t.Fatal("default value not added", err)
		return
	}

	val, err := GetValue(req, k)
	if err != nil {
		t.Fatal("could not get default value", err)
		return
	}
	if val != v {
		t.Fatal()
	}

	//extend
	err = Extend(w, req)
	if err != nil {
		t.Fatal("could not extend default", err)
		return
	}

	//destroy
	err = Destroy(w, req)
	if err != nil {
		t.Fatal("could not destroy", err)
		return
	}

}
