
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