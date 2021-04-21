package aws

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/spf13/cobra"
)

type DestroyIAMOptions struct {
	Region             string
	AWSCredentialsFile string
	InfraID            string
}

func NewDestroyIAMCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws",
		Short: "Destroys AWS instance profile for workers",
	}

	opts := DestroyIAMOptions{
		Region:             "us-east-1",
		AWSCredentialsFile: "",
		InfraID:            "",
	}

	cmd.Flags().StringVar(&opts.AWSCredentialsFile, "aws-creds", opts.AWSCredentialsFile, "Path to an AWS credentials file (required)")
	cmd.Flags().StringVar(&opts.InfraID, "infra-id", opts.InfraID, "Infrastructure ID to use for AWS resources.")
	cmd.Flags().StringVar(&opts.Region, "region", opts.Region, "Region where cluster infra lives")

	cmd.MarkFlagRequired("aws-creds")
	cmd.MarkFlagRequired("infra-id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)
		go func() {
			<-sigs
			cancel()
		}()
		return opts.DestroyIAM(ctx)
	}

	return cmd
}

func (o *DestroyIAMOptions) DestroyIAM(ctx context.Context) error {
	var err error
	iamClient, err := IAMClient(o.AWSCredentialsFile, o.Region)
	if err != nil {
		return err
	}
	s3Client, err := S3Client(o.AWSCredentialsFile, o.Region)
	if err != nil {
		return err
	}
	err = o.DestroyOIDCResources(ctx, iamClient, s3Client)
	if err != nil {
		return err
	}
	err = o.DestroyWorkerInstanceProfile(iamClient)
	if err != nil {
		return err
	}
	return nil
}

func (o *DestroyIAMOptions) DestroyOIDCResources(ctx context.Context, iamClient iamiface.IAMAPI, s3Client s3iface.S3API) error {
	bucketName := o.InfraID

	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(discoveryURI),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != s3.ErrCodeNoSuchBucket &&
				aerr.Code() != s3.ErrCodeNoSuchKey {
				log.Error(aerr, "Error deleting OIDC discovery document", "bucket", bucketName, "key", discoveryURI)
				return aerr
			}
		} else {
			log.Error(err, "Error deleting OIDC discovery document", "bucket", bucketName, "key", discoveryURI)
			return err
		}
	} else {
		log.Info("Deleted OIDC discovery document", "bucket", bucketName, "key", discoveryURI)
	}

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(jwksURI),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != s3.ErrCodeNoSuchBucket &&
				aerr.Code() != s3.ErrCodeNoSuchKey {
				log.Error(aerr, "Error deleting JWKS document", "bucket", bucketName, "key", jwksURI)
				return aerr
			}
		} else {
			log.Error(err, "Error deleting JWKS document", "bucket", bucketName, "key", jwksURI)
			return err
		}
	} else {
		log.Info("Deleted JWKS document", "bucket", bucketName, "key", jwksURI)
	}

	_, err = s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != iam.ErrCodeNoSuchEntityException &&
				aerr.Code() != s3.ErrCodeNoSuchBucket {
				log.Error(aerr, "Error deleting OIDC discovery endpoint", "bucket", bucketName)
				return aerr
			}
		} else {
			log.Error(err, "Error deleting OIDC discovery endpoint", "bucket", bucketName)
			return err
		}
	} else {
		log.Info("Deleted OIDC discovery endpoint", "bucket", bucketName)
	}

	oidcProviderList, err := iamClient.ListOpenIDConnectProviders(&iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return err
	}

	for _, provider := range oidcProviderList.OpenIDConnectProviderList {
		if strings.Contains(*provider.Arn, bucketName) {
			_, err := iamClient.DeleteOpenIDConnectProvider(&iam.DeleteOpenIDConnectProviderInput{
				OpenIDConnectProviderArn: provider.Arn,
			})
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					if aerr.Code() != iam.ErrCodeNoSuchEntityException {
						log.Error(aerr, "Error deleting OIDC provider", "providerARN", provider.Arn)
						return aerr
					}
				} else {
					log.Error(err, "Error deleting OIDC provider", "providerARN", provider.Arn)
					return err
				}
			} else {
				log.Info("Deleted OIDC provider", "providerARN", provider.Arn)
			}
			break
		}
	}
	err = o.DestroyOIDCRole(iamClient, "openshift-ingress")
	err = o.DestroyOIDCRole(iamClient, "openshift-image-registry")
	err = o.DestroyOIDCRole(iamClient, "aws-ebs-csi-driver-operator")

	cloudControllerUserName := fmt.Sprintf("%s-%s", o.InfraID, "cloud-controller")
	nodePoolUserName := fmt.Sprintf("%s-%s", o.InfraID, "node-pool")
	if err := o.DestroyUser(ctx, iamClient, cloudControllerUserName); err != nil {
		return err
	}
	if err := o.DestroyUser(ctx, iamClient, nodePoolUserName); err != nil {
		return err
	}
	if err := o.DestroyPolicy(ctx, iamClient, cloudControllerUserName); err != nil {
		return err
	}
	if err := o.DestroyPolicy(ctx, iamClient, nodePoolUserName); err != nil {
		return err
	}

	return nil
}

// CreateOIDCRole create an IAM Role with a trust policy for the OIDC provider
func (o *DestroyIAMOptions) DestroyOIDCRole(client iamiface.IAMAPI, name string) error {
	roleName := fmt.Sprintf("%s-%s", o.InfraID, name)
	_, err := client.DeleteRolePolicy(&iam.DeleteRolePolicyInput{
		PolicyName: aws.String(roleName),
		RoleName:   aws.String(roleName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != iam.ErrCodeNoSuchEntityException {
				log.Error(aerr, "Error deleting role policy", "role", roleName)
				return aerr
			}
		} else {
			log.Error(err, "Error deleting role policy", "role", roleName)
			return err
		}
	} else {
		log.Info("Deleted role policy", "role", roleName)
	}

	_, err = client.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != iam.ErrCodeNoSuchEntityException {
				log.Error(aerr, "Error deleting role", "role", roleName)
				return aerr
			}
		} else {
			log.Error(err, "Error deleting role", "role", roleName)
			return err
		}
	} else {
		log.Info("Deleted role", "role", roleName)
	}
	return nil
}

func (o *DestroyIAMOptions) DestroyWorkerInstanceProfile(client iamiface.IAMAPI) error {
	profileName := DefaultProfileName(o.InfraID)
	instanceProfile, err := existingInstanceProfile(client, profileName)
	if err != nil {
		return fmt.Errorf("cannot check for existing instance profile: %w", err)
	}
	if instanceProfile != nil {
		for _, role := range instanceProfile.Roles {
			_, err := client.RemoveRoleFromInstanceProfile(&iam.RemoveRoleFromInstanceProfileInput{
				InstanceProfileName: aws.String(profileName),
				RoleName:            role.RoleName,
			})
			if err != nil {
				return fmt.Errorf("cannot remove role %s from instance profile %s: %w", aws.StringValue(role.RoleName), profileName, err)
			}
			log.Info("Removed role from instance profile", "profile", profileName, "role", aws.StringValue(role.RoleName))
		}
		_, err := client.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
		})
		if err != nil {
			return fmt.Errorf("cannot delete instance profile %s: %w", profileName, err)
		}
		log.Info("Deleted instance profile", "profile", profileName)
	}
	roleName := fmt.Sprintf("%s-role", profileName)
	policyName := fmt.Sprintf("%s-policy", profileName)
	role, err := existingRole(client, roleName)
	if err != nil {
		return fmt.Errorf("cannot check for existing role: %w", err)
	}
	if role != nil {
		hasPolicy, err := existingRolePolicy(client, roleName, policyName)
		if err != nil {
			return fmt.Errorf("cannot check for existing role policy: %w", err)
		}
		if hasPolicy {
			_, err := client.DeleteRolePolicy(&iam.DeleteRolePolicyInput{
				PolicyName: aws.String(policyName),
				RoleName:   aws.String(roleName),
			})
			if err != nil {
				return fmt.Errorf("cannot delete role policy %s from role %s: %w", policyName, roleName, err)
			}
			log.Info("Deleted role policy", "role", roleName, "policy", policyName)
		}
		_, err = client.DeleteRole(&iam.DeleteRoleInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			return fmt.Errorf("cannot delete role %s: %w", roleName, err)
		}
		log.Info("Deleted role", "role", roleName)
	}
	return nil
}

func (o *DestroyIAMOptions) DestroyUser(ctx context.Context, client iamiface.IAMAPI, userName string) error {
	// Tear down any access keys for the user
	if output, err := client.ListAccessKeysWithContext(ctx, &iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	}); err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil
			}
		}
		return fmt.Errorf("failed to list access keys: %w", err)
	} else {
		for _, key := range output.AccessKeyMetadata {
			if _, err := client.DeleteAccessKeyWithContext(ctx, &iam.DeleteAccessKeyInput{
				AccessKeyId: key.AccessKeyId,
				UserName:    key.UserName,
			}); err != nil {
				if awsErr, ok := err.(awserr.Error); ok {
					if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
						continue
					}
				}
				return fmt.Errorf("failed to delete access key: %w", err)
			} else {
				log.Info("Deleted access key", "id", key.AccessKeyId, "user", userName)
			}
		}
	}

	// Detach any policies from the user
	if output, err := client.ListAttachedUserPoliciesWithContext(ctx, &iam.ListAttachedUserPoliciesInput{
		UserName: aws.String(userName),
	}); err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() != iam.ErrCodeNoSuchEntityException {
				return fmt.Errorf("failed to list user policies: %w", err)
			}
		} else {
			return fmt.Errorf("failed to list user policies: %w", err)
		}
	} else {
		for _, policy := range output.AttachedPolicies {
			if _, err := client.DetachUserPolicyWithContext(ctx, &iam.DetachUserPolicyInput{
				PolicyArn: policy.PolicyArn,
				UserName:  aws.String(userName),
			}); err != nil {
				return fmt.Errorf("failed to detach policy from user: %w", err)
			} else {
				log.Info("Detached user policy", "user", userName, "policyArn", aws.StringValue(policy.PolicyArn), "policyName", aws.StringValue(policy.PolicyName))
			}
		}
	}

	// Now the user can be deleted
	if _, err := client.DeleteUserWithContext(ctx, &iam.DeleteUserInput{UserName: aws.String(userName)}); err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil
			}
		}
		return fmt.Errorf("failed to delete user: %w", err)
	} else {
		log.Info("Deleted user")
	}
	return nil
}

func (o *DestroyIAMOptions) DestroyPolicy(ctx context.Context, client iamiface.IAMAPI, name string) (result error) {
	return client.ListPoliciesPagesWithContext(ctx, &iam.ListPoliciesInput{}, func(output *iam.ListPoliciesOutput, _ bool) bool {
		for _, policy := range output.Policies {
			if aws.StringValue(policy.PolicyName) != name {
				continue
			}
			if _, err := client.DeletePolicyWithContext(ctx, &iam.DeletePolicyInput{
				PolicyArn: policy.Arn,
			}); err != nil {
				if awsErr, ok := err.(awserr.Error); ok {
					if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
						return true
					}
				}
				result = fmt.Errorf("failed to delete policy: %w", err)
				return true
			} else {
				log.Info("Deleted policy", "name", name, "arn", aws.StringValue(policy.Arn))
				return true
			}
		}
		return true
	})
}
