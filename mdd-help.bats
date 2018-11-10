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

@test "mdd help templates" {
  run $BATS_CWD/mdd help templates
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd templates lists the templates available" ]
}

@test "mdd help new" {
  run $BATS_CWD/mdd help new
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd new creates a new document from a template" ]
}

@test "mdd help edit" {
  run $BATS_CWD/mdd help edit
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd edit opens a document in your editor" ]
}

@test "mdd help rm" {
  run $BATS_CWD/mdd help rm
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd rm deletes documents, and cleans up any links to them" ]
}

@test "mdd help ls" {
  run $BATS_CWD/mdd help ls
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd ls lists all the documents created" ]
}

@test "mdd help link" {
  run $BATS_CWD/mdd help link
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd link links a parent and child document" ]
}

@test "mdd help unlink" {
  run $BATS_CWD/mdd help unlink
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd unlink breaks the link between a parent and child document" ]
}

@test "mdd help tag" {
  run $BATS_CWD/mdd help tag
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd tag adds tags to a document" ]
}

@test "mdd help untag" {
  run $BATS_CWD/mdd help untag
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd untag removes tags from a document" ]
}

@test "mdd help verify" {
  run $BATS_CWD/mdd help verify
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd verify checks the integrity of the documents" ]
}

@test "mdd help publish" {
  run $BATS_CWD/mdd help publish
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "mdd publish creates a static website for the mdd repository" ]
}
