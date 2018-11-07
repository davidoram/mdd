#!/usr/bin/env bats
#
# Test script for 'mdd new' command
#

setup() {
  rm -rf ./tmp/.mdd
  rm -rf ./.mdd
}

@test "mdd new, missing template arg" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd new
  [ "$status" -eq 1 ]
}

@test "mdd new, invalid template arg" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd new not-a-template
  [ "$status" -eq 1 ]
}

@test "mdd new, no title" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd new adr
  [ "$status" -eq 0 ]
}

@test "mdd new, no title, file created ok" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  new_file=$( $BATS_CWD/mdd new adr )
  run ls -1 ${new_file}
  [ "$status" -eq 0 ]
}

@test "mdd new, with title, file created ok" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  new_file=$( $BATS_CWD/mdd new adr 'Important architectural record')
  run ls -1 ${new_file}
  [ "$status" -eq 0 ]
}

@test "mdd new, many files" {
  $BATS_CWD/mdd init
  for i in {1..100}; do
    $BATS_CWD/mdd new adr
  done

  # Count documents created
  file_count=$( ls ./.mdd/documents/* | wc -l)
  [ $(expr "${file_count}" : "^ *100") -ne 0 ]
}
