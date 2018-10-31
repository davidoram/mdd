#!/usr/bin/env bats
#
# Test script for 'mdd init' command
#

@test "mdd init" {
  rm -rf ./.mdd
  run $BATS_CWD/mdd init
  [ "$status" -eq 0 ]
}

@test "mdd init, creates correct directory structure" {
  rm -rf ./.mdd
  run $BATS_CWD/mdd init
  run ls -a -1 ./.mdd
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "." ]
  [ "${lines[1]}" = ".." ]
  [ "${lines[2]}" = "documents" ]
  [ "${lines[3]}" = "project.data" ]
  [ "${lines[4]}" = "publish" ]
  [ "${lines[5]}" = "templates" ]
}

@test "mdd init -p, saves project meta-data" {
  rm -rf ./.mdd
  run $BATS_CWD/mdd init -p my-project
  run cat ./.mdd/project.data
  [ "${lines[1]}" = "project: my-project" ]
}

@test "mdd init -o, directory must exist" {
  rm -rf ./foo
  run $BATS_CWD/mdd init -o $BATS_CWD/foo
  [ "$status" -eq 1 ]
}

@test "mdd init -o, saves to a different directory" {
  rm -rf ./foo
  mkdir -p ./foo
  run $BATS_CWD/mdd init -o $BATS_CWD/foo
  run ls -a -1 ./foo/.mdd
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "." ]
  [ "${lines[1]}" = ".." ]
  [ "${lines[2]}" = "documents" ]
  [ "${lines[3]}" = "project.data" ]
  [ "${lines[4]}" = "publish" ]
  [ "${lines[5]}" = "templates" ]
}

