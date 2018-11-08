#!/usr/bin/env bats
#
# Test script for 'mdd verify' command
#

setup() {
  rm -rf ./tmp/.mdd
  rm -rf ./.mdd
}


@test "mdd verify, missing project" {
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "No project found" ]
}

@test "mdd verify, empty project valid" {
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd verify
  [ "$status" -eq 0 ]
}

@test "mdd verify, remove project.data file" {
  $BATS_CWD/mdd init
  rm ./.mdd/project.data
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
}

