package money

import "math"

type calculator struct{}

func (c *calculator) add(a, b *Amount) *Amount {
	return &Amount{a.Val + b.Val}
}

func (c *calculator) subtract(a, b *Amount) *Amount {
	return &Amount{a.Val - b.Val}
}

func (c *calculator) multiply(a *Amount, m int64) *Amount {
	return &Amount{a.Val * m}
}

func (c *calculator) divide(a *Amount, d int64) *Amount {
	return &Amount{a.Val / d}
}

func (c *calculator) modulus(a *Amount, d int64) *Amount {
	return &Amount{a.Val % d}
}

func (c *calculator) allocate(a *Amount, r, s int) *Amount {
	return &Amount{a.Val * int64(r) / int64(s)}
}

func (c *calculator) absolute(a *Amount) *Amount {
	if a.Val < 0 {
		return &Amount{-a.Val}
	}

	return &Amount{a.Val}
}

func (c *calculator) negative(a *Amount) *Amount {
	if a.Val > 0 {
		return &Amount{-a.Val}
	}

	return &Amount{a.Val}
}

func (c *calculator) round(a *Amount, e int) *Amount {
	if a.Val == 0 {
		return &Amount{0}
	}

	absam := c.absolute(a)
	exp := int64(math.Pow(10, float64(e)))
	m := absam.Val % exp

	if m > (exp / 2) {
		absam.Val += exp
	}

	absam.Val = (absam.Val / exp) * exp

	if a.Val < 0 {
		a.Val = -absam.Val
	} else {
		a.Val = absam.Val
	}

	return &Amount{a.Val}
}
