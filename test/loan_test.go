package test

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pranish23/mini-aspire-app/controllers"
	"github.com/pranish23/mini-aspire-app/db"
	"github.com/pranish23/mini-aspire-app/middleware"
	"github.com/pranish23/mini-aspire-app/models"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()

	dbConn, _ := db.Init()
	lc := controllers.NewLoanController(dbConn)

	loanService := app.Group("/loan")
	loanService.Patch("/:loanID/approve", middleware.CheckAdminAuthorization(), lc.Approve)
	loanService.Post("/:customerID/submit", middleware.CheckCustomerAuthorization(), lc.Submit)
	loanService.Get("/:customerID/view", middleware.CheckCustomerAuthorization(), lc.View)
	loanService.Patch("/:customerID/repay", middleware.CheckCustomerAuthorization(), lc.AddRepayment)

	return app
}

func TestLoanView(t *testing.T) {
	app := setupTestApp()
	// T1 => Without basic auth
	req := httptest.NewRequest(http.MethodGet, "/loan/customer1/view", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal("Error occurred:", err)
	}
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// T2 => With basic auth
	req = httptest.NewRequest(http.MethodGet, "/loan/customer1/view", nil)
	getAuth(req)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal("Error occurred:", err)
	}
	respStr, err := getResponseData(res)
	if err != nil {
		t.Fatal(err, respStr)
	}
	assert.Equal(t, "no data found", respStr)

	// T3 => Without basic auth, different Customer ID
	req = httptest.NewRequest(http.MethodGet, "/loan/customer2/view", nil)
	getAuth(req)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal("Error occurred:", err, res.StatusCode)
	}
	respStr, err = getResponseData(res)
	if err != nil {
		t.Fatal(err, respStr)
	}
	assert.Equal(t, "user unauthorized to perform this operation", respStr)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	// T4 => With loan data
	_, _, err = createLoan(t, app, "1000", "3")
	if err != nil {
		t.Fatal(err)
	}
	req = httptest.NewRequest(http.MethodGet, "/loan/customer1/view", nil)
	getAuth(req)
	res, err = app.Test(req)
	if err != nil {
		t.Fatal("Error occurred:", err, res.StatusCode)
	}

	loanData, err := getParsedLoanDataFromList(res)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "customer1", loanData.CustomerID)

	// Stop the API server gracefully
	if err := app.Shutdown(); err != nil {
		t.Fatalf("Failed to stop API server gracefully: %v", err)
	}
}

func TestLoanSubmit(t *testing.T) {
	app := setupTestApp()
	// T5 => Without basic auth
	req := httptest.NewRequest(http.MethodPost, "/loan/customer1/submit", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal("Error occurred:", err)
	}
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// T6 => With basic auth
	loanData, res, err := createLoan(t, app, "1000", "3")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "customer1", loanData.CustomerID)

	// Stop the API server gracefully
	if err := app.Shutdown(); err != nil {
		t.Fatalf("Failed to stop API server gracefully: %v", err)
	}
}

func TestLoanApprove(t *testing.T) {
	app := setupTestApp()
	// T7 => Approve a loan, with admin credentials
	loanData, _, err := createLoan(t, app, "1000", "3")
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPatch, "/loan/"+loanData.LoanID+"/approve", nil)
	getAdminAuth(req)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal("Error occurred:", err)
	}
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Stop the API server gracefully
	if err := app.Shutdown(); err != nil {
		t.Fatalf("Failed to stop API server gracefully: %v", err)
	}
}

func TestLoanRepay(t *testing.T) {
	app := setupTestApp()
	// T8 => Repay a loan term
	loanData, _, err := createLoan(t, app, "1000", "3")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/loan/"+loanData.LoanID+"/approve", nil)
	getAdminAuth(req)
	_, err = app.Test(req)
	if err != nil {
		t.Fatal("Error occurred:", err)
	}

	payload := `{"loan_id":"` + loanData.LoanID + `","repayment_amount":333.33}`
	req = httptest.NewRequest(http.MethodPatch, "/loan/"+loanData.CustomerID+"/repay", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	getAuth(req)
	res, err := app.Test(req)
	if err != nil {
		t.Fatal("Error occurred:", err, res)
	}

	newloanData, err := getParsedLoanData(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "PAID", newloanData.Repayments[0].Status)

	// Stop the API server gracefully
	if err := app.Shutdown(); err != nil {
		t.Fatalf("Failed to stop API server gracefully: %v", err)
	}
}

// *======================== HELPERS ========================*

func getAuth(req *http.Request) {
	username := "customer1"
	password := "customer123#1"
	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", basicAuth)
}

func getAdminAuth(req *http.Request) {
	username := "admin1"
	password := "admin123#1"
	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", basicAuth)
}

func getResponseData(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func createLoan(t *testing.T, app *fiber.App, amount string, loanTerm string) (*models.Loan, *http.Response, error) {
	payload := `{"loan_amount":` + amount + `,"loan_term":` + loanTerm + `}`
	req := httptest.NewRequest(http.MethodPost, "/loan/customer1/submit", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	getAuth(req)
	res, err := app.Test(req)
	if err != nil {
		return nil, nil, err
	}
	var loanData *models.Loan
	respStr, err := getResponseData(res)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal([]byte(respStr), &loanData)
	if err != nil {
		return nil, nil, err
	}
	return loanData, res, nil
}

func getParsedLoanDataFromList(res *http.Response) (*models.Loan, error) {
	var loanData *models.Loan
	respStr, err := getResponseData(res)
	if err != nil {
		return loanData, nil
	}
	var loanDataList []*models.Loan
	err = json.Unmarshal([]byte(respStr), &loanDataList)
	if err != nil {
		return loanData, nil
	}
	return loanDataList[0], nil
}

func getParsedLoanData(res *http.Response) (*models.Loan, error) {
	var loanData *models.Loan
	respStr, err := getResponseData(res)
	if err != nil {
		return loanData, nil
	}
	var loanDataList *models.Loan
	err = json.Unmarshal([]byte(respStr), &loanDataList)
	if err != nil {
		return loanData, nil
	}
	return loanDataList, nil
}
