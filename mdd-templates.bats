#!/usr/bin/env bats
#
# Test script for 'mdd templates' command
#

@test "mdd templates" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd templates
  [ "$status" -eq 0 ]
}

@test "mdd templates, lists the markdown template files" {
  rm -rf ./.mdd
  $BATS_CWD/mdd init
  run $BATS_CWD/mdd templates
  [ "${lines[0]}" = "   adr: Architecture Decision Record" ]
  [ "${lines[1]}" = "   att: Automated test" ]
  [ "${lines[2]}" = "  itst: Inspection test" ]
  [ "${lines[3]}" = "   mtg: Meeting" ]
  [ "${lines[4]}" = "   nfr: Non Functional Requirement" ]
  [ "${lines[5]}" = "   req: Functional Requirement" ]
}
