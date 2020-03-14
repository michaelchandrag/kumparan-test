package main

import (
	"fmt"
	"os"
	"time"
	"sync"
	"encoding/json"
    
    "github.com/adjust/rmq"
    gin "github.com/gin-gonic/gin"
	database "bitbucket.org/michaelchandrag/kumparan-test/database"
	model "bitbucket.org/michaelchandrag/kumparan-test/model"
	helper "bitbucket.org/michaelchandrag/kumparan-test/helper"
)


type Consumer struct {
	name   string
	count  int
	before time.Time
}

func (consumer *Consumer) Consume(delivery rmq.Delivery) {
    var taskBody model.News
    if err := json.Unmarshal([]byte(delivery.Payload()), &taskBody); err != nil {
        // handle error
        delivery.Reject()
        return
    }

    // perform task
    var newNews model.News
	resultNews, err := newNews.Create(taskBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resultNews)
	if err = helper.EsPost(resultNews); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fmt.Sprintf("%d created on Elastic Search", resultNews.ID))
    delivery.Ack()
}

var queue rmq.Queue

func main() {
	connection := rmq.OpenConnection("kumparan", "tcp", fmt.Sprintf("localhost:%s", os.Getenv("Q_PORT")), 1)
	queue = connection.OpenQueue("kumparan")
	queue.StartConsuming(1000, 500*time.Millisecond)

	consumer := Consumer{
		name: "First Consumer",
		count: 0,
		before: time.Now(),
	}
	queue.AddConsumer("First Consumer", &consumer)

	err := database.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	fmt.Println("Database connected")


	router := SetupRouter()
	router.Run()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello",
		})
		return
	})

	// GET localhost:{PORT}/news
	r.GET("/news", func(c *gin.Context) {
		q := c.Request.URL.Query()

		news := helper.EsGet(q)
		var wg sync.WaitGroup
		wg.Add(len(news.HitsResult.HitsHits))
		for index, value := range news.HitsResult.HitsHits {
			newsID := value.Source.ID
			go func(indexLoop int, id int) {
				var whereNews model.News
				whereNews.FindByID(id)
				news.HitsResult.HitsHits[indexLoop].News = whereNews
				fmt.Println(whereNews)
				defer wg.Done()
			}(index, newsID)
		}

		wg.Wait()

		c.JSON(200, gin.H{
			"success": true,
			"news": news.HitsResult.HitsHits,
		})
		return
	})

	// POST localhost:{PORT}/news
	r.POST("/news", func(c *gin.Context) {
		var reqBodyNews model.News

		if err := c.BindJSON(&reqBodyNews); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{
				"success": false,
				"message": "Request error. Check your body request.",
			})
			return
		}

		taskBody := model.News{
			Author: reqBodyNews.Author,
			Body: reqBodyNews.Body,
		}
		testJson, err := json.Marshal(taskBody)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		queue.PublishBytes(testJson)

		c.JSON(200, gin.H{
			"success": true,
			"news": reqBodyNews,
		})
		return
	})

	return r
}