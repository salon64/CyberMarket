# CyberMarket

A market website made in the course d0018e

## Starting the application

Each part is a self contained container, to launch them we have provided a docker compose file for your convenience.
To run the docker file type

```bash
docker compose up --build
```

The structure of the services are as shown below

```mermaid
architecture-beta

    group compose(server)

    service db(database)[Database] in compose
    service backend(server)[Backend] in compose
    service webserver(server)[Web server] in compose
    service web(disk)[Web storage] in compose
    service storage(disk)[Storage] in compose

    db:B --> T:backend
    db:R --> L:storage

    backend:R --> L:web
    webserver:T --> B:web

```

Note that you ether have run as superuser or be in the docker group

## Backend API

This section goes trough the REST api that is used to communicate to the backend server.

Note many paths require you to set the authorization header, when missing the server would respond with 404 to avoid leakage of data

### User Login

This creates a token that the user would use to authenticate its actions across other api calls.
The POST request takes the field ``username`` and ``pswd``, this returns the token as well as the associated user id.
This function might return a simple string indicating a error if one happened

```curl
> POST /login HTTP/1.1
> Host: example.org
> User-Agent: curl/7.81.0
> Accept: */*
> Content-Length: 55
> Content-Type: application/json
>
{
    "name" : "john_doe"
    "pswd" : "dog1234!"
}

< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: *
< Date: Tue, 28 Jan 2025 12:25:02 GMT
< Content-Length: 55
< Content-Type: text/plain; charset=utf-8
<
{
    "Token": "00112233445566778899aabbccddeeff",
    "Userid" : 6
}
```

### Listing users

Calling ``GET http://example.org/users`` returns the public info off all users.
Note this query might be slow

```curl
> GET /users HTTP/1.1
> Host: example.org
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
    },
    {
        "Id": 2,
        "Name": "SEAL_MAN"
    }
]
```

### Add a user

To add an user use the ``POST http://example.org/user``.
The body of this request, should contain ``name`` and ``pswd``.
Note that the max length of both pswd and name is 45 bytes.

```curl
> POST /user HTTP/1.1
> Host: example.org
> User-Agent: curl/7.81.0
> Accept: */*
> Content-Length: 23
> Content-Type: application/json
>
{
    "name": "john_doe",
    "pswd": "dog1234!"
}

< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: *
< Date: Tue, 28 Jan 2025 14:31:20 GMT
< Content-Length: 1
< Content-Type: text/plain; charset=utf-8
<
{
    "UserID": 7,
    "Token": "c6e88ef3-533a-40c1-b08c-ee2074a3a5dc"
}
```

## Update user info

to update the user send a PATCH request to ``/users/{id}``.
This requires the fields token, where the fields new_name or new_pswd are optional.
This will then return the new username as well as the old one

```curl
> PATCH /users/6 HTTP/1.1
> Host: example.org
> User-Agent: curl/7.81.0
> Authorization: Bearer c6e88ef3-533a-40c1-b08c-ee2074a3a5dc
> Accept: */*
> Content-Length: 23
> Content-Type: application/json
>
{
    "new_name": "svenne_bannan",
    "new_pswd": "cat4321?"
}

< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: *
< Date: Tue, 28 Jan 2025 14:31:20 GMT
< Content-Length: 64
< Content-Type: text/plain; charset=utf-8
<
{
    "old_name": "john_doe",
    "new_name": "svenne_bannan",
}
```

## List user Items

NOTE this will later change to a POST request which would require the token to be passed along for verification

The way to list a users items, GET a post request to ``/inventory/{id}`` where the id is the users id.
Possible return values are 404 if the user is not found. 400 is returned if the token is not valid

```curl
GET /inventory/1 HTTP/1.1
Host: ronstad.se
User-Agent: curl/7.81.0
Accept: */*

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Date: Tue, 11 Feb 2025 13:53:28 GMT
Content-Length: 426
Content-Type: text/plain; charset=utf-8

[
    {
        "ItemID": 1,
        "TypeID": 1,
        "ItemName": "MRE",
        "ItemDescription": null,
        "ImgURL":null
    },
    {
        "ItemID": 2,
        "TypeID": 2,
        "ItemName": "Cyberarm",
        "ItemDescription": null,
        "ImgURL": null
    },
    {
        "ItemID": 3,
        "TypeID": 3,
        "ItemName": "Techtool",
        "ItemDescription": null,
        "ImgURL": null
    }
]
```

## List Market

NOTE, This will be rewrite as a GET request when i force Malcolm to follow the REST standard.

NOTE 2, currently search is not implemented

NOTE 3, if ``SortBy`` is not one of the valid strings an unrelated error is returned, will be fixed

This ``POST /Marketplace/displayMarket`` method takes the arguments ``SortBy`` and ``Search`` in json format and returns the market place items, the valid strings for OrderBy are

- Newest
- Oldest
- Price_Ascending
- Price_Descending
- Alphabetically_Ascending
- Alphabetically_Descending

Example

```curl
POST /Marketplace/displayMarket HTTP/1.1
Host: example.org
User-Agent: curl/7.81.0
Accept: */*
Content-Length: 32
Content-Type: application/x-www-form-urlencoded

{"SortBy":"Newest", "Search":""}

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Date: Tue, 18 Feb 2025 14:25:54 GMT
Content-Length: 561
Content-Type: text/plain; charset=utf-8

[    
    {
        "ItemID": 2,
        "TypeID": 2,
        "UserID": 1,
        "ItemName": "Cyberarm",
        "ItemDescription": null,
        "ImgURL": null,
        "OfferID": 2,
        "Price": 20,
        "CreationDate": "2025-02-12T00:00:00Z",
        "Username": "sha"
    },
    {
        "ItemID": 1,
        "TypeID": 1,
        "UserID": 1,
        "ItemName": "MRE",
        "ItemDescription": null,
        "ImgURL": null,
        "OfferID": 6,
        "Price": 10,
        "CreationDate": "2025-02-12T00:00:00Z",
        "Username": "sha"
    }
]
```

## Adding a listing

Note, this will later be changed to follow REST

By sending a POST to ``/Marketplace/addListing`` with the fields ItemID and Price in json format, a listing is created

```curl
POST /Marketplace/addListing HTTP/1.1
Host: ronstad.se
User-Agent: curl/7.88.1
Accept: */*
Content-Length: 26
Content-Type: application/x-www-form-urlencoded

{"ItemID":4, "Price":2000}

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Date: Tue, 18 Feb 2025 14:48:07 GMT
Content-Length: 1
Content-Type: text/plain; charset=utf-8

8
```

## Removing a listing

Note, this will be rewritten as a DELETE

To remove an listing Send the following get request ``GET /Marketplace/removeListing/{ItemID}``
where itemID is the intended itemid of the listing to be removed

```curl
GET /Marketplace/removeListing/1 HTTP/1.1
Host: example.org
User-Agent: curl/7.81.0
Accept: */*

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Date: Tue, 18 Feb 2025 14:42:33 GMT
Content-Length: 16
Content-Type: text/plain; charset=utf-8

removed listing
```

## Adding a comment on an itemtype

This adds a comment to the specified itemtype use ``POST /ItemType/{TypeID}``.
The required json fields are as following ``Grade``, ``UserID`` and ``Comment``.
To be allowed to post you need to have bought this itemType before. This is checked by looking in the transaction log.
Currently only buying from marketplace inserts an entry in this log, but in the future adding an item will be the same as buying from null with price null.

Note that there are currently no security, There are no checks that the userid is you.

```curl
POST /ItemType/1 HTTP/1.1
Host: example.org
User-Agent: curl/7.88.1
Accept: */*
Content-Length: 53
Content-Type: application/x-www-form-urlencoded

{
    "UserID":12,
    "Grade":0,
    "Comment":"must discusting"
}

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Date: Mon, 03 Mar 2025 11:37:11 GMT
Content-Length: 17
Content-Type: text/plain; charset=utf-8

Added the comment
```

## Get information and comments on an itemtype

To get information about an itemtype send a get request to ``GET /ItemType/{ItemTypeId}``
Note that there is no transaction used when getting this information,
this means that the first part of the transaction, getting mane and such can complete then, the itemType and its comments are deleted.
Then the method tries to get the Comments.
This was deemed acceptable since its meant that item types are more or less constant.

```curl
> GET /ItemType/1 HTTP/1.1
> Host: ronstad.se:5687
> User-Agent: curl/7.88.1
> Accept: */*

< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: *
< Date: Mon, 03 Mar 2025 10:54:08 GMT
< Content-Length: 498
< Content-Type: text/plain; charset=utf-8
< 
{
    "Name": "MRE",
    "ImgURL": null,
    "ShortDesc": null,
    "DescURL": null,
    "Comments": [
        {
            "UserName": "pelle",
            "UserID": 10,
            "Grade": 5,
            "Comment": "Tastes good mmm",
            "PostedOn": "2025-02-27T14:54:26Z"
        },
        {
            "UserName": "salonguy",
            "UserID": 12,
            "Grade": 0,
            "Comment": "must discusting",
            "PostedOn": "2025-03-03T10:19:03Z"
        }
    ]
}
```

## Delete comment

To delete a comment send ``DELETE /comment/{commentID}``

```curl
> DELETE /comment/3 HTTP/1.1
> Host: example.org
> User-Agent: curl/7.88.1
> Accept: */*
 
< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: *
< Date: Tue, 04 Mar 2025 09:54:05 GMT
< Content-Length: 15
< Content-Type: text/plain; charset=utf-8
 
Deleted comment
```
