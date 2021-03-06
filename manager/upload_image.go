package manager

import (
	"bytes"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/upload-image-service/data"
)

var (
	awsScretKey = os.Getenv("S3_SECRET_ACCESS_KEY")
	token       = ""
	pathImage   = "https://s3.amazonaws.com/"
)

// UploadImageToS3 is upload to AWS S3
func UploadImageToS3(model data.UploadImage) (error, string) {
	creds := credentials.NewStaticCredentials(model.AccessKey, awsScretKey, token)
	_, err := creds.Get()
	if err != nil {
		return err, ""
	}
	cfg := aws.NewConfig().WithRegion(model.Region).WithCredentials(creds)
	s, err := session.NewSession(cfg)
	if err != nil {
		return err, ""
	}

	err = addFileToS3(s, model)
	if err != nil {
		return err, ""
	}
	return nil, (pathImage + model.Bucket + "/" + model.ImageName)
}

func addFileToS3(s *session.Session, model data.UploadImage) error {
	file, err := os.Create(model.ImageName)
	if err != nil {
		return err
	}

	_, err = file.Write(model.ImageByte)

	defer file.Close()
	if err != nil {
		return err
	}

	fileInfo, _ := file.Stat()
	fileName := fileInfo.Name()
	fileBytes := bytes.NewReader(model.ImageByte)
	fileType := http.DetectContentType(model.ImageByte)

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(model.Bucket),
		Key:           aws.String(fileName),
		Body:          fileBytes,
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String(fileType),
	})
	os.Remove("./" + fileName)
	if err != nil {
		return err
	}
	return err
}
