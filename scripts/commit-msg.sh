#!/usr/bin/bash

MIN_LENGTH=10
MAX_LENGTH=52
TYPES="feat|fix|docs|style|refactor|perf|test|chore|build|ci|revert"
PATTERN='^('"$TYPES"')(\([a-zA-Z0-9._-]+\))?!?: .{'"$MIN_LENGTH"','"$MAX_LENGTH"'}$'

commit_msg=$(head -1 "$1")

if [[ ! $commit_msg =~ $PATTERN ]]; then
  echo -e "\n\e[1m\e[31m[INVALID COMMIT MESSAGE]"
  echo -e "------------------------\033[0m\e[0m"
  echo -e "\e[1mValid types:\e[0m \e[34m${TYPES}\033[0m"
  echo -e "\e[1mMin length (first line):\e[0m \e[34m$MIN_LENGTH\033[0m"
  echo -e "\e[1mMax length (first line):\e[0m \e[34m$MAX_LENGTH\033[0m\n"
  exit 1
fi
