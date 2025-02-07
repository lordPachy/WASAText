#!/bin/bash

printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n" 
printf "USERS CREATION\n"
pippoauth="$(curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Pippo"}' http://0.0.0.0:3000/session | jq '.identifier')"
pippoauth=${pippoauth//\"}
topolinoauth="$(curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Topolino"}' http://0.0.0.0:3000/session | jq '.identifier')"
topolinoauth=${topolinoauth//\"}
printf "$pippoauth\n$topolinoauth\n"
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT CREATION\n"
groupid="$(curl -s --header "Content-Type: application/json" --header "Authorization: ${pippoauth}" --request PUT --data '{"isgroup":true, "members":[{"name":"Topolino"}], "groupname": "Eccehomo"}' http://0.0.0.0:3000/conversations | jq '.id')"
printf "$groupid\n"
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT MESSAGE 1\n"
curl -s --header "Content-Type: application/json" --header "Authorization: ${pippoauth}" --request POST --data '{"content":"Hi, Im Pippo", "replyingto":-1}' http://0.0.0.0:3000/conversations/$groupid
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT MESSAGE 2\n"
curl -s --header "Content-Type: application/json" --header "Authorization: ${topolinoauth}" --request POST --data '{"content":"Hi, Im Topolino", "replyingto":-1}' http://0.0.0.0:3000/conversations/$groupid
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT RETRIEVAL FROM TOPOLINO\n"
curl -s --header "Content-Type: application/json" --header "Authorization: ${topolinoauth}" --request GET http://0.0.0.0:3000/conversations/$groupid
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT RETRIEVAL FROM PIPPO \n"
curl -s --header "Content-Type: application/json" --header "Authorization: ${pippoauth}" --request GET http://0.0.0.0:3000/conversations/$groupid

# Saving Pippo's message
pippomessid="$(curl -s --header "Content-Type: application/json" --header "Authorization: ${pippoauth}" --request GET http://0.0.0.0:3000/conversations/$groupid | jq -r '.messages[] | select(.username=="Pippo") | .messageid')"

# Saving Topolino's message
topolinomessid="$(curl -s --header "Content-Type: application/json" --header "Authorization: ${pippoauth}" --request GET http://0.0.0.0:3000/conversations/$groupid | jq -r '.messages[] | select(.username=="Topolino") | .messageid')"
printf "Retrieved Pippo message id : $pippomessid\n"
printf "Retrieved Topolino message id : $topolinomessid\n"
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "FAILED MESSAGE DELETING FROM PIPPO \n"
curl -s --header "Content-Type: application/json" --header "Authorization: ${pippoauth}" --request DELETE http://0.0.0.0:3000/conversations/$groupid/messages/$topolinomessid -v
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "MESSAGE DELETING FROM PIPPO \n"
curl -s --header "Content-Type: application/json" --header "Authorization: ${pippoauth}" --request DELETE http://0.0.0.0:3000/conversations/$groupid/messages/$pippomessid -v
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT COMMMENT TO 2\n"
commentid="$(curl -s --header "Content-Type: application/json" --header "Authorization: ${topolinoauth}" --request PUT --data '{"reaction":"thumbs_up"}' http://0.0.0.0:3000/conversations/$groupid/messages/$topolinomessid | jq -r '.commentid')"
printf "Comment id is $commentid\n"
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT RETRIEVAL FROM PIPPO \n"
curl -s --header "Content-Type: application/json" --header "Authorization: ${pippoauth}" --request GET http://0.0.0.0:3000/conversations/$groupid
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT COMMMENT DELETING TO 2\n"
curl -s --header "Content-Type: application/json" --header "Authorization: ${topolinoauth}" --request DELETE http://0.0.0.0:3000/conversations/$groupid/messages/$topolinomessid/comments/$commentid -v
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT RETRIEVAL FROM PIPPO \n"
curl -s --header "Content-Type: application/json" --header "Authorization: ${pippoauth}" --request GET http://0.0.0.0:3000/conversations/$groupid
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
