## Description

Home assistant server and single page application for managing smart devices.

Written in Go, for personal use by @wumbabum.

## Getting started

Before running the application you will need a working PostgreSQL installation and a valid DSN (data source name) for connecting to the database.

Please open the `cmd/web/main.go` file and edit it to include your valid DSN as the default value.

```
cfg.db.dsn = env.GetString("DB_DSN", "YOUR DEFAULT DSN GOES HERE")
```

Note that this DSN must be in the format `user:pass@localhost:port/db` and **not** be prefixed with `postgres://`.

Make sure that you're in the root of the project directory, fetch the dependencies with `go mod tidy`, then run the application using `go run ./cmd/web`:

```
$ go mod tidy
$ go run ./cmd/web
```

Then visit [http://localhost:5749](http://localhost:5749) in your browser.

You can also start the application with live reload support by using the `run` task in the `Makefile`:

```
$ make run
```

## Project structure

Everything in the codebase is designed to be editable. Feel free to change and adapt it to meet your needs.

|     |     |
| --- | --- |
| **`assets`** | Contains the non-code assets for the application. |
| `↳ assets/migrations/` | Contains SQL migrations. |
| `↳ assets/static/` | Contains static UI files (images, CSS etc). |
| `↳ assets/templates/` | Contains HTML templates. |
| `↳ assets/efs.go` | Declares an embedded filesystem containing all the assets. |

|     |     |
| --- | --- |
| **`cmd/web`** | Your application-specific code (handlers, routing, middleware, helpers) for dealing with HTTP requests and responses. |
| `↳ cmd/web/errors.go` | Contains helpers for managing and responding to error conditions. |
| `↳ cmd/web/handlers.go` | Contains your application HTTP handlers. |
| `↳ cmd/web/helpers.go` | Contains helper functions for common tasks. |
| `↳ cmd/web/main.go` | The entry point for the application. Responsible for parsing configuration settings initializing dependencies and running the server. Start here when you're looking through the code. |
| `↳ cmd/web/middleware.go` | Contains your application middleware. |
| `↳ cmd/web/routes.go` | Contains your application route mappings. |
| `↳ cmd/web/server.go` | Contains a helper functions for starting and gracefully shutting down the server. |

|     |     |
| --- | --- |
| **`internal`** | Contains various helper packages used by the application. |
| `↳ internal/database/` | Contains your database-related code (setup, connection and queries). |
| `↳ internal/env` | Contains helper functions for reading configuration settings from environment variables. |
| `↳ internal/funcs/` | Contains custom template functions. |
| `↳ internal/request/` | Contains helper functions for decoding HTML forms, JSON requests, and URL query strings. |
| `↳ internal/response/` | Contains helper functions for rendering HTML templates and sending JSON responses. |
| `↳ internal/validator/` | Contains validation helpers. |
| `↳ internal/version/` | Contains the application version number definition. |

## Configuration settings

Configuration settings are managed via environment variables, with the environment variables read into your application in the `run()` function in the `main.go` file.

You can try this out by setting a `HTTP_PORT` environment variable to configure the network port that the server is listening on:

```
$ export HTTP_PORT="9999"
$ go run ./cmd/web
```

Feel free to adapt the `run()` function to parse additional environment variables and store their values in the `config` struct. The application uses helper functions in the `internal/env` package to parse environment variable values or return a default value if no matching environment variable is set. It includes `env.GetString()`, `env.GetInt()` and `env.GetBool()` functions for reading string, integer and bool values from environment variables. Again, you can add any additional helper functions that you need.

## Creating new handlers

Handlers are defined as `http.HandlerFunc` methods on the `application` struct. They take the pattern:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    // Your handler logic...
}
```

Handlers are defined in the `cmd/web/handlers.go` file. For small applications, it's fine for all handlers to live in this file. For larger applications (10+ handlers) you may wish to break them out into separate files.

## Handler dependencies

Any dependencies that your handlers have should be initialized in the `run()` function `cmd/web/main.go` and added to the `application` struct. All of your handlers, helpers and middleware that are defined as methods on `application` will then have access to them.

You can see an example of this in the `cmd/web/main.go` file where we initialize a new `logger` instance and add it to the `application` struct.

## Creating new routes

[chi](https://github.com/go-chi/chi) version 5 is used for routing. Routes are defined in the `routes()` method in the `cmd/web/routes.go` file. For example:

```
func (app *application) routes() http.Handler {
    mux := chi.NewRouter()

    mux.Get("/your/path", app.yourHandler)

    return mux
}
```

For more information about chi and example usage, please see the [official documentation](https://github.com/go-chi/chi).

## Adding middleware

Middleware is defined as methods on the `application` struct in the `cmd/web/middleware.go` file. Feel free to add your own. They take the pattern:

```
func (app *application) yourMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Your middleware logic...
        next.ServeHTTP(w, r)
    })
}
```

You can then register this middleware with the router using the `Use()` method:

```
func (app *application) routes() http.Handler {
    mux := chi.NewRouter()
    mux.Use(app.yourMiddleware)

    mux.Get("/your/path", app.yourHandler)

    return mux
}
```

It's possible to use middleware on specific routes only by creating route 'groups':

```
func (app *application) routes() http.Handler {
    mux := chi.NewRouter()
    mux.Use(app.yourMiddleware)

    mux.Get("/your/path", app.yourHandler)

    mux.Group(func(mux chi.Router) {
        mux.Use(app.yourOtherMiddleware)

        mux.Get("/your/other/path", app.yourOtherHandler)
    })

    return mux
}
```

Note: Route 'groups' can also be nested.

## Rendering HTML templates

HTML templates are stored in the `assets/templates` directory and use the standard library `html/template` package. The structure looks like this:

|     |     |
| --- | --- |
| `assets/templates/base.tmpl` | The 'base' template containing the shared HTML markup for all your web pages. |
| `assets/templates/pages/` | Directory containing files with the page-specific content for your web pages. See `assets/templates/pages/home.tmpl` for an example. |
| `assets/templates/partials/` | Directory containing files with 'partials' to embed in your web pages or base template. See `assets/templates/partials/footer.tmpl` for an example. |

The HTML for web pages can be sent using the `response.Page()` function. For convenience, an `app.newTemplateData()` method is provided which returns a `map[string]any` map. You can add data to this map and pass it on to your templates.

For example, to render the HTML in a `assets/templates/pages/example.tmpl` file:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData()
    data["hello"] = "world"

    err := response.Page(w, http.StatusOK, data, "pages/example.tmpl")
    if err != nil {
        app.serverError(w, r, err)
    }
}
```

Specific HTTP headers can optionally be sent with the response too:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData()
    data["hello"] = "world"

    headers := make(http.Header)
    headers.Set("X-Server", "Go")

    err := response.PageWithHeaders(w, http.StatusOK, data, headers, "pages/example.tmpl")
    if err != nil {
        app.serverError(w, r, err)
    }
}
```

Note: All the files in the `assets/templates` directory are embedded into your application binary and can be accessed via the `EmbeddedFiles` variable in `assets/efs.go`.

## Adding default template data

If you have data that you want to display or use on multiple web pages, you can adapt the `newTemplateData()` helper in the `helpers.go` file to include this by default. For example, if you wanted to include the current year value you could adapt it like this:

```
func (app *application) newTemplateData() map[string]any {
    data := map[string]any{
        "CurrentYear": time.Now().Year(),
    }

    return data
}
```

## Custom template functions

Custom template functions are defined in `internal/funcs/funcs.go` and are automatically made available to your

HTML templates when you use `response.Page()`
.

The following custom template functions are already included by default:

|     |     |
| --- | --- |
| `now` | Returns the current time. |
| `timeSince arg1` | Returns the time elapsed since arg1. |
| `timeUntil arg2` | Returns the time until arg1. |
| `formatTime arg1 arg2` | Returns the time arg2 as formatted using the pattern arg1. |
| `approxDuration arg1` | Returns the approximate duration of arg1 in a 'human-friendly' format ("3 seconds", "2 months", "5 years") etc. |
| `uppercase arg1` | Returns arg1 converted to uppercase. |
| `lowercase arg1` | Returns arg1 converted to lowercase. |
| `pluralize arg1 arg2 arg3` | If arg1 equals 1 then return arg2, otherwise return arg3. |
| `slugify arg1` | Returns the lowercase of arg1 with all non-ASCII characters and punctuation removed (expect underscores and hyphens). Whitespaces are also replaced with a hyphen. |
| `safeHTML arg1` | Output the verbatim value of arg1 without escaping the content. This should only be used when arg1 is from a trusted source. |
| `join arg1 arg2` | Returns the values in slice arg1 joined using the separator arg2. |
| `incr arg1` | Increments arg1 by 1. |
| `decr arg1` | Decrements arg1 by 1. |
| `formatInt arg1` | Returns arg1 formatted with commas as the thousands separator. |
| `formatFloat arg1 arg2` | Returns arg1 rounded to arg2 decimal places and formatted with commas as the thousands separator. |
| `yesNo arg1` | Returns "Yes" if arg1 is true, or "No" if arg1 is false. |
| `urlSetParam arg1 arg2 arg3` | Returns the URL arg1 with the key arg2 and value arg3 added to the query string parameters. |
| `urlDelParam arg1 arg2` | Returns the URL arg1 with the key arg2 (and corresponding value) removed from the query string parameters. |

To add another custom template function, define the function in `internal/funcs/funcs.go` and add it to the `TemplateFuncs` map. For example:

```
var TemplateFuncs = template.FuncMap{
    ...
    "yourFunction": yourFunction,
}

func yourFunction(s string) (string, error) {
    // Do something...
}
```

## Static files

By default, the files in the `assets/static` directory are served using Go's `http.Fileserver` whenever the application receives a `GET` request with a path beginning `/static/`. So, for example, if the application receives a `GET /static/css/main.css` request it will respond with the contents of the `assets/static/css/main.css` file.

If you want to change or remove this behavior you can by editing the `routes.go` file.

Note: The files in `assets/static` directory are embedded into your application binary and can be accessed via the `EmbeddedFiles` variable in `assets/efs.go`.

## Working with forms

The codebase includes a `request.DecodePostForm()` function for automatically decoding HTML form data from a POST request into a struct, and `request.DecodeQueryString()` for decoding URL query strings into a struct. Behind the scenes decoding is managed using the [go-playground/form](https://github.com/go-playground/form) package.

As an example, let's say you have a page with the following HTML form for creating a 'person' record and routing rule:

```
<form action="/person/create" method="POST">
    <div>
        <label>Your name:</label>
        <input type="text" name="Name" value="{{.Form.Name}}">
    </div>
    <div>
        <label>Your age:</label>
        <input type="number" name="Age" value="{{.Form.Age}}">
    </div>
    <button>Submit</button>
</form>
```

```
func (app *application) routes() http.Handler {
    mux := flow.New()

    mux.HandleFunc("/person/create", app.createPerson, "GET", "POST")

    return mux
}
```

Then you can display and parse this form with a `createPerson` handler like this:

```
package main

import (
    "net/http"

    "github.com/wumbabum/home_assist/internal/request"
    "github.com/wumbabum/home_assist/internal/response"
)

func (app *application) createPerson(w http.ResponseWriter, r *http.Request) {
    type createPersonForm struct {
        Name string `form:"Name"`
        Age  int    `form:"Age"`
    }

    switch r.Method {
    case http.MethodGet:
        data := app.newTemplateData()

        // Add any default values to the form.
        data["Form"] = createPersonForm{
            Age: 21,
        }

        err := response.Page(w, http.StatusOK, data, "/path/to/page.tmpl")
        if err != nil {
            app.serverError(w, r, err)
        }

    case http.MethodPost:
        var form createPersonForm

        err := request.DecodePostForm(r, &form)
        if err != nil {
            app.badRequest(w, r, err)
            return
        }

        // Do something with the data in the form variable...
    }
}
```

## Validating forms

The `internal/validator` package includes a simple (but powerful) `validator.Validator` type that you can use to carry out validation checks.

Extending the example above:

```
package main

import (
    "net/http"

    "github.com/wumbabum/home_assist/internal/request"
    "github.com/wumbabum/home_assist/internal/response"
    "github.com/wumbabum/home_assist/internal/validator"
)

func (app *application) createPerson(w http.ResponseWriter, r *http.Request) {
    type createPersonForm struct {
        Name      string              `form:"Name"`
        Age       int                 `form:"Age"`
        Validator validator.Validator `form:"-"`
    }

    switch r.Method {
    case http.MethodGet:
        data := app.newTemplateData()

        // Add any default values to the form.
        data["Form"] = createPersonForm{
            Age: 21,
        }

        err := response.Page(w, http.StatusOK, data, "/path/to/page.tmpl")
        if err != nil {
            app.serverError(w, r, err)
        }

    case http.MethodPost:
        var form createPersonForm

        err := request.DecodePostForm(r, &form)
        if err != nil {
            app.badRequest(w, r, err)
            return
        }

        form.Validator.CheckField(form.Name != "", "Name", "Name is required")
        form.Validator.CheckField(form.Age != 0, "Age", "Age is required")
        form.Validator.CheckField(form.Age >= 21, "Age", "Age must be 21 or over")

        if form.Validator.HasErrors() {
            data := app.newTemplateData()
            data["Form"] = form

            err := response.Page(w, http.StatusUnprocessableEntity, data, "/path/to/page.tmpl")
            if err != nil {
                app.serverError(w, r, err)
            }
            return
        }

        // Do something with the form information, like adding it to a database...
    }
}
```

And you can display the error messages in your HTML form like this:

```
<form action="/person/create" method="POST">
    {{if .Form.Validator.HasErrors}}
        <p>Something was wrong. Please correct the errors below and try again.</p>
    {{end}}
    <div>
        <label>Your name:</label>
        {{with .Form.Validator.FieldErrors.Name}}
            <span class='error'>{{.}}</span>
        {{end}}
        <input type="text" name="Name" value="{{.Form.Name}}">
    </div>
    <div>
        <label>Your age:</label>
        {{with .Form.Validator.FieldErrors.Age}}
            <span class='error'>{{.}}</span>
        {{end}}
        <input type="number" name="Age" value="{{.Form.Age}}">
    </div>
    <button>Submit</button>
</form>
```

In the example above we use the `CheckField()` method to carry out validation checks for specific fields. You can also use the `Check()` method to carry out a validation check that is _not related to a specific field_. For example:

```
input.Validator.Check(input.Password == input.ConfirmPassword, "Passwords do not match")
```

The `validator.AddError()` and `validator.AddFieldError()` methods also let you add validation errors directly:

```
input.Validator.AddFieldError("Email", "This email address is already taken")
input.Validator.AddError("Passwords do not match")
```

The `internal/validator/helpers.go` file also contains some helper functions to simplify validations that are not simple comparison operations.

|     |     |
| --- | --- |
| `NotBlank(value string)` | Check that the value contains at least one non-whitespace character. |
| `MinRunes(value string, n int)` | Check that the value contains at least n runes. |
| `MaxRunes(value string, n int)` | Check that the value contains no more than n runes. |
| `Between(value, min, max T)` | Check that the value is between the min and max values inclusive. |
| `Matches(value string, rx *regexp.Regexp)` | Check that the value matches a specific regular expression. |
| `In(value T, safelist ...T)` | Check that a value is in a 'safelist' of specific values. |
| `AllIn(values []T, safelist ...T)` | Check that all values in a slice are in a 'safelist' of specific values. |
| `NotIn(value T, blocklist ...T)` | Check that the value is not in a 'blocklist' of specific values. |
| `NoDuplicates(values []T)` | Check that a slice does not contain any duplicate (repeated) values. |
| `IsEmail(value string)` | Check that the value has the formatting of a valid email address. |
| `IsURL(value string)` | Check that the value has the formatting of a valid URL. |

For example, to use the `Between` check your code would look similar to this:

```
input.Validator.CheckField(validator.Between(input.Age, 18, 30), "Age", "Age must between 18 and 30")
```

Feel free to add your own helper functions to the `internal/validator/helpers.go` file as necessary for your application.

## Sending JSON responses

JSON responses and a specific HTTP status code can be sent using the `response.JSON()` function. The `data` parameter can be any JSON-marshalable type.

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"hello":  "world"}

    err := response.JSON(w, http.StatusOK, data)
    if err != nil {
        app.serverError(w, r, err)
    }
}
```

Specific HTTP headers can optionally be sent with the response too:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"hello":  "world"}

    headers := make(http.Header)
    headers.Set("X-Server", "Go")

    err := response.JSONWithHeaders(w, http.StatusOK, data, headers)
    if err != nil {
        app.serverError(w, r, err)
    }
}
```

## Parsing JSON requests

HTTP requests containing a JSON body can be decoded using the `request.DecodeJSON()` function. For example, to decode JSON into an `input` struct:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name string `json:"Name"`
        Age  int    `json:"Age"`
    }

    err := request.DecodeJSON(w, r, &input)
    if err != nil {
        app.badRequest(w, r, err)
        return
    }

    ...
}
```

Note: The target decode destination passed to `request.DecodeJSON()` (which in the example above is `&input`) must be a non-nil pointer.

The `request.DecodeJSON()` function returns friendly, well-formed, error messages that are suitable to be sent directly to the client using the `app.badRequest()` helper.

There is also a `request.DecodeJSONStrict()` function, which works in the same way as `request.DecodeJSON()` except it will return an error if the request contains any JSON fields that do not match a name in the the target decode destination.

## Validating JSON requests

The `internal/validator` package includes a simple (but powerful) `validator.Validator` type that you can use to carry out validation checks.

Extending the example above:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name      string              `json:"Name"`
        Age       int                 `json:"Age"`
        Validator validator.Validator `json:"-"`
    }

    err := request.DecodeJSON(w, r, &input)
    if err != nil {
        app.badRequest(w, r, err)
        return
    }

    input.Validator.CheckField(input.Name != "", "Name", "Name is required")
    input.Validator.CheckField(input.Age != 0, "Age", "Age is required")
    input.Validator.CheckField(input.Age >= 21, "Age", "Age must be 21 or over")

    if input.Validator.HasErrors() {
        app.failedValidation(w, r, input.Validator)
        return
    }

    ...
}
```

The `app.failedValidation()` helper will send a `422` status code along with any validation error messages. For the example above, the JSON response will look like this:

```
{
    "FieldErrors": {
        "Age": "Age must be 21 or over",
        "Name": "Name is required"
    }
}
```

## Working with the database

This codebase is set up to use PostgreSQL with the [lib/pq](https://github.com/lib/pq) driver. You can control which database you connect to using the `DB_DSN` environment variable to pass in a DSN, or by adapting the default value in `run()`.

The codebase is also configured to use [jmoiron/sqlx](https://github.com/jmoiron/sqlx), so you have access to the whole range of sqlx extensions as well as the standard library `Exec()`, `Query()` and `QueryRow()` methods .

The database is available to your handlers, middleware and helpers via the `application` struct. If you want, you can access the database and carry out queries directly. For example:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    _, err := app.db.Exec("INSERT INTO people (name, age) VALUES ($1, $2)", "Alice", 28)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    ...
}
```

Generally though, it's recommended to isolate your database logic in the `internal/database` package and extend the `DB` type to include your own methods. For example, you could create a `internal/database/people.go` file containing code like:

```
type Person struct {
    ID    int    `db:"id"`
    Name  string `db:"name"`
    Age   int    `db:"age"`
}

func (db *DB) NewPerson(name string, age int) error {
    _, err := db.Exec("INSERT INTO people (name, age) VALUES ($1, $2)", name, age)
    return err
}

func (db *DB) GetPerson(id int) (Person, error) {
    var person Person
    err := db.Get(&person, "SELECT * FROM people WHERE id = $1", id)
    return person, err
}
```

And then call this from your handlers:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    err := app.db.NewPerson("Alice", 28)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    ...
}
```

## Managing SQL migrations

The `Makefile` in the project root contains commands to easily create and work with database migrations:

|     |     |
| --- | --- |
| `$ make migrations/new name=add_example_table` | Create a new database migration in the `assets/migrations` folder. |
| `$ make migrations/up` | Apply all up migrations. |
| `$ make migrations/down` | Apply all down migrations. |
| `$ make migrations/goto version=N` | Migrate up or down to a specific migration (where N is the migration version number). |
| `$ make migrations/force version=N` | Force the database to be specific version without running any migrations. |
| `$ make migrations/version` | Display the currently in-use migration version. |

Hint: You can run `$ make help` at any time for a reminder of these commands.

These `Makefile` tasks are simply wrappers around calls to the `github.com/golang-migrate/migrate/v4/cmd/migrate` tool. For more information, please see the [official documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

By default all 'up' migrations are automatically run on application startup using embeded files from the `assets/migrations` directory. You can disable this by setting the `DB_AUTOMIGRATE` environment variable to `false`.

## Logging

Leveled logging is supported using the [slog](https://pkg.go.dev/log/slog) and [tint](https://github.com/lmittmann/tint) packages.

By default, a logger is initialized in the `main()` function. This logger writes all log messages above `Debug` level to `os.Stdout`.

```
logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
```

Feel free to customize this further as necessary.

Also note: Any messages that are automatically logged by the Go `http.Server` are output at the `Warn` level.

## Authentication

The application uses Auth0 for authentication via OAuth2/OIDC, enabling secure single sign-on (SSO) with multiple identity providers (Google, GitHub, email/password, etc.).

### Configuration

Set the following environment variables:
```
AUTH0_DOMAIN='your-tenant.us.auth0.com'
AUTH0_CLIENT_ID='your-client-id'
AUTH0_CLIENT_SECRET='your-client-secret'
AUTH0_CALLBACK_URL='http://localhost:5749/callback'
BASE_URL='http://localhost:5749'
```

### Auth0 Login
- Navigate to `/login` to authenticate
- Protected routes automatically redirect unauthenticated users to login
- Access user profile at `/profile` after authentication
- Logout at `/logout`

The `requireAuth` middleware in `cmd/web/middleware.go` protects routes requiring authentication.

## Using sessions

The codebase is set up so that server-side sessions (using the [SCS](https://github.com/alexedwards/scs) package) work out-of-the-box.

You can use them in your handlers like this:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    // Add the key "foo" and value "bar" to the session.
    app.sessionManager.Put(r.Context(), "foo", "bar")

    // Retrieve the value for the key "baz" from the session.
    baz := app.sessionManager.GetString(r.Context(), "baz")
    ...
}
```

By default sessions are set to expire after 1 week. You can configure this along with other settings in the `cmd/web/main.go` file by changing the `sessionManager` struct values. For example:

```
sessionManager := scs.New()
...
sessionManager.Lifetime = 3 * time.Hour
sessionManager.IdleTimeout = 20 * time.Minute
sessionManager.Cookie.HttpOnly = false
sessionManager.Cookie.Persist = false
sessionManager.Cookie.SameSite = http.SameSiteStrictMode
sessionManager.Cookie.Partitioned = true
```

The session cookie name is randomized on a per-application basis, in order to reduce the risk of cookie conflicts if you are developing multiple applications on the same machine. The default session cookie name takes the form `"session_{8 random characters}"`, and this can be configured to a different value at runtime using the `SESSION_COOKIE_NAME` environment variable

For more information please see the [documentation for the SCS package](https://github.com/alexedwards/scs).

## Admin tasks

The `Makefile` in the project root contains commands to easily run common admin tasks:

|     |     |
| --- | --- |
| `$ make tidy` | Format all code using `go fmt` and tidy the `go.mod` file. |
| `$ make audit` | Run `go vet`, `staticheck`, `govulncheck`, execute all tests and verify required modules. |
| `$ make test` | Run all tests. |
| `$ make test/cover` | Run all tests and outputs a coverage report in HTML format. |
| `$ make build` | Build a binary for the `cmd/web` application and store it in the `/tmp/bin` folder. |
| `$ make run` | Build and then run a binary for the `cmd/web` application. |
| `$ make run/live` | Build and then run a binary for the `cmd/web` application (uses live reloading). |

## Live reload

When you use `make run/live` to run the application, the application will automatically be rebuilt and restarted whenever you make changes to any files with the following extensions:

```
.go
.tpl, .tmpl, .html
.css, .js, .sql
.jpeg, .jpg, .gif, .png, .bmp, .svg, .webp, .ico
```

Behind the scenes the live reload functionality uses the [cosmtrek/air](https://github.com/cosmtrek/air) tool. You can configure how it works (including which file extensions and folders are watched for changes) by editing the `Makefile` file.

## Running background tasks

A `backgroundTask()` helper is included in the `cmd/web/helpers.go` file. You can call this in your handlers, helpers and middleware to run any logic in a separate background goroutine. This useful for things like sending emails, or completing slow-running jobs.

You can call it like so:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    app.backgroundTask(r, func() error {
        // The logic you want to execute in a background task goes here.
        // It should return an error, or nil.
        err := doSomething()
        if err != nil {
            return err
        }

        return nil
    })

    ...
}
```

Using the `backgroundTask()` helper will automatically recover any panics in the background task logic, and when performing a graceful shutdown the application will wait for any background tasks to finish running before it exits.

## Application version

The application version number is defined in a `Get()` function in the `internal/version/version.go` file. Feel free to change this as necessary.

```
package version

func Get() string {
    return "0.0.1"
}
```

## Changing the module path

The module path is currently set to `github.com/wumbabum/home_assist`. If you want to change this please find and replace all instances of `github.com/wumbabum/home_assist` in the codebase with your own module path.
