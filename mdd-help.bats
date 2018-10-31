#!/usr/bin/env bats
#
# Test script for 'mdd init' command
#

@test "mdd help" {
  run $BATS_CWD/mdd help
  [ "$status" -eq 0 ]
}

@test "mdd help init" {
  run $BATS_CWD/mdd help init
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd init creates a new mdd document repository" ]
}
