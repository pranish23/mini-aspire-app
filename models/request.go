package models

import "github.com/go-playground/validator/v10"

type LoanRequest struct {
	LoanAmount float64 `json:"loan_amount,omitempty" validate:"required,number"`
	LoanTerm   int     `json:"loan_term,omitempty" validate:"required,number"`
}

type LoanRepaymentRequest struct {
	LoanID          string  `json:"loan_id,omitempty" validate:"required,number"`
	RepaymentAmount float64 `json:"repayment_amount,omitempty" validate:"required,number"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func (loanReq LoanRequest) ValidateLoanRequest() []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(loanReq)
	checkErrors(errors, err)
	return errors
}

func (loanRepayReq LoanRepaymentRequest) ValidateRepaymentRequest() []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(loanRepayReq)
	checkErrors(errors, err)
	return errors
}

// Helpers
func checkErrors(errors []*ErrorResponse, err error) {
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
}
