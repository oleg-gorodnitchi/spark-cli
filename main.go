package main

import (
	"encoding/json"
	cli "gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	_ "reflect"
)

func main() {
	const baseHistoryApiUrl = "http://localhost:18080/api/v1/"
	cliApp := &cli.App{
		Name:        "spark-cli",
		Usage:       "CLI for Apache Spark REST API",
		Version:     "0.1.0",
		Description: "Fetches data from the Spark History Server REST API.",
		Authors: []*cli.Author{
			{
				Name:  "Aravind R. Yarram",
				Email: "yaravind@gmail.com",
			},
		},
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "apps",
				Usage: "Lists all Spark applications",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "completed",
						Aliases: []string{"c"},
						Usage:   "Lists all 'completed' spark applications",
					},
					&cli.BoolFlag{
						Name:    "running",
						Aliases: []string{"r"},
						Usage:   "Lists all 'running' spark applications",
					},
				},
				Action: func(c *cli.Context) error {

					log.Printf("Total Args = %d, Args=%s", c.NArg(), c.Args())

					log.Printf("IsSet(Completed) = %t, IsSet(Running) = %t", c.IsSet("completed"), c.IsSet("running"))

					var url string = baseHistoryApiUrl + "applications"

					if c.IsSet("completed") {
						log.Println("Listing all 'completed' applications")

						url = url + "?status=completed"
						respStr := getAsStr(url)
						log.Println(respStr)
					} else if c.IsSet("running") {
						log.Println("Listing all 'running' applications")

						url = url + "?status=running"
						respStr := getAsStr(url)
						log.Println(respStr)
					} else {
						log.Println("Listing all applications")
						var apps []Apps
						if respBuff, err := get(url); err == nil {
							//log.Println(string(respBuff))
							if jsonErr := json.Unmarshal(respBuff, &apps); jsonErr == nil {
								//log.Println(apps)
								cntTot, cntCompleted, cntIncomplete := Summary(apps)
								log.Printf("Total: %d (Completed: %d, Incomplete: %d)", cntTot, cntCompleted, cntIncomplete)
							} else {
								log.Fatal(jsonErr)
							}
						}
					}
					return nil
				},
			},
		},
	}

	cliApp.Run(os.Args)
}
func Summary(apps []Apps) (cntTot, cntCompleted, cntIncomplete int) {
	cntTot = len(apps)
	for _, app := range apps {
		if app.Attempts[0].IsCompleted {
			cntCompleted++
		} else {
			cntIncomplete++
		}
	}
	return
}

type Attempt struct {
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
	LastUpdated      string `json:"lastUpdated"`
	Duration         uint32 `json:"duration"`
	SparkUser        string `json:"sparkUser"`
	IsCompleted      bool   `json:"completed"`
	LastUpdatedEpoch int64  `json:"lastUpdatedEpoch"`
	StartTimeEpoch   int64  `json:"startTimeEpoch"`
	EndTimeEpoch     int64  `json:"EndTimeEpoch"`
}

type Apps struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Attempts []Attempt `json:"attempts"`
}

func get(url string) ([]byte, error) {
	log.Printf("GET %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		defer resp.Body.Close()
		respBuff, _ := ioutil.ReadAll(resp.Body)
		return respBuff, nil
	}
}

func getAsStr(url string) string {
	if respBuff, err := get(url); err != nil {
		return string(respBuff)
	}
	return "{}"
}
