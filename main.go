package main

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	views    = jet.NewHTMLSet("./templates")
	jsonData []byte
	videos   map[string][]byte
	prefix   []string
)

func main() {
	videos = make(map[string][]byte, 5000)
	jsonFile, err := os.Open("./DatabaseBuild/output.server.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jsonData, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = jsonparser.ObjectEach(jsonData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		videos[string(key)] = value
		return nil
	})
	if err != nil {
		return
	}

	rand.Seed(time.Now().Unix())
	prefix = make([]string, 0)
	prefix = append(prefix,
		"decent",
		"light",
		"strong")

	r := gin.Default()

	r.GET("/", index)
	r.GET("/404", error404)
	r.GET("/album/:id", album)
	r.StaticFS("/static", http.Dir("./Static"))
	r.StaticFile("/data.client.json", "./DatabaseBuild/output.client.json")

	r.NoRoute(error404)

	//start server
	err = r.Run(":" + getEnv("PORT", "8080"))
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func index(c *gin.Context) {
	t, _ := views.GetTemplate("index.jet.html")
	vars := make(jet.VarMap)
	c.Writer.WriteHeader(200)
	if err := t.Execute(c.Writer, vars, nil); err != nil {
		log.Println(err)
	}
}

type track struct {
	name     string
	duration string
	id       string
	explict  bool
	fav      bool
	notFav   bool
}

func album(c *gin.Context) {
	t, _ := views.GetTemplate("album.jet.html")
	id := c.Param("id")
	vars := make(jet.VarMap)
	if data, ok := videos[id]; ok {
		albumName, err := jsonparser.GetString(data, "spotify_name")
		if err != nil {
			albumName, err = jsonparser.GetString(data, "album")
			if err != nil {
				albumName, _ = jsonparser.GetString(data, "title")
			}
		}
		artistId := ""
		artistName, err := jsonparser.GetString(data, "spotify_artists", "[0]", "name")
		if err != nil {
			artistName, err = jsonparser.GetString(data, "artist")
			if err != nil {
				artistName = "Could not find artist name"
			}
		} else {
			artistId, err = jsonparser.GetString(data, "spotify_artists", "[0]", "id")
		}
		image, _ := jsonparser.GetString(data, "spotify_obj", "images", "[0]", "url")
		label, err := jsonparser.GetString(data, "spotify_obj", "label")
		if err != nil {
			label, err = jsonparser.GetString(data, "fantano_genre")
		}
		spotifyId, _ := jsonparser.GetString(data, "spotify_obj", "id")
		ratingString, _ := jsonparser.GetString(data, "rating")
		var ratings = strings.Split(strings.Split(ratingString, ",")[0], "/")[0]
		var ratingNumber = -1
		var ratingUrl = ""
		if val, err := strconv.Atoi(ratings); err == nil {
			ratingNumber = val
		}
		if ratingNumber == 10 {
			ratingUrl = "/static/10.png"
		} else if ratingNumber < 10 && ratingNumber > -1 {
			ratingUrl = fmt.Sprintf("/static/%s%d.png", prefix[rand.Intn(len(prefix))], ratingNumber)
		}
		tracks := make([]track, 0)
		_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			name, _ := jsonparser.GetString(value, "name")
			duration, _ := jsonparser.GetInt(value, "duration_ms")
			id, _ := jsonparser.GetString(value, "id")
			explicit, _ := jsonparser.GetBoolean(value, "explicit")
			fav, _ := jsonparser.GetBoolean(value, "is_fav")
			leastFav, _ := jsonparser.GetBoolean(value, "least_fav")
			tracks = append(tracks, track{
				name,
				time.Duration(1000000 * duration).Round(1000000000).String(),
				id,
				explicit,
				fav,
				leastFav,
			})
		}, "spotify_obj", "tracks", "items")
		if err != nil {
			tracks = nil
		}

		vars.Set("videoId", id)
		vars.Set("albumName", albumName)
		vars.Set("artistName", artistName)
		vars.Set("image", image)
		vars.Set("data", data)
		vars.Set("spotifyId", spotifyId)
		vars.Set("ratingString", ratingString)
		vars.Set("ratingUrl", ratingUrl)
		vars.Set("tracks", tracks)
		vars.Set("artistId", artistId)
		vars.Set("label", label)
	} else {
		c.Redirect(301, "/404")
	}
	c.Writer.WriteHeader(200)
	if err := t.Execute(c.Writer, vars, nil); err != nil {
		log.Println(err)
	}
}

func error404(c *gin.Context) {
	t, _ := views.GetTemplate("404.jet.html")
	//id := c.Param("id")
	vars := make(jet.VarMap)
	c.Writer.WriteHeader(404)
	if err := t.Execute(c.Writer, vars, nil); err != nil {
		log.Println(err)
	}
}
