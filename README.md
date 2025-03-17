# CyberMarket

- [Github Link](https://github.com/salon64/CyberMarket)
- [Github Backlog](https://github.com/users/salon64/projects/3/views/1)

## Executive summary

The goal of this website is to provide a marketplace for users to buy and sell items inspired by the tabletop game Cyberpunk RED.
The site allows for users to enter items to the marketplace and set their price.
These items on the market would be able to be bought by other users transferring the item to the new users inventory and updating their funds.

## Table of Contents

- [Executive summary](#executive-summary)
- [Table of Contents](#table-of-contents)
- [Requirements, Assumptions, Limitations](#requirements-assumptions-limitations)
  - [Requirements](#requirements)
    - [Functional Requirements](#functional-requirements)
    - [Non-Functional Requirements](#non-functional-requirements)
  - [Assumptions](#assumptions)
  - [Limitations](#limitations)
- [Changelog](#changelog)
  - [Sprint 1 \& 2](#sprint-1--2)
- [Upcoming work (Prioritized backlog)](#upcoming-work-prioritized-backlog)
- [System Architecture](#system-architecture)
  - [Backend](#backend)
  - [Database](#database)
    - [Schema](#schema)
    - [Transaction Log](#transaction-log)
    - [Comments and grading](#comments-and-grading)
    - [Shopping basket](#shopping-basket)
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
    - [HTTP](#http)
    - [SQL](#sql)
    - [JSON](#json)
    - [Error handling](#error-handling)
      - [Rust](#rust)
      - [Go](#go)
      - [Python](#python)
- [Starting the application](#starting-the-application)
- [Website Documentation](#website-documentation)
  - [Login Page](#login-page)
  - [Register](#register)
  - [Marketplace](#marketplace)
  - [Inventory](#inventory)
  - [Profile](#profile)
  - [Admin](#admin)
  - [Cart](#cart)
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
  - [Olle](#olle)
  - [Shaya:](#shaya)
  - [Malcolm:](#malcolm)
    - [Important:](#important)
    - [Informal](#informal)
  - [General](#general)

## Requirements, Assumptions, Limitations

### Requirements

#### Functional Requirements

Some functional requirements would be allowing the user to login, view their inventory, purchase items from the marketplace, sell items they own, add items to a shopping cart and checkout.

#### Non-Functional Requirements

Non-functional requirements would be;

- Usability, the site needs to be easy to use
- ACID, a marketplace is one of the prime examples of where the principles of Atomicity, Consistency, Isolation and Durability should be followed
- Security/Data Integrity, user information should be safely stored and encrypted
- Availability, the site should be consistently available to the user

### Assumptions

The system assumes that there are no man in the middle attacks or other attacks on non encrypted communications. (due to ssl cert cost money)
Internet traffic on the website is also assumed to be low as to not put too much of a load on the backend server.

### Limitations

Many relations slow down the database query's and executions. Items are not grouped by quantity, so if 20 MREs are put up for sale, there exists 20 different listing, one for each MRE.

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
    Marketplace ||--|| Inventory: listing

    Inventory {
        int ItemID "PK, auto_inc"
        int UserID "not null"
        int TypeID "not null"
    }
    Inventory }|--|| Users: owner
    Inventory }|--|| ItemTypes: Type

    ItemTypes {
        int TypeID "PK, auto_inc"
        varchar(45) ItemName "not null"
        varchar(255) ShortDescription "null"
        varchar(45) ItemDescriptionURL "null"
        varchar(45) ImgURL "null"
    }

    TokenTable {
        binary(16) Token "PK"
        int UserID  "not null"
        datetime CreatedOn "not null"
    }
    TokenTable }|--|| Users: Bearer

    TransactionLog {
        int TransID "PK auto_inc"
        int Price "not null"
        datetime Date "not null"
        int ItemID "null"
        int Buyer "null"
        int Seller "null"
    }
    TransactionLog }o--o| Inventory: Item
    TransactionLog }o--o| Users: Buyer
    TransactionLog }o--o| Users: Seller

    TypeComments {
        int CommentID "PK auto_inc"
        int ParentComment "null"
        int TypeID "not null"
        int UserID "not null"
        int Grade "not null"
        varchar(255) comment "not null"
        datetime CreatedOn "not null"
    }
    TypeComments }|--|| Users: PostedBy
    TypeComments }|--|| ItemTypes: Comment
    TypeComments |o--o| TypeComments: "Parent comment"

    ShoppingBasket {
        int CartID "PK auto_inc"
        int UserID "NOT NULL"
        int OfferID "NOT NULL"
    }
    ShoppingBasket }|--|| Users: Basket
    ShoppingBasket }|--|| Marketplace: Item

```

#### Transaction Log

The transaction log is used to keep track of all transaction, such as when a user buys an item from the marketplace.
The log also allows for tracking when an item is created, buy setting the seller and price to null.

#### Comments and grading

There exist the ability to comment on itemTypes, This allows user to add comments and a grade to itemtype.
For a user to be allowed to comment on an itemtype, they need to have been in possession of them item, hover briefly.
This is accomplished by looking att the transaction log, this allows us to test if that itemtype have ever been sold to the user,
this includes when an item is created since that is also tracked in the log

Comments can also be related to another comment as to achieve hierarchical comment structure,
note that for this to work the parent comment must have the same Type ID as this is the column that is quired on,
there currently exist no constraint to enforce this in the backend or database, the behavior of breaking this has not been studied.

#### Shopping basket

THe shopping basket keeps tracks of a users shopping basket, it acts as a many to many between users and marketplace listings. The foreign keys relation to marketplace is set to cascading delete, so that if an item on the marketplace is bought by another user or removed, the shopping cart reference is removed.

### WebServer

As react and vite was chosen to create the fronted, we need a way to host the static files vite generates,
An option would be to host them on the same application as the backend, this would have been the easier option
but would have some cons. One of with would have been the close link between the website which would have made
the system harder to link to other application, in this case Discord.
So the decision was made to host the website separate from the backend.
The final solution was set up a docker build system that builds the vite application and then sets up a busybox container as the webserver.

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

6. As a Admin I want to be able to add money into a account, to have sufficient funds.

7. As a Admin I want to be able to create items that i can list on the marketplace, to interact with the market.

8. As a User I want to be able to sell an item from my inventory at a price i set to get rid off items from my inventory.

9. As a User I want to be able to retract a listing from the marketplace I have made, In order to get that item back solely to my inventory.

10. As a User I want to be able to see items listed in them market with details about it (itemtype, price, description, seller). So that i can see the items currently listed.

11. As a User I want to be able to sort the market by different metrics (time, price, alphabetically). To display the market in the way i prefer.

12. As a Users I want to be able to buy items from the marketplace, In order to get that item to my inventory.

13. As A Admin I want to be able to create a new itemtype that can be added as an item to users inventories.

14. As A User I want to be able to make comments on items types i have owned.

15. As A User I want to be able to remove a comment I have made.

16. As A User I want to be able to view comments made on an itemtype.

17. As A Admin I want to be able to remove comments made by users.

18. As A Admin I want to be able to remove marketplace listings made by users.

19. As a Admin I want to be able to view the transaction log, and filter it based on Users.

20. As a Admin i should be able to spoof a user id, allowing me to complete actions as that users.

21. As a Users I want to be able to buy items from the marketplace, In order to get that item to my inventory.

22. As a user, I want to be able to add items/remove items from the marketplace to my cart.

23. As a user, I want to be able to remove items from inside my cart.

24. As a user, I want to be able to checkout all items from my cart for them to then be added to my inventory.


### Testing

1. To register An account The User should press on the "register your account button" that should take the user to a different page.
The user should then fill in a user name and password and press the register button. When the button is pressed a http call to the backend to register a new user with the details the user filled in. The backend should take the request and make a new account in the db and return a userID and token.\
**TODO** add checks and handle non unique usernames \
**NOTE** No need to log in the user

2. To log into an account the user should enter username and password in their input felids and press the button or press enter on their keyboard.
That should send a http call with the information to the backend that should return a userId and a token if matching and an error if not.
If the backend call was successful to user should be redirected to marketplace.

3. When logged in a user should be able to press the "profile" button in the navbar, it should take them to the profile page and should display the users ID.
It should also request the backend for transactions where user was seller of buyer, if user token is invalid the view of transactions should be empty.

4. When on the profile website the User should be able to update their user information by entering the new data into the update username or update password
and pressing the "submit username" or "submit password" button.
That should send a http call to the backend that should check the token to see if the token has permissions for that userID,
If it doesn't it should return an error and if it does it should update the information in the DB.

5. As a user I should be able to press the inventory button in the navbar. That should display the users wallet.
It should also display the users items in a table with a distinction if they are in the marketplace or not.
If they are not in the marketplace it should display a sell option with a input prompt and a sell button,
and if it is in the marketplace it should display a button to recall the item from the marketplace.
Before data is sent to the frontend the token is validated to see if the user is the owner or an admin.
if not an error is returned.

6. To add money to to their wallet the admin should type the users ID and amount in the input fields under "add money to wallet" and press the "add money" button.
That should send a http call to the backend where it should check if the active user is an Admin, if it can't the backend should return an error.
If it can it should add the money (in the DB) to the wallet of the user and return 200 ok.

7. To add a new item to a users inventory the admin should add the users id and the itemType in their respective fields under "create item" and press the button "create item".
That should send a http call to the backend, where it should check if the one who sent the request is an admin.
It should also check if the itemtype exists. And return an error if they don't.
If it passes the checks it should send to the DB to add that the item to that users inventory and return the ItemID, userId and itemtype.

8. A user should be able to sell an item from their inventory that is not already on the marketplace by entering the price in items row and press sell.
That should send a http call to the backend, that should check if the caller can sell that users items (the user is the owner or an admin), and return an error if it shouldn't.
If it can it will add the item to the marketplace table in the db and return the OfferID.

9. A user should be able to remove a listed item from the marketplace that they have access (either is the owner or an admin) to by pressing the "remove listing" button.
That should send a http call to the backend where it checks if the caller can recall that persons items and return an error if not.
If it should it should update the DB by removing the item from the marketplace table.

10. When a user is on the marketplace it should display all the items on the marketplace in the sorted order, on the left side it should display what they are sorted by.
This data should be gathered with an http call from when the user loads the site. that call should return all the items in json format.

11. A user should be able to change the order the items are displayed in by pressing the section bar and selecting another method. That should send a call to the backend where that calls the db for the marketplace information to be returned in that order. it should then update the table with all the items to the marketplace sorted in the new order.
    - Newest
    - Oldest
    - Price_Ascending
    - Price_Descending
    - Alphabetically_Ascending
    - Alphabetically_Descending

12. A user should be able to buy an item by pressing the "buy" button on the row that item is displayed in.
That should send a http call to the backend. The backend will perform checks to se if the sender of the http is the same of the buyer or an admin, and that the user have sufficient funds.
And return an error if it fails a check. if not the item should be transferred to the new owner, funds should be updated and and a transaction should be added to the transaction log in the db.

13. In the Admin page an admin should be able to add an new itemtype, by filling in information like name, descriptionUrl, imageUrl and short description, using the create itemtype button an http request should be sent to the backend, the request should be checked at the backend to see if the user is an admin otherwise return an error. If the user is an admin, the type should be added to the data base.

14. A user should be able when in market place page to add a comment to ether an itemtype or another comment, by filling out the fields grade and comment.
When user presses post button the a request should be sent to the backend where its checked if the user have ever owned the item.
Another check is done to determine if the token matches the user or the token belongs to an admin.
if the checks pass a new comment should be added to the database.

15. A user should be able to remove a comment they have made bu pressing a button next to their comment, a request should be sent the the backend where its the token is checked to se if the user is the commenter or an admin. if so the comment is remove, as well as all its children comments.

16. A user should be able to view comments on a itemtype by clicking it in the marketplace page, this should be facilitated by an request to the backend.
The comments should be hierarchically ordered with some comments att the top level and each child comment should be more indented than the parent.

17. A Admin should be able to remove a comment they have made buy pressing a button next to the comment, a request should be sent the the backend where its the token is checked to se if the user is an admin. if so the comment is remove, as well as all its children comments.

18. An admin should be able to remove marketplace listings made by other users, \
**TODO** currently this is only possible by the Spoof user feature, (changing userid while keeping the admin token, available in the admin panel)

19. an admin should be able the search the transaction log in the admin page for specif user or all, this should be supported by an request to the backend with the desired userid or "all". The backend should check if the token matches the requested user or if the user is an admin, if not an error should be returned.

20. An admin should be able to spoof another user by entering their userid in the admin page. This should change all outgoing request to use that userid.

## Reflection

TODO *from task description* A reflection of the system's limitations and the possibilities for improvements.

### frontend specific reflection

The frontend framework that was decided on was React
as it is the most popular framework today. The code was written in TypeScript to get more reliability from the typed nature of the language.

This might have been a downside since non off the developers had worked with neither react ot ts before causing a lot of headaches in development.
Most of them stem from not knowing the right way to work with the tools leading to some hacks having to stay in place to hold it together in some places.
Other then that the language offered a lot of tools to do anything you wanted with most of our problems already having documented solutions
on stackOverflow or other websites.

### backend specific reflection

As mentioned in [Backend](#backend) go was chosen as the language, due to its standard library having http servers and sql interfaces.
This meant that only two external modules had to be used, the driver code for mysql which is invisible to the programer
and is dealt with the standard library sql implementation and for our token implementation an UUID module was needed.
The small mount of modules required did make the backend easier to work with and since the uuid module was written by google there was an unified philosophy.

#### HTTP

The HTTP module that is included in go is extensive and provide large about of the functnialty needed for any web application.
This project encountered some problems regarding *CORS* policy and Complex http methods (**Delete**, **patch** and others),
the cors problem was solved by simply setting the response header for each request to allow *CORS*.
As well as a simple blanket preflight check that allows the authorization header.
While we currently avoid complex http methods such as delete and simply stuck with **post** and **get**.
In the future it would be wise to implement the better pre-flight checks that are required to use these functions.

#### SQL

Go sql implantation provided method that make it easy to use transactions, as instead of calling query or exec on the database object,
you use the same methods on a transaction object. There exist places where an transaction might benefit but it was decided against,
one example is ItemTypes Where first common information is retrieved (name, description, ect) then comments are retried.
There exist the possibility that the itemtype would have been deleted between these two operations.
This was okay as ItemTypes are not meant to be deleted, and even if they were deleted in the middle of the operation the effect would just be that the user saw no comets on that Itemtype
and the item type would be gone when the site refreshed.
For where the result of an race condition would be "fatal" transactions are used, example [Buy](#buy)

For more specific documentation about transactions see the backend api [documentation](#backend-api)

#### JSON

But Go json implantation caused us some trouble, First of for an field to be decoded or encoded the field needed to be public,
This was easy to fixed but also easy to miss as an public field starts with an capitalized letter compared to an private which starts with an lowercase.
Another annoyance is how default values and missing fields are handled, Since decoding to a struct dont return the struct but only updates an existing one.
When a struct is created to be filed with values decoded from json, all field are initializes to their default values, string are set to empty and ints to zero,
this causes problem for the decoders does not cause an error if an felid is not found in the json.
So if for example name is left out we cant determine if the filed was missing or intently left blank (empty string).
A way to avoid this is to use another json library which dont have these problems, but that would require a rewrite of most backed functions

#### Error handling

Another annoyance in go was error handling and Nil, Compared to rust which granites that no uninitialize values exist (no nil)
and uses the enum types result and option to conway the same meaning while allowing for error handling methods on these types.
Gu uses a more verbose system which forces you to think about error in the same way as rust but eventing is more verbose and laking
inbuilt operators but still not preventing nil deref errors.
This puts Go in an awkward position were error can be ignored as in Python but not dealt with later in a try catch block.

Personally (@spooky_firefox) i feel like go error handling tries t get you to do the right ting but ends up as verbose mess,
where each function that can return an error needs three lines to deal with it.

The following code blocks show the verbosely in go error handling.

##### Rust

```rust
fn i_can_error(arg: i32) -> Result<&i32,dyn Error> {
    if i < 0 {
        return Err()
    } else {
        return Ok(arg+1)
    }
}

fn main() -> Result<(),dyn Error>{
    println!("{}", i_can_error(-1).unwrap()) // Unwrap panic if the enum is Err

    if let Ok(output) = i_can_error(-1) {
        // this branch will only be evaluated if i_can_error return the enum ok,
    }

    println("{}", i_can_error()?) // the ? operator returns the value if its ok, otherwise it exits the function (this case main) and return the error
}
```

##### Go

```go
func iCanError(arg int) (int ,error) {
    if i<0 {
        return 0, error
    } else {
        return arg + 1, nil
    }
}

func main() error{
    val, _ := iCanError(-1) //ignore the error, our result might be wrong, might also be nil and cause nil deref error later on
    fmt.PrintLn(val)

    val, err := iCanError(-1) // this is the same as rust unwrap, but incapably verbose and often clutters up your codebase
    if err != nil {
        panic("i can error errored")
    }
    fmt.PrintLn(val)

    val, err = iCanError(-1) // this is the same as rust if let
    if err == nil {
        // do stuff
    }

    val, err = iCanError(-1) // this is the same as rust ? operator, but incapably verbose and often clutters up your codebase
    if err != nil {
        return err
    }
    fmt.PrintLn(val)
}
```

##### Python

```python
def iCanError(arg):
    if arg < 0:
        raise ValueError
    else:
        return arg+1

print(iCanError(-1)) # this is the same as rust ? operator, we dont have to del with error let some else do

try:
    print(iCanError(-1)) # this is the same as rust if let which gives us nice control if its an error
except:
    # do something
```

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

### Login Page

Allows users to enter a username and password to log in,
also directs Users to a register page

### Register

Allows users to enter a username and password and register that as a account.

### Marketplace

Shows the entire marketplace, with all its listed items with the data:

Itemname | Price | ShortDescription | Seller | BuyButton

The marketplace is sorted by postdate and can be sorted by
newest, oldest, price acending, price decending, alphabetically a-ö, alphabetically ö-a. in the top left

There are also tags but they are not yet interactive.
TODO implement tags

### Inventory

Shows the current balance of the user
and alsow shows the users items in the same way as on the marketplace.
But now having a remove/add listing (with a price input feild)

### Profile

Shows the users ID. Allows the user to update username and password

TODO fix profile

### Admin

Allows and admin to chage UserID to check that users view,
allows them to create an item into a users inventory,
allows them to add money to a users wallet,
allows them to create a new itemtype

### Cart

Shows the items in the users cart. Can remove items from cart. Can checkout (buy) all items in cart

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
Authorization: Bearer c6e88ef3-533a-40c1-b08c-ee2074a3a5dc
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
Authorization: Bearer c6e88ef3-533a-40c1-b08c-ee2074a3a5dc
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
Authorization: Bearer c6e88ef3-533a-40c1-b08c-ee2074a3a5dc
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

NOTE, currently search is not implemented

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
Authorization: Bearer c6e88ef3-533a-40c1-b08c-ee2074a3a5dc
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
Authorization: Bearer c6e88ef3-533a-40c1-b08c-ee2074a3a5dc
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
Authorization: Bearer c6e88ef3-533a-40c1-b08c-ee2074a3a5dc
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
> Authorization: Bearer c6e88ef3-533a-40c1-b08c-ee2074a3a5dc
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
To avoid that the user does purchases in quick sequence to trick the wallet checking by having both checks complete before any of the purchases have happened,
the whole of checking funds, transferring ownership, updating funds and removing the listing is done in a singe sql transaction

**DANGEROUS FLAW, CURRENTLY BEING FIXED.**\
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
Authorization: Bearer c6e88ef3-533a-40c1-b08c-ee2074a3a5dc
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

#### DisplayCart

Here we re-use and tweak code from the marketplace function that displays items by just changing around the query a bit to only select items from the shopping cart associated with the user requesting the cart.

#### Add To Cart

This function exists in the marketplace where a user clicks "Add to Cart" on an item to send an http request ```POST /Marketplace/addToCart/{OfferID}``` to the backend where we just insert the item with the buyers id into the shopping cart table.

#### Remove from cart

Similar to the Add to Cart function, except we simply remove the item from the cart table.

#### Cart Checkout

Once the user decides to checkout, they send an http request ```GET /Marketplace/cartCheckout/{UserID}```. The backend then checks if the user has enough money for the total purchase by querying the total price from the DB. We then begin an SQL transaction to ensure consistency and start looping through every item in the shopping cart and do something similar to the buy function on every item. This does mean one shopping cart checkout will generate one transaction for every item in the cart. Finally, if no errors occur, we commit the SQL transaction.

## References

### Olle

Important

- [example dockerize vite](https://dev.to/code42cate/how-to-dockerize-vite-44d3)
- [go database/sql](https://go.dev/doc/tutorial/database-access)

Informal

- [sql, overview and examples](https://www.w3schools.com/sql/) \
    contains most of the necessary information about sql, describing syntax and functions
- [go http documentation](https://pkg.go.dev/net/http)
- [http preflight, mdn](https://developer.mozilla.org/en-US/docs/Glossary/Preflight_request)
- [self signed ssl cert, https (not used in final version)](forgot_to_write_down_address)
- [go json documentation](https://pkg.go.dev/encoding/json)
- [go uuid documentation](https://pkg.go.dev/github.com/google/UUID)
- [go database driver](https://github.com/go-sql-driver/mysql/)
- [golang docker container](https://hub.docker.com/_/golang)
- [mysql docker container](https://hub.docker.com/_/mysql/)
- [docker env documentation](https://docs.docker.com/compose/how-tos/environment-variables/set-environment-variables/)
- [vite env documentation](https://vite.dev/guide/env-and-mode)
- [go sql null](https://medium.com/@dual/ways-to-handle-sql-null-values-in-go-e237f31b82fc)
- *old projects from when i had a webserver course in a swedish gymnasium*

### Shaya

- [http connections] <https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch>
- [general react syntax] <https://react.dev/reference/react>
- [private routes] <https://www.robinwieruch.de/react-router-private-routes/>
- [navigation/redirects] <https://api.reactrouter.com/v7/functions/react_router.useNavigate.html>
- [chatgpt] (For implementing hashroutes, due to browserRoutes causing the site to crash on reloads)

### Malcolm

#### Important

- [cyberpunk css style] <https://www.cssscript.com/demo/cyberpunk-2077/>
- [crt screen affects] <https://medium.com/@dovid11564/using-css-animations-to-mimic-the-look-of-a-crt-monitor-3919de3318e2>
- <https://codepen.io/lbebber/pen/XJRdrV/>
- <https://codesandbox.io/p/sandbox/crt-terminal-in-css-js-tlijm?file=%2Findex.html>
- <https://stackoverflow.com/questions/70498819/retro-crt-curved-screen-effect-for-website-ccs>
- <https://aleclownes.com/2017/02/01/crt-display.html>
- And AI (to assist in putting things together)
<https://www.youtube.com/watch?v=Gy4G68WoHq4>

#### Informal

- [react tutorial] <https://www.youtube.com/watch?v=SqcY0GlETPk>
- [CSS examples] <https://getbootstrap.com/>
- [user input] <https://stackoverflow.com/questions/53071851/getting-the-value-from-input-element-in-typescript>
- <https://stackoverflow.com/questions/57302715/how-to-get-input-field-value-on-button-click-in-react>
- <https://www.youtube.com/watch?v=zR5FoKMAJp4>

### General

- <https://stackoverflow.com/questions/54300334/react-change-state-based-on-select-option>
- <https://www.w3schools.com/react/react_useeffect.asp>
- <https://stackoverflow.com/questions/3304014/how-to-interpolate-variables-in-strings-in-javascript-without-concatenation>
- <https://stackoverflow.com/questions/29051368/how-do-i-make-a-button-fill-in-a-td-completely>
- <https://stackoverflow.com/questions/47270595/how-to-parse-json-string-to-struct>
- <https://medium.com/@prashant2018/parse-json-string-to-struct-object-in-golang-59d3cf6bce7>
- <https://www.w3schools.com/css/css_text.asp>
- <https://dev.to/writech/how-to-handle-the-mouse-hover-event-in-react-h7h>
- <https://stackoverflow.com/questions/27827234/how-to-handle-the-onkeypress-event-in-reactjs>
- <https://dev.to/silvenleaf/fetch-api-easiest-explanation-part-4-4-delete-by-silvenleaf-4376>
- <https://go.dev/blog/routing-enhancements>
- <https://git-scm.com/docs/git-merge>
