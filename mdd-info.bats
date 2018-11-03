#!/usr/bin/env bats
#
# Test script for 'mdd info' command
#

setup() {
  rm -rf ./tmp/.mdd
  rm -rf ./.mdd
}

teardown() {
  rm -rf ./tmp/.mdd
  rm -rf ./.mdd
}

@test "mdd info, no project" {
  run $BATS_CWD/mdd info
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "No projects found" ]
}

@test "mdd info, no files" {
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd info
  [ "$status" -eq 0 ]
}

@test "mdd info, 1 file" {
  $BATS_CWD/mdd init
  $BATS_CWD/mdd new adr
  run $BATS_CWD/mdd info
  [ "$status" -eq 0 ]
  [ "${lines[2]}" = "path      : .mdd" ]
  [ "${lines[3]}" = "templates : 6" ]
  [ "${lines[4]}" = "documents : 1" ]
}

@test "mdd info, different path" {
  rm -rf ./tmp/.mdd
  mkdir -p tmp
  $BATS_CWD/mdd init -o tmp
  $BATS_CWD/mdd new adr
  run $BATS_CWD/mdd info
  [ "$status" -eq 0 ]
  [ "${lines[2]}" = "path      : tmp/.mdd" ]
  [ "${lines[3]}" = "templates : 6" ]
  [ "${lines[4]}" = "documents : 1" ]
}

@test "mdd info, 10 files" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  for i in {1..10}; do
    $BATS_CWD/mdd new adr
  done
  run $BATS_CWD/mdd info
  [ "$status" -eq 0 ]
  [ "${lines[2]}" = "path      : .mdd" ]
  [ "${lines[3]}" = "templates : 6" ]
  [ "${lines[4]}" = "documents : 10" ]
}
