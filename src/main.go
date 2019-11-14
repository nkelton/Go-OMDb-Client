package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "log"
    "encoding/json"
)

type Movie struct {
    Title string `json:"title"`
    Year string `json: "year"`
    Rated string `json: "rated"`
    Released string `json: "released"`
    Runtime string `json: "runtime"`
    Genre string `json: "genre"`
    Director string `json: "director"`
}

type OmdbClient struct {
    Api struct {
        Key string `json: "key"`
        Host string `json: "host"`
    }
    Params struct {
        Id string `json: "id"`
        Title string `json: "Title"`
        Year string `json: "Year"`
        Plot string `json: "Plot"`
    }
    Types struct {
        Movie string `json: "movie"`
        Series string `json: "movie"`
        Episode string `json: "movie"`
    }
}

func (c *OmdbClient) Init(apiKey string) {
    c.Api.Key = apiKey
    c.Api.Host = "http://www.omdbapi.com"

    c.Types.Movie = "movie"
    c.Types.Series = "series"
    c.Types.Episode = "episode"
}

func buildRequestByIdOrTitle (c OmdbClient, contentType string) string {
    req, err := http.NewRequest("GET", c.Api.Host, nil)

    if err != nil {
        log.Fatalf("Creating new Request: %v", err)
    }

    q := req.URL.Query()
    q.Add("apikey", c.Api.Key)
    q.Add("type", contentType)

    if c.Params.Id != "" {
        q.Add("i", c.Params.Id)
    }

    if c.Params.Title != "" {
        q.Add("t", c.Params.Title)
    }

    if c.Params.Year != "" {
        q.Add("y", c.Params.Year)
    }

    if c.Params.Plot != "" {
        q.Add("plot", c.Params.Plot)
    }

    req.URL.RawQuery = q.Encode()
    return req.URL.String()
}

func (c *OmdbClient) getMovieByParams() Movie {
    if c.Params.Id != "" && c.Params.Title != "" {
        panic("Cannot search by Id and Title")
    } else if c.Params.Id == "" && c.Params.Title == "" {
        panic("Must provide search param Id or Title")
    }

    var request string = buildRequestByIdOrTitle(*c, c.Types.Movie)

    response, err := http.Get(request)

    if err != nil {
        panic(err)
    }

    body, _ := ioutil.ReadAll(response.Body)
    movie, err := toMovie([]byte(body))

    if err != nil {
        panic(err)
    }

    return movie
}

func toMovie(body []byte) (Movie, error) {
    m := new(Movie)
    err := json.Unmarshal(body, &m)

    return *m, err
}


func main() {
    var apiKey string = "6d4e8018"
    client := new(OmdbClient)
    client.Init(apiKey)

    client.Params.Title = "Jaws"
    m := client.getMovieByParams()

    fmt.Println(m)
}

