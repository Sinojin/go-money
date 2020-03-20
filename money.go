package money

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// Injection points for backward compatibility.
// If you need to keep your JSON marshal/unmarshal way, overwrite them like below.
//   money.UnmarshalJSON = func (m *Money, b []byte) error { ... }
//   money.MarshalJSON = func (m Money) ([]byte, error) { ... }
var (
	// UnmarshalJSONFunc is injection point of json.Unmarshaller for money.Money
	UnmarshalJSON = defaultUnmarshalJSON
	// MarshalJSONFunc is injection point of json.Marshaller for money.Money
	MarshalJSON = defaultMarshalJSON
)

func defaultUnmarshalJSON(m *Money, b []byte) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	ref := New(int64(data["AmountData"].(float64)), data["CurrencyData"].(string))
	*m = *ref
	return nil
}

func defaultMarshalJSON(m Money) ([]byte, error) {
	buff := bytes.NewBufferString(fmt.Sprintf(`{"AmountData": %d, "CurrencyData": "%s"}`, m.Amount(), m.Currency().Code))
	return buff.Bytes(), nil
}

// AmountData is a datastructure that stores the AmountData being used for calculations.
type Amount struct {
	Val int64
}

// Money represents monetary Value information, stores
// CurrencyData and AmountData Value.
type Money struct {
	AmountData   *Amount
	CurrencyData *Currency
}

// New creates and returns new instance of Money.
func New(amount int64, code string) *Money {
	return &Money{
		AmountData:   &Amount{Val: amount},
		CurrencyData: newCurrency(code).get(),
	}
}

// CurrencyData returns the CurrencyData used by Money.
func (m *Money) Currency() *Currency {
	return m.CurrencyData
}

// AmountData returns a copy of the internal monetary Value as an int64.
func (m *Money) Amount() int64 {
	return m.AmountData.Val
}

// SameCurrencyData check if given Money is equals by CurrencyData.
func (m *Money) SameCurrencyData(om *Money) bool {
	return m.CurrencyData.equals(om.CurrencyData)
}

func (m *Money) assertSameCurrencyData(om *Money) error {
	if !m.SameCurrencyData(om) {
		return errors.New("currencies don't match")
	}

	return nil
}

func (m *Money) compare(om *Money) int {
	switch {
	case m.AmountData.Val > om.AmountData.Val:
		return 1
	case m.AmountData.Val < om.AmountData.Val:
		return -1
	}

	return 0
}

// Equals checks equality between two Money types.
func (m *Money) Equals(om *Money) (bool, error) {
	if err := m.assertSameCurrencyData(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// GreaterThan checks whether the Value of Money is greater than the other.
func (m *Money) GreaterThan(om *Money) (bool, error) {
	if err := m.assertSameCurrencyData(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual checks whether the Value of Money is greater or equal than the other.
func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrencyData(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// LessThan checks whether the Value of Money is less than the other.
func (m *Money) LessThan(om *Money) (bool, error) {
	if err := m.assertSameCurrencyData(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// LessThanOrEqual checks whether the Value of Money is less or equal than the other.
func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrencyData(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero returns boolean of whether the Value of Money is equals to zero.
func (m *Money) IsZero() bool {
	return m.AmountData.Val == 0
}

// IsPositive returns boolean of whether the Value of Money is positive.
func (m *Money) IsPositive() bool {
	return m.AmountData.Val > 0
}

// IsNegative returns boolean of whether the Value of Money is negative.
func (m *Money) IsNegative() bool {
	return m.AmountData.Val < 0
}

// Absolute returns new Money struct from given Money using absolute monetary Value.
func (m *Money) Absolute() *Money {
	return &Money{AmountData: mutate.calc.absolute(m.AmountData), CurrencyData: m.CurrencyData}
}

// Negative returns new Money struct from given Money using negative monetary Value.
func (m *Money) Negative() *Money {
	return &Money{AmountData: mutate.calc.negative(m.AmountData), CurrencyData: m.CurrencyData}
}

// Add returns new Money struct with Value representing sum of Self and Other Money.
func (m *Money) Add(om *Money) (*Money, error) {
	if err := m.assertSameCurrencyData(om); err != nil {
		return nil, err
	}

	return &Money{AmountData: mutate.calc.add(m.AmountData, om.AmountData), CurrencyData: m.CurrencyData}, nil
}

// Subtract returns new Money struct with Value representing difference of Self and Other Money.
func (m *Money) Subtract(om *Money) (*Money, error) {
	if err := m.assertSameCurrencyData(om); err != nil {
		return nil, err
	}

	return &Money{AmountData: mutate.calc.subtract(m.AmountData, om.AmountData), CurrencyData: m.CurrencyData}, nil
}

// Multiply returns new Money struct with Value representing Self multiplied Value by multiplier.
func (m *Money) Multiply(mul int64) *Money {
	return &Money{AmountData: mutate.calc.multiply(m.AmountData, mul), CurrencyData: m.CurrencyData}
}

// Round returns new Money struct with Value rounded to nearest zero.
func (m *Money) Round() *Money {
	return &Money{AmountData: mutate.calc.round(m.AmountData, m.CurrencyData.Fraction), CurrencyData: m.CurrencyData}
}

// Split returns slice of Money structs with split Self Value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later.
func (m *Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, errors.New("split must be higher than zero")
	}

	a := mutate.calc.divide(m.AmountData, int64(n))
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{AmountData: a, CurrencyData: m.CurrencyData}
	}

	l := mutate.calc.modulus(m.AmountData, int64(n)).Val

	// Add leftovers to the first parties.
	for p := 0; l != 0; p++ {
		ms[p].AmountData = mutate.calc.add(ms[p].AmountData, &Amount{1})
		l--
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self Value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m *Money) Allocate(rs ...int) ([]*Money, error) {
	if len(rs) == 0 {
		return nil, errors.New("no ratios specified")
	}

	// Calculate sum of ratios.
	var sum int
	for _, r := range rs {
		sum += r
	}

	var total int64
	ms := make([]*Money, 0, len(rs))
	for _, r := range rs {
		party := &Money{
			AmountData:   mutate.calc.allocate(m.AmountData, r, sum),
			CurrencyData: m.CurrencyData,
		}

		ms = append(ms, party)
		total += party.AmountData.Val
	}

	// Calculate leftover Value and divide to first parties.
	lo := m.AmountData.Val - total
	sub := int64(1)
	if lo < 0 {
		sub = -sub
	}

	for p := 0; lo != 0; p++ {
		ms[p].AmountData = mutate.calc.add(ms[p].AmountData, &Amount{sub})
		lo -= sub
	}

	return ms, nil
}

// Display lets represent Money struct as string in given CurrencyData Value.
func (m *Money) Display() string {
	c := m.CurrencyData.get()
	return c.Formatter().Format(m.AmountData.Val)
}

// AsMajorUnits lets represent Money struct as subunits (float64) in given CurrencyData Value
func (m *Money) AsMajorUnits() float64 {
	c := m.CurrencyData.get()
	return c.Formatter().ToMajorUnits(m.AmountData.Val)
}

// UnmarshalJSON is implementation of json.Unmarshaller
func (m *Money) UnmarshalJSON(b []byte) error {
	return UnmarshalJSON(m, b)
}

// MarshalJSON is implementation of json.Marshaller
func (m Money) MarshalJSON() ([]byte, error) {
	return MarshalJSON(m)
}
