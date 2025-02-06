#!/bin/bash

printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n" 
printf "USER CREATION:\n"
auth="$(curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Pippo"}' http://0.0.0.0:3000/session | jq '.identifier')"
auth=${auth//\"}
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "AUTHENTICATION PRINTING:\n"
printf $auth
printf "\n++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"