#!/bin/bash

printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n" 
printf "USERS CREATION\n"
auth="$(curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Pippo"}' http://0.0.0.0:3000/session | jq '.identifier')"
auth=${auth//\"}
curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Topolino"}' http://0.0.0.0:3000/session
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "PRIVATE CHAT CREATION\n"
curl -s --header "Content-Type: application/json" --header "Authentication: ${auth}" --request PUT --data '{"isgroup":false, "members":[{"name":"Topolino"}], "groupname":null}' http://0.0.0.0:3000/conversations
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT CREATION\n"
curl -s --header "Content-Type: application/json" --header "Authentication: ${auth}" --request PUT --data '{"isgroup":true, "members":[{"name":"Topolino"}], "groupname": "Eccehomo"}' http://0.0.0.0:3000/conversations
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "PRIVATE CHAT CREATION: ALREADY CREATED\n"
curl -s --header "Content-Type: application/json" --header "Authentication: ${auth}" --request PUT --data '{"isgroup":false, "members":[{"name":"Topolino"}], "groupname":null}' http://0.0.0.0:3000/conversations
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"