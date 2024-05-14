package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"

	"github.com/kai-thompson/toronto-humane-society-alerts/utils/appointment"
	"github.com/kai-thompson/toronto-humane-society-alerts/utils/sms"
	"github.com/kai-thompson/toronto-humane-society-alerts/utils/logger"
)

const (
	baseCalendarURL = "https://ths.use1.ezyvet.com/external/portal/calendar"

	appointmentsFoundMessage = "There are %d appointment(s) available for %s!"
)

func checkAppointment(id appointment.ID) (int, error) {
	res, err := http.Get(fmt.Sprintf("%s/firstAvailableForType?appointmenttypeid=%d", baseCalendarURL, id))
	if err != nil {
		logger.Error("failed to get appointment data", zap.Error(err))
		return 0, err
	}

	defer res.Body.Close()

	var data []interface{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		logger.Error("failed to decode appointment list", zap.Error(err))
		return 0, err
	}

	return len(data), nil
}

func HandleRequest(ctx context.Context, event interface{}) (error) {
	ids := appointment.IDs()

	doneErr := make(chan error, len(ids))

	for _, id := range ids {
		go func(id appointment.ID) {
			var err error
			defer func() {
				doneErr <- err
			}()

			apmt, err := appointment.New(id)
			if err != nil {
				logger.Error("failed to create appointment", zap.Error(err))
				return
			}

			numAppointments, err := checkAppointment(apmt.ID())
			if err != nil {
				logger.Error("failed to check appointment", zap.Error(err))
				return
			}

			if numAppointments == 0 {
				logger.Info("no appointment available", zap.String("appointment", apmt.Name()))
				return
			}

			err = sms.SendMessage(apmt.ID(), fmt.Sprintf(appointmentsFoundMessage, numAppointments, apmt.Name()))
			if err != nil {
				logger.Error("failed to send text message", zap.Error(err))
			}
		}(id)
	}

	for range ids {
		select {
		case err := <-doneErr:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
			return fmt.Errorf("timed out")
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}