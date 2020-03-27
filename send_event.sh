#!/bin/env bash

rm tmp_deployment-finish.event
export UUID=$(uuidgen)
export TESTID=$RANDOM
envsubst < deployment-finish.event > tmp_deployment-finish.event
keptn send event -f tmp_deployment-finish.event
