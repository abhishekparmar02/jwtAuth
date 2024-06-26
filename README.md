#To start the API server
  go run main.go handlers.go

#SignUp
    curl -v -X POST -d '{"Username":"user0","Password":"password0"}' http://localhost:8080/signup

    Note: If trying to create user with previously used username then 400 BadRequest response will be received Note: The usernames user1 and user2 are already stored in the file store

#Login
    curl -v -X POST -d '{"Username":"user0","Password":"password0"}' http://localhost:8080/login

    cookie will be received (lets call it received_cookie)

#Authorization
    curl -v -X GET --cookie "received_cookie" http://localhost:8080/home

    Note: The expiry time has been set for 5 mins from the time the token was generated.

#Revoke
    curl -v -X POST --cookie "received_cookie" http://localhost:8080/logout

#Refresh
    curl -v -X POST --cookie "received_cookie" http://localhost:8080/refresh
