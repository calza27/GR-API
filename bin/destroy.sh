#!/usr/bin/env bash
die() { echo "${1:-urgh}" >&2; exit "${2:-1}"; }

hash sam 2>/dev/null || die "missing dep: sam"

profile=$1
[[ -z $profile ]] && die "Usage: $0 <profile>"

STACK_NAME="gr-api"

echo "~~~ Destroy processors stack"
sam delete \
  --stack-name "$STACK_NAME-processors" \
  --region "ap-southeast-2" \
  --profile "$profile" || die "failed to destroy stack $STACK_NAME-processors"

echo "~~~ Destroy gateway stack"
sam delete \
  --stack-name "$STACK_NAME-gateway" \
  --region "ap-southeast-2" \
  --profile "$profile" || die "failed to destroy stack $STACK_NAME-gateway"

die "~~ cleaning up" 0