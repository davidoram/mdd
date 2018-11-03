#!/usr/bin/env bats
#
# Test script for 'mdd ls' command
#

@test "mdd ls, no files" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd ls
  [ "$status" -eq 0 ]
}

@test "mdd ls, 1 file" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  $BATS_CWD/mdd new adr
  run $BATS_CWD/mdd ls
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "adr-b7-0001.md   Architecture Decision Record  " ]
}

@test "mdd ls, 1 file" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  $BATS_CWD/mdd new adr
  run $BATS_CWD/mdd ls
  [ "$status" -eq 0 ]
  [ $(expr "$output" : "adr.*md.*Architecture Decision Record") -ne 0 ]
}

@test "mdd ls, 10 files" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  for i in {1..10}; do
    $BATS_CWD/mdd new adr
  done
  run $BATS_CWD/mdd ls
  [ "$status" -eq 0 ]
  [ ${#lines[@]} -eq 10 ]
}
