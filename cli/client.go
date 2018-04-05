package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prantoran/gcd-grpc/pb"
	"google.golang.org/grpc"
)

/*

Frontend service uses gin web framework to serve a REST API and calls the GCD service
for the actual calculation.
*/

func main() {
	/*
		Create a client to communicate with the GCD service inside main function.
		Because Kubernetes (v1.3+) has a built-in DNS service, you can refer to the GCD service
		with the name "gcd-service", defined later on.
		Will not run on its own, need to define gcd-service
	*/
	// conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	conn, err := grpc.Dial("gcd-service:3000", grpc.WithInsecure())
	fmt.Println("conn:", conn, " err:", err)
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}
	gcdClient := pb.NewGCDServiceClient(conn)
	fmt.Printf("gcdclient type: %T \ngcdclient val: %v\n", gcdClient, gcdClient)
	/*
		After that, declare a handler for /gcd/:a/:b endpoint which reads parameters A and B,
		and then calls the GCD service.
	*/
	r := gin.Default()
	r.GET("/gcd/:a/:b", func(c *gin.Context) {
		// Parse parameters
		a, err := strconv.ParseUint(c.Param("a"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter A"})
			return
		}
		b, err := strconv.ParseUint(c.Param("b"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter B"})
			return
		}
		// Call GCD service
		req := &pb.GCDRequest{A: a, B: b}
		if res, err := gcdClient.Compute(c, req); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprint(res.Result),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	// run the server
	if err := r.Run(":3030"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
