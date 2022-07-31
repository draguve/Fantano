package main

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"math/rand"
	"time"
	"os"
	"strconv"
	"strings"
)

var (
	views = jet.NewHTMLSet("./templates")
	jsonData []byte
	videos map[string][]byte
	prefix []string
)

func main(){
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
		//fmt.Printf("Key: '%s'\n Type: %s\n", string(key), dataType)
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

	r.GET("/",index)
	r.GET("/404",error404)
	r.GET("/album/:id",album)
	r.StaticFS("/static", http.Dir("./ScoreImages"))
	r.StaticFile("/data.client.json", "./DatabaseBuild/output.client.json")

	r.NoRoute(error404)

	//start server
	err = r.Run(":"+getEnv("PORT","8080"))
	if err != nil{
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

func index(c *gin.Context){
	t, _ := views.GetTemplate("index.jet.html")
	vars := make(jet.VarMap)
	c.Writer.WriteHeader(200)
	if err := t.Execute(c.Writer, vars, nil); err != nil {
		log.Println(err)
	}
}

func album(c *gin.Context){
	t, _ := views.GetTemplate("album.jet.html")
	id := c.Param("id")
	vars := make(jet.VarMap)
	if data, ok := videos[id]; ok {

		albumName, err := jsonparser.GetString(data, "spotify_name")
		if err!=nil {
			albumName,err = jsonparser.GetString(data,"album")
			if err!=nil{
				albumName, _ = jsonparser.GetString(data,"title")
			}
		}
		artistName, err := jsonparser.GetString(data, "spotify_artists","[0]","name")
		if err!=nil {
			artistName,err = jsonparser.GetString(data,"artist")
			if err!=nil{
				artistName = "Could not find artist name"
			}
		}
		image, _ := jsonparser.GetString(data, "spotify_obj","images","[0]","url")
		spotifyId, _ := jsonparser.GetString(data, "spotify_obj","id")
		ratingString, _ := jsonparser.GetString(data, "rating")
		var ratings = strings.Split(strings.Split(ratingString, ",")[0],"/")[0]
		var ratingNumber = -1
		var ratingUrl = ""
		if val, err := strconv.Atoi(ratings); err == nil {
			ratingNumber = val
		}
		if ratingNumber == 10{
			ratingUrl = "/static/10.png"
		}else if ratingNumber<10 && ratingNumber>-1{
			ratingUrl = fmt.Sprintf("/static/%s%d.png",prefix[rand.Intn(len(prefix))],ratingNumber)
		}
		vars.Set("videoId",id)
		vars.Set("albumName",albumName)
		vars.Set("artistName",artistName)
		vars.Set("image",image)
		vars.Set("data",data)
		vars.Set("spotifyId",spotifyId)
		vars.Set("ratingString",ratingString)
		vars.Set("ratingUrl",ratingUrl)
	}else{
		c.Redirect(301,"/404")
	}
	c.Writer.WriteHeader(200)
	if err := t.Execute(c.Writer, vars, nil); err != nil {
		log.Println(err)
	}
}

func error404(c *gin.Context){
	t, _ := views.GetTemplate("404.jet.html")
	//id := c.Param("id")
	vars := make(jet.VarMap)
	c.Writer.WriteHeader(404)
	if err := t.Execute(c.Writer, vars, nil); err != nil {
		log.Println(err)
	}
}