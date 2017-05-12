
# My Cool API
My API usually works as expected.

Table of Contents

1. [login and get token](#sign-in)
1. [get all users](#users)

<a name="sign-in"></a>

## sign-in

| Specification | Value |
|-----|-----|
| Resource Path | /sign-in |
| API Version | 1.0.0 |
| BasePath for the API | http://localhost:3000 |
| Consumes | application/json |
| Produces |  |



### Operations


| Resource Path | Operation | Description |
|-----|-----|-----|
| /sign-in | [POST](#signIn) | login and get token |



<a name="signIn"></a>

#### API: /sign-in (POST)


login and get token



| Param Name | Param Type | Data Type | Description | Required? |
|-----|-----|-----|-----|-----|
| username | query | string | username should not less than 6 characters |  |
| password | query | string | password should not less than 6 characters |  |


| Code | Type | Model | Message |
|-----|-----|-----|-----|
| 200 | object | [LoginForm](#bitbucket.org.huayun321.test-heroku.LoginForm) |  |
| 400 | object | [ErrorResp](#bitbucket.org.huayun321.test-heroku.ErrorResp) | Customer ID must be specified |




### Models

<a name="bitbucket.org.huayun321.test-heroku.ErrorResp"></a>

#### ErrorResp

| Field Name (alphabetical) | Field Type | Description |
|-----|-----|-----|
| Code | int |  |
| Msg | string |  |

<a name="bitbucket.org.huayun321.test-heroku.LoginForm"></a>

#### LoginForm

| Field Name (alphabetical) | Field Type | Description |
|-----|-----|-----|
| password | string |  |
| username | string |  |


<a name="users"></a>

## users

| Specification | Value |
|-----|-----|
| Resource Path | /users |
| API Version | 1.0.0 |
| BasePath for the API | http://localhost:3000 |
| Consumes |  |
| Produces |  |



### Operations


| Resource Path | Operation | Description |
|-----|-----|-----|
| /users | [GET](#get users) | get all users |



<a name="get users"></a>

#### API: /users (GET)


get all users



