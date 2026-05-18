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


aws cognito-idp admin-create-user \
  --user-pool-id eu-west-1_GQ2jZGoqJ \
  --username admin@admin.com \
  --user-attributes Name=email,Value=admin@admin.com Name=email_verified,Value=true \
  --message-action SUPPRESS \
  --region eu-west-1

aws cognito-idp admin-set-user-password \
  --user-pool-id eu-west-1_GQ2jZGoqJ \
  --username admin@admin.com \
  --password Admin1234! \
  --permanent \
  --region eu-west-1

aws cognito-idp initiate-auth \
  --auth-flow USER_PASSWORD_AUTH \
  --client-id 3c9p2gol9ti304rcdsm4907nak \
  --auth-parameters USERNAME=admin@admin.com,PASSWORD=Admin1234! \
  --region eu-west-1

curl -X POST http://localhost:8080/api/auth/insert-admin \
  -H "Authorization: Bearer <IdToken>"

  
eyJraWQiOiJzSmtwSk42SmE5XC9XZWRjUVJucU1QUDNJNXpIRkFaWDdEbDRLdXQ0cFdsWT0iLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiI5MjQ1OTQ0NC00MDUxLTcwNWMtNWM1ZS1mYWU1OGU4ZTY5OTYiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfR1EyalpHb3FKIiwiY29nbml0bzp1c2VybmFtZSI6IjkyNDU5NDQ0LTQwNTEtNzA1Yy01YzVlLWZhZTU4ZThlNjk5NiIsIm9yaWdpbl9qdGkiOiIwOTE0OTg4NC0zNGEwLTQyODQtOWQ2NC1kMWVmYzFkYjkxODAiLCJhdWQiOiIzYzlwMmdvbDl0aTMwNHJjZHNtNDkwN25hayIsImV2ZW50X2lkIjoiOGU5MTM2NzgtMGRhYy00YTYwLTg5NWItZWY4ZGFjYWQ5MGRlIiwidG9rZW5fdXNlIjoiaWQiLCJhdXRoX3RpbWUiOjE3NzkxMjEyNTIsImV4cCI6MTc3OTEyNDg1MiwiaWF0IjoxNzc5MTIxMjUyLCJqdGkiOiI0NGY2ZTgxNS00YzE2LTQwNjctYTZjNC03OTcyMDU0N2I5YjUiLCJlbWFpbCI6ImFkbWluQGFkbWluLmNvbSJ9.c7QPDGhMdIkGwUigcOMwKz-YiRkECjRVMnN1-f_dy1xn7FIdLh74k5jrjk0WqL6diov5GnMD-impZDKVMAvu9a9YbhCIZwpLDR_WboEZOP9N6cbrxAntO5nEvxwcm7Vxeu6Nsq9HlORRBjgKRQFcXtrcHVuoTIe0aV2cJmS0dbHhrsaWC_z_5xVadv-NkjMuU28Rhrj-lPc-YUfDGzHVMwJB6niNrm_jjM-5sfVPS9MGqgGjE5BQtSze7oD8TqqVnoNrUBuP5Rp6gqvfrm-OX9qJ1ZGPN2VYPbFsf8D0jHfn4CFKOfagJ8OgaeYJcdwGDjT23LBgC7-7ykBpHXZdvA