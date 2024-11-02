package webcontrol

//gs://

import (
	//internet related

	"fmt"

	// "fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	//general purpose related

	"encoding/json"
	"io"
	"os"
	"regexp"

	// local packages
	fbcon "veles/db/FirebaseConnection"
	prbl "veles/db/Problem"

	// external packages

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

var App *firebase.App
var Err error

var firebaseConnection fbcon.FirebaseConnection

func InitDBconnection(path_to_key string, bucketName string) {
	firebaseConnection = *fbcon.New(path_to_key, bucketName)
}

func WhyPageHandler(c *gin.Context) {
	data := gin.H{
		"PageName": "why",
	}
	c.HTML(http.StatusOK, "why.html", data)
}

func NewPageHandler(c *gin.Context) {
	problems_json, _ := os.ReadFile("../frontend/static/precomputed/problems_list.json")

	problems_list := firebaseConnection.GetProblemsList()
	done_problems_json, _ := json.Marshal(prbl.GetFilteredIdList(problems_list, "done"))
	awaiting_problems_json, _ := json.Marshal(prbl.GetFilteredIdList(problems_list, "inProgress"))

	data := gin.H{
		"PageName":        "new",
		"ProblemsList":    string(problems_json),
		"ReadyProblems":   string(done_problems_json),
		"WaitingProblems": string(awaiting_problems_json),
	}
	c.HTML(http.StatusOK, "new.html", data)
}

func CataloguePageHandler(c *gin.Context) {

	problems_list := firebaseConnection.GetProblemsList()

	done_problems_json, _ := json.Marshal(prbl.GetFilteredIdList(problems_list, "done"))
	awaiting_problems_json, _ := json.Marshal(prbl.GetFilteredIdList(problems_list, "inProgress"))

	data := gin.H{
		"PageName":        "catalogue",
		"DoneProblems":    string(done_problems_json),
		"WaitingProblems": string(awaiting_problems_json),
	}

	c.HTML(http.StatusOK, "catalogue.html", data)
}

func HowPageHandler(c *gin.Context) {
	data := gin.H{
		"PageName": "how",
	}
	c.HTML(http.StatusOK, "how.html", data)
}

func FormSumbissionHandler(c *gin.Context) {

	// Get text values
	number := c.PostForm("number")
	text := c.PostForm("text")
	answer := c.PostForm("answer")
	comment := c.PostForm("comment")

	// Check number for existance
	r, _ := regexp.Compile(`\d{1,2}\.\d{1,2}\.\d{1,3}`)
	if !r.MatchString(number) {
		c.JSON(http.StatusBadRequest, "Bad number value")
		return
	}

	// Parse form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Get files from form
	files := form.File["images"]
	var filenames []string
	var files_to_save []multipart.File
	var files_ext []string

	// Save each file in storage
	for _, file := range files {
		filereader, _ := file.Open()
		files_to_save = append(files_to_save, filereader)
		files_ext = append(files_ext, filepath.Ext(file.Filename))
	}
	// Save text information
	filenames = firebaseConnection.SaveProblemSubmission(number, text, answer, files_ext, comment)

	firebaseConnection.SaveFiles(files_to_save, filenames)

	c.JSON(http.StatusOK, gin.H{"message": "Form handled succesfully"})
}

func ReadyProblemsHandler(c *gin.Context) {
	number, _ := c.GetQuery("problem")
	filereader := firebaseConnection.GetReadyMdFile(number)

	b, err := io.ReadAll(filereader)

	if err != nil {
		fmt.Println(err)
	}

	data := gin.H{
		"MDFile": string(b[:]),
	}

	c.HTML(http.StatusAccepted, "ready.html", data)
}

func PotentialyFullHandler(c *gin.Context) {
	number, _ := c.GetQuery("problemNum")
	num := firebaseConnection.GetNumsOfSubmissions(number)
	if num > 5 {
		c.JSON(http.StatusOK, gin.H{"permission": false})
	} else {
		c.JSON(http.StatusOK, gin.H{"permission": true})
	}
}
