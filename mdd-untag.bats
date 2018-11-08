#!/usr/bin/env bats
#
# Test script for 'mdd untag' command
#

setup() {
  rm -rf ./tmp/.mdd
  rm -rf ./.mdd
}


@test "mdd untag missing both args" {
  run $BATS_CWD/mdd untag
  [ "$status" -eq 1 ]
}

@test "mdd untag missing tag arg" {
  run $BATS_CWD/mdd untag document
  [ "$status" -eq 1 ]
}

@test "mdd untag, missing project" {
  run $BATS_CWD/mdd untag document tag
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "No project found" ]
}

@test "mdd untag, file doesnt exist" {
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd untag document tag
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Cant find document 'document.md'" ]
}

@test "mdd untag, tag doesnt exist is ignored" {
  $BATS_CWD/mdd init
  document=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd untag ${document} tag
  [ "$status" -eq 0 ]
}

@test "mdd untag, valid" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  document=$(basename $file_path)
  $BATS_CWD/mdd tag ${document} tag-1
  run $BATS_CWD/mdd untag ${document} tag-1
  [ "$status" -eq 0 ]
}

@test "mdd untag, updates metadata" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  document=$(basename $file_path)
  $BATS_CWD/mdd tag ${document} tag-1

  # Test for metadata added
  run grep "mdd-tag: tag-1" ${file_path}
  [ "$status" -eq 0 ]

  $BATS_CWD/mdd untag ${document} tag-1

  # Test for metadata removed
  run grep "mdd-tag: tag-1" ${file_path}
  [ "$status" -eq 1 ]
}

@test "mdd untag, many tags" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  document=$(basename $file_path)

  for i in {1..10}; do
    $BATS_CWD/mdd tag ${document} tag-${i}
  done

  # Count number of mdd-tag metadata entries
  run grep -c "mdd-tag:" $file_path
  [ "${lines[0]}" = "10" ]

  # Untag all
  for i in {1..10}; do
    $BATS_CWD/mdd untag ${document} tag-${i}
  done

  # Count number of mdd-tag metadata entries
  run grep -c "mdd-tag:" $file_path
  [ "${lines[0]}" = "0" ]
}