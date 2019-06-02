package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	yuque = make(map[string]string)

	// listRecode
	lr sync.Map

	// detailRecode
	dr sync.Map

	// repoRecode
	rr sync.Map

	// ListClient -
	ListClient *http.Client

	// DetailClient -
	DetailClient *http.Client

	// ClassifyClient -
	ClassifyClient *http.Client
)

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
	fmt.Println(string(body))

	check := bytes.Contains(body, []byte("ERROR"))
	if check {
		fmt.Println("2222")
		log.Println("Requst Failed")
		val, ok := lr.Load(list.RepoID)
		if !ok {
			c.JSON(http.StatusRequestTimeout, gin.H{"status": http.StatusRequestTimeout})
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Details": val})
		return
	}
	fmt.Println("111")
	val, ok := lr.Load(list.RepoID)
	fmt.Println(ok, "----")
	if !ok {
		lr.Store(list.RepoID, string(body))
	}
	val, ok = lr.Load(list.RepoID)
	fmt.Println(val, ok, "====")
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "List": val})
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

	url := "http://127.0.0.1:9090/detail?RepoID=" + detail.RepoID + "&ID=" + detail.ID
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
		val, ok := lr.Load(detail.RepoID)
		if !ok {
			c.JSON(http.StatusRequestTimeout, gin.H{"status": http.StatusRequestTimeout})
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Details": val})

	}

	info := bytes.Split(body, []byte("****"))

	yuque = map[string]string{
		con[0]: string(info[1]),
		con[1]: string(info[3]),
		con[2]: string(info[5]),
		con[3]: string(info[7]),
		con[4]: string(info[9]),
	}

	val, ok := dr.Load(detail.RepoID)
	if !ok {
		dr.Store(detail.ID, yuque)
	}

	val, ok = dr.Load(detail.RepoID)
	fmt.Println(val, ok)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Detail": val})

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

	url := "http://127.0.0.1:9090/repo?RepoID=" + Repo.RepoID
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

		val, ok := rr.Load(Repo.RepoID)
		if !ok {
			c.JSON(http.StatusRequestTimeout, gin.H{"status": http.StatusRequestTimeout})
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Details": val})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "List": string(body)})
}
