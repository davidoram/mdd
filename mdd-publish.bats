#!/usr/bin/env bats
#
# Test script for 'mdd publish' command
#

setup() {
  rm -rf ./tmp/.mdd
  rm -rf ./.mdd
}

@test "mdd publish, missing project" {
  run $BATS_CWD/mdd publish
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "No project found" ]
}

@test "mdd publish, empty project valid" {
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd publish
  [ "$status" -eq 0 ]
}

@test "mdd publish, converts md to html" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  file=$(basename ${file_path})
  run $BATS_CWD/mdd publish
  [ "$status" -eq 0 ]
  file_no_suffix=$(basename ${file_path} .md)
  run ls ./.mdd/publish/${file_no_suffix}.html
  [ "$status" -eq 0 ]
  run ls ./.mdd/publish/index.html
  [ "$status" -eq 0 ]
}

