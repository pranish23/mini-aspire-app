## Description

Mini apsire app by Pranish Ajagekar
 

## Run project

```bash
Run locally     :   make dev
Build           :   make build
Run             :   make run
Build & Run     :   make build && make run
```
## Test
```bash
# make test
```

## Installation

```bash
(If you dont have golang installed)
```
Refer this link : https://go.dev/doc/install

## Postman collection
Please check ```./postman_collection``` folder in the project

- ```/``` Healthcheck API
- ```/loan/:customerID/view``` View customer loans  
- ```/loan/:customerID/submit``` Apply for a new customer loan  
- ```/loan/:loanID/approve``` Approve loan (On admins allowed)  
- ```/loan/:customerID/repay``` Repay loan term and mark them as paid

## Next steps
- The API will run on http://localhost:3000/
- Use postman_collection for testing
- Protected APIs can be accesed with ```Basic Auth``` , check credentials in ```user/users.go```
- Run ```make test``` to run unit tests

### Technologies Used
- Golang (Backend)
- Fiber https://github.com/gofiber/fiber (Web Framework)
- go-memdb https://github.com/hashicorp/go-memdb (In memory Database) 

### Design
The mini aspire app is minimalistic web app where have two type of roles:
- Admins
- Customers
Customers can Read , Submit & Repay a loan.

Admins can approve the loan submitted by customers.

*NOTE:*

Customer can only repay loans terms for loan which are approved to admins.



## API Reference

#### Health check API

```http
  GET /
```


#### Get loans

```http
  GET /loan/:customerID/view
```
List all the loan applied by customer.
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `customerID`      | `string` | **Required**. e.g: `customer1` |

#### Submit loan
```http
  POST /loan/:customerID/view
```
Submit a new loan request for customer
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `customerID`      | `string` | **Required**. e.g: `customer1` |

#### Approve loan
```http
  PATCH /loan/:loanID/view
```
Approve loan request from customer
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `loanID`      | `string` | **Required**. e.g: `loan-e0f7f400-6948-45f9-a5ec-59303b0f0a28` |

*Use the loan ID is returned by `/submit` api*

**Request Body:**
`
{
    "loan_amount": 1000,
    "loan_term": 3
}
`

#### Repay loan
```http
  PATCH /loan/:customerID/repay
```
Repay an approved loan request
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `loanID`      | `string` | **Required**. e.g: `customer1` |

**Request Body:**
`
{
    "loan_id":"loan-e0f7f400-6948-45f9-a5ec-59303b0f0a28",
    "repayment_amount":333.33
}
`

### Demo
Here is a quick demo of the app.

https://drive.google.com/file/d/1cFznlNftnmLMSvgYtHHnaaGI2ENGqZRv/view?usp=drive_link



