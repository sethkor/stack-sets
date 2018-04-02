include env.mk

create: _validate-stack-name _validate-profile
	$(call cfn,create-stack --enable-termination-protection)

update: _validate-stack-name _validate-profile
	$(call cfn,update-stack)

delete-master: _validate-stack-name _validate-profile
	aws cloudformation delete-stack --profile $(profile) --region $(region) --stack-name $(STACK-NAME)

delete-slaves: _validate-stack-name _validate-profile
	for account in $(accounts) ; do \
    aws cloudformation delete-stack --profile $$accounts --region $(region) --stack-name $(STACK-NAME-SLAVE) \
  done

delete: _validate-stack-name _validate-profile delete-slaves delete-master
	aws cloudformation delete-stack --profile $(profile) --region $(region) --stack-name $(STACK-NAME-SLAVE)
	
query: _validate-stack-name _validate-profile
	aws cloudformation --profile $(profile) --region $(region) describe-stacks --stack-name $(STACK-NAME) --query 'Stacks[].StackStatus' --output text
	
master:
	aws cloudformation $1 --profile $(profile) \
	      --capabilities CAPABILITY_NAMED_IAM \
        --region $(region)  \
        --stack-name $(STACK-NAME) \
        --template-body file://master.yml

slaves:
	for account in $(accounts) ; do \
    	aws cloudformation create-stack --profile $$account \
          --capabilities CAPABILITY_NAMED_IAM \
          --region $(region)  \
          --stack-name $(STACK-NAME-SLAVE) \
          --template-body file://slave.yml --parameters \
	        ParameterKey=MasterRegion,ParameterValue=$(MasterRegion) \
	        ParameterKey=MasterAccountId,ParameterValue=$(MasterAccountId); \
  done
        
define cfn
	aws cloudformation $1 --profile $(profile) \
	      --capabilities CAPABILITY_NAMED_IAM \
        --region $(region)  \
        --stack-name $(STACK-NAME) \
        --template-body file://master.yml
endef

_validate-stack-name:
ifndef STACK-NAME
	$(error STACK-NAME is required)
endif

_validate-profile:
ifndef profile
	$(error profile is required)
endif
