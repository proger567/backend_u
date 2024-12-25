# REST API for Users

## Application API

`GET` **/user** `Get user by name from JWT token`

`GET` **/users** `Get all users<->role`

`POST` **/user** `Add user to local database`

`PUT` **/user** `Update user's role`

`DELETE` **/user/{username}** `Delete user by name`

## Environment Variables:

| Variable    | Default value | Description                                      |
|-------------|---------------|--------------------------------------------------|
| LOG_LEVEL   | INFO          | this word level logger(INFO, DEBUG, ERROR, WARN) |
| LISTEN_PORT | :80           | it is listen port                                |
| SECRET_KEY  | *empty*       | it is secret key for check JWT token             |
| DB_HOST     | *empty*       | this is a URL where database is hosted           |
| DB_NAME     | *empty*       | gggg                                             |
