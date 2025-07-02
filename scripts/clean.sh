#!/usr/bin/env -S bash -x

rm -rf builds completions ./fsize
find . -name "*.mo" -delete
find . -name "*.po~" -delete
find . -name "*.log" -delete
