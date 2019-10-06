# The Ping The Pong
this app is a simulated ping pong web game in the best Hattrick fashion. It exposes a REST API, to store known players on database and simulate an exciting ping pong match.


## How to run
Clone the project, put on the project root folder and run:

```zsh
go run main.go
```

## How to consume
The application provide the following APIs

* Players
  
  * Get all players
  
    ```
    curl -X GET http://localhost:8287/players
    ```

  * Get a player with a given Id
  
    ```
    curl -X GET http://localhost:8287/players/{playerid}
    ```

  * Create a player
    
    Here you are required to generate the token through SignIn capability.

    ```
    curl -d '{"names":"Fan Zhendong", "wins":10, "losses": 2}' -H "Content-Type: application/json" -H "Authorization: Bearer ${TOKEN}" -X POST http://localhost:8287/players
    ```
  
* Sign in
  
  To sign in you must provide user name and password. The current users are :

  * **user1**:password1
  * **user2**:password2
  
  The capability generates a token and return it as a cookie.

  ```
  curl -d '{"username":"user1", "password":"password1"}' -H "Content-Type: application/json" -X POST http://localhost:8287/signin
  ```

* Play
  
  To play a match between two players you have to sign in and consume the match API as follows :

  ```
  curl -d '{"player1ID":"", "player2ID":""}' -H "Content-Type: application/json" -H "Authorization: Bearer ${TOKEN}" -X POST http://localhost:8287/matches
  ```

## HTTP Client
In the root of the project was added a **insonmina** script to consume the API 

## Thirdparty libraries

* [Viper](https://github.com/spf13/viper) for configuration purposes.
* [Gorilla](https://github.com/gorilla/mux) to take advantage of its powerful router.
* [Logrus](https://github.com/sirupsen/logrus) for logging mechanism.
