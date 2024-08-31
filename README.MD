# Reader blog
This is simple pure REST API backend application in GOLANG using beautiful framework a.k.a. Fiber + GORM. Implemented all the CRUD operations for `User` model. Hope, that my resource will be useful for you.

## Stack
1. Golang -> Fiber
2. PostgreSQL -> via GORM postgres driver

## Requirements
 - You must have Docker installed in your computer.
 - Create .env file in the root of project with following fields:
  1. DB_HOST=
  2. DB_NAME=
  3. DB_USER=
  4. DB_PASSWORD=
  5. DB_PORT=
  6. JWT_SECRET=
  7. ADMIN_ROLE=admin
  8. WRITER_ROLE=writer

## Routes
1. GET `api/` -> Should get a message "API is up and running"
2. REGISTERING A NEW USER
POST `api/register`
```JSON
{
    "username": "exampleuser",
    "email": "example@mail.com",
    "password": "somePassword",
    "password_confirmation": "somePassword"
}
```
> Output: You will get the user that registered just, and the JWT token!
3. Login route
POST `api/login`
```JSON
{
    "username": "exampleuser",
    "password": "somePassword",
    "password_confirmation": "somePassword"
}
```
> Output: In success case you get the JWT token.
4. Profile page
GET `api/profile`
> Requires authentication in Header section add the following line:
`Authorization: Bearer your_jwt_token`

5. To fetch all users 
GET `api/users`
> Requires authorization: Role => Admin. Output: you get the list of existing users

6. To fetch all soft deleted users
GET `api/users?deleted=true`
> Requires authorization: Role => Admin. Output: you get the list of existing users

7. To get one specified user by id
GET `api/users/{:userId}`
> Requires authorization: Role => Admin. Output: you get the user if is exists

8. To create a new user
POST `api/users`
> Requires authorization: Role => Admin. Output: you get the user just created
```JSON
{
    "username": "exampleuser",
    "email": "usermail@mail.com",
    "password": "somePassword",
    "password_confirmation": "somePassword"
}
```

9. To update user
> Requires authorization: Role => Admin. Output: you get the user just updated
> You can change specified fields if you want, or all fields such as: email, password, role
```JSON
{
    "email": "updatedemail@mail.com",
    "password": "updatedPassword01",
    "password_confirmation": "updatedPassword01",
    "role": "writer"
}
```

10. Deleting the user
> Requires authorization: Role => Admin.
> Here you have two options like: Soft deleting and Force deleting
- In order to softly delete you have to perform request like:
DELETE `api/users/{:userId}`
- In order to force delete you have to perform following request:
DELETE `api/users/{:userId}?force=true`

11. You can restore the soft deleted user by:
> Requires authorization: Role => Admin.
PUT `api/users/{:userId}/restore`
