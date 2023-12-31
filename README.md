# Go-url-shortener

## Description

Go url shortener is a project that offers 3 major features:
1. API to shorten URL
2. API to redirect the shortened URL to original URL
3. A metrics API to get the top N most shortened URL


## Architecture
![](img/architecture.png)

## Getting Started

### Dependencies

- AWS account
- AWS CLI
- Golang
- Windows OS


### Deploy the AWS resources

There are multiple ways to deploy AWS resources . In this we'll cover the 2 basic methods :
1. Via AWS console.
2. Via AWS cli

Resources :
- Lambda function:
    - Build all 3 Go program 
    ```
    go build -o bin/main
    ```
    - Create a ZIP of the build files and upload to AWS via console.
    ![](img/lambda.png)
    OR

    - Deploy via aws cli
    ```
    aws lambda create-function --function-name my-function \
    --zip-file fileb://main.zip --handler main --runtime go1.x \
    --role arn:aws:iam::123456789012:role/lambda-ex
    ```
- API Gateway:
    - Create API via AWS console
    - Create routes for the API
    - Attach the lambda integration to the route
    ![](img/routes-intergration.png)
    - Add necessary permission for the API to invoke the lambda function
    - Dont forget to deploy the API after creation.

- s3 Bucket :
    - Create the S3 bucket via AWS console 
    - Provide the bucket name and region(Bucket name must be unique. Change the bucket name in `constants.go` as per the input)
    - Provide the lambda function necessary permission to read , write and list bucket in s3.

    OR

    - Create bucket via aws cli
    ```
    aws s3api create-bucket --bucket {bucketName} --region {regionName}
    ```
 
### Executing the APIs

To Execute the url-shortener API :
- Execute via postman :
    - Method : **POST**
    - API : https:/{{api-id}}.execute-api.{{region}}.amazonaws.com/url-shortener
    - Body : 
    ```json
    Example
    {
    "url":"https://www.google.com"
    }
    ```
    ![](img/url-shortener.png)

To Execute the redirect API :
- Execute via postman :
    - Method : **GET**
    - API : https:/{{api-id}}.execute-api.{{region}}.amazonaws.com/redirect/{shortUrl}
    - PathParametr : **shortUrl**
    ![](img/redirect-api.png)

To Execute the metrics API :
- Execute via postman :
    - Method : **GET**
    - API : https:/{{api-id}}.execute-api.{{region}}.amazonaws.com/metrics
    - QueryParameter (optional) : **limit**
    ![](img/metrics-api.png)

## Authors

Contributors names and contact info

- Name : Shobhit Gupta
- Email : shobhit00gupta@gmail.com


