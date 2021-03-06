package base

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func GetLatestCloudwatchLogs(region string, logGroupName string) ([]string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	svc := cloudwatchlogs.New(sess)

	streamsIn := new(cloudwatchlogs.DescribeLogStreamsInput).SetDescending(
		true).SetLimit(1).SetLogGroupName(logGroupName).SetOrderBy("LastEventTime")
	if err = streamsIn.Validate(); err != nil {
		return nil, err
	}
	streamsOut, err := svc.DescribeLogStreams(streamsIn)
	if err != nil {
		return nil, err
	}

	if len(streamsOut.LogStreams) != 1 {
		return nil, fmt.Errorf("unable to find log groups. Found %d streams", len(streamsOut.LogStreams))
	}

	streamName := streamsOut.LogStreams[0].LogStreamName
	if streamName == nil {
		return nil, fmt.Errorf("Unable to find valid stream %+v", streamsOut.LogStreams[0])
	}

	eventIn := new(cloudwatchlogs.GetLogEventsInput).SetLogGroupName(
		logGroupName).SetLogStreamName(*streamName)
	if err := eventIn.Validate(); err != nil {
		return nil, err
	}
	eventOut, err := svc.GetLogEvents(eventIn)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(eventOut.Events))
	for _, event := range eventOut.Events {
		if event == nil || event.Message == nil || *event.Message == "" {
			continue
		}
		res = append(res, *event.Message)
	}
	return res, nil
}
