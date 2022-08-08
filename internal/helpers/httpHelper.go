package helpers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	s "github.com/suprnova/src-webhook-fetcher/internal/structs"
)

var client http.Client = http.Client{
	Timeout: time.Second * 5,
}

func New(c chan s.Response) {
	r, err := client.Do(createReq("?status=new&orderby=submitted&direction=desc"))
	if err != nil {
		log.Panic(err)
	}
	c <- parseRes(r)
}

func Verified(c chan s.Response) {
	r, err := client.Do(createReq("?status=verified&orderby=verify-date&direction=desc"))
	if err != nil {
		log.Panic(err)
	}
	c <- parseRes(r)
}

func Rejected(c chan s.Response) {
	r, err := client.Do(createReq("?status=rejected&orderby=verify-date&direction=desc"))
	if err != nil {
		log.Panic(err)
	}
	c <- parseRes(r)
}

func createReq(queries string) *http.Request {
	url := "https://speedrun.com/api/v1/runs"
	// this nanosecond formatting has to be done because of src's aggressive caching behavior
	req, err := http.NewRequest("GET", url+queries+"&vary="+strconv.FormatInt(int64(time.Now().Nanosecond()), 10), nil)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("User-Agent", "SRCStats Webhook")
	return req
}

func parseRes(r *http.Response) s.Response {
	if r.StatusCode == 400 || r.StatusCode == 404 {
		// this doesnt *have* to be a panic, could just return an empty Response
		log.Panic("Server returned a failure!")
	}
	result, err := io.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}
	r.Body.Close()
	var res s.Response
	json.Unmarshal(result, &res)
	return res
}
