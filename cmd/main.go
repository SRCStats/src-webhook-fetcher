package main

import (
	"fmt"
	"sync"

	h "github.com/srcstats/src-webhook-fetcher/internal/helpers"
	s "github.com/srcstats/src-webhook-fetcher/internal/structs"
)

var wg sync.WaitGroup

func main() {
	chNew, chVerified, chRejected := make(chan s.Response), make(chan s.Response), make(chan s.Response)
	go h.New(chNew)
	go h.Verified(chVerified)
	go h.Rejected(chRejected)
	for {
		if chNew == nil && chVerified == nil && chRejected == nil {
			break
		}
		select {
		case res := <-chNew:
			go handleRuns(&res, "new")
			chNew = nil
		case res := <-chVerified:
			go handleRuns(&res, "verified")
			chVerified = nil
		case res := <-chRejected:
			go handleRuns(&res, "rejected")
			chRejected = nil
		}
	}
	wg.Wait()
}

func handleRuns(r *s.Response, scope string) {
	wg.Add(1)
	runs := s.MakeRuns(*r, scope)
	runs = *update(&runs, scope)
	if len(runs.Data) != 0 {
		h.Delete(scope)
		h.Create(runs)
	}
	wg.Done()
}

func update(runs *s.Response, scope string) *s.Response {
	exists := false
	var newRuns []s.Data
	existingRuns := h.List(scope).Data
	for _, run := range runs.Data {
		for _, xRun := range existingRuns {
			if run.ID == xRun.ID {
				exists = true
				break
			}
		}
		if exists {
			break
		}
		run.New = true
		newRuns = append(newRuns, run)
	}
	fmt.Printf("New runs found for scope %v: %v\n", scope, len(newRuns))
	i := len(newRuns)
	if i == 0 {
		return &s.Response{Data: newRuns, Scope: scope}
	}
	for i < 20 {
		(existingRuns)[i].Order = i
		(existingRuns)[i].New = false
		newRuns = append(newRuns, (existingRuns)[i])
		i++
	}
	return &s.Response{Data: newRuns, Scope: scope}
}
