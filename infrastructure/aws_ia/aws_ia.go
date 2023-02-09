package aws_ia

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"worker-validation-identity/infrastructure/logger"
)

const (
	AwsRegion = "us-east-1"
)

func CompareFaces(face1, face2 []byte) (bool, error) {

	sess, err := sessionAws()
	if err != nil {
		return false, err
	}

	svc := rekognition.New(sess)
	input := &rekognition.CompareFacesInput{
		SimilarityThreshold: aws.Float64(90.000000),
		SourceImage: &rekognition.Image{
			Bytes: face1,
		},
		TargetImage: &rekognition.Image{
			Bytes: face2,
		},
	}

	result, err := svc.CompareFaces(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rekognition.ErrCodeInvalidParameterException:
				logger.Error.Println("InvalidParameterException:", aerr.Error())
			case rekognition.ErrCodeInvalidS3ObjectException:
				logger.Error.Println("ErrCodeInvalidS3ObjectException:", aerr.Error())
			case rekognition.ErrCodeImageTooLargeException:
				logger.Error.Println("ErrCodeImageTooLargeException:", aerr.Error())
			case rekognition.ErrCodeAccessDeniedException:
				logger.Error.Println("ErrCodeAccessDeniedException:", aerr.Error())
			case rekognition.ErrCodeInternalServerError:
				logger.Error.Println("ErrCodeInternalServerError:", aerr.Error())
			case rekognition.ErrCodeThrottlingException:
				logger.Error.Println("ErrCodeThrottlingException:", aerr.Error())
			case rekognition.ErrCodeProvisionedThroughputExceededException:
				logger.Error.Println("ErrCodeProvisionedThroughputExceededException:", aerr.Error())
			case rekognition.ErrCodeInvalidImageFormatException:
				logger.Error.Println("ErrCodeInvalidImageFormatException:", aerr.Error())
			default:
				logger.Error.Println("Other error:", aerr.Error())
			}
		} else {
			logger.Error.Println("Other error:", aerr.Error())
		}
		return false, err
	}
	var similarity float64

	if result.FaceMatches == nil || len(result.FaceMatches) == 0 {
		return false, fmt.Errorf("faces are not similar enough: %f", similarity)
	}

	similarity = *result.FaceMatches[0].Similarity

	if similarity <= 90 {
		return false, fmt.Errorf("faces are not similar enough: %f", similarity)
	}

	return true, nil
}

func sessionAws() (client.ConfigProvider, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AwsRegion),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		fmt.Println("error iniciando sesión con aws: #{err}")
		return sess, err
	}
	return sess, nil
}
