#!/usr/bin/env bats
#
# Test script for 'mdd tag' command
#

setup() {
  rm -rf ./tmp/.mdd
  rm -rf ./.mdd
}

@test "mdd tag missing both args" {
  run $BATS_CWD/mdd tag
  [ "$status" -eq 1 ]
}

@test "mdd tag missing tag arg" {
  run $BATS_CWD/mdd tag file
  [ "$status" -eq 1 ]
}

@test "mdd tag, file doest exist" {
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd tag file tag
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Cant find document 'file.md'" ]
}

@test "mdd tag, invalid tag characters" {
  $BATS_CWD/mdd init
  file=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd tag ${file} '^%%$#@'
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Tags must be 3-20 chars long, made up of the following characters: '0-9A-Za-z_-'" ]
}

@test "mdd tag, invalid tag too short" {
  $BATS_CWD/mdd init
  file=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd tag ${file} ab
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Tags must be 3-20 chars long, made up of the following characters: '0-9A-Za-z_-'" ]
}

@test "mdd tag, invalid tag too long" {
  $BATS_CWD/mdd init
  file=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd tag ${file} abcefgbhighlomopqrstu
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Tags must be 3-20 chars long, made up of the following characters: '0-9A-Za-z_-'" ]
}

@test "mdd tag, valid" {
  $BATS_CWD/mdd init
  file=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd tag ${file} abcefgbhighlomopqrst abc
  [ "$status" -eq 0 ]
}

@test "mdd tag, duplicate ignored" {
  $BATS_CWD/mdd init
  file=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd tag ${file} abc abc abc
  [ "$status" -eq 0 ]
  run $BATS_CWD/mdd ls -l
  [ "$status" -eq 0 ]
  [ $(expr "${lines[0]}" : "^${file}.*#abc *$") -ne 0 ]
}

@test "mdd tag, updates metadata" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  file=$(basename ${file_path})
  $BATS_CWD/mdd tag ${file} tag-1 tag-2
  # Test for metadata in the file
  run grep "mdd-tag: tag-1" ${file_path}
  [ "$status" -eq 0 ]
  run grep "mdd-tag: tag-2" ${file_path}
  [ "$status" -eq 0 ]
}

@test "mdd tag, many tags" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  file=$(basename ${file_path})
  for suffix in {1..100}; do
    run $BATS_CWD/mdd tag ${file} tag-${suffix}
    [ "$status" -eq 0 ]
  done
  # Check all links saved
  for suffix in {1..100}; do
    run grep "mdd-tag: tag-${suffix}" ${file_path}
    [ "$status" -eq 0 ]
  done
}