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
### Step 2: run the following command to build and start the Docker containers:

```bash
docker-compose up --build   
```

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
