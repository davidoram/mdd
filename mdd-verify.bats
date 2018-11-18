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
  [ "${lines[0]}" = "File '.mdd/project.data' doesnt exist" ]
}

@test "mdd verify, remove documents directory" {
  $BATS_CWD/mdd init
  rm -r ./.mdd/documents
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Directory '.mdd/documents' doesnt exist" ]
}

@test "mdd verify, remove templates directory" {
  $BATS_CWD/mdd init
  rm -r ./.mdd/templates
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Directory '.mdd/templates' doesnt exist" ]
}

@test "mdd verify, remove publish directory" {
  $BATS_CWD/mdd init
  rm -r ./.mdd/publish
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Directory '.mdd/publish' doesnt exist" ]
}

@test "mdd verify, cant read template" {
  $BATS_CWD/mdd init
  chmod -r ./.mdd/templates/adr.md
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "open .mdd/templates/adr.md: permission denied" ]
}

@test "mdd verify, missing title" {
  $BATS_CWD/mdd init
  echo "" > ./.mdd/templates/adr.md
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Template '.mdd/templates/adr.md', is missing a title" ]
}

@test "mdd verify, cant read document" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  chmod -r $file_path
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "open ${file_path}: permission denied" ]
}

@test "mdd verify, document doesnt match filename regex" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  mv $file_path ./.mdd/documents/not-a-valid-document-name.md
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Document 'not-a-valid-document-name.md' doesnt match mdd filename regex" ]
}

@test "mdd verify, no template for document" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  mv $file_path ./.mdd/documents/xyz-b7-0001.md
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Document 'xyz-b7-0001.md' no template matching shortcode 'xyz'" ]
}

@test "mdd verify, document has no title" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  file=$(basename ${file_path})
  echo "Moo" > $file_path
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Document '${file}' has no title" ]
}

@test "mdd verify, missing metadata block" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  file=$(basename ${file_path})
  echo "# title" > $file_path
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Document '${file}' missing metadata block" ]
}

@test "mdd verify, incomplete metadata block, invalid entry" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  file=$(basename ${file_path})
  echo "# title" > $file_path
  echo "<!-- mdd" >> $file_path
  echo "The end" >> $file_path
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Document '${file}' expected 2 values, found 1 from metadata 'The end'" ]
}

@test "mdd verify, invalid metadata tag" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  file=$(basename ${file_path})
  echo "# title" > $file_path
  echo "<!-- mdd" >> $file_path
  echo "mdd-unknown: moo" >> $file_path
  echo "-->" >> $file_path
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Document '${file}' unrecognised metadata tag 'mdd-unknown'" ]
}

@test "mdd verify, invalid metadata tag - key without value" {
  $BATS_CWD/mdd init
  file_path=$($BATS_CWD/mdd new adr)
  file=$(basename ${file_path})
  echo "# title" > $file_path
  echo "<!-- mdd" >> $file_path
  echo "mdd-empty:" >> $file_path
  echo "-->" >> $file_path
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Document '${file}' metadata value for key 'mdd-empty' is empty" ]
}

@test "mdd verify, invalid link" {
  $BATS_CWD/mdd init
  parent_path=$($BATS_CWD/mdd new adr)
  parent=$(basename ${parent_path})
  child_path=$($BATS_CWD/mdd new adr)
  child=$(basename ${child_path})
  $BATS_CWD/mdd link ${parent} ${child}
  rm ${child_path}
  run $BATS_CWD/mdd verify
  [ "$status" -eq 1 ]
  [ "${lines[0]}" = "Document '${parent}' has child '${child}' which doesnt exist" ]
}

@test "mdd verify, circular link" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  parent=$(basename $($BATS_CWD/mdd new adr))
  child=$(basename $($BATS_CWD/mdd new adr))
  $BATS_CWD/mdd link ${parent} ${child}
  $BATS_CWD/mdd link ${child} ${parent}
  run $BATS_CWD/mdd verify
  [ "$status" -eq 0 ]
}
