package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	logFile, err := os.OpenFile(time.Now().Format("./logs/2006-01-02 15.04.05")+".log", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	mw := io.MultiWriter(os.Stderr, logFile)
	log.SetOutput(mw)

	log.Println("Loading configuration...")
	loadConfig()

	log.Println("Starting server...")
	srv := startServer(":7000")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "stop" {
			ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

			if err = srv.Shutdown(ctx); err != nil {
				log.Fatalln("[http] Shutdown failed", err)
			}
			break
		} else if line == "reload" {
			log.Println("Re-loading configuration")
			loadConfig()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
