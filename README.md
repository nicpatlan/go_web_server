# go_web_server
A project to explore building a web server with RESTful functionality in Go.
Includes a basic database stored in JSON at the project root file: database.json

## Installation & Setup
Inside a go module:
>     go get github.com/nicpatlan/go_web_server

The server runs on the localhost port 8080 but can be adjusted in main.go

## API
Typical access pattern will be http://localhost:8080 followed by an endpoint
detailed below.

### Server Metric Endpoints
- GET /api/healthz    "returns the server status"
- GET /admin/metrics  "returns the number of hits the server has received"
- /api/reset          "resets the counter tracking server hits"

### User Endpoints
- POST /api/users  "creates a new user from email and password given in body"
- POST /api/login  "verifies user credentials and returns a JWT at login"
- PUT /api/users   "updates the authorized user's email and password"

The request body to the endpoint for creating a user or to log a user in should
conform to the example below:
> {
>    "email": "name@example.com",
>    "password: "01234"
> }

The request header to the endpoint for updating a user should occur only after a
user has logged in and follow this pattern:
> {
>     "Authorization": "Bearer example_jwt"
> }    

The response body of the login endpoint will include a javascript web token or JWT
used for authentication for a period of 1 hour. This JWT is used for verifying posts
and updating user information.

### Refresh Token Endpoints
- POST /api/refresh "refreshes a JWT that has expired"
- POST /api/revoke  "revokes a refresh token"

These endpoints can be used to refresh an expired JWT or revoke a refresh token as 
needed. As previously mentioned a JWT is valid for 1 hour. A refresh token is 
returned in the response body of a login request and is valid for 60 days or until 
revoked.

The request header of these endpoints should contain the following pattern:
>{
    "Authorization": "Bearer example_refresh_token"
>}

A new refresh token can be acquired in the response body of a login request if it 
has expired or is no longer known.

### Post Create/Delete Endpoints
- POST /api/posts            "creates a new post by the user"
- DELETE /api/posts/{postID} "deletes a post from the user" 

Creating or deleting a post is an authenticated endpoint and requires a header conforming
to the following:
>{
    "Authorization": "Bearer example_jwt" 
>}

PostIDs can be obtained from in responses to requests made to the endpoints detailed 
below and are not authenticated. All posts are visible to all users.

### Post Retrieval Endpoints
- GET /api/posts      "returns a listing of all posts currently stored in the database"
- GET /api/posts/{id} "returns a listing of all posts created by the given userID"

UserIDs are included with each post and can be obtained in the response body of these
endpoints. These endpoints are not authenticated and do not require a logged in user.

### Webhook Endpoint
- POST /api/polka/webhooks "sets a users premium user status"

This endpoint currently sets a users premium status but can be expanded to include more
webhook functionality as needed. This is an authenticated endpoint that utilizes an 
API key to authorize the premium user setting. The header of the request to this endpoint
should include the following:
>{
    "Authorization": "ApiKey example_APIkey" 
>}