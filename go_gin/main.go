package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// アルバム情報
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// アルバムデータのスライス
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func UpdateCache() {
	var lock sync.Mutex
	timer1 := time.NewTicker(time.Second * 10)
	defer timer1.Stop()
	timer2 := time.NewTicker(time.Second * 5)
	defer timer2.Stop()
	for {
		/* run forever */
		select {
		case <-timer1.C:
			go func() {
				lock.Lock()
				defer lock.Unlock()
				/* do things I need done every 10 seconds */
				currentTime := time.Now()
				fmt.Printf("%s  10 seconds\n", currentTime.Format("2006-1-2 15:04:05"))
			}()
		case <-timer2.C:
			go func() {
				lock.Lock()
				defer lock.Unlock()
				/* do things I need done every 5 seconds */
				currentTime := time.Now()
				zone, offset := currentTime.Zone()
				fmt.Printf("%s  5 seconds %s %d\n", currentTime.Format("2006-1-2 15:04:05"), zone, (offset / 60 / 60))
			}()
		}
	}
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
	fmt.Println("getAlbums end")
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id") // string
	// 非同期処理
	go func() {
		// simulate a long task with time.Sleep(). 5 seconds
		time.Sleep(7 * time.Second)

		// idを参照しても特に問題なし
		fmt.Printf("Done! %s\n", id)
	}()

	for _, a := range albums {
		if a.ID == id { // string同士の比較
			c.IndentedJSON(http.StatusOK, a)
			fmt.Println("getAlbums finished1")
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	fmt.Println("getAlbums finished2")
}

// albumでも*albumでも問題なかった
func asyncProcess(newAlbum *album) {
	// simulate a long task with time.Sleep(). 5 seconds
	time.Sleep(7 * time.Second)

	// 変数newAlbumを参照しても特に問題なし
	fmt.Printf("Done! %v\n", newAlbum)
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		fmt.Println(err)
		return
	}
	// 非同期処理
	go asyncProcess(&newAlbum)

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
	fmt.Println("postAlbum finished")
}

func main() {
	// 非同期処理ループ開始
	go UpdateCache()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}
