# Order-Nest
Welcome to the Order-Nest System. In the following sections, you will get an idea about the whole system in brief:

[![Postman Documentation](https://img.shields.io/badge/Postman-Documentation-orange)](https://documenter.getpostman.com/view/29231217/2sB3QCRYoa)

## Getting Started

To run the project on your machine, follow these steps:

### Step 1: Clone the repository and navigate to the project directory.

```bash
git clone git@github.com:SaifulJnU/order-nest.git
cd order-nest
```
### Step 2 (with Docker): build and start the containers

```bash
docker-compose up --build   
```

### Run without Docker (local setup)

Prerequisites:
- Go 1.25+
- PostgreSQL running locally (or use the `postgres` service from docker-compose)

Steps:
1) Ensure DB is available and matches `order-nest-config.yaml` (default compose values):
   - host: 127.0.0.1
   - port: 5432
   - user: ordernest
   - password: adminsecret
   - dbname: ordernest
2) Install deps and run the app (use local config):
```bash
go mod tidy
go run ./main.go --config=order-nest-config-local.yaml serve
```
The API will listen on http://localhost:8080

## Problem Statement

This project involves building a backend service in Go that enables basic order management. The system must provide authentication for users and support essential order operations such as creating new orders, listing existing orders, and cancelling orders. The backend should also ensure secure access using JWT tokens, handle validation for user inputs (such as phone number and required fields), and return structured responses for both success and error cases.

### Key Features

* Login User – Authenticate a user and generate an access token for secure API access.
* Logout User – Invalidate the user session and revoke the access token.
* Create an Order – Place a new delivery request with recipient and order details.
* List Orders – Retrieve all orders associated with a store or user.
* Cancel Order – Cancel an existing delivery order by consignment ID.

## WorkFlow Design

![img.png](docs/img-system-flow.png)

### High Level System flow

1. **Front-End Client:** User-facing interface (CLI, web, or Postman) that interacts with the API for login, order creation, listing, cancellation, and logout.
2. **API Layer (Go Application):** Exposes REST endpoints (/login, /orders, /orders/all, /orders/{id}/cancel, /logout). It validates requests, applies business logic, and ensures secure access with JWT.
3. **Authentication Service:** Verifies email & password, issues JWT tokens, and validates them on each protected request. Handles login and logout operations.
4. **Order Service:** Core module that manages creating new orders, fetching order lists, and cancelling orders. Implements business rules (e.g., delivery fee, COD fee, validation of phone numbers).
5. **Database (PostgreSQL):** Stores user credentials, order details, and related metadata (status, consignment IDs, fees, amounts).
6. **Validation & Pricing:** Validates user inputs (phone, required fields), and calculates delivery charges & COD fee based on city, weight, and amount.
7. **Transaction Handling:** Ensures atomic operations while creating or cancelling orders so that no partial updates occur.
8. **Error Handling Layer:** Returns clear, structured responses for success and failure (200, 400, 401, 422), ensuring predictable behavior.
9. **Order Listing & Pagination:** Provides paginated results when fetching orders, with filters like transfer_status and archive.
10. **Order Cancellation:** Uses the consignment_id to cancel specific orders securely, updating their status consistently.

## Database schema

![img.png](docs/img-schema-design.png)

## Used Packages and Tools

- **Web Framework**: GIN
- **Logger**: Logrus
- **Testing**: Testify
- **Database**: Postgres
- **Auth Token**: JWT
- **ORM**: GORM
- **Validation**: Go-playground/validator
- **Environment Variables**: Viper
- **Password Hashing**: Bcrypt
- **API Documentation**: Postman
- **Containerization**: Docker

## cURL examples

Base URL: `http://localhost:8080`

### Login 
```bash
curl --location 'http://localhost:8080/api/v1/login' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "username": "01901901901@mailinator.com",
    "password": "321dsa"
}'
```

Copy `access_token` from the response as TOKEN below.

### Create order (required)
```bash
curl --location 'http://localhost:8080/api/v1/orders' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer {{TOKEN}}' \
  --data '{
    "store_id": 131172,
    "merchant_order_id": "my-order-001",
    "recipient_name": "Karim",
    "recipient_phone": "01626692267",
    "recipient_address": "Jamuna Future Park, Vatara, Dhaka",
    "recipient_city": 1,
    "recipient_zone": 1,
    "recipient_area": 1,
    "delivery_type": 48,
    "item_type": 2,
    "special_instruction": "Call on arrival",
    "item_quantity": 1,
    "item_weight": 0.5,
    "amount_to_collect": 1200,
    "item_description": "Books"
}'
```

### List orders (required)
```bash
curl --location 'http://localhost:8080/api/v1/orders/all?transfer_status=1&archive=0&limit=10&page=1' \
  --header 'Authorization: Bearer {{TOKEN}}'
```

### Cancel order 
```bash
curl --location --request PUT 'http://localhost:8080/api/v1/orders/{{CONSIGNMENT_ID}}/cancel' \
  --header 'Authorization: Bearer {{TOKEN}}'
```

## API Doc

### Login

- URL: `{{HOST}}/api/v1/login`
- Method: POST

```bash
curl --location '{{HOST}}/api/v1/login' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "username": "01901901901@mailinator.com",
    "password": "321dsa"
}'
```

Success 200:
```json
{
   "token_type": "Bearer",
   "expires_in": 432000,
   "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxIiwiZXhwIjoxNzU4ODAwNjMyLCJpYXQiOjE3NTg3OTk3MzIsImp0aSI6Ijk1OGVkNjBjLTJlYjgtNDUyNC1iMDhlLTgyY2U3YjJjYTk1NiIsIm5hbWUiOiIwMTkwMTkwMTkwMUBtYWlsaW5hdG9yLmNvbSIsInRva2VuX3R5cGUiOiJhY2Nlc3MifQ.0gN5xn1kNyaxlMURjgL5Wh3zjbjkzVRJRqIgQXqW5RY",
   "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIxIiwiZXhwIjoxNzU5MjMxNzMyLCJpYXQiOjE3NTg3OTk3MzIsImp0aSI6Ijk1OGVkNjBjLTJlYjgtNDUyNC1iMDhlLTgyY2U3YjJjYTk1NiIsIm5hbWUiOiIwMTkwMTkwMTkwMUBtYWlsaW5hdG9yLmNvbSIsInRva2VuX3R5cGUiOiJyZWZyZXNoIn0.zbuKAcPXxPeDD3xjh6mjtnz0OapbgsZnLUUAsHF35oc"
}
```

Error 400:
```json
{
  "message": "The user credentials were incorrect.",
  "type": "error",
  "code": 400
}
```

### Logout 

- URL: `{{HOST}}/api/v1/logout`
- Method: POST

```bash
curl --location --request POST '{{HOST}}/api/v1/logout' \
  --header 'authorization: Bearer {{TOKEN}}'
```

Success 200:
```json
{
  "message": "Successfully logged out",
  "type": "success",
  "code": 200
}
```

Error 401:
```json
{
  "message": "Unauthorized",
  "type": "error",
  "code": 401
}
```

### Create New Order (required)

- URL: `{{HOST}}/api/v1/orders`
- Method: POST

```bash
curl --location '{{HOST}}/api/v1/orders' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer {{TOKEN}}' \
  --data '{
    "store_id": 131172,
    "merchant_order_id": "DD1134498",
    "recipient_name": "Karim",
    "recipient_phone": "01626692267",
    "recipient_address": "Jamuna Future Park, Vatara, Dhaka",
    "recipient_city": 1,
    "recipient_zone": 1,
    "recipient_area": 1,
    "delivery_type": 48,
    "item_type": 2,
    "special_instruction": "call on arrival",
    "item_quantity": 1,
    "item_weight": 0.5,
    "amount_to_collect": 1300,
    "item_description": "Books"
}'
```

Success 200:
```json
{
  "consignment_id": "{{CONSIGNMENT_ID}}",
  "merchant_order_id": "DD1134498",
  "order_status": "Pending",
  "delivery_fee": 60
}
```

Error 401:
```json
{
  "message": "Unauthorized",
  "type": "error",
  "code": 401
}
```

Error 422:
```json
{
  "message": "Please fix the given errors",
  "type": "error",
  "code": 422,
  "errors": {
    "recipient_name": ["The recipient name field is required."],
    "recipient_phone": ["The recipient phone field is required."],
    "recipient_address": ["The recipient address field is required."],
    "delivery_type": ["The delivery type field is required."],
    "amount_to_collect": ["The amount to collect field is required."],
    "item_quantity": ["The item quantity field is required."],
    "item_weight": ["The item weight field is required."],
    "item_type": ["The item type field is required."]
  }
}
```

### List All Orders (required)

- URL: `{{HOST}}/api/v1/orders/all`
- Method: GET

```bash
curl --location '{{HOST}}/api/v1/orders/all?limit=10&page=1&transfer_status=1&archive=0' \
  --header 'Authorization: Bearer {{TOKEN}}'
```

Success 200 (with data):
```json
{
   "code": 200,
   "data": {
      "data": [
         {
            "order_consignment_id": "b54bb151-7cc9-4724-974f-765e15476385",
            "order_created_at": "2025-09-25T11:29:11.42301Z",
            "order_description": "Books",
            "merchant_order_id": "1444498",
            "recipient_name": "Karim",
            "recipient_address": "Jamuna Future Park, Dhaka, Bangladesh",
            "recipient_phone": "01626692267",
            "order_amount": 1200,
            "total_fee": 72,
            "instruction": "Call before delivery",
            "order_type_id": 1,
            "cod_fee": 12,
            "promo_discount": 0,
            "discount": 0,
            "delivery_fee": 60,
            "order_status": "Pending",
            "order_type": "Delivery",
            "item_type": "Parcel",
            "transfer_status": 1,
            "archive": 0,
            "updated_at": "2025-09-25T11:29:11.42301Z",
            "created_by": 1,
            "updated_by": 1
         }
      ],
      "total": 1,
      "current_page": 1,
      "per_page": 10,
      "total_in_page": 1,
      "last_page": 1
   },
   "message": "Orders successfully fetched.",
   "type": "success"
}
```

### Cancel Order 

- URL: `{{HOST}}/api/v1/orders/{{CONSIGNMENT_ID}}/cancel`
- Method: PUT

```bash
curl --location --request PUT '{{HOST}}/api/v1/orders/{{CONSIGNMENT_ID}}/cancel' \
  --header 'Authorization: Bearer {{TOKEN}}'
```

Success 200:
```json
{
  "message": "Order Cancelled Successfully",
  "type": "success",
  "code": 200
}
```
