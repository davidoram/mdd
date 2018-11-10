#!/usr/bin/env bats
#
# Test script for 'mdd rm' command
#

setup() {
  rm -rf ./tmp/.mdd
  rm -rf ./.mdd
}


@test "mdd rm missing both args" {
  run $BATS_CWD/mdd rm
  [ "$status" -eq 1 ]
}

@test "mdd rm, missing project" {
  run $BATS_CWD/mdd rm file
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "No project found" ]
}

@test "mdd rm, no such file" {
  $BATS_CWD/mdd init
  $BATS_CWD/mdd new adr
  run $BATS_CWD/mdd rm no-such-file
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "No such file: 'no-such-file'" ]
}


@test "mdd rm, valid unlinked file" {
  $BATS_CWD/mdd init
  parent=$(basename $($BATS_CWD/mdd new adr))
  child=$(basename $($BATS_CWD/mdd new adr))

  # Delete the parent
  run ls ./.mdd/documents/${parent}
  [ "$status" -eq 0 ]
  run $BATS_CWD/mdd rm ${parent}
  [ "$status" -eq 0 ]
  run ls ./.mdd/documents/${parent}
  [ "$status" -eq 1 ]

  # Delete the child
  run ls ./.mdd/documents/${child}
  [ "$status" -eq 0 ]
  run $BATS_CWD/mdd rm ${child}
  [ "$status" -eq 0 ]
  run ls ./.mdd/documents/${child}
  [ "$status" -eq 1 ]
}

@test "mdd rm, valid linked file" {
  $BATS_CWD/mdd init
  parent_path=$($BATS_CWD/mdd new adr)
  parent=$(basename ${parent_path})
  child=$(basename $($BATS_CWD/mdd new adr))
  run $BATS_CWD/mdd link ${parent} ${child}
  [ "$status" -eq 0 ]

  # Check metadata exists
  run grep "mdd-child: ${child}" $parent_path
  [ "$status" -eq 0 ]

  # Delete the child
  run $BATS_CWD/mdd rm ${child}
  [ "$status" -eq 0 ]

  # Check metadata updated
  run grep "mdd-child: ${child}" $parent_path
  [ "$status" -eq 1 ]
}

@test "mdd rm, many files" {
  $BATS_CWD/mdd init
  file_paths=()
  files=()
  for i in {1..10}; do
    fpath=$($BATS_CWD/mdd new adr)
    file_pathss+=(${fpath})
    files+=($(basename ${fpath}))
  done

  # Check all files created
  run $BATS_CWD/mdd ls -1
  [ "$status" -eq 0 ]
  [ "${#lines[@]}" -eq 10 ]

  for f in ${files[@]}; do
    run $BATS_CWD/mdd rm ${f}
    [ "$status" -eq 0 ]
  done

  # Check all files deleted
  run $BATS_CWD/mdd ls
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "" ]
}