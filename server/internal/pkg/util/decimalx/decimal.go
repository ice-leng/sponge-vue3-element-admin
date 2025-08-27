package decimalx

import "github.com/shopspring/decimal"

func Ratio(today, last decimal.Decimal) decimal.Decimal {
	if last.IsZero() {
		return decimal.Zero
	}
	hundred := decimal.NewFromFloat(100)
	if today.IsZero() {
		return hundred
	}
	return last.Sub(today).Mul(hundred).Div(last).Round(2)
}
