package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	database   string
	collection string
)

var client http.Client = http.Client{
	Timeout: time.Second * 15,
}

var wg sync.WaitGroup

const (
	connectionStringEnv = "SRC_WEBHOOK_MONGODB_CONNECTION_STRING"
	databaseNameEnv     = "SRC_WEBHOOK_DATABASE"
	collectionNameEnv   = "SRC_WEBHOOK_COLLECTION"
	postRunUri          = "SRC_WEBHOOK_RUN_URI"
)

type Data struct {
	ID      string `json:"id,omitempty"`
	Order   int    `json:"order,omitempty"`
	New     bool   `json:"new,omitempty"`
	Weblink string `json:"weblink,omitempty"`
	Game    struct {
		Data struct {
			ID    string `json:"id,omitempty"`
			Names struct {
				International string `json:"international,omitempty"`
				Japanese      string `json:"japanese,omitempty"`
				Twitch        string `json:"twitch,omitempty"`
			} `json:"names,omitempty"`
			Name         string   `json:"name,omitempty"`
			Abbreviation string   `json:"abbreviation,omitempty"`
			Platforms    []string `json:"platforms,omitempty"`
			Regions      []string `json:"regions,omitempty"`
			// im not sure this works
			Moderators []string `json:"moderators,omitempty"`
			Assets     struct {
				Trophy1st struct {
					URI string `json:"uri,omitempty"`
				} `json:"trophy-1st,omitempty"`
				Trophy2nd struct {
					URI string `json:"uri,omitempty"`
				} `json:"trophy-2nd,omitempty"`
				Trophy3rd struct {
					URI string `json:"uri,omitempty"`
				} `json:"trophy-3rd,omitempty"`
				Trophy4th struct {
					URI string `json:"uri,omitempty"`
				} `json:"trophy-4th,omitempty"`
			} `json:"assets,omitempty"`
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
			Variables     struct {
				Data []struct {
					ID       string      `json:"id,omitempty"`
					Name     string      `json:"name,omitempty"`
					Category interface{} `json:"category,omitempty"`
					Scope    struct {
						Type string `json:"type,omitempty"`
					} `json:"scope,omitempty"`
					Mandatory     bool                   `json:"mandatory,omitempty"`
					UserDefined   bool                   `json:"user-defined,omitempty"`
					Obsoletes     bool                   `json:"obsoletes,omitempty"`
					Values        map[string]interface{} `json:"values,omitempty"`
					IsSubcategory bool                   `json:"is-subcategory,omitempty"`
				} `json:"data,omitempty"`
			} `json:"variables,omitempty"`
		} `json:"data,omitempty"`
	} `json:"category,omitempty"`
	Times struct {
		Primary          float64 `json:"primary_t,omitempty"`
		RealTime         float64 `json:"realtime_t,omitempty"`
		RealTimeLoadless float64 `json:"realtime_noloads_t,omitempty"`
		InGameTime       float64 `json:"ingame_t,omitempty"`
	}
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
	Players struct {
		Data []struct {
			ID      string `json:"id,omitempty"`
			Weblink string `json:"weblink,omitempty"`
			Names   struct {
				International string `json:"international,omitempty"`
				Japanese      string `json:"japanese,omitempty"`
			} `json:"names,omitempty"`
			Assets struct {
				Icon struct {
					URI string `json:"uri,omitempty"`
				} `json:"icon,omitempty"`
				Image struct {
					URI string `json:"uri,omitempty"`
				} `json:"image,omitempty"`
			} `json:"assets,omitempty"`
		} `json:"data,omitempty"`
	} `json:"players,omitempty"`
	Date      string    `json:"date,omitempty"`
	Submitted time.Time `json:"submitted,omitempty"`
	System    struct {
		Platform string `json:"platform,omitempty"`
		Emulated bool   `json:"emulated,omitempty"`
		Region   string `json:"region,omitempty"`
	} `json:"system,omitempty"`
	Values map[string]string `json:"values,omitempty"`
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

func Main() {
	chNew, chVerified, chRejected := make(chan Response), make(chan Response), make(chan Response)
	go New(chNew)
	go Verified(chVerified)
	go Rejected(chRejected)
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

func handleRuns(r *Response, scope string) {
	runs := MakeRuns(*r, scope)
	runs = *update(&runs, scope)
	if len(runs.Data) != 0 {
		Delete(scope)
		Create(runs)
	}
	wg.Done()
}

func update(runs *Response, scope string) *Response {
	exists := false
	var newRuns []Data
	existingRuns := List(scope).Data
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
		return &Response{Data: newRuns, Scope: scope}
	}
	body, _ := json.Marshal(newRuns)
	http.Post(os.Getenv(postRunUri), "application/json", bytes.NewBuffer(body))
	for i < 20 {
		(existingRuns)[i].Order = i
		(existingRuns)[i].New = false
		newRuns = append(newRuns, (existingRuns)[i])
		i++
	}
	return &Response{Data: newRuns, Scope: scope}
}

func Connect() *mongo.Client {
	connectionString := os.Getenv(connectionStringEnv)
	if connectionString == "" {
		log.Fatal("The database connection string variable is missing!")
	}
	database = os.Getenv(databaseNameEnv)
	if database == "" {
		log.Fatal("The database name variable is missing!")
	}
	collection = os.Getenv(collectionNameEnv)
	if collection == "" {
		log.Fatal("The collection name variable is missing!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	options := options.Client().ApplyURI(connectionString).SetDirect(true)
	c, _ := mongo.NewClient(options)

	err := c.Connect(ctx)
	if err != nil {
		log.Fatalf("Connection couldn't be initialized!\n%v", err)
	}
	err = c.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Connection didn't respond back!\n%v", err)
	}
	return c
}

func Create(runs Response) {
	c := Connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)

	runCollection := c.Database(database).Collection(collection)
	r, err := runCollection.InsertOne(ctx, runs)
	if err != nil {
		log.Fatalf("The runs couldn't be added!\n%v", err)
	}
	fmt.Println("Added runs", r.InsertedID)
}

func List(scope string) *Response {
	c := Connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)
	var runs Response
	runCollection := c.Database(database).Collection(collection)
	err := runCollection.FindOne(ctx, bson.D{{Key: "scope", Value: scope}}).Decode(&runs)

	if err != nil {
		log.Printf("The runs couldn't be listed!\n%v", err)
	}
	return &runs
}

func Delete(scope string) {
	c := Connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)

	runCollection := c.Database(database).Collection(collection)
	res, err := runCollection.DeleteMany(ctx, bson.D{{Key: "scope", Value: scope}})
	if err != nil {
		log.Fatalf("The runs couldn't be deleted!\n%v", err)
	}
	fmt.Printf("Deleted %v entry for scope %v\n", res.DeletedCount, scope)
}

// todo: implement pagination to all of these
func New(c chan Response) {
	wg.Add(1)
	r, err := client.Do(createReq("?status=new&orderby=submitted&direction=desc&embed=category.variables,game,level,region,platform,players"))
	if err != nil {
		log.Panic(err)
	}
	c <- parseRes(r)
}

func Verified(c chan Response) {
	wg.Add(1)
	r, err := client.Do(createReq("?status=verified&orderby=verify-date&direction=desc&embed=category.variables,game,level,region,platform,players"))
	if err != nil {
		log.Panic(err)
	}
	c <- parseRes(r)
}

func Rejected(c chan Response) {
	wg.Add(1)
	r, err := client.Do(createReq("?status=rejected&orderby=verify-date&direction=desc&embed=category.variables,game,level,region,platform,players"))
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

func parseRes(r *http.Response) Response {
	if r.StatusCode == 400 || r.StatusCode == 404 {
		// this doesnt *have* to be a panic, could just return an empty Response
		log.Panic("Server returned a failure!")
	}
	result, err := io.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}
	r.Body.Close()
	var res Response
	json.Unmarshal(result, &res)
	return res
}
