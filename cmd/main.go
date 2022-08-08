package main

import (
	"fmt"
	"sync"

	h "github.com/suprnova/src-webhook-fetcher/internal/helpers"
	s "github.com/suprnova/src-webhook-fetcher/internal/structs"
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
	runs = update(runs, scope)
	if len(*runs) != 0 {
		// this is inefficient, probably better to shift the orders instead of deleting all of them
		h.Delete(scope)
		fmt.Println(runs)
		h.Create(runs)
	}
	wg.Done()
}

func update(runs *[]s.Run, scope string) *[]s.Run {
	exists := false
	var newRuns []s.Run
	existingRuns := h.List(scope)
	for _, run := range *runs {
		for _, xRun := range *existingRuns {
			if run.SiteId == xRun.SiteId {
				exists = true
				break
			}
		}
		if exists {
			break
		}
		newRuns = append(newRuns, run)
	}
	fmt.Printf("New runs found for scope %v: %v\n", scope, len(newRuns))
	i := len(newRuns)
	if i == 0 {
		return &newRuns
	}
	for i < 20 {
		(*existingRuns)[i].Order = i
		newRuns = append(newRuns, (*existingRuns)[i])
		i++
	}
	return &newRuns
}
