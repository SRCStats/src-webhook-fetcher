package helpers

import (
	"time"
)

type Data struct {
	ID    string `json:"id,omitempty"`
	Order int    `json:"order,omitempty"`
	New   bool   `json:"new,omitempty"`
	Game  struct {
		Data struct {
			ID    string `json:"id,omitempty"`
			Names struct {
				International string `json:"international,omitempty"`
				Japanese      string `json:"japanese,omitempty"`
				Twitch        string `json:"twitch,omitempty"`
			} `json:"names,omitempty"`
			Abbreviation string   `json:"abbreviation,omitempty"`
			Platforms    []string `json:"platforms,omitempty"`
			Regions      []string `json:"regions,omitempty"`
			// im not sure this works
			Moderators []string `json:"moderators,omitempty"`
		} `json:"data,omitempty"`
	} `json:"game,omitempty"`
	Level struct {
		Data struct {
			ID   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"data,omitempty"`
	} `json:"level,omitempty"`
	Category struct {
		Data struct {
			ID            string `json:"id,omitempty"`
			Name          string `json:"name,omitempty"`
			Type          string `json:"type,omitempty"`
			Miscellaneous bool   `json:"miscellaneous,omitempty"`
		} `json:"data,omitempty"`
	} `json:"category,omitempty"`
	Videos struct {
		Links []struct {
			URI string `json:"uri,omitempty"`
		} `json:"links,omitempty"`
	} `json:"videos,omitempty"`
	Comment string `json:"comment,omitempty"`
	Status  struct {
		Status     string    `json:"status,omitempty"`
		Examiner   string    `json:"examiner,omitempty"`
		Reason     string    `json:"reason,omitempty"`
		VerifyDate time.Time `json:"verify-date,omitempty"`
	} `json:"status,omitempty"`
	Players []struct {
		Rel string `json:"rel,omitempty"`
		ID  string `json:"id,omitempty"`
	} `json:"players,omitempty"`
	Date      string    `json:"date,omitempty"`
	Submitted time.Time `json:"submitted,omitempty"`
	System    struct {
		Platform string `json:"platform,omitempty"`
		Emulated bool   `json:"emulated,omitempty"`
		Region   string `json:"region,omitempty"`
	} `json:"system,omitempty"`
	Region struct {
		Data struct {
			ID   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"data,omitempty"`
	} `json:"region,omitempty"`
	Platform struct {
		Data struct {
			ID   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"data,omitempty"`
	} `json:"platform,omitempty"`
}

type Response struct {
	Data       []Data `json:"data,omitempty"`
	Pagination struct {
		Offset int `json:"offset,omitempty"`
		Max    int `json:"max,omitempty"`
		Size   int `json:"size,omitempty"`
	} `json:"pagination,omitempty"`
	Scope string `json:"scope,omitempty"`
}

func MakeRuns(res Response, scope string) Response {
	for i := range res.Data {
		res.Data[i].Order = i
	}
	res.Scope = scope
	return res
}
