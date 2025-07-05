# Welcome to the goChirpy - "twitter"-type server written in Golang!

This project is a part of the boot.dev course (guided project).
It allows you to host sever on your local machine and post/delete/read 140 symbols long messages (chirps) in/from PostgreSQL database. This project also uses user registration and autentication (through JWT access and refresh tokens).

To run this program you need:
1. [go](https://webinstall.dev/golang/);
2. [postgreSQL database](https://webinstall.dev/postgres/);
The database should be migrated to the version 005_userss. I used [Goose command line tool](https://github.com/pressly/goose#install). To run a migration with it, run "goose postgres *connection_string* up" in folder "sql/schema "
3. .env file in the root of the project with the following environment variables:
```
DB_URL="*connection string*" - db connection string
PLATFORM="dev" - to access *admin/* endpoints
SECRET="*some random string*" - secret string for your website JWTs
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e" - as it is. Used to verify "Subscription webhook"
```

To install goGator do "go install" in the goGator folder. Now the goGator can be run!


The website/app is available using **http://localhost:8080/app/** link.
Other than that, the following endpoints are available (using **http://localhost:8080/app/***endpoint*):

0. **/api/healthz** - GET request - returns code 200 if the server is up;

1. **/admin/reset** - POST requst - requires no body and deletes all the users from the *users* table;
Requires correct **PLATFORM** .env variable

2. **/api/users** - POST request - to creat Chirpy user;
Requires JSON with the following fields: 
```
string `json:"email"`
string `json: "password"`
```
Returns json with the user data.

3. **/api/users** - PUT request - to update Chirpy user (new password and email);
Requires JSON with the following fields: 
```
string `json:"email"`
string `json: "password"`
```
Also requires access token in the "Authorization" header in the form *"Bearer token"* for changes authorization. Returns json with the updated user data.

4. **/api/login** - POST request - to login Chirpy user;
Requires JSON with the following fields: 
```
string `json:"email"`
string `json: "password"`
```
The responce returns JWT access token (valid 1 hour) and refresh tokens (valid 60 days) in the `json:"token"` and `json:"refresh_token"` fields, respectively.

5. **/api/refresh** - POST request - returns a new JWT access token;
Requires non-expired/revoked refresh token in the "Authorization" header in the form *"Bearer token"* for access token generation. The responce returns access token in the `json:"token"` field.

6. **/api/revoke** - POST request - returns a new JWT access token;
Requires refresh token in the "Authorization" header in the form *"Bearer token"*. Revokes this token if it exist in database.

7. **/api/chirps** - POST request - to post a chirp;
Requires JSON with the following fields: 
```
string `json:"body"` (message)
string `json:"user_id" (user ID)
```
Also requires access token in the "Authorization" header in the form "Bearer token". Returns json with the chirp data.

7. **/api/chirps** - GET request - to return the list of chirps;
Returns all the chirps in the database if not asked otherwise. Accepts two optional parameters in the reluest link:
* "author_id" (*api/chirps?author_id=1*) - returns chirps for the specific user_id 
* "sort" (*api/chirps?sort=desc*) - return chirps in the descending order of their creation (ascending by default)
Both of the parameters can be combined.

7. **/api/chirps/{chirpID}** - GET request - to return the chirp with the given chirpID;

8. **/api/chirps/{chirpID}** - DELETE request - to delete the chirp with the given chirpID;
Also requires access token in the "Authorization" header in the form *"Bearer token"* for chirp deletion. Will only delete the chirp if the access token's user_id is the same as the user_id of the chirp creator.

9. **/api/polka/webhooks** - POST request - auxiliary endpoint to validate "chirpy red subscription" from the webhook;
Requires API key in the "Authorization" header in the form *"ApiKey token"*. The key should be equal to the .env POLKA_KEY variable.


Have a pleasant goChirpy use! If you have any issues with running this program or have suggestions - don't hesitate to let me know!


