## Introduction:
This package wraps around the `gorilla/sessions` package to provide some addition tooling used in typical use cases and simplify the implementation of sessions for simple apps.

## Details:
The point of this package is to provide the boilerplate functionality for starting and interacting with sessions. While `gorrilla/sessions` itself provides the underlying tools to work with sessions, this package further simplifies the usage of sessions while also reducing the ability to customize the implementation of your sessions store. This package was designed for a very simple use case: to get started with sessions quickly.

The session data is stored in a cookie. The data is encrypted and hashed to prevent viewing and tampering of the data outside of your server that is managing the sessions. If you use a load balancer in front of your website/app, you *will* need to make sure it is session-aware.

There is no dependency on a filesystem or database; all session information is stored in your website/app's memory and in the browser cookie. This provides pros in that it is very simple to get started; however the cons are that there is no server side validation outside of ensure the cookie's contents haven't changed (i.e.: the cookie expiration hasn't been changed). Furthermore, there is a limit to how much data you would want to store in a cookie and ideally you simply store a session ID in the session and refer back to a database for further information.

## Getting Started:
1) Get a default configuraiton with `NewConfig()` or `DefaultConfig()`. `NewConfig()` requires you to store the configuration elsewhere in your app and pass it around as needed while `DefaultConfig()` stores the configuration, and thus your session store, globally so you can access the session configuration without needing to pass around a variable.
2) Initialize the session store using `Init()`. This will allow sessions to be created and data to be stored.
3) Now you can add and read data from the session based on the request using `AddValue()` and `GetValue()`.

## Typical Use Case:
1) Call `DefaultConfig()`, modify the configuration as needed, and call `Init()` in your `main.go's init()` prior to initializing any HTTP routers.
2) Within your HTTP handler function for handling user logins, after successfully authenticating a user, use `AddValue(w, r, "user_session_id", "2554")` where 2554 references a session in your database table. A cookie will now exist for the user.
3) Within other HTTP handler functions, use `GetValue(r, "user_session_id")` to look up the session ID, and use it to look up your user's data in your database.
4) Typically you would call `Extend(w, r)` on each HTTP endpoint as well to "keep a user logged in". This would most likely be performed in some middleware after validating the session is still active and valid.
