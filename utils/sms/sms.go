package sms

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"go.uber.org/zap"

	"github.com/kai-thompson/toronto-humane-society-alerts/utils/appointment"
	"github.com/kai-thompson/toronto-humane-society-alerts/utils/logger"
)

const subscribeConfirmationMessage = "You are now subscribed to receive notifications when new appointments are available for %s. Send 'STOP' to unsubscribe."

var client *sns.SNS

func init() {
	sess, err := session.NewSession(aws.NewConfig().WithRegion("ca-central-1"))
	if err != nil {
		logger.Fatal("failed to create new session", zap.Error(err))
	}

	client = sns.New(sess)
}

func SendMessage(appointmentID appointment.ID, message string) error {
	appointment, err := appointment.New(appointmentID)
	if err != nil {
		return err
	}

	params := &sns.PublishInput{
		Message: aws.String(message),
		TargetArn: aws.String(appointment.TopicARN()),
	}

	_, err = client.Publish(params)
	if err != nil {
		return err
	}

	return nil
}

func Subscribe(apmt appointment.ID, phoneNumber string) error {
	appointment, err := appointment.New(apmt)
	if err != nil {
		return err
	}

	params := &sns.SubscribeInput{
		Protocol: aws.String("sms"),
		TopicArn: aws.String(appointment.TopicARN()),
		Endpoint: aws.String(phoneNumber),
	}

	_, err = client.Subscribe(params)
	if err != nil {
		return err
	}

	return sendSubscribeConfirmation(appointment)
}

func sendSubscribeConfirmation(appointment *appointment.Appointment) error {
	confirmationMsg := fmt.Sprintf(subscribeConfirmationMessage, appointment.Name())
	return SendMessage(appointment.ID(), confirmationMsg)
}
