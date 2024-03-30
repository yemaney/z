package aws

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// ec2 states
var (
	RUNNING = "running"
	STOPPED = "stopped"
)

// list searches for ec2 instances and prints out their name and current state
func list() {

	result, _ := getInstances()

	// Print information about each instance
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceName := getInstanceName(instance)
			fmt.Printf("%s: %s\n", instanceName, *instance.State.Name)
		}
	}

}

// Create a session using the default AWS configuration and credentials
// Create an EC2 service client
// Call the DescribeInstances API to list instances
func getInstances() (*ec2.DescribeInstancesOutput, *ec2.EC2) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		os.Exit(1)
	}

	svc := ec2.New(sess)

	result, err := svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error describing instances:", err)
		os.Exit(1)
	}
	return result, svc
}

// get searches for an instance and prints information for it.
// if the instance has a publicIP then the information printed will be in format of an SSH config section
// otherwise just prints the name and state
func get(searchString string) {
	result, _ := getInstances()

	instance := filterInstances(result, searchString)

	if instance != nil {
		if *instance.State.Name != RUNNING {
			fmt.Printf("Instance: %s: State%s\n", getInstanceName(instance), *instance.State.Name)
		} else {
			fmt.Printf("Host %s\n", getInstanceName(instance))
			fmt.Printf("	HostName %s\n", *instance.PublicIpAddress)
			fmt.Printf("	IdentityFile %s\n", *instance.KeyName)
			fmt.Printf("	User ubuntu\n")
		}
	}

}

// getInstanceName gets a name for an instance.
// If an instance doesn't have a name, then the instanceId is returned
func getInstanceName(instance *ec2.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return *instance.InstanceId
}

// start searches for an ec2 instance and, if the instance is not running, starts it
func start(searchString string) {
	result, svc := getInstances()

	instance := filterInstances(result, searchString)

	if instance != nil {
		// Check the instance state
		if *instance.State.Name != RUNNING {
			// Start the instance if its state is not running
			fmt.Println("Instance state is not running. Starting instance...")
			_, err := svc.StartInstances(&ec2.StartInstancesInput{
				InstanceIds: []*string{instance.InstanceId},
			})
			if err != nil {
				fmt.Println("Error starting instance:", err)
				return
			}
			fmt.Println("Instance started successfully.")
		} else {
			fmt.Println("Instance is already running.")
		}
	}

}

// stop searches for an ec2 instance and, if the instance is not in a stopped state, stops it
func stop(searchString string) {
	result, svc := getInstances()

	instance := filterInstances(result, searchString)
	if instance != nil {
		if *instance.State.Name != STOPPED {
			// Stop the instance if its state is not "stopped"
			fmt.Println("Instance state is not stopped. Stopping instance...")
			_, err := svc.StopInstances(&ec2.StopInstancesInput{
				InstanceIds: []*string{instance.InstanceId},
			})
			if err != nil {
				fmt.Println("Error stopping instance:", err)
				return
			}
			fmt.Println("Instance stopped successfully.")
		} else {
			fmt.Println("Instance is already stopped.")
		}

	}

}

// filterInstances searches for an instance that has a name or instanceID that matches the searchString
func filterInstances(result *ec2.DescribeInstancesOutput, searchString string) *ec2.Instance {
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if *instance.InstanceId == searchString || getInstanceName(instance) == searchString {
				return instance
			}
		}
	}
	return nil
}

func getLatestImage() *ec2.Image {
	_, svc := getInstances()

	owner := "amazon"
	// Describe Ubuntu AMIs
	input := &ec2.DescribeImagesInput{
		Owners: []*string{&owner},
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: []*string{aws.String("ubuntu/images/hvm-ssd/ubuntu-*amd64*")},
			}},
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		fmt.Println("Error describing images:", err)
		os.Exit(1)
	}
	// Sort the images by creation date
	sort.Slice(result.Images, func(i, j int) bool {
		return parseCreationDate(*result.Images[i].CreationDate).Before(parseCreationDate(*result.Images[j].CreationDate))
	})

	return result.Images[len(result.Images)-1]
}

// Function to parse creation date string into time.Time object
func parseCreationDate(dateStr string) time.Time {
	layout := "2006-01-02T15:04:05.000Z" // ISO 8601 format
	t, _ := time.Parse(layout, dateStr)
	return t
}

func create(instanceName, ami string) {
	_, svc := getInstances()

	// Specify the parameters for the new instance
	key := "new_ssh"
	sg := "sg-03718963ed05bc2f9"
	runInput := &ec2.RunInstancesInput{
		ImageId:      aws.String(ami),        // Replace with the AMI ID of your choice
		InstanceType: aws.String("t2.micro"), // Replace with the instance type of your choice
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		KeyName:      &key,
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(instanceName),
					},
				},
			},
		},
		SecurityGroupIds: []*string{&sg},
	}

	// Run the instance
	runResult, err := svc.RunInstances(runInput)
	if err != nil {
		fmt.Println("Error creating instance:", err)
		os.Exit(1)
	}

	// Print the instance ID of the newly created instance
	fmt.Println("Instance created successfully with ID:", *runResult.Instances[0].InstanceId)
}

func delete(searchString string) error {
	// Create a new AWS session
	result, svc := getInstances()

	instance := filterInstances(result, searchString)

	if instance == nil {
		return nil
	}

	// Terminate the instance
	instanceID := instance.InstanceId
	terminateInput := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{instanceID},
	}

	_, err := svc.TerminateInstances(terminateInput)
	if err != nil {
		return err
	}
	fmt.Printf("Terminating Instance: %s\n", *instanceID)

	return nil
}
