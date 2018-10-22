package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"
	"github.com/go-redis/redis"
	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

var (
	s3Client *minio.Client
	rClient  *redis.Client
	s3Bucket string
)

func main() {
	fmt.Println("Archivist booting up")

	var err error
	s3Client, err = minio.NewV2(os.Getenv("S3_ENDPOINT"), os.Getenv("S3_ACCESS_ID"), os.Getenv("S3_ACCESS_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	s3Bucket = os.Getenv("S3_BUCKET")

	rClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = rClient.Ping().Result()
	if err != nil {
		panic("Redis unavailable: " + err.Error())
	}

	for {
		persistMessages()
		time.Sleep(time.Millisecond * 1200)
	}
}

func persistMessages() {
	toPersist := make(map[string][]string)

	for i := 0; i < 250; i++ {
		message, err := rClient.LPop("messages").Result()
		if err != nil {
			break
		}

		_, user, msg := twitch.ParseMessage(message)

		objectName := fmt.Sprintf("%s_%d_%d", user.UserID, msg.Time.Year(), msg.Time.Month())

		if value, ok := toPersist[objectName]; ok {
			toPersist[objectName] = append(value, msg.Raw)
		} else {
			toPersist[objectName] = []string{msg.Raw}
		}
	}

	for objectName, messages := range toPersist {
		object, err := s3Client.GetObject(s3Bucket, objectName, minio.GetObjectOptions{})

		if err != nil {
			log.Error("Failure getting previous logs", err)
			return
		}
		defer object.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(object)

		result := buf.String()

		for _, message := range messages {
			result += "\n" + message
		}

		fmt.Println("persisting:")
		fmt.Println(result)

		reader := strings.NewReader(result)

		_, err = s3Client.PutObject(s3Bucket, objectName, reader, reader.Size(), minio.PutObjectOptions{})
		if err != nil {
			log.Error("Failure putting logs", err)
		}
	}
}
