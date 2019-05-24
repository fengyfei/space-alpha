package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	yuque = make(map[string]string)

	// list cache
	listRecode = make(map[string]string)

	// detail cache
	detailRecode = make(map[string](map[string]string))

	// classify cache
	classifyRecode = make(map[string]string)

	// ListClient -
	ListClient *http.Client

	// DetailClient -
	DetailClient *http.Client

	// ClassifyClient -
	ClassifyClient *http.Client
)

func main() {
	router := gin.Default()

	RegisterRouter(router)
	router.Run()
}

// RegisterRouter -
func RegisterRouter(r gin.IRouter) {

	r.GET("/getlist", getList)
	r.GET("/getdetails", getDetails)
	r.GET("/getrepo", getRepo)

}

func getList(c *gin.Context) {
	var list struct {
		RepoID string `json:"repoid"`
	}

	err := c.ShouldBind(&list)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	url := "http://127.0.0.1:9090/list?RepoID=" + list.RepoID
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	Token := c.Request.Header
	request.Header.Add("X-Auth-Token", Token["X-Auth-Token"][0])

	ListClient = &http.Client{}
	response, err := ListClient.Do(request)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	check := bytes.Contains(body, []byte("ERROR"))
	if check {
		log.Println("Requst Failed")
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Details": listRecode[list.RepoID]})
		return
	}

	_, isPresent := listRecode[list.RepoID]
	if !isPresent {
		listRecode[list.RepoID] = string(body)
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "List": listRecode[list.RepoID]})
}

func getDetails(c *gin.Context) {
	var (
		detail struct {
			RepoID string `json:"repoid"`
			ID     string `json:"id"`
		}

		con = make([]string, 5)
	)

	con[0] = "封面"
	con[1] = "作者"
	con[2] = "内容简介"
	con[3] = "捐助者"
	con[4] = "时间"

	err := c.ShouldBind(&detail)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	url := "http://127.0.0.1:9090/list?RepoID=" + detail.RepoID + "&ID=" + detail.ID
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	Token := c.Request.Header
	request.Header.Add("X-Auth-Token", Token["X-Auth-Token"][0])

	DetailClient = &http.Client{}
	response, err := DetailClient.Do(request)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	check := bytes.Contains(body, []byte("ERROR"))
	if check {
		log.Println("Requst Failed")
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Details": detailRecode[detail.ID]})
		return
	}

	info := bytes.Split(body, []byte("****"))

	yuque = map[string]string{
		con[0]: string(info[1]),
		con[1]: string(info[3]),
		con[2]: string(info[5]),
		con[3]: string(info[7]),
		con[4]: string(info[9]),
	}

	_, isPresent := detailRecode[detail.ID]
	if !isPresent {
		detailRecode[detail.ID] = yuque
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Details": detailRecode[detail.ID]})

}

func getRepo(c *gin.Context) {
	var Repo struct {
		RepoID string `json:"repoid"`
	}

	err := c.ShouldBind(&Repo)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	url := "http://127.0.0.1:9090/list?RepoID=" + Repo.RepoID
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	Token := c.Request.Header
	request.Header.Add("X-Auth-Token", Token["X-Auth-Token"][0])

	ClassifyClient = &http.Client{}
	response, err := ClassifyClient.Do(request)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	check := bytes.Contains(body, []byte("ERROR"))
	if check {
		log.Println("Requst Failed")
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Details": classifyRecode[Repo.RepoID]})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "List": string(body)})
}
