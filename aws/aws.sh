aws cloudformation create-stack --stack-name doclane-vpc --template-body file://vpc.yml --parameters ParameterKey=AZ1,ParameterValue=eu-west-1a ParameterKey=AZ2,ParameterValue=eu-west-1b ParameterKey=MyPublicCidr,ParameterValue=/32

aws cloudformation create-stack --stack-name doclane-rds-params --template-body file://rds-params.yml --parameters ParameterKey=DBPassword,ParameterValue=robert ParameterKey=DBUsername,ParameterValue=robertrobert

aws cloudformation create-stack --stack-name doclane-rds --template-body file://rds.yml 

aws cloudformation create-stack --stack-name doclane-bastion --template-body file://bastion-host.yml