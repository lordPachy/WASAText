#!/bin/bash

printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n" 
printf "USERS CREATION\n"
pippoauth="$(curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Pippo"}' http://0.0.0.0:3000/session | jq '.identifier')"
pippoauth=${pippoauth//\"}
curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Topolino"}' http://0.0.0.0:3000/session
curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Titti"}' http://0.0.0.0:3000/session
curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Minni"}' http://0.0.0.0:3000/session
printf "$pippoauth\n"
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "USERS QUERYING\n"
curl -s --header "Content-Type: application/json" --header "Authorization: $pippoauth" --request GET --data '{"name":""}' http://0.0.0.0:3000/users -v
curl -s --header "Content-Type: application/json" --header "Authorization: $pippoauth" --request GET --data '{"name":"Minni"}' http://0.0.0.0:3000/users -v
curl -s --header "Content-Type: application/json" --header "Authorization: $pippoauth" --request GET --data '{"name":"T"}' http://0.0.0.0:3000/users -v
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"