# go_web_server
A project to explore building a web server with RESTful functionality in Go.
Includes a basic database stored in JSON at the project root file: database.json

## Installation & Setup
Inside a go module:
> go get github.com/nicpatlan/go_web_server

The server runs on the localhost port 8080 but can be adjusted in main.go

## API
Typical access pattern will be http://localhost:8080 followed by an endpoint
detailed below.

### Server Metric Endpoints
- GET /api/healthz    "returns the server status"
- GET /admin/metrics  "returns the number of hits the server has received"
-     /api/reset      "resets the counter tracking server hits"

### User Endpoints
The header of a user endpoint http request should follow this pattern:
> {
>     "Authorization": "Bearer example_token"
> }

The body of a user endpoint http request should conform to the example below:
> {
>    "email": "name@example.com",
>    "password: "01234"
> }

- POST /api/users  "creates a new user from email and password given in body"
- PUT /api/users   "updates the authorized user's email and password"
- POST /api/login  "verifies user credentials and returns a JWT at login"    