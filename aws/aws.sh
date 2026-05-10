curl ifconfig.me

aws cloudformation create-stack --stack-name doclane --template-body file://aws.yml --parameters ParameterKey=MyPublicCidr,ParameterValue=109.166.139.204/32 --capabilities CAPABILITY_NAMED_IAM

aws cloudformation create-stack --stack-name doclane-vpc --template-body file://vpc.yml --parameters ParameterKey=AZ1,ParameterValue=eu-west-1a ParameterKey=AZ2,ParameterValue=eu-west-1b ParameterKey=MyPublicCidr,ParameterValue=/32

aws cloudformation update-stack --stack-name doclane-vpc --template-body file://vpc.yml --parameters ParameterKey=AZ1,UsePreviousValue=true ParameterKey=AZ2,UsePreviousValue=true ParameterKey=MyPublicCidr,UsePreviousValue=true

aws cloudformation create-stack --stack-name doclane-rds-params --template-body file://rds-params.yml --parameters ParameterKey=DBPassword,ParameterValue=robert ParameterKey=DBUsername,ParameterValue=robertrobert

aws cloudformation create-stack --stack-name doclane-rds --template-body file://rds.yml 

aws cloudformation create-stack --stack-name doclane-bastion --template-body file://bastion-host.yml

aws cloudformation create-stack --stack-name doclane-lambda-backend --template-body file://lambda-backend.yml --capabilities CAPABILITY_NAMED_IAM

aws cloudformation create-stack --stack-name doclane-apigw --template-body file://apigw.yml

psql -h <host> -p 5432 -U robert -d doclane

aws cloudformation create-stack --stack-name doclane-asg --template-body file://frontend.yml

