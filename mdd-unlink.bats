#!/usr/bin/env bats
#
# Test script for 'mdd unlink' command
#

setup() {
  rm -rf ./tmp/.mdd
  rm -rf ./.mdd
}


@test "mdd unlink missing both args" {
  run $BATS_CWD/mdd unlink
  [ "$status" -eq 1 ]
}

@test "mdd unlink missing child arg" {
  run $BATS_CWD/mdd unlink parent
  [ "$status" -eq 1 ]
}

@test "mdd unlink, missing project" {
  run $BATS_CWD/mdd unlink parent child
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "No project found" ]
}

@test "mdd unlink, missing parent" {
  $BATS_CWD/mdd init
  child=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd unlink parent ${child}
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Cant find parent 'parent.md'" ]
}

@test "mdd unlink, missing child" {
  $BATS_CWD/mdd init
  parent=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd unlink ${parent} child
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Cant find child 'child.md'" ]
}

@test "mdd unlink, to self" {
  $BATS_CWD/mdd init
  parent=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd unlink ${parent} ${parent}
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Cant unlink from self" ]
}

@test "mdd unlink, duplicate is a no-op" {
  $BATS_CWD/mdd init
  parent=$(basename $($BATS_CWD/mdd new adr))
  child=$(basename $($BATS_CWD/mdd new adr))
  $BATS_CWD/mdd link ${parent} ${child}
  $BATS_CWD/mdd unlink ${parent} ${child}
  run $BATS_CWD/mdd unlink ${parent} ${child}
  [ "$status" -eq 0 ]
}

@test "mdd unlink, valid" {
  $BATS_CWD/mdd init
  parent=$(basename $($BATS_CWD/mdd new adr))
  child=$(basename $($BATS_CWD/mdd new adr))
  $BATS_CWD/mdd link ${parent} ${child}
  run $BATS_CWD/mdd unlink ${parent} ${child}
  [ "$status" -eq 0 ]
}

@test "mdd unlink, updates metadata" {
  $BATS_CWD/mdd init
  parent_path=$($BATS_CWD/mdd new adr)
  parent=$(basename ${parent_path})
  child=$(basename $($BATS_CWD/mdd new adr))
  $BATS_CWD/mdd link ${parent} ${child}
  run grep "mdd-child: ${child}" $parent_path
  [ "$status" -eq 0 ]
  $BATS_CWD/mdd unlink ${parent} ${child}
  run grep "mdd-child: ${child}" $parent_path
  [ "$status" -eq 1 ]
}

@test "mdd unlink, many files" {
  $BATS_CWD/mdd init
  children=()
  for i in {1..10}; do
    children+=($(basename $($BATS_CWD/mdd new adr)))
  done
  parent_path=$($BATS_CWD/mdd new adr)
  parent=$(basename ${parent_path})
  # Link all
  for child in ${children[@]}; do
    run $BATS_CWD/mdd link ${parent} ${child}
    [ "$status" -eq 0 ]
  done

  # Count number of mdd-child metadata entries
  run grep -c "mdd-child:" $parent_path
  [ "${lines[0]}" = "10" ]

  # Unlink all
  for child in ${children[@]}; do
    run $BATS_CWD/mdd unlink ${parent} ${child}
    [ "$status" -eq 0 ]
  done

  # Count number of mdd-child metadata entries
  run grep -c "mdd-child:" $parent_path
  [ "${lines[0]}" = "0" ]
}