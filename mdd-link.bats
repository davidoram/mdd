#!/usr/bin/env bats
#
# Test script for 'mdd link' command
#

setup() {
  rm -rf ./tmp/.mdd
  rm -rf ./.mdd
}

# teardown() {
#   rm -rf ./tmp/.mdd
#   rm -rf ./.mdd
# }

@test "mdd link missing both args" {
  run $BATS_CWD/mdd link
  [ "$status" -eq 1 ]
}

@test "mdd link missing child arg" {
  run $BATS_CWD/mdd link parent
  [ "$status" -eq 1 ]
}

@test "mdd link, missing project" {
  run $BATS_CWD/mdd link parent child
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "No project found" ]
}

@test "mdd link, missing parent" {
  $BATS_CWD/mdd init
  child=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd link parent ${child}
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Cant find parent 'parent.md'" ]
}

@test "mdd link, missing child" {
  $BATS_CWD/mdd init
  parent=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd link ${parent} child
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Cant find child 'child.md'" ]
}

@test "mdd link, to self" {
  $BATS_CWD/mdd init
  parent=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd link ${parent} ${parent}
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Cant link to self" ]
}

@test "mdd link, valid" {
  $BATS_CWD/mdd init
  parent=$(basename $($BATS_CWD/mdd new adr))
  child=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd link ${parent} ${child}
  [ "$status" -eq 0 ]
}
