package main

import (
	"github.com/CloudyKit/jet"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	views = jet.NewHTMLSet("./templates")
)

func main(){
	r := gin.Default()

	r.GET("/",index)
	r.StaticFile("/data.client.json", "./DatabaseBuild/output.client.json")

	//start server
	err := r.Run(":"+getEnv("PORT","8080"))
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
	//id := c.Param("id")
	//vars.Set("playlist",playlist)
	c.Writer.WriteHeader(200)
	if err := t.Execute(c.Writer, vars, nil); err != nil {
		log.Println(err)
	}
}