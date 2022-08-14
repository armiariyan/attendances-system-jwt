# Documentation Attendances System with JWT
This is simple Restful API personal project

## Healthcheck 
Check the server

<b>GET</b>
```
https://attendances-system-jwt.herokuapp.com/api/v1/health
```
  
Response (Status: 200)
```
{
    "status": true,
    "message": "healthcheck successfull",
    "errors": null,
    "data": null
}
```

# Auth Routes
## Register User
Create new user

<b>POST</b>
```
https://attendances-system-jwt.herokuapp.com/api/v1/register
```

Request Body
```
{
    "name" : "Armia Riyan",
    "email" : "armiariyan@gmail.com",
    "password" : "password"
}
```

Response success (Status: 201 Created)
```
{
    "status": true,
    "message": "OK!",
    "errors": null,
    "data": {
        "id": 4,
        "name": "Armia Riyan",
        "email": "armiariyan@gmail.com",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNCIsInN0YW5kYXJkX2NsYWltcyI6eyJleHAiOjE2NjA1NzM3ODksImlhdCI6MTY2MDQ4NzM4OSwiaXNzIjoiYXJtaWFyaXlhbiJ9fQ.6GWyhjusqNBX2CrkHbRAiUYRVGQkw0Wp-cRXF9G5zf4"
    }
}
```

Response error (Status : 422 Unprocessable Entity) Validation for body error
```
{
    "status": false,
    "message": "register failed",
    "errors": [
        "Key: 'RegisterDTO.Email' Error:Field validation for 'Email' failed on the 'email' tag"
    ],
    "data": {}
}
```

## Login
Login using email and password

<b>POST</b>
```
https://attendances-system-jwt.herokuapp.com/api/v1/login
```

Request body
```
{
    "email" : "armiariyan@gmail.com",
    "password" : "password"
}
```

Response Success (Status: 200)
```
{
    "status": true,
    "message": "successfully logged in",
    "errors": null,
    "data": {
        "id": 4,
        "name": "Armia Riyan",
        "email": "armiariyan@gmail.com",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNCIsInN0YW5kYXJkX2NsYWltcyI6eyJleHAiOjE2NjA1NzUyNjksImlhdCI6MTY2MDQ4ODg2OSwiaXNzIjoiYXJtaWFyaXlhbiJ9fQ.eB2dWsTWBcOpYF-WdG_Z58TjVhGNCzXpWXua_-8KklY"
    }
}
```

Response error for invalid credential (Status: 401 Unauthorized)
```
{
    "status": false,
    "message": "failed to process request",
    "errors": [
        "invalid email or password"
    ],
    "data": {}
}
```


# Attendances Routes
Check in or Check out only can be done if user already login (provide valid jwt token)  

Attendance ID is ATD + 5 random alphanumeric string

## Check In
Create check in attendance

<b>POST</b>
```
https://attendances-system-jwt.herokuapp.com/api/v1/checkin
```

Headers
```
Authorization: JWT Token
```

Response success (status: 200)
```
{
    "status": true,
    "message": "successfully check in!",
    "errors": null,
    "data": {
        "id": "ATD-8AO58",
        "id_user": 4,
        "label": "check in",
        "date": "2022-08-14",
        "time": "15:06:54"
    }
}
```
Response if token not provide (status: 400)
```
{
    "status": false,
    "message": "failed to process request",
    "errors": [
        "no token found"
    ],
    "data": null
}
```
Response if token expired (status: 401)
```
{
    "status": false,
    "message": "token is not valid",
    "errors": [
        "token expired!"
    ],
    "data": null
}
```


## Check Out
Create check out attendance

<b>POST</b>
```
https://attendances-system-jwt.herokuapp.com/api/v1/checkout
```

Headers
```
Authorization: JWT Token
```

Response success (status: 200)
```
{
    "status": true,
    "message": "successfully check in!",
    "errors": null,
    "data": {
        "id": "ATD-8AO58",
        "id_user": 4,
        "label": "check in",
        "date": "2022-08-14",
        "time": "15:06:54"
    }
}
```
Response if have not check in today (status: 400)
```
{
    "status": false,
    "message": "failed to process request",
    "errors": [
        "please check in first"
    ],
    "data": {}
}
```

Refer another error to Check In docs above


## Get Attendances History
Get all attendances history by user id  

For this documentation, i already create couple of attendances before

<b>GET</b>
```
https://attendances-system-jwt.herokuapp.com/api/v1/activity/1?startDate=2022-08-01&endDate=2022-08-31
```

Headers
```
Authorization: JWT Token
```

Response success (Status: 200)
```
{
    "status": true,
    "message": "successfully get attendances history!",
    "errors": null,
    "data": [
        {
            "id": "ATD-9ZMGZ",
            "id_user": 4,
            "label": "check in",
            "date": "2022-08-14",
            "time": "15:41:55"
        },
        {
            "id": "ATD-GNI8y",
            "id_user": 4,
            "label": "check in",
            "date": "2022-08-14",
            "time": "15:41:49"
        },
        {
            "id": "ATD-Zc0TB",
            "id_user": 4,
            "label": "check out",
            "date": "2022-08-14",
            "time": "15:41:53"
        }
    ]
}
```

Response error empty (Status: 200)
```
{
    "status": false,
    "message": "failed to process request",
    "errors": [
        "attendances history is empty"
    ],
    "data": {}
}
```


# Activities Routes
Manage activity only can be done if user already login and check in  

Activity ID is ACT + 5 random alphanumeric string

## Create Activity
Create activity

<b>POST</b>
```
https://attendances-system-jwt.herokuapp.com/api/v1/activity
```

Headers
```
Authorization: JWT Token
```

Request Body
```
{
    "description" : "Test Activity 0"
}
```

Response success (Status: 200)
```
{
    "status": true,
    "message": "successfully created activity!",
    "errors": null,
    "data": {
        "id": "ACT-SbuxR",
        "id_user": 4,
        "description": "Test Activity 0",
        "date_created": "2022-08-14",
        "time_created": "15:25:37"
    }
}
```

## Update Activity
Update activity with activity_id

<b>PUT</b>
ACT-SbuxR is ID taken from Create Activity above
```
https://attendances-system-jwt.herokuapp.com/api/v1/activity/ACT-SbuxR
```

Headers
```
Authorization: JWT Token
```

Request Body
```
{
    "description" : "Update Test Activity 0"
}
```

Response success (Status: 200)
```
{
    "status": true,
    "message": "successfully update activity!",
    "errors": null,
    "data": {
        "id": "ACT-SbuxR",
        "id_user": 4,
        "description": "Update Test Activity 0",
        "date_created": "2022-08-14",
        "time_created": "15:25:37"
    }
}
```

## Delete Activity
Delete activity with activity_id

<b>DEL</b>
ACT-SbuxR is ID taken from Create Activity above
```
https://attendances-system-jwt.herokuapp.com/api/v1/activity/ACT-SbuxR
```

Headers
```
Authorization: JWT Token
```

Response success (Status: 200)
```
{
    "status": true,
    "message": "activity deleted!",
    "errors": null,
    "data": {}
}
```

## Get Activity History By Date
Get all activity history between two dates.  

For this documentation, i already create couple of activity before

<b>GET</b>
```
https://attendances-system-jwt.herokuapp.com/api/v1/activity/1?startDate=2022-08-01&endDate=2022-08-31
```

Headers
```
Authorization: JWT Token
```

Params
```
startDate: valid date with format (YYYY-MM-DD) example from url above --> 2022-08-01
endDate: valid date with format (YYYY-MM-DD) example from url above --> 2022-08-31
```


Response success (Status: 200)
```
{
    "status": true,
    "message": "successfully get activity history!",
    "errors": null,
    "data": [
        {
            "id": "ACT-1jDHR",
            "id_user": 4,
            "description": "Test Activity 2",
            "date_created": "2022-08-14",
            "time_created": "15:34:12"
        },
        {
            "id": "ACT-bbZrs",
            "id_user": 4,
            "description": "Test Activity 1",
            "date_created": "2022-08-14",
            "time_created": "15:34:10"
        }
    ]
}
```

Response error empty (Status: 200)
```
{
    "status": false,
    "message": "failed to process request",
    "errors": [
        "activities in that range date is empty"
    ],
    "data": {}
}
```
