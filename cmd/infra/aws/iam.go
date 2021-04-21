package aws

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	jose "gopkg.in/square/go-jose.v2"

	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
)

const (
	// requires Provider ARN and Issuer URL
	oidcTrustPolicyTemplate = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"Federated": "%s"
				},
					"Action": "sts:AssumeRoleWithWebIdentity",
				"Condition": {
					"StringEquals": {
						"%s:aud": "openshift"
					}
				}
			}
		]
	}`

	ingressPermPolicy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": [
				"elasticloadbalancing:DescribeLoadBalancers",
				"route53:ListHostedZones",
				"route53:ChangeResourceRecordSets",
				"tag:GetResources"
			],
			"Resource": "*"
		}
	]
}`

	imageRegistryPermPolicy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": [
				"s3:CreateBucket",
				"s3:DeleteBucket",
				"s3:PutBucketTagging",
				"s3:GetBucketTagging",
				"s3:PutBucketPublicAccessBlock",
				"s3:GetBucketPublicAccessBlock",
				"s3:PutEncryptionConfiguration",
				"s3:GetEncryptionConfiguration",
				"s3:PutLifecycleConfiguration",
				"s3:GetLifecycleConfiguration",
				"s3:GetBucketLocation",
				"s3:ListBucket",
				"s3:GetObject",
				"s3:PutObject",
				"s3:DeleteObject",
				"s3:ListBucketMultipartUploads",
				"s3:AbortMultipartUpload",
				"s3:ListMultipartUploadParts"
			],
			"Resource": "*"
		}
	]
}`

	awsEBSCSIPermPolicy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": [
				"ec2:AttachVolume",
				"ec2:CreateSnapshot",
				"ec2:CreateTags",
				"ec2:CreateVolume",
				"ec2:DeleteSnapshot",
				"ec2:DeleteTags",
				"ec2:DeleteVolume",
				"ec2:DescribeInstances",
				"ec2:DescribeSnapshots",
				"ec2:DescribeTags",
				"ec2:DescribeVolumes",
				"ec2:DescribeVolumesModifications",
				"ec2:DetachVolume",
				"ec2:ModifyVolume"
			],
			"Resource": "*"
		}
	]
}`

	discoveryURI      = ".well-known/openid-configuration"
	jwksURI           = "openid/v1/jwks"
	discoveryTemplate = `{
	"issuer": "%s",
	"jwks_uri": "%s/%s",
	"response_types_supported": [
		"id_token"
	],
	"subject_types_supported": [
		"public"
	],
	"id_token_signing_alg_values_supported": [
		"RS256"
	],
	"claims_supported": [
		"aud",
		"exp",
		"sub",
		"iat",
		"iss",
		"sub"
	]
}`

	cloudControllerPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:DescribeInstances",
        "ec2:DescribeImages",
        "ec2:DescribeRegions",
        "ec2:DescribeRouteTables",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeSubnets",
        "ec2:DescribeVolumes",
        "ec2:CreateSecurityGroup",
        "ec2:CreateTags",
        "ec2:CreateVolume",
        "ec2:ModifyInstanceAttribute",
        "ec2:ModifyVolume",
        "ec2:AttachVolume",
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:CreateRoute",
        "ec2:DeleteRoute",
        "ec2:DeleteSecurityGroup",
        "ec2:DeleteVolume",
        "ec2:DetachVolume",
        "ec2:RevokeSecurityGroupIngress",
        "ec2:DescribeVpcs",
        "elasticloadbalancing:AddTags",
        "elasticloadbalancing:AttachLoadBalancerToSubnets",
        "elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
        "elasticloadbalancing:CreateLoadBalancer",
        "elasticloadbalancing:CreateLoadBalancerPolicy",
        "elasticloadbalancing:CreateLoadBalancerListeners",
        "elasticloadbalancing:ConfigureHealthCheck",
        "elasticloadbalancing:DeleteLoadBalancer",
        "elasticloadbalancing:DeleteLoadBalancerListeners",
        "elasticloadbalancing:DescribeLoadBalancers",
        "elasticloadbalancing:DescribeLoadBalancerAttributes",
        "elasticloadbalancing:DetachLoadBalancerFromSubnets",
        "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
        "elasticloadbalancing:ModifyLoadBalancerAttributes",
        "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
        "elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer",
        "elasticloadbalancing:AddTags",
        "elasticloadbalancing:CreateListener",
        "elasticloadbalancing:CreateTargetGroup",
        "elasticloadbalancing:DeleteListener",
        "elasticloadbalancing:DeleteTargetGroup",
        "elasticloadbalancing:DescribeListeners",
        "elasticloadbalancing:DescribeLoadBalancerPolicies",
        "elasticloadbalancing:DescribeTargetGroups",
        "elasticloadbalancing:DescribeTargetHealth",
        "elasticloadbalancing:ModifyListener",
        "elasticloadbalancing:ModifyTargetGroup",
        "elasticloadbalancing:RegisterTargets",
        "elasticloadbalancing:SetLoadBalancerPoliciesOfListener",
        "iam:CreateServiceLinkedRole",
        "kms:DescribeKey"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    }
  ]
}`

	nodePoolPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:AllocateAddress",
        "ec2:AssociateRouteTable",
        "ec2:AttachInternetGateway",
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:CreateInternetGateway",
        "ec2:CreateNatGateway",
        "ec2:CreateRoute",
        "ec2:CreateRouteTable",
        "ec2:CreateSecurityGroup",
        "ec2:CreateSubnet",
        "ec2:CreateTags",
        "ec2:DeleteInternetGateway",
        "ec2:DeleteNatGateway",
        "ec2:DeleteRouteTable",
        "ec2:DeleteSecurityGroup",
        "ec2:DeleteSubnet",
        "ec2:DeleteTags",
        "ec2:DescribeAccountAttributes",
        "ec2:DescribeAddresses",
        "ec2:DescribeAvailabilityZones",
        "ec2:DescribeInstances",
        "ec2:DescribeInternetGateways",
        "ec2:DescribeNatGateways",
        "ec2:DescribeNetworkInterfaces",
        "ec2:DescribeNetworkInterfaceAttribute",
        "ec2:DescribeRouteTables",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeSubnets",
        "ec2:DescribeVpcs",
        "ec2:DescribeVpcAttribute",
        "ec2:DescribeVolumes",
        "ec2:DetachInternetGateway",
        "ec2:DisassociateRouteTable",
        "ec2:DisassociateAddress",
        "ec2:ModifyInstanceAttribute",
        "ec2:ModifyNetworkInterfaceAttribute",
        "ec2:ModifySubnetAttribute",
        "ec2:ReleaseAddress",
        "ec2:RevokeSecurityGroupIngress",
        "ec2:RunInstances",
        "ec2:TerminateInstances",
        "tag:GetResources",
        "ec2:CreateLaunchTemplate",
        "ec2:CreateLaunchTemplateVersion",
        "ec2:DescribeLaunchTemplates",
        "ec2:DescribeLaunchTemplateVersions",
        "ec2:DeleteLaunchTemplate",
        "ec2:DeleteLaunchTemplateVersions"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Condition": {
        "StringLike": {
          "iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"
        }
      },
      "Action": [
        "iam:CreateServiceLinkedRole"
      ],
      "Resource": [
        "arn:*:iam::*:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "iam:PassRole"
      ],
      "Resource": [
        "arn:*:iam::*:role/*-worker-role"
      ],
      "Effect": "Allow"
    }
  ]
}`
)

type KeyResponse struct {
	Keys []jose.JSONWebKey `json:"keys"`
}

func DefaultProfileName(infraID string) string {
	return infraID + "-worker"
}

// inputs: none
// outputs rsa keypair
func (o *CreateIAMOptions) CreateOIDCResources(iamClient iamiface.IAMAPI, s3Client s3iface.S3API) (*CreateIAMOutput, error) {
	output := &CreateIAMOutput{
		Region:  o.Region,
		InfraID: o.InfraID,
	}
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	output.ServiceAccountSigningKey = pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privKey),
	})

	pubKey := &privKey.PublicKey
	pubKeyDERBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	hasher := crypto.SHA256.New()
	hasher.Write(pubKeyDERBytes)
	pubKeyDERHash := hasher.Sum(nil)
	kid := base64.RawURLEncoding.EncodeToString(pubKeyDERHash)

	var keys []jose.JSONWebKey
	keys = append(keys, jose.JSONWebKey{
		Key:       pubKey,
		KeyID:     kid,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	})

	jwks, err := json.MarshalIndent(KeyResponse{Keys: keys}, "", "  ")
	if err != nil {
		return nil, err
	}

	bucketName := o.InfraID
	issuerURL := fmt.Sprintf("s3.%s.amazonaws.com/%s", o.Region, bucketName)
	issuerURLWithProto := fmt.Sprintf("https://%s", issuerURL)
	output.IssuerURL = issuerURLWithProto

	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		var aerr awserr.Error
		if errors.As(err, &aerr) {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				log.Info("Bucket already exists and is owned by us", "bucket", bucketName)
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		log.Info("Bucket created", "bucket", bucketName)
	}

	discoveryJSON := fmt.Sprintf(discoveryTemplate, issuerURLWithProto, issuerURLWithProto, jwksURI)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Body:   aws.ReadSeekCloser(strings.NewReader(discoveryJSON)),
		Bucket: aws.String(bucketName),
		Key:    aws.String(discoveryURI),
	})
	if err != nil {
		return nil, err
	}
	log.Info("OIDC discovery document updated", "bucket", bucketName)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader(jwks),
		Bucket: aws.String(bucketName),
		Key:    aws.String(jwksURI),
	})
	if err != nil {
		return nil, err
	}
	log.Info("JWKS document updated", "bucket", bucketName)

	oidcProviderList, err := iamClient.ListOpenIDConnectProviders(&iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return nil, err
	}

	var providerARN string
	for _, provider := range oidcProviderList.OpenIDConnectProviderList {
		if strings.Contains(*provider.Arn, bucketName) {
			providerARN = *provider.Arn
			log.Info("OIDC provider already exists", "provider", providerARN)
			break
		}
	}

	if len(providerARN) == 0 {
		oidcOutput, err := iamClient.CreateOpenIDConnectProvider(&iam.CreateOpenIDConnectProviderInput{
			ClientIDList: []*string{
				aws.String("openshift"),
			},
			ThumbprintList: []*string{
				aws.String("A9D53002E97E00E043244F3D170D6F4C414104FD"), // root CA thumbprint for s3 (DigiCert)
			},
			Url: aws.String(issuerURLWithProto),
		})
		if err != nil {
			return nil, err
		}

		providerARN = *oidcOutput.OpenIDConnectProviderArn
		log.Info("OIDC provider created", "provider", providerARN)
	}

	oidcTrustPolicy := fmt.Sprintf(oidcTrustPolicyTemplate, providerARN, issuerURL)

	// TODO: The policies and secrets for these roles can be extracted from the
	// release payload, avoiding this current hardcoding.
	arn, err := o.CreateOIDCRole(iamClient, "openshift-ingress", oidcTrustPolicy, ingressPermPolicy)
	if err != nil {
		return nil, err
	}
	output.Roles = append(output.Roles, hyperv1.AWSRoleCredentials{
		ARN:       arn,
		Namespace: "openshift-ingress-operator",
		Name:      "cloud-credentials",
	})

	arn, err = o.CreateOIDCRole(iamClient, "openshift-image-registry", oidcTrustPolicy, imageRegistryPermPolicy)
	if err != nil {
		return nil, err
	}
	output.Roles = append(output.Roles, hyperv1.AWSRoleCredentials{
		ARN:       arn,
		Namespace: "openshift-image-registry",
		Name:      "installer-cloud-credentials",
	})

	arn, err = o.CreateOIDCRole(iamClient, "aws-ebs-csi-driver-operator", oidcTrustPolicy, awsEBSCSIPermPolicy)
	if err != nil {
		return nil, err
	}
	output.Roles = append(output.Roles, hyperv1.AWSRoleCredentials{
		ARN:       arn,
		Namespace: "openshift-cluster-csi-drivers",
		Name:      "ebs-cloud-credentials",
	})

	return output, nil
}

// CreateOIDCRole create an IAM Role with a trust policy for the OIDC provider
func (o *CreateIAMOptions) CreateOIDCRole(client iamiface.IAMAPI, name, trustPolicy, permPolicy string) (string, error) {
	roleName := fmt.Sprintf("%s-%s", o.InfraID, name)
	role, err := existingRole(client, roleName)
	var arn string
	if err != nil {
		return "", err
	}
	if role == nil {
		output, err := client.CreateRole(&iam.CreateRoleInput{
			AssumeRolePolicyDocument: aws.String(trustPolicy),
			RoleName:                 aws.String(roleName),
		})
		if err != nil {
			return "", err
		}
		log.Info("Created role", "name", roleName)
		arn = *output.Role.Arn
	} else {
		log.Info("Found existing role", "name", roleName)
		arn = *role.Arn
	}

	rolePolicyName := roleName
	hasPolicy, err := existingRolePolicy(client, roleName, rolePolicyName)
	if err != nil {
		return "", err
	}
	if !hasPolicy {
		_, err = client.PutRolePolicy(&iam.PutRolePolicyInput{
			PolicyName:     aws.String(rolePolicyName),
			PolicyDocument: aws.String(permPolicy),
			RoleName:       aws.String(roleName),
		})
		if err != nil {
			return "", err
		}
		log.Info("Created role policy", "name", rolePolicyName)
	}

	return arn, nil
}

func (o *CreateIAMOptions) CreateWorkerInstanceProfile(client iamiface.IAMAPI, profileName string) error {
	const (
		assumeRolePolicy = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}`
		workerPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeInstances",
        "ec2:DescribeRegions"
      ],
      "Resource": "*"
    }
  ]
}`
	)
	roleName := fmt.Sprintf("%s-role", profileName)
	role, err := existingRole(client, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		_, err := client.CreateRole(&iam.CreateRoleInput{
			AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
			Path:                     aws.String("/"),
			RoleName:                 aws.String(roleName),
		})
		if err != nil {
			return fmt.Errorf("cannot create worker role: %w", err)
		}
		log.Info("Created role", "name", roleName)
	} else {
		log.Info("Found existing role", "name", roleName)
	}
	instanceProfile, err := existingInstanceProfile(client, profileName)
	if err != nil {
		return err
	}
	if instanceProfile == nil {
		result, err := client.CreateInstanceProfile(&iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			Path:                aws.String("/"),
		})
		if err != nil {
			return fmt.Errorf("cannot create instance profile: %w", err)
		}
		instanceProfile = result.InstanceProfile
		log.Info("Created instance profile", "name", profileName)
	} else {
		log.Info("Found existing instance profile", "name", profileName)
	}
	hasRole := false
	for _, role := range instanceProfile.Roles {
		if aws.StringValue(role.RoleName) == aws.StringValue(role.RoleName) {
			hasRole = true
		}
	}
	if !hasRole {
		_, err = client.AddRoleToInstanceProfile(&iam.AddRoleToInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			RoleName:            aws.String(roleName),
		})
		if err != nil {
			return fmt.Errorf("cannot add role to instance profile: %w", err)
		}
		log.Info("Added role to instance profile", "role", roleName, "profile", profileName)
	}
	rolePolicyName := fmt.Sprintf("%s-policy", profileName)
	hasPolicy, err := existingRolePolicy(client, roleName, rolePolicyName)
	if err != nil {
		return err
	}
	if !hasPolicy {
		_, err = client.PutRolePolicy(&iam.PutRolePolicyInput{
			PolicyName:     aws.String(rolePolicyName),
			PolicyDocument: aws.String(workerPolicy),
			RoleName:       aws.String(roleName),
		})
		if err != nil {
			return fmt.Errorf("cannot create profile policy: %w", err)
		}
		log.Info("Created role policy", "name", rolePolicyName)
	}
	return nil
}

func (o *CreateIAMOptions) CreateCredentialedUserWithPolicy(ctx context.Context, client iamiface.IAMAPI, userName, policyDocument string) (*iam.AccessKey, error) {
	var user *iam.User
	if output, err := client.CreateUserWithContext(ctx, &iam.CreateUserInput{
		UserName: aws.String(userName),
		Tags:     iamTags(o.InfraID, userName),
	}); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	} else {
		user = output.User
	}
	log.Info("Created user", "user", userName)

	var policy *iam.Policy
	if output, err := client.CreatePolicyWithContext(ctx, &iam.CreatePolicyInput{
		PolicyName:     aws.String(userName),
		PolicyDocument: aws.String(policyDocument),
	}); err != nil {
		return nil, fmt.Errorf("failed to create policy: %w", err)
	} else {
		policy = output.Policy
	}
	log.Info("Created policy", "name", aws.StringValue(policy.PolicyName), "arn", aws.StringValue(policy.Arn))

	if _, err := client.AttachUserPolicyWithContext(ctx, &iam.AttachUserPolicyInput{
		UserName:  user.UserName,
		PolicyArn: policy.Arn,
	}); err != nil {
		return nil, fmt.Errorf("failed to attach user policy: %w", err)
	}
	log.Info("Attached user to policy", "user", aws.StringValue(user.UserName), "policy", aws.StringValue(policy.Arn))

	if output, err := client.CreateAccessKeyWithContext(ctx, &iam.CreateAccessKeyInput{
		UserName: user.UserName,
	}); err != nil {
		return nil, fmt.Errorf("failed to create access key: %w", err)
	} else {
		log.Info("Created access key", "user", aws.StringValue(user.UserName))
		return output.AccessKey, nil
	}
}

func existingRole(client iamiface.IAMAPI, roleName string) (*iam.Role, error) {
	result, err := client.GetRole(&iam.GetRoleInput{RoleName: aws.String(roleName)})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("cannot get existing role: %w", err)
	}
	return result.Role, nil
}

func existingInstanceProfile(client iamiface.IAMAPI, profileName string) (*iam.InstanceProfile, error) {
	result, err := client.GetInstanceProfile(&iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("cannot get existing instance profile: %w", err)
	}
	return result.InstanceProfile, nil
}

func existingRolePolicy(client iamiface.IAMAPI, roleName, policyName string) (bool, error) {
	result, err := client.GetRolePolicy(&iam.GetRolePolicyInput{
		RoleName:   aws.String(roleName),
		PolicyName: aws.String(policyName),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				return false, nil
			}
		}
		return false, fmt.Errorf("cannot get existing role policy: %w", err)
	}
	return aws.StringValue(result.PolicyName) == policyName, nil
}

func iamTags(infraID, name string) []*iam.Tag {
	tags := []*iam.Tag{
		{
			Key:   aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", infraID)),
			Value: aws.String("owned"),
		},
	}
	if len(name) > 0 {
		tags = append(tags, &iam.Tag{
			Key:   aws.String("Name"),
			Value: aws.String(name),
		})
	}
	return tags
}
