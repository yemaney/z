package aws

import (
	"fmt"
	"os"

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
			fmt.Printf("Host %s:\n", getInstanceName(instance))
			fmt.Printf("    HostName: %s\n", *instance.PublicIpAddress)
			fmt.Printf("    IdentityFile: %s\n", *instance.KeyName)
			fmt.Printf("    User ubuntu\n")
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
