package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"paymentgateway/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	database = mockdb{}
}

func TestGetPayment(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)

	ctx.Request.Method = "GET"

	u := url.Values{}
	u.Add("id", "1")
	// set query params
	ctx.Request.URL.RawQuery = u.Encode()

	GetPayment(ctx)

	var p model.Payment
	t.Log(w.Body.String())
	json.Unmarshal(w.Body.Bytes(), &p)
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.EqualValues(t, "XXXXXXXXXXXX7567", p.CardNumber)

}

func TestSendPayment(t *testing.T) {
	p := model.Payment{
		CardName:        "processout test",
		CardNumber:      "1234567891234567",
		CardType:        "visa",
		CardExpiryMonth: 11,
		CardExpiryYear:  2025,
		Cvv:             "123",
		Amount:          1000,
		Description:     "testing payment",
	}

	w := MockJsonPost("/payment", SendPayment, p)

	type response struct {
		Id int `json:"paymentID"`
	}

	var resp response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, resp.Id)

}

func TestRegister(t *testing.T) {
	m := model.Merchant{
		Name: "testing",
	}

	w := MockJsonPost("/register", Register, m)

	type response struct {
		Name string `json:"merchant"`
	}

	var resp response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, m.Name, resp.Name)
}

func TestLogin(t *testing.T) {

	m := model.Merchant{
		Name: "test1",
	}

	w := MockJsonPost("/login", Login, m)

	type response struct {
		Id    int    `json:"merchantID"`
		Token string `json:"token"`
	}

	var resp response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, resp.Id)
	assert.NotNil(t, resp.Token)
}

func TestLoginFail(t *testing.T) {
	m := model.Merchant{
		Name: "testfauk",
	}

	w := MockJsonPost("/login", Login, m)

	type response struct {
		Id    int    `json:"merchantID"`
		Token string `json:"token"`
	}

	var resp response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "", resp.Token)
}

// mock gin context
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

func MockJsonPost(endpoint string, f func(*gin.Context), d interface{}) *httptest.ResponseRecorder {
	r := SetUpRouter()
	r.POST(endpoint, f)

	jsonValue, _ := json.Marshal(d)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// mock getrequest
func MockJsonGet(c *gin.Context, params gin.Params, u url.Values) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("id", 1)

	// set path params
	c.Params = params

	// set query params
	c.Request.URL.RawQuery = u.Encode()
}

type mockdb struct {
}

func (d mockdb) SaveUser(m model.Merchant) error {
	if m.Name == "fail_test" {
		return errors.New("failed to save user")
	}

	return nil
}

func (d mockdb) VerifyMerchant(m model.Merchant) (int, bool, error) {
	switch m.Name {
	case "test1":
		return 1, true, nil
	case "test2":
		return 2, true, nil
	default:
		return 0, false, nil
	}

}

func (d mockdb) SavePayment(p model.Payment) error {
	if p.Id == 2 {
		return errors.New("failed to save payment")
	}

	return nil
}

func (d mockdb) GetPayment(id int) (model.Payment, error) {
	if id == 1 {
		return model.Payment{
			Currency:        "GBP",
			CardNumber:      "2323565666547567",
			CardType:        "visa",
			CardExpiryMonth: 11,
			CardExpiryYear:  2025,
		}, nil
	}

	return model.Payment{}, errors.New("couldn't find payment")
}
