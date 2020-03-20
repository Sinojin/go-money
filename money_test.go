package money

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	m := New(1, "EUR")

	if m.AmountData.Val != 1 {
		t.Errorf("Expected %d got %d", 1, m.AmountData.Val)
	}

	if m.CurrencyData.Code != "EUR" {
		t.Errorf("Expected CurrencyData %s got %s", "EUR", m.CurrencyData.Code)
	}

	m = New(-100, "EUR")

	if m.AmountData.Val != -100 {
		t.Errorf("Expected %d got %d", -100, m.AmountData.Val)
	}
}

func TestCurrencyData(t *testing.T) {
	code := "MOCK"
	decimals := 5
	AddCurrency(code, "M$", "1 $", ".", ",", decimals)
	m := New(1, code)
	c := m.Currency().Code
	if c != code {
		t.Errorf("Expected %s got %s", code, c)
	}
	f := m.Currency().Fraction
	if f != decimals {
		t.Errorf("Expected %d got %d", decimals, f)
	}
}

func TestMoney_SameCurrencyData(t *testing.T) {
	m := New(0, "EUR")
	om := New(0, "USD")

	if m.SameCurrencyData(om) {
		t.Errorf("Expected %s not to be same as %s", m.CurrencyData.Code, om.CurrencyData.Code)
	}

	om = New(0, "EUR")

	if !m.SameCurrencyData(om) {
		t.Errorf("Expected %s to be same as %s", m.CurrencyData.Code, om.CurrencyData.Code)
	}
}

func TestMoney_Equals(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		AmountData int64
		expected   bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.AmountData, "EUR")
		r, err := m.Equals(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equals %d == %t got %t", m.AmountData.Val,
				om.AmountData.Val, tc.expected, r)
		}
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		AmountData int64
		expected   bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.AmountData, "EUR")
		r, err := m.GreaterThan(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Greater Than %d == %t got %t", m.AmountData.Val,
				om.AmountData.Val, tc.expected, r)
		}
	}
}

func TestMoney_GreaterThanOrEqual(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		AmountData int64
		expected   bool
	}{
		{-1, true},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.AmountData, "EUR")
		r, err := m.GreaterThanOrEqual(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equals Or Greater Than %d == %t got %t", m.AmountData.Val,
				om.AmountData.Val, tc.expected, r)
		}
	}
}

func TestMoney_LessThan(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		AmountData int64
		expected   bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		om := New(tc.AmountData, "EUR")
		r, err := m.LessThan(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Less Than %d == %t got %t", m.AmountData.Val,
				om.AmountData.Val, tc.expected, r)
		}
	}
}

func TestMoney_LessThanOrEqual(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		AmountData int64
		expected   bool
	}{
		{-1, false},
		{0, true},
		{1, true},
	}

	for _, tc := range tcs {
		om := New(tc.AmountData, "EUR")
		r, err := m.LessThanOrEqual(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equal Or Less Than %d == %t got %t", m.AmountData.Val,
				om.AmountData.Val, tc.expected, r)
		}
	}
}

func TestMoney_IsZero(t *testing.T) {
	tcs := []struct {
		AmountData int64
		expected   bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, "EUR")
		r := m.IsZero()

		if r != tc.expected {
			t.Errorf("Expected %d to be zero == %t got %t", m.AmountData.Val, tc.expected, r)
		}
	}
}

func TestMoney_IsNegative(t *testing.T) {
	tcs := []struct {
		AmountData int64
		expected   bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, "EUR")
		r := m.IsNegative()

		if r != tc.expected {
			t.Errorf("Expected %d to be negative == %t got %t", m.AmountData.Val,
				tc.expected, r)
		}
	}
}

func TestMoney_IsPositive(t *testing.T) {
	tcs := []struct {
		AmountData int64
		expected   bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, "EUR")
		r := m.IsPositive()

		if r != tc.expected {
			t.Errorf("Expected %d to be positive == %t got %t", m.AmountData.Val,
				tc.expected, r)
		}
	}
}

func TestMoney_Absolute(t *testing.T) {
	tcs := []struct {
		AmountData int64
		expected   int64
	}{
		{-1, 1},
		{0, 0},
		{1, 1},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, "EUR")
		r := m.Absolute().AmountData.Val

		if r != tc.expected {
			t.Errorf("Expected absolute %d to be %d got %d", m.AmountData.Val,
				tc.expected, r)
		}
	}
}

func TestMoney_Negative(t *testing.T) {
	tcs := []struct {
		AmountData int64
		expected   int64
	}{
		{-1, -1},
		{0, -0},
		{1, -1},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, "EUR")
		r := m.Negative().AmountData.Val

		if r != tc.expected {
			t.Errorf("Expected absolute %d to be %d got %d", m.AmountData.Val,
				tc.expected, r)
		}
	}
}

func TestMoney_Add(t *testing.T) {
	tcs := []struct {
		AmountData1 int64
		AmountData2 int64
		expected    int64
	}{
		{5, 5, 10},
		{10, 5, 15},
		{1, -1, 0},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData1, "EUR")
		om := New(tc.AmountData2, "EUR")
		r, err := m.Add(om)

		if err != nil {
			t.Error(err)
		}

		if r.Amount() != tc.expected {
			t.Errorf("Expected %d + %d = %d got %d", tc.AmountData1, tc.AmountData2,
				tc.expected, r.AmountData.Val)
		}
	}
}

func TestMoney_Add2(t *testing.T) {
	m := New(100, "EUR")
	dm := New(100, "GBP")
	r, err := m.Add(dm)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Subtract(t *testing.T) {
	tcs := []struct {
		AmountData1 int64
		AmountData2 int64
		expected    int64
	}{
		{5, 5, 0},
		{10, 5, 5},
		{1, -1, 2},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData1, "EUR")
		om := New(tc.AmountData2, "EUR")
		r, err := m.Subtract(om)

		if err != nil {
			t.Error(err)
		}

		if r.AmountData.Val != tc.expected {
			t.Errorf("Expected %d - %d = %d got %d", tc.AmountData1, tc.AmountData2,
				tc.expected, r.AmountData.Val)
		}
	}
}

func TestMoney_Subtract2(t *testing.T) {
	m := New(100, "EUR")
	dm := New(100, "GBP")
	r, err := m.Subtract(dm)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Multiply(t *testing.T) {
	tcs := []struct {
		AmountData int64
		multiplier int64
		expected   int64
	}{
		{5, 5, 25},
		{10, 5, 50},
		{1, -1, -1},
		{1, 0, 0},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, "EUR")
		r := m.Multiply(tc.multiplier).AmountData.Val

		if r != tc.expected {
			t.Errorf("Expected %d * %d = %d got %d", tc.AmountData, tc.multiplier, tc.expected, r)
		}
	}
}

func TestMoney_Round(t *testing.T) {
	tcs := []struct {
		AmountData int64
		expected   int64
	}{
		{125, 100},
		{175, 200},
		{349, 300},
		{351, 400},
		{0, 0},
		{-1, 0},
		{-75, -100},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, "EUR")
		r := m.Round().AmountData.Val

		if r != tc.expected {
			t.Errorf("Expected rounded %d to be %d got %d", tc.AmountData, tc.expected, r)
		}
	}
}

func TestMoney_RoundWithExponential(t *testing.T) {
	tcs := []struct {
		AmountData int64
		expected   int64
	}{
		{12555, 13000},
	}

	for _, tc := range tcs {
		AddCurrency("CUR", "*", "$1", ".", ",", 3)
		m := New(tc.AmountData, "CUR")
		r := m.Round().AmountData.Val

		if r != tc.expected {
			t.Errorf("Expected rounded %d to be %d got %d", tc.AmountData, tc.expected, r)
		}
	}
}

func TestMoney_Split(t *testing.T) {
	tcs := []struct {
		AmountData int64
		split      int
		expected   []int64
	}{
		{100, 3, []int64{34, 33, 33}},
		{100, 4, []int64{25, 25, 25, 25}},
		{5, 3, []int64{2, 2, 1}},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, "EUR")
		var rs []int64
		split, _ := m.Split(tc.split)

		for _, party := range split {
			rs = append(rs, party.AmountData.Val)
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected split of %d to be %v got %v", tc.AmountData, tc.expected, rs)
		}
	}
}

func TestMoney_Split2(t *testing.T) {
	m := New(100, "EUR")
	r, err := m.Split(-10)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Allocate(t *testing.T) {
	tcs := []struct {
		AmountData int64
		ratios     []int
		expected   []int64
	}{
		{100, []int{50, 50}, []int64{50, 50}},
		{100, []int{30, 30, 30}, []int64{34, 33, 33}},
		{200, []int{25, 25, 50}, []int64{50, 50, 100}},
		{5, []int{50, 25, 25}, []int64{3, 1, 1}},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, "EUR")
		var rs []int64
		split, _ := m.Allocate(tc.ratios...)

		for _, party := range split {
			rs = append(rs, party.AmountData.Val)
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected allocation of %d for ratios %v to be %v got %v", tc.AmountData, tc.ratios,
				tc.expected, rs)
		}
	}
}

func TestMoney_Allocate2(t *testing.T) {
	m := New(100, "EUR")
	r, err := m.Allocate()

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Format(t *testing.T) {
	tcs := []struct {
		AmountData int64
		code       string
		expected   string
	}{
		{100, "GBP", "£1.00"},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.AmountData, tc.expected, r)
		}
	}
}

func TestMoney_Display(t *testing.T) {
	tcs := []struct {
		AmountData int64
		code       string
		expected   string
	}{
		{100, "AED", "1.00 .\u062f.\u0625"},
		{1, "USD", "$0.01"},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.AmountData, tc.expected, r)
		}
	}
}

func TestMoney_AsMajorUnits(t *testing.T) {
	tcs := []struct {
		AmountData int64
		code       string
		expected   float64
	}{
		{100, "AED", 1.00},
		{1, "USD", 0.01},
	}

	for _, tc := range tcs {
		m := New(tc.AmountData, tc.code)
		r := m.AsMajorUnits()

		if r != tc.expected {
			t.Errorf("Expected Value as major units of %d to be %f got %f", tc.AmountData, tc.expected, r)
		}
	}
}

func TestMoney_Allocate3(t *testing.T) {
	pound := New(100, "GBP")
	parties, err := pound.Allocate(33, 33, 33)

	if err != nil {
		t.Error(err)
	}

	if parties[0].Display() != "£0.34" {
		t.Errorf("Expected %s got %s", "£0.34", parties[0].Display())
	}

	if parties[1].Display() != "£0.33" {
		t.Errorf("Expected %s got %s", "£0.33", parties[1].Display())
	}

	if parties[2].Display() != "£0.33" {
		t.Errorf("Expected %s got %s", "£0.33", parties[2].Display())
	}
}

func TestMoney_Comparison(t *testing.T) {
	pound := New(100, "GBP")
	twoPounds := New(200, "GBP")
	twoEuros := New(200, "EUR")

	if r, err := pound.GreaterThan(twoPounds); err != nil || r {
		t.Errorf("Expected %d Greater Than %d == %t got %t", pound.AmountData.Val,
			twoPounds.AmountData.Val, false, r)
	}

	if r, err := pound.LessThan(twoPounds); err != nil || !r {
		t.Errorf("Expected %d Less Than %d == %t got %t", pound.AmountData.Val,
			twoPounds.AmountData.Val, true, r)
	}

	if r, err := pound.LessThan(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.GreaterThan(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.Equals(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.LessThanOrEqual(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.GreaterThanOrEqual(twoEuros); err == nil || r {
		t.Error("Expected err")
	}
}

func TestMoney_CurrencyData(t *testing.T) {
	pound := New(100, "GBP")

	if pound.Currency().Code != "GBP" {
		t.Errorf("Expected %s got %s", "GBP", pound.Currency().Code)
	}
}

func TestMoney_AmountData(t *testing.T) {
	pound := New(100, "GBP")

	if pound.Amount() != 100 {
		t.Errorf("Expected %d got %d", 100, pound.Amount())
	}
}

func TestDefaultMarshal(t *testing.T) {
	given := New(12345, "IQD")
	expected := `{"AmountData":12345,"CurrencyData":"IQD"}`

	b, err := json.Marshal(given)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}
}

func TestCustomMarshal(t *testing.T) {
	given := New(12345, "IQD")
	expected := `{"AmountData":12345,"CurrencyData_code":"IQD","CurrencyData_fraction":3}`
	MarshalJSON = func(m Money) ([]byte, error) {
		buff := bytes.NewBufferString(fmt.Sprintf(`{"AmountData": %d, "CurrencyData_code": "%s", "CurrencyData_fraction": %d}`, m.Amount(), m.Currency().Code, m.Currency().Fraction))
		return buff.Bytes(), nil
	}

	b, err := json.Marshal(given)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}
}

func TestDefaultUnmarshal(t *testing.T) {
	given := `{"AmountData": 10012, "CurrencyData":"USD"}`
	expected := "$100.12"
	var m Money
	err := json.Unmarshal([]byte(given), &m)
	if err != nil {
		t.Error(err)
	}

	if m.Display() != expected {
		t.Errorf("Expected %s got %s", expected, m.Display())
	}
}

func TestCustomUnmarshal(t *testing.T) {
	given := `{"AmountData": 10012, "CurrencyData_code":"USD", "CurrencyData_fraction":2}`
	expected := "$100.12"
	UnmarshalJSON = func(m *Money, b []byte) error {
		data := make(map[string]interface{})
		err := json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
		ref := New(int64(data["AmountData"].(float64)), data["CurrencyData_code"].(string))
		*m = *ref
		return nil
	}

	var m Money
	err := json.Unmarshal([]byte(given), &m)
	if err != nil {
		t.Error(err)
	}

	if m.Display() != expected {
		t.Errorf("Expected %s got %s", expected, m.Display())
	}
}
