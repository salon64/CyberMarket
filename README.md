# CyberMarket

[Github Link](https://github.com/salon64/CyberMarket)

## Executive summary

The goal of this website is to provide a marketplace for users to buy and sell items inspired by the tabletop game Cyberpunk RED.
The site allows for users to enter items to the marketplace and set their price.
These items on the market would be able to be bought by other users transferring the item to the new users inventory and updating their funds.

## Table of Contents

- [Executive summary](#executive-summary)
- [Table of Contents](#table-of-contents)
- [requirements, assumptions, limitations](#requirements-assumptions-limitations)
- [Changelog](#changelog)
  - [Sprint 1 \& 2](#sprint-1--2)
- [Upcoming work (Prioritized backlog)](#upcoming-work-prioritized-backlog)
- [System Architecture](#system-architecture)
  - [Backend](#backend)
  - [Database](#database)
    - [Schema](#schema)
    - [Transaction Log](#transaction-log)
    - [Comments and grading](#comments-and-grading)
  - [WebServer](#webserver)
  - [Frontend](#frontend)
  - [Debugging frontend](#debugging-frontend)
  - [Discord Frontend](#discord-frontend)
- [User stories and Testing](#user-stories-and-testing)
  - [User Stories](#user-stories)
  - [Testing](#testing)
- [Reflection](#reflection)
  - [frontend specific reflection](#frontend-specific-reflection)
  - [backend specific reflection](#backend-specific-reflection)
- [Starting the application](#starting-the-application)
- [Website Documentation](#website-documentation)
- [Backend API](#backend-api)
  - [User Login](#user-login)
  - [Listing users](#listing-users)
  - [Add a user](#add-a-user)
  - [Update user info](#update-user-info)
  - [Add funds](#add-funds)
  - [Show wallet](#show-wallet)
  - [List user Items](#list-user-items)
  - [Create Item](#create-item)
  - [List Market](#list-market)
  - [Adding a listing](#adding-a-listing)
  - [Removing a listing](#removing-a-listing)
  - [Adding a comment on an itemtype](#adding-a-comment-on-an-itemtype)
  - [Get information and comments on an itemtype](#get-information-and-comments-on-an-itemtype)
  - [Buy](#buy)
- [References](#references)

## requirements, assumptions, limitations

TODO *from task description*

## Changelog

changes and added features, the full view exist on
[GithuB CyberMarket Project](https://github.com/users/salon64/projects/3/views/1)

### Sprint 1 & 2

reflections and result of sprint 1 and 2,
see [GithuB CyberMarket Project](https://github.com/users/salon64/projects/3/views/1)

## Upcoming work (Prioritized backlog)

A Github project is used to track issues and features that are wanted.
[GithuB CyberMarket Project](https://github.com/users/salon64/projects/3/views/1)

TODO *from task description*: Prioritized backlog (or reference to it in a tool): Each item should have ID, description, priority(importance), effort(small, medium, large), sprint no (when assigned), status (e.g. when finished).

## System Architecture

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

### Backend

To create this website, a backend application designed in go handles access to the database.
The task delegated to the backend is handling of request that interact with the data of the system.
Good examples of this are Login, buying, or viewing a users inventory.
Some of these required a token be passed along in the header for authorization, so the backend is also the major source of security for the system.
In the current architecture the backend also has access to the webserver data, this is for the possibility to add images to the webserver,
by handling images this way it allows the webserver to  handel image request which frees upp recourses from the backend.

Go was chosen as the language as its included standard library has http support as well an standard sql interface (note, driver is installed separately).
Another benefit of go is familiarity as members of this project have worked with it before

### Database

The chosen database for this project is mysql, the choice was made as its one of the most common open source databases

#### Schema

- TokenTable, this table represents the authorization token given out to each user on login,
- User, Represents the users, with their money and role
- Marketplace, each listing is uniquely related to an item for sale with price and the date the listing was created
- Inventory, represents items and their owner
- ItemTypes, contain descriptive data of item types, such as image, name and description

```mermaid
erDiagram
    Users {
        int UserID "PK auto_inc"
        varchar(45) Username "uniq, not null"
        varchar(45) Password "not null"
        int role "not null, uniq, default=0"
    }

    Marketplace {
        int OfferID "PK, auto_inc"
        int ItemID "uniq, not null"
        int Price "not null"
        datetime CreationDate "not null"
    }

    Inventory {
        int ItemID "PK, auto_inc"
        int UserID "not null"
        int TypeID "not null"
    }

    ItemTypes {
        int TypeID "PK, auto_inc"
        varchar(45) ItemName "not null"
        varchar(45) ItemDescription "null"
        varchar(45) ImgURL "null"
    }

    TokenTable {
        binary(16) Token "PK"
        int UserID  "not null"
        datetime CreatedOn "not null"
    }

    TransactionLog {
        int TransID "PK auto_inc"
        int Price "not null"
        datetime Date "not null"
        int ItemID "null"
        int Buyer "null"
        int Seller "null"
    }
    TokenTable }|--|| Users: Bearer

    Inventory }|--|| Users: owner
    Inventory }|--|| ItemTypes: Type

    Marketplace ||--|| Inventory: listing

    TransactionLog }o--o| Inventory: Item
    TransactionLog }o--o| Users: Buyer
    TransactionLog }o--o| Users: Seller

    %% TODO UPDATE SCHEMA WITH COMMENTS
```

#### Transaction Log

The transaction log is used to keep track of all transaction, such as when a user buys an item from the marketplace.
The log also allows for tracking when an item is created, buy setting the seller and price to null.

#### Comments and grading

TODO *From task description*: Write about your implementation of grading and comments on particular assets.  

### WebServer

The webserver hosting the files are currently not finalized,
currently to run the frontend the included nodejs from vite is used.
Possible plans for the future are to use vite for build and using the docker container busybox as an http server

### Frontend

The frontend library chosen was react as its common and great recourses exist to assist development.
A full frontend library may result in extra work which is not relevant for d0018e (uni course this application is developed in)
But this applications goals is not only for the course as it will see use after d0018e has concluded.

To build the react application Vite was offered to us as an viable tool from @voffiedev

### Debugging frontend

A simple rust terminal interface is implemented as an client that is used for debugging, see github repo

### Discord Frontend

Currently not implemented, work is planed after sprint 4

## User stories and Testing

TODO *from task description*: Test case specifications (manual testing, related to the roles and user stories) Write step by step how you will execute the demo for TA. The test are manual and are based on the user stories bellow

### User Stories

1. As a User I want to be able to register an account on the website. So that i have access to the website.

2. As a User I want to be able to log into my account, In order to view my profile/inventory/marketplace.

3. As a User I want to be able to view my profile, so that i can see account details about my account.

4. As a User I want to be able to update my username/password in my profile page, in order to have account the account information i prefer.

5. As a User I want to be able to view my inventory in order to see my wallet and items.

6. As a User I want to be able to add money into my account, to have sufficient funds.

7. As a User I want to be able to create items that i can list on the marketplace, to interact with the market.

8. As a User I want to be able to sell an item from my inventory at a price i set to get rid off items from my inventory.

9. As a User I want to be able to retract a listing from the marketplace I have made, In order to get that item back solely to my inventory.

10. As a User I want to be able to see items listed in them market with details about it (itemtype, price, description, seller). So that i can see the items currently listed.

11. As a User I want to be able to sort the market by different metrics (time, price, alphabetically). To display the market in the way i prefer.

12. As a Users I want to be able to buy items from the marketplace, In order to get that item to my inventory.

### Testing

1. To register An account The User should press on the "register your account button" that should take the user to a different page.The user should then fill in a user name and password and press the register button. When the button is pressed a http call to the backend to register a new user with the details the user filled in. The backend should take the request and make a new account in the db and return a userID and token.
TODO add checks and handle non unique usernames

2. To log into an account the user should enter username and password in their input felids and press the button/ press enter on their keyboard. That should send a http call with the information to the backend that should return a userId and a token if matching and an error if not

3. When logged in a user should be able to press the "profile" button in the navbar, it should take them to the profile page and should display the users ID.

4. When on the profile website the User should be able to update their user information by entering the new data into the fields and pressing the "save changes" button. That should send a http call to the backend that should check the token to see if the token has permissions for that userID, If it doesn't it should return an error and if it does it should update the information in the DB and return oldname + newname or old psw + new psw.

5. As a user I should be able to press the inventory button in the navbar. That should display the users wallet. It should also display the users items in a table with a distinction if they are in the marketplace or not. If they are not in the marketplace it should display a sell option with a input prompt and a sell button, and if it is in the marketplace it should display a button to recall the item from the marketplace.

6. To add money to to their wallet the user should type the users ID and amount in the input fields under "add money to wallet" and press the "add money" button. That should send a http call to the backend where it should check if the active user can add money to that wallet, if it cant it should return an error. If it can it should add the money (in the DB) to the wallet of the user and return the total amount of money now in that persons wallet.

7. To add a new item to a users inventory the user should add the users id and the itemType in their respective fields under "create item" and press the button "create item". That should send a http call to the backend, where it should check if the one who sent the request can add items to that inventory. It should also check if the itemtype exists. and return an error if they dont. If it passes the checks it should send to the DB to add that the item to that users inventory and return the ItemID, userId and itemtype.

8. A user should be able to sell an item from their inventory that is not already on the marketplace by entering the price in items row and press sell. That should send a http call to the backend, that should check if the caller can sell that users items, and return an error if it shouldn't. if it should it will add the item to the marketplace table in the db and return the OfferID.
TODO

9. A user should be able to remove a listed item from the marketplace that they have access to by pressing the "recall" button. That should send a http call to the backend where it checks if the caller can recall that persons items and return an error if not. If it should it should update the DB by removing the item from the marketplace table and returning the offerID
TODO

10. When a user is on the marketplace it should display all the items on the marketplace in the sorted order, on the left side it should display what they are sorted by. This data should be gathered with an http call from when the user loads the site. that call should return all the items in json format.

11. A user should be able to change the order the items are displayed in by pressing the section bar and selecting another method. That should send a call to the backend where that calls the db for the marketplace information to be returned in that order. it should then update the table with all the items to the marketplace sorted in the new order.

12. A user should be able to buy an item by pressing the "buy" button on the row that item is displayed in. That should send a http call to the backend. The backend will perform checks to see that you are not buying your own item, you have sufficient funds etc. and return an error if it fails a check. if not the item should be transferred to the new owner, funds should be updated and and a transaction should be added to the transaction log in the db. then return the OfferID
TODO

## Reflection

TODO *from task description* A reflection of the system's limitations and the possibilities for improvements.

### frontend specific reflection

### backend specific reflection

## Starting the application

Each part is a self contained container, to launch them we have provided a docker compose file for your convenience.
To run the docker file type

```bash
docker compose up --build
```

To start the go backend application only and connect to known database.

```sh
docker run --rm -e DBHOST=database.org:3306 -e DBUSER=root -e DBPASS=pswd -p 5687:5687 $(docker build -q ./backend)
```

If Go is installed, access to environment variables and privileges to host on a specific ports this command could be run (faster than to build a docker image)

```sh
TODO "@spookyfirefox needs to write this"
```

Note that you ether have run as superuser or be in the docker group

## Website Documentation

TODO *from sprint 2 review*

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

### Update user info

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

### Add funds

The request ``POST /user/AddMoney`` with the additional fields UserID and Amount can any amount of money to an existing user.

```curl
POST /user/AddMoney HTTP/1.1
Host: example.org
User-Agent: curl/7.81.0
Accept: */*
Content-Length: 27
Content-Type: application/x-www-form-urlencoded

{"UserID":9, "Amount":1000}

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Date: Fri, 21 Feb 2025 10:00:52 GMT
Content-Length: 0
```

### Show wallet

A user can send the request ``GET /user/getMoney/{uid}`` in order to view the amount of currency they hold.
As of writing this, anyone can send a request to view anyone's wallet.
Bearer tokens will be implemented (in most, if not all http requests) in a later sprint to verify which user has access to this information.

```curl
GET /user/getMoney/12 HTTP/1.1
Host: example.org
User-Agent: curl/7.81.0
Accept: */*

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Date: Fri, 21 Feb 2025 10:28:31 GMT
Content-Length: 38
Content-Type: text/plain; charset=utf-8

[
    {
        "Amount": 8788
    }
]
```

### List user Items

NOTE this will later change to a POST request which would require the token to be passed along for verification

The way to list a users items, GET a post request to ``/inventory/{id}`` where the id is the users id.
Possible return values are 404 if the user is not found. 400 is returned if the token is not valid

```curl
GET /inventory/1 HTTP/1.1
Host: example.org
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
        "ItemName": "CyberArm",
        "ItemDescription": null,
        "ImgURL": null
    },
    {
        "ItemID": 3,
        "TypeID": 3,
        "ItemName": "TechTool",
        "ItemDescription": null,
        "ImgURL": null
    }
]
```

### Create Item

Sending the http request ``POST /Marketplace/CreateItem`` allows a user to create an item and place it into any user's inventory with the fields UserID and ItemType.

```curl
POST /Marketplace/CreateItem HTTP/1.1
Host: example.org
User-Agent: curl/7.81.0
Accept: */*
Content-Length: 26
Content-Type: application/json

{"UserID":9, "ItemType":3}

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Date: Fri, 21 Feb 2025 09:50:31 GMT
Content-Length: 38
Content-Type: text/plain; charset=utf-8

{
    "UserID": 9,
    "ItemType": 3
}
```

### List Market

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
        "ItemName": "CyberArm",
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

### Adding a listing

Note, this will later be changed to follow REST

By sending a POST to ``/Marketplace/addListing`` with the fields ItemID and Price in json format, a listing is created

```curl
POST /Marketplace/addListing HTTP/1.1
Host: example.org
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

### Removing a listing

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

### Adding a comment on an itemtype

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
    "Comment":"must disgusting"
}

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Date: Mon, 03 Mar 2025 11:37:11 GMT
Content-Length: 17
Content-Type: text/plain; charset=utf-8

Added the comment
```

### Get information and comments on an itemtype

To get information about an itemtype send a get request to ``GET /ItemType/{ItemTypeId}``
Note that there is no transaction used when getting this information,
this means that the first part of the transaction, getting mane and such can complete then, the itemType and its comments are deleted.
Then the method tries to get the Comments.
This was deemed acceptable since its meant that item types are more or less constant.

```curl
> GET /ItemType/1 HTTP/1.1
> Host: example.org
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
            "Comment": "must disgusting",
            "PostedOn": "2025-03-03T10:19:03Z"
        }
    ]
}
```

### Buy

This Call would by the item defined by the path ``/Marketplace/buy/{itemID}``, the passed JSON ``UserID`` would be the new owner of that item.
The funds from teh listing is transferred from the buyer (passed UserID) ot the owner of the item. The bought objects owner is updated to the buyer.
If the buyer does not have enough funds the plaintext "not enough funds" is returned with the status code 200 ok.

Notes on security.
Currently there is no authentication on this. This means that another user can buy stuff as another user,
will be fixed when token authentication is implemented in this method.
To avoid that the user does purchases in quick sequence to trick the found checking by having both checks complete before any of the purchases have happened,
the whole of checking funds, transferring ownership, updating funds and removing the listing is done in a singe sql transaction

DANGEROUS FLAW, CURRENTLY BEING FIXED.\
imagine the following scenario

1. User **A** adds a listing för item **1** on the marketplace with the price **100 EUR**
2. User **B** loads the website with the listing of item **1** with price **100 EUR**
3. **A** removes the listing
4. **A** add the same item **1** on the market place for **250 EUR**
5. **B** haven't reloaded the website still **sees price **100**
6. **B** buys item **1** thinking it buys it for **100 EUR**
7. **B** actually buys the item for **250 EUR**

This can happen since we use itemID as the identifier, by changing to marketID this is prevented, since marketplace ids are unique

```curl
POST /Marketplace/buy/5 HTTP/1.1
Host: example.org
User-Agent: curl/7.81.0
Accept: */*
Content-Length: 13
Content-Type: application/json

{"UserID":12}

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Date: Fri, 21 Feb 2025 10:25:31 GMT
Content-Length: 7
Content-Type: text/plain; charset=utf-8

Success
```

## References

TODO move references from google doc to here
