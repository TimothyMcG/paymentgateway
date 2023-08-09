# Tech Test : Processout

This is my solution to the scenario of creating a payment gateway

## Command Options

To use this solution, you can use the following terminal
commands from the repository root:

- `make up` - bring up the `Postgres` database and setup the
schema
- `make down` - bring down the `Postgres` database
- `make run` - run the service (listens on port `8080`)


## Implementation

### Schema
For this solution, you can review the database schema I designed
in the following location: `db/init.sql`

I've kept this fairly simple. Only using a merchant and payments table. 

### Endpoints

You can interact with my solution with the following endpoints all listening on localhost:8080 

- POST api/register: Requires a json body with a name parameter 
{
    "name":"tim"
}
- POST api/login: Requires a json body with the name parameter as above. Returns a JWT token for auth on the payment endpoints and a merchantId
{
    "merchantID": 1,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2OTE1NzI3MzksInVzZXJfaWQiOjB9.eJ7kkffpKDb_dRhsv2XJVfN-kejHjFbeFGBBkEbsqb4"
}


- POST /api/v1/payment: Handles a users payment expects a jwt token passed in the header and a json body with payment info. Will return a paymentId 
{
    "card_name": "tim",
    "card_number": "1234567891234567",
    "card_type": "visa",
    "currency": "GBP",
    "card_expiry_month": 11,
    "card_expiry_year": 2024, 
    "cvv": "123",
    "description": "work",
    "amount": 1000,
}

example call: curl --location --request POST 'localhost:8080/api/v1/payment' \
--header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2OTE1NzI3MzksInVzZXJfaWQiOjB9.eJ7kkffpKDb_dRhsv2XJVfN-kejHjFbeFGBBkEbsqb4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "card_name": "tim",
    "card_number": "1234567891234567",
    "card_type": "visa",
    "currency": "GBP",
    "card_expiry_month": 11,
    "card_expiry_year": 2024, 
    "cvv": "123",
    "description": "work",
    "amount": 1000,
    "merchant_id": 1
}'


- GET /api/v1/payment: Returns payment details. Expects an Id passed in the URL and a token in the header. 

example call: curl --location --request GET 'localhost:8080/api/v1/payment?id=1' \
--header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2OTE1NzExNDMsInVzZXJfaWQiOjB9.ZcADT6HMyqszz9LLAQqDbF0MPrS3csR58EB579Ld834'

### Future Features 
- Add a foreign key into the payments table for merchant_id. This would be a one to many relationship, one merchant can be related with many payments. This could then be extended to provide the merchant with endpoints to see an overview of their payments.
- Create a banking table. This would contain a list of banks and endpoints to call for each bank. The benefits being a user can select from multiple banks (if the merchant allows) and our payment gateway will facilitate that. This could be cached on startup. 
- More stats endpoint. The service could provide the merchant with endpoints that provide more granular data. A break down of payments by type or product for example 
- 
