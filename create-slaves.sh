#!/bin/bash

#
#    This script creates salve account CFN roles for stack sets
#


function usage
{
    echo "usage: create-slaves.sh [-h --help]
                                      --profiles profile_file
                                      --region aws-region
                                      --master master-account
                                      --stack_name stack-name
"
}

delete=false
profiles=""
region=""
master=""
stack=""
accounts=()



while [ "$1" != "" ]; do
    case $1 in
        -p | --profiles )       shift
                                profiles=$1
                                ;; 
        -r | --region )         shift
                                region=$1
                                ;; 
        -d | --delete )         shift
                                delete=true 
                                ;; 
        -m | --master )         shift
                                master=$1
                                ;; 
        -s | --stack-name )     shift
                                stack=$1
                                ;;
        -h | --help )           usage
                                exit
                                ;;
    esac
    shift
done

cat $profiles

if [ $delete == 'false' ]
then
  set -x
  for account in $(cat $profiles) ; do
      aws cloudformation create-stack --profile $account \
          --capabilities CAPABILITY_NAMED_IAM \
          --region $region  \
          --stack-name $stack \
          --template-body file://slave.yml --parameters \
          ParameterKey=MasterRegion,ParameterValue=$region \
          ParameterKey=MasterAccountId,ParameterValue=$master;
  done
else
  set -x
  for account in $(cat $profiles) ; do
      aws cloudformation delete-stack --profile $account \
          --region $region  \
          --stack-name $stack
  done  
fi