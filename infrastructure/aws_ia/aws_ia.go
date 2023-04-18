package aws_ia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"worker-validation-identity/infrastructure/env"
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

func CompareFacesV2(face1, face2 []byte) (bool, error) {
	e := env.NewConfiguration()
	face1Reader := bytes.NewReader(face1)
	face2Reader := bytes.NewReader(face2)
	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)

	face1Writer, err := multipartWriter.CreateFormFile("selfie", "selfie"+GetExtensionFromBytes(face1))
	if err != nil {
		return false, err
	}
	if _, err := io.Copy(face1Writer, face1Reader); err != nil {
		return false, err
	}

	face2Writer, err := multipartWriter.CreateFormFile("document", "document"+GetExtensionFromBytes(face2))
	if err != nil {
		return false, err
	}
	if _, err := io.Copy(face2Writer, face2Reader); err != nil {
		return false, err
	}

	multipartWriter.Close()

	request, err := http.NewRequest(http.MethodPost, e.FaceApi.CompareFace, body)
	request.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	httpClient := &http.Client{}

	resp, err := httpClient.Do(request)
	if err != nil {
		logger.Error.Printf("no se  puedo enviar la petición: %v  -- log: ", err)
		return false, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error.Printf("no se pudo ejecutar defer body close: %v  -- log: ", err)
		}
	}(resp.Body)

	rsBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("no se  puedo obtener respuesta: %v  -- log: ", err)
		return false, err
	}

	resFace := CompareFaceResponse{}

	err = json.Unmarshal(rsBody, &resFace)
	if err != nil {
		logger.Error.Printf("no se pudo parsear la respuesta del servio de comparacion de rostros: %v  -- log: ", err)
		return false, err
	}

	if resFace.Error {
		logger.Error.Printf("error al consumir el servicio ocr: " + resFace.Msg)
		return false, fmt.Errorf(resFace.Msg)
	}

	if resFace.Data.Verified == "true" {
		return true, nil
	}

	return false, nil
}

func GetExtensionFromBytes(file []byte) string {
	mimeType := http.DetectContentType(file)
	switch mimeType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpeg"
	case "image/tiff":
		return ".tif"
	}
	return ".txt"
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
