#!/bin/bash

CHAIN_ID=${CHAIN_ID:-Oraichain}
DS_RAW=${1:-classification}
TC_RAW=${2:-classification_testcase}
OS=${3:-classification_oscript}
DS_INPUT=${4:-''}
TC_INPUT=${5:-''}
OS_INPUT=${6:-''}
OS_DS=''
OS_TC=''
NONCE=${7:-1}
PASS=${8:-123456789}

# deploy data sources
bash $PWD/scripts/deploy_datasource.sh $DS_RAW "$DS_INPUT" $NONCE $PASS

# deploy test cases
bash $PWD/scripts/deploy_testcase.sh $TC_RAW "$TC_INPUT" $NONCE $PASS

# deploy oracle scripts
bash $PWD/scripts/deploy_oscript.sh $DS_RAW $TC_RAW $OS "$OS_INPUT" $NONCE $PASS