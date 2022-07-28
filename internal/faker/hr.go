package faker

import "math/rand"

var (
	hrNums      = "0123456789"
	hrAreaCodes = []string{"1", "20", "21", "22", "23", "31", "32", "33", "34", "35", "40", "42", "43", "44", "47", "48", "49", "51", "52", "53"}
)

func hrAreaCode() string {
	return hrAreaCodes[rand.Intn(len(hrAreaCodes))]
}

func HrPhoneNumber() string {
	phone := hrAreaCode()
	for i := 0; i < 7; i++ {
		phone += string(hrNums[rand.Intn(len(hrNums))])
	}
	return phone
}
