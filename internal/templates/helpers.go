package templates

import "strconv"

func itoa(id int32) string {
	return strconv.Itoa(int(id))
}

func formatPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}
