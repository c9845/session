/*
Package session handles managing user sessions. This provides some tooling around
gorilla/sessions to simplify use.

This file defines some helper functions for working with typical values stored in
sessions.
*/

package session

import (
	"net/http"
	"strconv"
)

//We define some typical fields stored in sessions with some helper funcs for retrieving
//these fields to reduce code required elsewhere.
const (
	keyUsername  = "username"
	keyUserID    = "user_id"
	keyToken     = "token"
	keySessionID = "session_id"
)

//AddUsername adds the username value to the session using the username key.
func (c *Config) AddUsername(w http.ResponseWriter, r *http.Request, value string) error {
	return c.AddValue(w, r, keyUsername, value)
}

//AddUsername adds the username value to the session using the username key and the default
//package level config.
func AddUsername(w http.ResponseWriter, r *http.Request, value string) error {
	return config.AddUsername(w, r, value)
}

//GetUsername looks up the username key from the session.
func (c *Config) GetUsername(r *http.Request) (value string, err error) {
	return c.GetValue(r, keyUsername)
}

//GetUsername looks up the username key from the session using the default package level config.
func GetUsername(r *http.Request) (value string, err error) {
	return config.GetUsername(r)
}

//----------------------------------------------------------------------------------------------

//AddUserID adds the user ID value to the session using the user ID key. We assume user IDs
//are provided as integers.
func (c *Config) AddUserID(w http.ResponseWriter, r *http.Request, value int64) error {
	return c.AddValue(w, r, keyUserID, strconv.FormatInt(value, 10))
}

//AddUserID adds the user ID value to the session using the user ID key and the default
//package level config. We assume user IDs are provided as integers.
func AddUserID(w http.ResponseWriter, r *http.Request, value int64) error {
	return config.AddUserID(w, r, value)
}

//GetUserID looks up the userID key from the session. We assume user IDs are integers
//and try to convert the value stored in the session accordingly.
func (c *Config) GetUserID(r *http.Request) (value int64, err error) {
	valStr, err := c.GetValue(r, keyUserID)
	if err != nil {
		return
	}

	value, err = strconv.ParseInt(valStr, 10, 64)
	return
}

//GetUserID looks up the userID key from the session using the default package level confing.
func GetUserID(r *http.Request) (value int64, err error) {
	return config.GetUserID(r)
}

//----------------------------------------------------------------------------------------------

//AddToken adds the token value to the session using the token key. We assume user IDs
//are provided as integers.
func (c *Config) AddToken(w http.ResponseWriter, r *http.Request, value string) error {
	return c.AddValue(w, r, keyToken, value)
}

//AddToken adds the token value to the session using the token key and the default
//package level config. We assume user IDs are provided as integers.
func AddToken(w http.ResponseWriter, r *http.Request, value string) error {
	return config.AddToken(w, r, value)
}

//GetToken looks up the token key from the session.
func (c *Config) GetToken(r *http.Request) (value string, err error) {
	return c.GetValue(r, keyToken)
}

//GetToken looks up the token key from the session using the default package level confing.
func GetToken(r *http.Request) (value string, err error) {
	return config.GetToken(r)
}

//----------------------------------------------------------------------------------------------

//AddSessionID adds the session ID value to the session using the session ID key. We assume session IDs
//are provided as integers.
func (c *Config) AddSessionID(w http.ResponseWriter, r *http.Request, value int64) error {
	return c.AddValue(w, r, keySessionID, strconv.FormatInt(value, 10))
}

//AddSessionID adds the session ID value to the session using the session ID key and the default
//package level config. We assume session IDs are provided as integers.
func AddSessionID(w http.ResponseWriter, r *http.Request, value int64) error {
	return config.AddSessionID(w, r, value)
}

//GetSessionID looks up the sessionID key from the session. We assume session IDs are
//integers and try to convert the value stored in the session accordingly.
func (c *Config) GetSessionID(r *http.Request) (value int64, err error) {
	valStr, err := c.GetValue(r, keySessionID)
	if err != nil {
		return
	}

	value, err = strconv.ParseInt(valStr, 10, 64)
	return
}

//GetSessionID looks up the sessionID key from the session using the default package level confing.
func GetSessionID(r *http.Request) (value int64, err error) {
	return config.GetSessionID(r)
}
