package controllers

import (
	"fmt"
	"strconv"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pranish23/mini-aspire-app/db"
	"github.com/pranish23/mini-aspire-app/models"
)

type loanController struct {
	db *db.LoanDB
}

func NewLoanController(dbConn *db.LoanDB) loanController {
	return loanController{
		db: dbConn,
	}
}

func (lc loanController) Approve(c *fiber.Ctx) error {
	loanID := c.Params("loanID")
	getLoanData, err := lc.db.ReadByID("id", loanID)
	if err != nil {
		return err
	}

	getLoanData[0].Status = models.LoanStatus_string[int32(models.LOAN_STATUS_APPROVED)]
	lc.db.Upsert(getLoanData[0])
	return c.JSON(getLoanData[0])
}

func (lc loanController) Submit(c *fiber.Ctx) error {
	customerID := c.Params("customerID")
	err := checkAccess(customerID, c)
	if err != nil {
		return err
	}

	loanRequest := new(models.LoanRequest)
	if err := c.BodyParser(loanRequest); err != nil {
		return err
	}
	errors := loanRequest.ValidateLoanRequest()
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	loanAmount := loanRequest.LoanAmount
	paymentTerms := 3

	newLoan := &models.Loan{
		LoanID:     "loan-" + uuid.NewString(),
		CustomerID: customerID,
		Status:     models.LoanStatus_string[int32(models.LOAN_STATUS_PENDING)],
	}
	addRepayments(loanAmount, paymentTerms, newLoan)

	err = lc.db.Upsert(newLoan)
	if err != nil {
		return err
	}

	return c.JSON(newLoan)
}

func (lc loanController) View(c *fiber.Ctx) error {
	customerID := c.Params("customerID")
	err := checkAccess(customerID, c)
	if err != nil {
		return err
	}

	getLoanData, err := lc.db.ReadByID("customer_id", customerID)
	if err != nil {
		return err
	}
	return c.JSON(getLoanData)
}

func (lc loanController) AddRepayment(c *fiber.Ctx) error {
	customerID := c.Params("customerID")
	err := checkAccess(customerID, c)
	if err != nil {
		return err
	}

	loanRepayRequest := new(models.LoanRepaymentRequest)
	if err := c.BodyParser(loanRepayRequest); err != nil {
		return err
	}
	errors := loanRepayRequest.ValidateRepaymentRequest()
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	getLoanData, err := lc.db.ReadByID("id", loanRepayRequest.LoanID)
	if err != nil {
		return err
	}

	_, err = processRepayment(getLoanData[0], loanRepayRequest, &lc)
	if err != nil {
		return err
	}

	return c.JSON(getLoanData[0])

}

// ++========================= Helpers =========================++
func getNextRepaymentDate(t time.Time, addDays int) string {
	newTime := t.AddDate(0, 0, addDays).Format("2006-January-02")
	return newTime
}

func addRepayments(loanAmount float64, paymentTerms int, newLoan *models.Loan) {
	timeNow := time.Now()
	repaymentAmount := loanAmount / float64(paymentTerms)
	for i := 1; i <= paymentTerms; i++ {
		newLoan.Repayments = append(newLoan.Repayments, &models.Repayment{
			ID:     i,
			Date:   getNextRepaymentDate(timeNow, i*int(models.LOAN_TERM_WEEKLY)),
			Amount: fmt.Sprintf("%.2f", repaymentAmount),
			Status: models.RepaymentStatus_string[int32(models.REPAYMENT_STATUS_UNPAID)],
		})
	}
}

func checkAccess(reqUserName string, c *fiber.Ctx) error {
	var err error
	contextUsername := c.Context().UserValue("username")
	if reqUserName != contextUsername {
		err = fmt.Errorf("user unauthorized to perform this operation")
	}
	return err
}

func processRepayment(loanData *models.Loan, repaymentReq *models.LoanRepaymentRequest, lc *loanController) (*models.Loan, error) {
	repayments := loanData.Repayments

	if loanData.Status != getLoanString(models.LOAN_STATUS_APPROVED) {
		var err error
		switch loanData.Status {
		case getLoanString(models.LOAN_STATUS_PENDING):
			err = fmt.Errorf("loan is still pending")
		case getLoanString(models.LOAN_STATUS_PAID):
			err = fmt.Errorf("loan is fully paid")
		}
		return nil, err
	}

	for _, v := range repayments {
		if v.Status == getRepaymentString(models.REPAYMENT_STATUS_UNPAID) {
			termAmount, err := strconv.ParseFloat(v.Amount, 64)
			if err != nil {
				return nil, err
			}

			if repaymentReq.RepaymentAmount >= termAmount {
				v.Status = getRepaymentString(models.REPAYMENT_STATUS_PAID)
				err := lc.db.Upsert(loanData)
				if err != nil {
					return nil, err
				}
				if v.ID == len(repayments) {
					err := markLoanAsPaid(loanData, lc)
					if err != nil {
						return nil, err
					}
				}
				return loanData, nil
			} else {
				return nil, fmt.Errorf("repayment amount less than loan term amount ")
			}

		}
	}
	return loanData, nil
}

func getLoanString(status models.LoanStatus) string {
	return models.LoanStatus_string[int32(status)]
}

func getRepaymentString(status models.LoanRepaymentStatus) string {
	return models.RepaymentStatus_string[int32(status)]
}

func markLoanAsPaid(loanData *models.Loan, lc *loanController) error {
	loanData.Status = getLoanString(models.LOAN_STATUS_PAID)
	err := lc.db.Upsert(loanData)
	if err != nil {
		return err
	}
	return nil
}
