package helpers

import (
	"time"
)

type Response struct {
	Data []struct {
		ID       string `json:"id,omitempty"`
		Weblink  string `json:"weblink,omitempty"`
		Game     string `json:"game,omitempty"`
		Level    string `json:"level,omitempty"`
		Category string `json:"category,omitempty"`
		Videos   struct {
			Links []struct {
				URI string `json:"uri,omitempty"`
			} `json:"links,omitempty"`
		} `json:"videos,omitempty"`
		Comment interface{} `json:"comment,omitempty"`
		Status  struct {
			Status string `json:"status,omitempty"`
		} `json:"status,omitempty"`
		Players []struct {
			Rel string `json:"rel,omitempty"`
			ID  string `json:"id,omitempty"`
			URI string `json:"uri,omitempty"`
		} `json:"players,omitempty"`
		Date      string    `json:"date,omitempty"`
		Submitted time.Time `json:"submitted,omitempty"`
		Times     struct {
			Primary          string      `json:"primary,omitempty"`
			PrimaryT         float64     `json:"primary_t,omitempty"`
			Realtime         string      `json:"realtime,omitempty"`
			RealtimeT        float64     `json:"realtime_t,omitempty"`
			RealtimeNoloads  interface{} `json:"realtime_noloads,omitempty"`
			RealtimeNoloadsT int         `json:"realtime_noloads_t,omitempty"`
			Ingame           interface{} `json:"ingame,omitempty"`
			IngameT          int         `json:"ingame_t,omitempty"`
		} `json:"times,omitempty"`
		System struct {
			Platform string      `json:"platform,omitempty"`
			Emulated bool        `json:"emulated,omitempty"`
			Region   interface{} `json:"region,omitempty"`
		} `json:"system,omitempty"`
		Splits interface{} `json:"splits,omitempty"`
		Values struct {
			E8M53Zxn string `json:"e8m53zxn,omitempty"`
			Gnxjxzgl string `json:"gnxjxzgl,omitempty"`
		} `json:"values,omitempty"`
		Links []struct {
			Rel string `json:"rel,omitempty"`
			URI string `json:"uri,omitempty"`
		} `json:"links,omitempty"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Max    int `json:"max"`
		Size   int `json:"size"`
		Links  []struct {
			Rel string `json:"rel"`
			URI string `json:"uri"`
		} `json:"links"`
	} `json:"pagination"`
}

type Run struct {
	SiteId string `bson:"siteId"`
	Scope  string `bson:"scope"`
	Order  int    `bson:"order"`
}

func MakeRuns(r Response, scope string) *[]Run {
	var list []Run
	for i, r := range r.Data {
		list = append(list, Run{SiteId: r.ID, Scope: scope, Order: i})
	}
	return &list
}
