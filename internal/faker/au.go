package faker

import (
	"fmt"
	"math/rand"
)

var (
	auFirst  = "238"
	auSecond = "6789"
	auThird  = "12345678"
	auRest   = "123456789"
)

func auPhoneFirst() string {
	return string(auFirst[rand.Intn(len(auFirst))])
}

func auPhoneSecond() string {
	return string(auSecond[rand.Intn(len(auSecond))])
}

func auPhoneThird() string {
	return string(auThird[rand.Intn(len(auThird))])
}

func auPhoneRest() string {
	part := ""

	for i := 0; i < 6; i++ {
		part += string(auRest[rand.Intn(len(auRest))])
	}

	return part

}

func AuPhoneNumber() string {
	return fmt.Sprintf("%s%s%s%s", auPhoneFirst(), auPhoneSecond(), auPhoneThird(), auPhoneRest())
}
