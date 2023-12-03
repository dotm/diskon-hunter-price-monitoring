package s3helper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/envhelper"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func CreateClientFromSession() *s3.S3 {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{},
		SharedConfigState: session.SharedConfigEnable,
	}))
	return s3.New(session)
}

const DefaultExpireDurationForPresignedURL = 24 * time.Hour

func DownloadFileToLambdaTemporaryDirectory(s3Client *s3.S3, bucketName, keyName string) (
	file *os.File,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	downloader := s3manager.NewDownloaderWithClient(s3Client)
	file, err = os.Create(fmt.Sprintf("/tmp/%s", keyName)) //Lambda has tmp directory
	if err != nil {
		err = fmt.Errorf("error creating file %s: %v", keyName, err)
		return file, createerror.InternalException(err), err
	}
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	if err != nil {
		err = fmt.Errorf("error downloading file %s/%s: %v", bucketName, keyName, err)
		return file, createerror.InternalException(err), err
	}

	return file, nil, nil
}

func UploadFile(
	s3Client *s3.S3, file io.Reader, bucketName, keyName string,
) (
	errObj *serverresponse.ErrorObj,
	err error,
) {
	uploader := s3manager.NewUploaderWithClient(s3Client)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
		Body:   file,
	})
	if err != nil {
		err = fmt.Errorf("error uploading file %s/%s: %v", bucketName, keyName, err)
		return createerror.InternalException(err), err
	}

	return nil, nil
}

func GeneratePresignedURLForGetObject(s3Client *s3.S3, bucketName, keyName string) (
	presignedUrl string,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	presignedUrl, err = req.Presign(DefaultExpireDurationForPresignedURL)
	if err != nil {
		err = fmt.Errorf("error presign GetObjectRequest for %s/%s: %v", bucketName, keyName, err)
		return presignedUrl, createerror.InternalException(err), err
	}

	return presignedUrl, nil, nil
}

func GeneratePresignedURLForPutObject(s3Client *s3.S3, bucketName, keyName string) (
	presignedUrl string,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	presignedUrl, err = req.Presign(DefaultExpireDurationForPresignedURL)
	if err != nil {
		err = fmt.Errorf("error presign PutObjectRequest for %s/%s: %v", bucketName, keyName, err)
		return presignedUrl, createerror.InternalException(err), err
	}

	return presignedUrl, nil, nil
}

func GetTemporaryBucketName() string {
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s-temporary-bucket",
		envhelper.GetEnvVar("aws_deployment_account_id"),
		envhelper.GetEnvVar("aws_deployment_region_short"),
		envhelper.GetEnvVar("deployment_environment_name"),
		envhelper.GetEnvVar("project_name_short"),
		"db-s3",
	)
}
