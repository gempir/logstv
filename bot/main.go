package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"
	_ "github.com/go-sql-driver/mysql"
	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

var (
	s3Client *minio.Client
	s3Bucket string
)

func main() {
	tClient := twitch.NewClient("justinfan123123", "oauth:123123123")

	var err error
	s3Client, err = minio.NewV2(os.Getenv("S3_ENDPOINT"), os.Getenv("S3_ACCESS_ID"), os.Getenv("S3_ACCESS_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	s3Bucket = os.Getenv("S3_BUCKET")

	tClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		go persistMessage(message.Tags["room-id"], message.Tags["user-id"], message.Time, message)
	})

	tClient.OnNewClearchatMessage(func(channel string, user twitch.User, message twitch.Message) {
		go persistMessage(message.Tags["room-id"], message.Tags["target-user-id"], message.Time, message)
	})

	tClient.OnNewUsernoticeMessage(func(channel string, user twitch.User, message twitch.Message) {
		go persistMessage(message.Tags["room-id"], message.Tags["user-id"], message.Time, message)
	})

	tClient.Join("gempir")
	// tClient.Join("pajlada")
	// tClient.Join("nymn")
	// tClient.Join("forsen")

	panic(tClient.Connect())
}

func persistMessage(channelid string, userid string, timestamp time.Time, message twitch.Message) {
	objectName := fmt.Sprintf("%s_%d_%d", userid, timestamp.Year(), timestamp.Month())

	object, err := s3Client.GetObject(s3Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Error("Failure getting previous logs", err)
		return
	}
	defer object.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(object)

	reader := strings.NewReader(buf.String() + "\n" + message.Raw)

	_, err = s3Client.PutObject(s3Bucket, objectName, reader, reader.Size(), minio.PutObjectOptions{})
	if err != nil {
		log.Error("Failure putting logs", err)
	}
}
