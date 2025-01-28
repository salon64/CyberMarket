# CyberMarket

A market website made in the course d0018e

## Backend API

This section goes trough the REST api that is used to communicate to the backend server

### Authentication users TODO

To authenticate a user send a request to `POST ronstad.se:5714/users/auth/{USERNAME}`
with the body `pswd={PSWD}`, example

```
POST /users/auth/myuser HTTP/2.0
...
pswd=MyPassword1234
```

This would return a token that could be used for authenticating the user

### Listing users

example

```
> GET /users HTTP/1.1
> Host: ronstad.se
> User-Agent: curl/7.81.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: *
< Date: Tue, 28 Jan 2025 12:25:02 GMT
< Content-Length: 55
< Content-Type: text/plain; charset=utf-8
<
[
    {
        "Id": 1,
        "Name": "myUser"
    }
]
```
