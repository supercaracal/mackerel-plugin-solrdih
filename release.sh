#!/bin/bash

set -e

latest_tag=$(git describe --abbrev=0 --tags)
goxz -d dist/$latest_tag -z -os windows,darwin,linux -arch amd64,386
ghr -u supercaracal -r mackerel-plugin-solrdih $latest_tag dist/$latest_tag
