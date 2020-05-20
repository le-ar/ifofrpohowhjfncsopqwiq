#!/bin/bash

OUT=participantSolution.go
TMP=tempNameForParticipantSolution.go

cat $filename > $TMP || exit 1
rm $filename
cat $TMP > $OUT || exit 1
rm $TMP