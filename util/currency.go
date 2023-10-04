package util

const (
	USD  = "USD"
	EURO = "EURO"
	CAD  = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EURO, CAD:
		return true
	}
	return false
}
