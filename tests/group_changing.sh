#!/bin/bash

printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n" 
printf "USERS CREATION\n"
pippoauth="$(curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Pippo"}' http://0.0.0.0:3000/session | jq '.identifier')"
pippoauth=${pippoauth//\"}
curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Topolino"}' http://0.0.0.0:3000/session
paperoneauth="$(curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Paperone"}' http://0.0.0.0:3000/session | jq '.identifier')"
paperoneauth=${paperoneauth//\"}
curl -s --header "Content-Type: application/json" --request PUT --data '{"name":"Paperoga"}' http://0.0.0.0:3000/session
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT CREATION BY PIPPO: THERE ARE PIPPO AND TOPOLINO\n"
groupid="$(curl -s --header "Content-Type: application/json" --header "Authentication: ${pippoauth}" --request PUT --data '{"isgroup":true, "members":[{"name":"Topolino"}], "groupname": "Eccehomo"}' http://0.0.0.0:3000/conversations)"
printf $groupid
printf "\n++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "GROUP CHAT NAME CHANGING BY PIPPO\n"
groupidnum=$(jq '.id' <<< $groupid) 
printf "Group id num is: $groupidnum\n"
curl -s --header "Content-Type: application/json" --header "Authentication: ${pippoauth}" --request PUT --data '{"value":"Eccepippo"}' http://0.0.0.0:3000/conversations/$groupidnum/settings/groupname
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
printf "CONVERSATIONS RETRIEVAL FROM PIPPO \n"
curl -s --header "Content-Type: application/json" --header "Authentication: ${pippoauth}" --request GET http://0.0.0.0:3000/conversations
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
