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
GROUP_ADDING=$(printf '{
    "username": {"name": "Paperoga"},
    "group": %s
}' "$groupid"| jq -c .)
printf "FAILED GROUP CHAT ADDING: PAPERONE TRIES ADDING PAPEROGA\n"
printf "${GROUP_ADDING}\n"
curl -s --header "Content-Type: application/json" --header "Authentication: ${paperoneauth}" --request PUT --data $GROUP_ADDING http://0.0.0.0:3000/groups -v
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
GROUP_ADDING=$(printf '{
    "username": {"name": "Paperone"},
    "group": %s
}' "$groupid"| jq -c .)
printf "GROUP CHAT ADDING: PIPPO ADDS PAPERONE\n"
printf "${GROUP_ADDING}\n"
curl -s --header "Content-Type: application/json" --header "Authentication: ${pippoauth}" --request PUT --data $GROUP_ADDING http://0.0.0.0:3000/groups -v
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
GROUP_ADDING=$(printf '{
    "username": {"name": "Paperoga"},
    "group": %s
}' "$groupid"| jq -c .)
printf "GROUP CHAT ADDING: PAPERONE ADDS PAPEROGA\n"
printf "${GROUP_ADDING}\n"
curl -s --header "Content-Type: application/json" --header "Authentication: ${paperoneauth}" --request PUT --data $GROUP_ADDING http://0.0.0.0:3000/groups -v
printf "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"