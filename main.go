package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	r.GET("/probe", probe)

	r.Run(":8080")
}

func probe(c *gin.Context) {

	// 解析 target 参数
	target := c.Query("target")
	// if target != "" {
	// 	// urldecoded
	// 	target, _ = url.QueryUnescape(target)
	// }
	// fmt.Println("target: ", target)s

	debug := c.Query("debug")

	fmt.Printf("target: %s \n", target)

	dsn, err := ParseRedisDSN(target)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	result, _ := CheckLen(*dsn)

	if debug == "true" {
		c.JSON(200, gin.H{
			"data":   dsn,
			"result": result,
		})

		return
	}

	registry := prometheus.NewRegistry()
	probeQueueLen := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis_exporter",
		Name:      "queue_len",
		Help:      "queue lenth",
	},
		[]string{"key", "db", "addr"},
	)

	registry.MustRegister(probeQueueLen)

	for k, v := range result {
		probeQueueLen.WithLabelValues(k, strconv.Itoa(dsn.DB), dsn.Addr).Set(float64(v))
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(c.Writer, c.Request)
}
