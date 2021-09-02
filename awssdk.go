///////////////////////////////////////////////////////////////////
//      (c) 2021 Fujitsu Services                                //
//       By: GyanPatel                                           //
//      Ref: https://pol-jira.atlassian.net/browse/BMP-4421      //
//     Date: 27-Jan-2021                                         //
//  Version: v01.001                                             //
///////////////////////////////////////////////////////////////////

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const forcePasswordChallengeName = "NEW_PASSWORD_REQUIRED"

// This is the username and password of a user from congnito user pool.

//AuthenticateUser to verify if input user is a valid user
func AuthenticateUser(username string, password string) (string, error) {
	log.Println("Info :", "awssdk-AuthenticateUser", "Authentication starting for ", username)

	conf := &aws.Config{Region: aws.String("eu-west-2")}
	sess := session.Must(session.NewSession(conf))
	mac := hmac.New(sha256.New, []byte(secretDetails.CognitoUserPoolClientSecret))
	_, err := mac.Write([]byte(username + secretDetails.CognitoUserPoolClientID))
	if err != nil {
		log.Println("ERROR:AuthenticateUser Error occured awssdk.go - mac.Write ", err)
	}

	secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	cognitoClient := cognitoidentityprovider.New(sess)

	authTry := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(username),
			"PASSWORD":    aws.String(password),
			"SECRET_HASH": aws.String(secretHash),
		},
		ClientId: aws.String(secretDetails.CognitoUserPoolClientID),
	}

	res, err := cognitoClient.InitiateAuth(authTry)
	challengeName := aws.StringValue(res.ChallengeName)
	AccessToken := ""
	if err != nil {
		log.Println("ERROR :", "awssdk-AuthenticateUser", username, err)
		return "N", err
	} else if strings.Compare(challengeName, forcePasswordChallengeName) == 0 {
		log.Println("Info :", "awssdk-AuthenticateUser", username, " authenticated - return access R  ", err)
		return "R", nil
	} else if AccessToken = aws.StringValue(res.AuthenticationResult.AccessToken); len(AccessToken) > 0 {
		log.Println("Info :", "awssdk-AuthenticateUser", username, " authenticated - return access Y ")
		return "Y", nil
	} else {
		log.Println("Info :", "awssdk-AuthenticateUser", username, " unable to perform authentication ")
		return "", nil
	}
}

/*
func getAWSSecret(secretName string) (SecretDetails, error) {
	log.Println("Info : getAWSSecret: Retrieving application secrets ")
	//secretName := secretName
	region := "eu-west-2"
	secretDetails := SecretDetails{}
	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(),
		aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				log.Println("ERROR : getAWSSecret: ", secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				log.Println("ERROR : getAWSSecret: ", secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				log.Println("ERROR : getAWSSecret: ", secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				log.Println("ERROR : getAWSSecret: ", secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				log.Println("ERROR : getAWSSecret: ", secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println("ERROR : getAWSSecret: ", err.Error())
		}
		return secretDetails, err
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		secret := secretString
		json.Unmarshal([]byte(secret), &secretDetails)
	} else { // This is not needed in if block return the secret
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			log.Println("ERROR : getAWSSecret: ", "Base64 Decode Error:", err)
			return secretDetails, err
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		secret := decodedBinarySecret
		json.Unmarshal([]byte(secret), &secretDetails)
	}
	return secretDetails, nil

}
*/
