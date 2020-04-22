# rest-api-creator/backend-sls #

This serverless project implements the backend that supports the rest-api-creator website. It is a Go module that runs on AWS Lambda and uses DynamoDB as a database. AWS APIGateway provides a REST API to invoke the Lambda functions.

## Package Structure

There are three common Go packages in the project: `auth`, `dao` and `errors`. Each of these packages is used in every Lambda function. `auth` implements session token generation and validation. `dao` implements database queries. `errors` implements error generation and handling.

In addition to these three packages, each Lambda function has its own package named after the API endpoint that it implements. For example, the `getproject` package implements the `GET /project/{id}` endpoint. Within each of these packages, there is a `handler.go` file and an `action.go` file. `action.go` implements the actual business logic of the API endpoint. `handler.go` is in charge of parsing the HTTP request that the Lambda function receives from APIGateway and forwarding the request parameters to `action.go`. `handler.go` then takes the response from `action.go`, converts it to a format that APIGateway understands, and returns the converted response to APIGateway. 

The only execption to the above rule is the `portal` package. This package combines the signup and login API endpoints, as the two are extremely similar. In the `portal` package, `handlers.go` implements the handlers for both the signup and login APIs, while `signup.go` implements the business logic for the signup endpoint and `login.go` implements the business logic for the login endpoint.
