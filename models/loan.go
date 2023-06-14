package models

type Loan struct {
	LoanID     string       `json:"loan_id,omitempty"`
	CustomerID string       `json:"customer_id,omitempty"`
	Status     string       `json:"status,omitempty"`
	Repayments []*Repayment `json:"repayments,omitempty"`
}

type Repayment struct {
	ID     int    `json:"repayment_id,omitempty"`
	Date   string `json:"date,omitempty"`
	Amount string `json:"amount,omitempty"`
	Status string `json:"repayment_status,omitempty"`
}

// Loan Term
type LoanTerm int32

const (
	LOAN_TERM_WEEKLY  LoanTerm = 7
	LOAN_TERM_MONTHLY LoanTerm = 30
)

// Loan status
type LoanStatus int32

const (
	LOAN_STATUS_PENDING  LoanStatus = 1
	LOAN_STATUS_APPROVED LoanStatus = 2
	LOAN_STATUS_PAID     LoanStatus = 3
)

var (
	LoanStatus_string = map[int32]string{
		1: "PENDING",
		2: "APPROVED",
		3: "PAID",
	}
	LoanStatus_int = map[string]int32{
		"PENDING":  1,
		"APPROVED": 2,
		"PAID":     3,
	}
)

// Repayment status
type LoanRepaymentStatus int32

const (
	REPAYMENT_STATUS_UNPAID LoanRepaymentStatus = 1
	REPAYMENT_STATUS_PAID   LoanRepaymentStatus = 2
)

var (
	RepaymentStatus_string = map[int32]string{
		1: "UNPAID",
		2: "PAID",
	}
	RepaymentStatus_int = map[string]int32{
		"UNPAID": 1,
		"PAID":   2,
	}
)
