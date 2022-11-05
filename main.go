package main

const (
	targetServiceUrl = "https://qrxjmztf2h.execute-api.us-west-2.amazonaws.com/prod"

	minPwdLen = 11 // inclusive
	maxPwdLen = 14 // inclusive
)

func main() {
	attacker := NewPasswordAttacker("0123456789", minPwdLen, maxPwdLen)
	attacker.Attack()
}
