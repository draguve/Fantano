package main

import (
	"encoding/json"
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
)

var (
	views = jet.NewHTMLSet("./templates")
	info map[string]interface{}
)

func main(){
	jsonFile, err := os.Open("./DatabaseBuild/output.server.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := json.Unmarshal(jsonData, &info); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := gin.Default()

	r.GET("/",index)
	r.GET("/404",error404)
	r.GET("/album/:id",album)
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
	if val, ok := info[id]; ok {
		println(val)
	}else{
		c.Redirect(404,"/404")
	}
	vars := make(jet.VarMap)
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