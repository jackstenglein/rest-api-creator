package ec2

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

type mockService struct {
	addIngressRuleInput *ec2.AuthorizeSecurityGroupIngressInput
	addIngressRuleErr   error

	createSecurityGroupInput *ec2.CreateSecurityGroupInput
	createSecurityGroupErr   error

	createTagsInput *ec2.CreateTagsInput
	createTagsErr   error

	describeInstanceInput  *ec2.DescribeInstancesInput
	describeInstanceOutput *ec2.DescribeInstancesOutput
	describeInstanceErr    error

	describeGroupInput  *ec2.DescribeSecurityGroupsInput
	describeGroupOutput *ec2.DescribeSecurityGroupsOutput
	describeGroupErr    error

	runInstanceInput  *ec2.RunInstancesInput
	runInstanceOutput *ec2.Reservation
	runInstanceErr    error
}

func (mock *mockService) AuthorizeSecurityGroupIngress(input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	if !reflect.DeepEqual(input, mock.addIngressRuleInput) {
		return nil, errors.NewServer("Incorrect input to AuthorizeSecurityGroupIngress mock")
	}
	return nil, mock.addIngressRuleErr
}

func (mock *mockService) CreateSecurityGroup(input *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
	if !reflect.DeepEqual(input, mock.createSecurityGroupInput) {
		return nil, errors.NewServer("Incorrect input to CreateSecurityGroup mock")
	}
	return nil, mock.createSecurityGroupErr
}

func (mock *mockService) CreateTags(input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	if !reflect.DeepEqual(input, mock.createTagsInput) {
		return nil, errors.NewServer("Incorrect input to CreateTags mock")
	}
	return nil, mock.createTagsErr
}

func (mock *mockService) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	if !reflect.DeepEqual(input, mock.describeInstanceInput) {
		return nil, errors.NewServer("Incorrect input to DescribeInstances mock")
	}
	return mock.describeInstanceOutput, mock.describeInstanceErr
}

func (mock *mockService) DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	if !reflect.DeepEqual(input, mock.describeGroupInput) {
		return nil, errors.NewServer("Incorrect input to DescribeSecurityGroup mock")
	}
	return mock.describeGroupOutput, mock.describeGroupErr
}

func (mock *mockService) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	if !reflect.DeepEqual(input, mock.runInstanceInput) {
		return nil, errors.NewServer("Incorrect input to RunInstances mock")
	}
	return mock.runInstanceOutput, mock.runInstanceErr
}

var describeInstanceTests = []struct {
	name         string
	instanceID   string
	mock         *mockService
	wantInstance *ec2.Instance
	wantErr      error
}{
	{
		name:       "ServiceError",
		instanceID: "instance",
		mock: &mockService{
			describeInstanceInput: &ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String("instance")},
			},
			describeInstanceErr: errors.NewServer("EC2 failure"),
		},
		wantErr: errors.Wrap(errors.NewServer("EC2 failure"), "Failed call to DescribeInstances"),
	},
	{
		name:       "NoReservations",
		instanceID: "instance",
		mock: &mockService{
			describeInstanceInput: &ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String("instance")},
			},
			describeInstanceOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{},
			},
		},
		wantErr: errors.NewServer("DescribeInstances returned no results"),
	},
	{
		name:       "NoInstances",
		instanceID: "instance",
		mock: &mockService{
			describeInstanceInput: &ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String("instance")},
			},
			describeInstanceOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					&ec2.Reservation{
						Instances: []*ec2.Instance{},
					},
				},
			},
		},
		wantErr: errors.NewServer("DescribeInstances returned no results"),
	},
	{
		name:       "InstanceExists",
		instanceID: "instance",
		mock: &mockService{
			describeInstanceInput: &ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String("instance")},
			},
			describeInstanceOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					&ec2.Reservation{
						Instances: []*ec2.Instance{
							&ec2.Instance{
								InstanceId:    aws.String("instance"),
								PublicDnsName: aws.String("instance.example.com"),
							},
						},
					},
				},
			},
		},
		wantInstance: &ec2.Instance{
			InstanceId:    aws.String("instance"),
			PublicDnsName: aws.String("instance.example.com"),
		},
	},
}

func TestDescribeInstance(t *testing.T) {
	for _, test := range describeInstanceTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			svc = test.mock
			defer func() {
				svc = defaultSvc
			}()

			// Execute
			instance, err := EC2.DescribeInstance(test.instanceID)

			// Verify
			if !reflect.DeepEqual(instance, test.wantInstance) {
				t.Errorf("Got err `%v`; want `%v`", err, test.wantErr)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err `%v`; want `%v`", err, test.wantErr)
			}
		})
	}
}

var getPublicURLTests = []struct {
	name       string
	instanceID string
	mock       *mockService
	wantURL    string
	wantErr    error
}{
	{
		name:       "ServiceError",
		instanceID: "instance",
		mock: &mockService{
			describeInstanceInput: &ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String("instance")},
			},
			describeInstanceErr: errors.NewServer("EC2 failure"),
		},
		wantErr: errors.Wrap(errors.Wrap(errors.NewServer("EC2 failure"), "Failed call to DescribeInstances"), "Failed to describe instance"),
	},
	{
		name:       "NoReservations",
		instanceID: "instance",
		mock: &mockService{
			describeInstanceInput: &ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String("instance")},
			},
			describeInstanceOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{},
			},
		},
		wantErr: errors.Wrap(errors.NewServer("DescribeInstances returned no results"), "Failed to describe instance"),
	},
	{
		name:       "NoInstances",
		instanceID: "instance",
		mock: &mockService{
			describeInstanceInput: &ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String("instance")},
			},
			describeInstanceOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					&ec2.Reservation{
						Instances: []*ec2.Instance{},
					},
				},
			},
		},
		wantErr: errors.Wrap(errors.NewServer("DescribeInstances returned no results"), "Failed to describe instance"),
	},
	{
		name:       "InstanceExists",
		instanceID: "instance",
		mock: &mockService{
			describeInstanceInput: &ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String("instance")},
			},
			describeInstanceOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					&ec2.Reservation{
						Instances: []*ec2.Instance{
							&ec2.Instance{
								InstanceId:    aws.String("instance"),
								PublicDnsName: aws.String("instance.example.com"),
							},
						},
					},
				},
			},
		},
		wantURL: "instance.example.com",
	},
}

func TestGetPublicURL(t *testing.T) {
	for _, test := range getPublicURLTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			svc = test.mock
			defer func() {
				svc = defaultSvc
			}()

			// Execute
			url, err := EC2.GetPublicURL(test.instanceID)

			// Verify
			if url != test.wantURL {
				t.Errorf("Got url `%v`; want `%v`", url, test.wantURL)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err `%v`; want `%v`", err, test.wantErr)
			}
		})
	}
}

var launchInstanceTests = []struct {
	name    string
	mock    *mockService
	wantID  string
	wantURL string
	wantErr error
}{
	{
		name: "DescribeSecurityGroupError",
		mock: &mockService{
			describeGroupInput: &ec2.DescribeSecurityGroupsInput{
				GroupNames: []*string{aws.String(securityGroupName)},
			},
			describeGroupErr: errors.NewServer("EC2 failure"),
		},
		wantErr: errors.Wrap(errors.Wrap(errors.NewServer("EC2 failure"), "Failed call to DescribeSecurityGroups"), "Failed to get security group"),
	},
	{
		name: "CreateSecurityGroupError",
		mock: &mockService{
			createSecurityGroupInput: &ec2.CreateSecurityGroupInput{
				Description: aws.String(securityGroupDescription),
				GroupName:   aws.String(securityGroupName),
			},
			createSecurityGroupErr: errors.NewServer("EC2 failure"),
			describeGroupInput: &ec2.DescribeSecurityGroupsInput{
				GroupNames: []*string{aws.String(securityGroupName)},
			},
			describeGroupOutput: &ec2.DescribeSecurityGroupsOutput{},
		},
		wantErr: errors.Wrap(errors.Wrap(errors.NewServer("EC2 failure"), "Failed call to CreateSecurityGroup"), "Failed to create security group"),
	},
	{
		name: "AuthorizeSecurityGroupIngressError",
		mock: &mockService{
			addIngressRuleInput: &ec2.AuthorizeSecurityGroupIngressInput{
				CidrIp:     aws.String(ipRangeAnywhere),
				FromPort:   aws.Int64(ingressPort),
				GroupName:  aws.String(securityGroupName),
				IpProtocol: aws.String(ipProtocol),
				ToPort:     aws.Int64(ingressPort),
			},
			addIngressRuleErr: errors.NewServer("EC2 failure"),
			createSecurityGroupInput: &ec2.CreateSecurityGroupInput{
				Description: aws.String(securityGroupDescription),
				GroupName:   aws.String(securityGroupName),
			},
			describeGroupInput: &ec2.DescribeSecurityGroupsInput{
				GroupNames: []*string{aws.String(securityGroupName)},
			},
			describeGroupOutput: &ec2.DescribeSecurityGroupsOutput{},
		},
		wantErr: errors.Wrap(errors.Wrap(errors.NewServer("EC2 failure"), "Failed call to AuthorizeSecurityGroupIngress"), "Failed to add ingress rule to security group"),
	},
	{
		name: "RunInstanceError",
		mock: &mockService{
			addIngressRuleInput: &ec2.AuthorizeSecurityGroupIngressInput{
				CidrIp:     aws.String(ipRangeAnywhere),
				FromPort:   aws.Int64(ingressPort),
				GroupName:  aws.String(securityGroupName),
				IpProtocol: aws.String(ipProtocol),
				ToPort:     aws.Int64(ingressPort),
			},
			createSecurityGroupInput: &ec2.CreateSecurityGroupInput{
				Description: aws.String(securityGroupDescription),
				GroupName:   aws.String(securityGroupName),
			},
			describeGroupInput: &ec2.DescribeSecurityGroupsInput{
				GroupNames: []*string{aws.String(securityGroupName)},
			},
			describeGroupOutput: &ec2.DescribeSecurityGroupsOutput{},
			runInstanceInput: &ec2.RunInstancesInput{
				ImageId:      aws.String(imageID),
				InstanceType: aws.String(instanceType),
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				UserData:     aws.String(encUserData),
			},
			runInstanceErr: errors.NewServer("EC2 failure"),
		},
		wantErr: errors.Wrap(errors.Wrap(errors.NewServer("EC2 failure"), "Failed call to RunInstances"), "Failed to run instance"),
	},
	{
		name: "NoInstances",
		mock: &mockService{
			addIngressRuleInput: &ec2.AuthorizeSecurityGroupIngressInput{
				CidrIp:     aws.String(ipRangeAnywhere),
				FromPort:   aws.Int64(ingressPort),
				GroupName:  aws.String(securityGroupName),
				IpProtocol: aws.String(ipProtocol),
				ToPort:     aws.Int64(ingressPort),
			},
			createSecurityGroupInput: &ec2.CreateSecurityGroupInput{
				Description: aws.String(securityGroupDescription),
				GroupName:   aws.String(securityGroupName),
			},
			describeGroupInput: &ec2.DescribeSecurityGroupsInput{
				GroupNames: []*string{aws.String(securityGroupName)},
			},
			describeGroupOutput: &ec2.DescribeSecurityGroupsOutput{},
			runInstanceInput: &ec2.RunInstancesInput{
				ImageId:      aws.String(imageID),
				InstanceType: aws.String(instanceType),
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				UserData:     aws.String(encUserData),
			},
			runInstanceOutput: &ec2.Reservation{},
		},
		wantErr: errors.Wrap(errors.NewServer("RunInstances returned no instances"), "Failed to run instance"),
	},
	{
		name: "NilInstanceId",
		mock: &mockService{
			addIngressRuleInput: &ec2.AuthorizeSecurityGroupIngressInput{
				CidrIp:     aws.String(ipRangeAnywhere),
				FromPort:   aws.Int64(ingressPort),
				GroupName:  aws.String(securityGroupName),
				IpProtocol: aws.String(ipProtocol),
				ToPort:     aws.Int64(ingressPort),
			},
			createSecurityGroupInput: &ec2.CreateSecurityGroupInput{
				Description: aws.String(securityGroupDescription),
				GroupName:   aws.String(securityGroupName),
			},
			describeGroupInput: &ec2.DescribeSecurityGroupsInput{
				GroupNames: []*string{aws.String(securityGroupName)},
			},
			describeGroupOutput: &ec2.DescribeSecurityGroupsOutput{},
			runInstanceInput: &ec2.RunInstancesInput{
				ImageId:      aws.String(imageID),
				InstanceType: aws.String(instanceType),
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				UserData:     aws.String(encUserData),
			},
			runInstanceOutput: &ec2.Reservation{
				Instances: []*ec2.Instance{
					&ec2.Instance{},
				},
			},
		},
		wantErr: errors.Wrap(errors.NewServer("RunInstances did not return InstanceId"), "Failed to run instance"),
	},
	{
		name: "EmptyInstanceId",
		mock: &mockService{
			addIngressRuleInput: &ec2.AuthorizeSecurityGroupIngressInput{
				CidrIp:     aws.String(ipRangeAnywhere),
				FromPort:   aws.Int64(ingressPort),
				GroupName:  aws.String(securityGroupName),
				IpProtocol: aws.String(ipProtocol),
				ToPort:     aws.Int64(ingressPort),
			},
			createSecurityGroupInput: &ec2.CreateSecurityGroupInput{
				Description: aws.String(securityGroupDescription),
				GroupName:   aws.String(securityGroupName),
			},
			describeGroupInput: &ec2.DescribeSecurityGroupsInput{
				GroupNames: []*string{aws.String(securityGroupName)},
			},
			describeGroupOutput: &ec2.DescribeSecurityGroupsOutput{},
			runInstanceInput: &ec2.RunInstancesInput{
				ImageId:      aws.String(imageID),
				InstanceType: aws.String(instanceType),
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				UserData:     aws.String(encUserData),
			},
			runInstanceOutput: &ec2.Reservation{
				Instances: []*ec2.Instance{
					&ec2.Instance{InstanceId: aws.String("")},
				},
			},
		},
		wantErr: errors.Wrap(errors.NewServer("RunInstances did not return InstanceId"), "Failed to run instance"),
	},
	{
		name: "CreateTagsError",
		mock: &mockService{
			addIngressRuleInput: &ec2.AuthorizeSecurityGroupIngressInput{
				CidrIp:     aws.String(ipRangeAnywhere),
				FromPort:   aws.Int64(ingressPort),
				GroupName:  aws.String(securityGroupName),
				IpProtocol: aws.String(ipProtocol),
				ToPort:     aws.Int64(ingressPort),
			},
			createSecurityGroupInput: &ec2.CreateSecurityGroupInput{
				Description: aws.String(securityGroupDescription),
				GroupName:   aws.String(securityGroupName),
			},
			createTagsInput: &ec2.CreateTagsInput{
				Resources: []*string{aws.String("instance")},
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("CRUD Creator Server"),
					},
				},
			},
			createTagsErr: errors.NewServer("EC2 failure"),
			describeGroupInput: &ec2.DescribeSecurityGroupsInput{
				GroupNames: []*string{aws.String(securityGroupName)},
			},
			describeGroupOutput: &ec2.DescribeSecurityGroupsOutput{},
			runInstanceInput: &ec2.RunInstancesInput{
				ImageId:      aws.String(imageID),
				InstanceType: aws.String(instanceType),
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				UserData:     aws.String(encUserData),
			},
			runInstanceOutput: &ec2.Reservation{
				Instances: []*ec2.Instance{
					{
						InstanceId: aws.String("instance"),
					},
				},
			},
		},
		wantErr: errors.Wrap(errors.Wrap(errors.NewServer("EC2 failure"), "Failed call to CreateTags"), "Failed to run instance"),
	},
	{
		name: "SuccessfulInvocation",
		mock: &mockService{
			addIngressRuleInput: &ec2.AuthorizeSecurityGroupIngressInput{
				CidrIp:     aws.String(ipRangeAnywhere),
				FromPort:   aws.Int64(ingressPort),
				GroupName:  aws.String(securityGroupName),
				IpProtocol: aws.String(ipProtocol),
				ToPort:     aws.Int64(ingressPort),
			},
			createSecurityGroupInput: &ec2.CreateSecurityGroupInput{
				Description: aws.String(securityGroupDescription),
				GroupName:   aws.String(securityGroupName),
			},
			createTagsInput: &ec2.CreateTagsInput{
				Resources: []*string{aws.String("instance")},
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("CRUD Creator Server"),
					},
				},
			},
			describeGroupInput: &ec2.DescribeSecurityGroupsInput{
				GroupNames: []*string{aws.String(securityGroupName)},
			},
			describeGroupOutput: &ec2.DescribeSecurityGroupsOutput{},
			runInstanceInput: &ec2.RunInstancesInput{
				ImageId:      aws.String(imageID),
				InstanceType: aws.String(instanceType),
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				UserData:     aws.String(encUserData),
			},
			runInstanceOutput: &ec2.Reservation{
				Instances: []*ec2.Instance{
					{
						InstanceId:    aws.String("instance"),
						PublicDnsName: aws.String("instance.example.com"),
					},
				},
			},
		},
		wantID:  "instance",
		wantURL: "instance.example.com",
	},
	{
		name: "ExistingSecurityGroup",
		mock: &mockService{
			createTagsInput: &ec2.CreateTagsInput{
				Resources: []*string{aws.String("instance")},
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("CRUD Creator Server"),
					},
				},
			},
			describeGroupInput: &ec2.DescribeSecurityGroupsInput{
				GroupNames: []*string{aws.String(securityGroupName)},
			},
			describeGroupOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						Description: aws.String(securityGroupDescription),
						GroupName:   aws.String(securityGroupName),
						IpPermissions: []*ec2.IpPermission{
							{
								IpProtocol: aws.String(ipProtocol),
								FromPort:   aws.Int64(ingressPort),
								ToPort:     aws.Int64(ingressPort),
								IpRanges: []*ec2.IpRange{
									{
										CidrIp: aws.String(ipRangeAnywhere),
									},
								},
							},
						},
					},
				},
			},
			runInstanceInput: &ec2.RunInstancesInput{
				ImageId:      aws.String(imageID),
				InstanceType: aws.String(instanceType),
				MaxCount:     aws.Int64(1),
				MinCount:     aws.Int64(1),
				UserData:     aws.String(encUserData),
			},
			runInstanceOutput: &ec2.Reservation{
				Instances: []*ec2.Instance{
					{
						InstanceId:    aws.String("instance"),
						PublicDnsName: aws.String("instance.example.com"),
					},
				},
			},
		},
		wantID:  "instance",
		wantURL: "instance.example.com",
	},
}

func TestLaunchInstance(t *testing.T) {
	for _, test := range launchInstanceTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			svc = test.mock
			defer func() {
				svc = defaultSvc
			}()

			// Execute
			id, url, err := EC2.LaunchInstance()

			// Verify
			if id != test.wantID {
				t.Errorf("Got id `%v`; want `%v`", id, test.wantID)
			}
			if url != test.wantURL {
				t.Errorf("Got url `%v`; want `%v`", url, test.wantURL)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err `%v`; want `%v`", err, test.wantErr)
			}
		})
	}
}

var shouldAddTests = []struct {
	name  string
	group *ec2.SecurityGroup
	want  bool
}{
	{
		name: "NilGroup",
		want: true,
	},
	{
		name:  "NoIpPermissions",
		group: &ec2.SecurityGroup{},
		want:  true,
	},
	{
		name: "WrongIpProtocol",
		group: &ec2.SecurityGroup{
			IpPermissions: []*ec2.IpPermission{
				&ec2.IpPermission{IpProtocol: aws.String("udp")},
			},
		},
		want: true,
	},
	{
		name: "WrongFromPort",
		group: &ec2.SecurityGroup{
			IpPermissions: []*ec2.IpPermission{
				&ec2.IpPermission{
					IpProtocol: aws.String("tcp"),
					FromPort:   aws.Int64(1337),
					ToPort:     aws.Int64(ingressPort),
				},
			},
		},
		want: true,
	},
	{
		name: "WrongToPort",
		group: &ec2.SecurityGroup{
			IpPermissions: []*ec2.IpPermission{
				&ec2.IpPermission{
					IpProtocol: aws.String("tcp"),
					FromPort:   aws.Int64(ingressPort),
					ToPort:     aws.Int64(1337),
				},
			},
		},
		want: true,
	},
	{
		name: "NoIpRanges",
		group: &ec2.SecurityGroup{
			IpPermissions: []*ec2.IpPermission{
				&ec2.IpPermission{
					IpProtocol: aws.String("tcp"),
					FromPort:   aws.Int64(ingressPort),
					ToPort:     aws.Int64(ingressPort),
				},
			},
		},
		want: true,
	},
	{
		name: "WrongIpRange",
		group: &ec2.SecurityGroup{
			IpPermissions: []*ec2.IpPermission{
				&ec2.IpPermission{
					IpProtocol: aws.String("tcp"),
					FromPort:   aws.Int64(ingressPort),
					ToPort:     aws.Int64(ingressPort),
					IpRanges: []*ec2.IpRange{
						&ec2.IpRange{
							CidrIp: aws.String("127.0.0.0/0"),
						},
					},
				},
			},
		},
		want: true,
	},
	{
		name: "CorrectSecurityGroup",
		group: &ec2.SecurityGroup{
			IpPermissions: []*ec2.IpPermission{
				&ec2.IpPermission{
					IpProtocol: aws.String("tcp"),
					FromPort:   aws.Int64(ingressPort),
					ToPort:     aws.Int64(ingressPort),
					IpRanges: []*ec2.IpRange{
						&ec2.IpRange{
							CidrIp: aws.String(ipRangeAnywhere),
						},
					},
				},
			},
		},
		want: false,
	},
}

func TestShouldAddIngressRule(t *testing.T) {
	for _, test := range shouldAddTests {
		t.Run(test.name, func(t *testing.T) {
			// Execute
			got := shouldAddIngressRule(test.group)

			// Verify
			if got != test.want {
				t.Errorf("Got %v; want %v", got, test.want)
			}
		})
	}
}
