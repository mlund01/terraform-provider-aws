// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package cognitoidp

import (
	"context"

	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newResourceUserPoolClient,
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceUserPoolClient,
			TypeName: "aws_cognito_user_pool_client",
		},
		{
			Factory:  DataSourceUserPoolClients,
			TypeName: "aws_cognito_user_pool_clients",
		},
		{
			Factory:  DataSourceUserPoolSigningCertificate,
			TypeName: "aws_cognito_user_pool_signing_certificate",
		},
		{
			Factory:  DataSourceUserPools,
			TypeName: "aws_cognito_user_pools",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceIdentityProvider,
			TypeName: "aws_cognito_identity_provider",
		},
		{
			Factory:  ResourceResourceServer,
			TypeName: "aws_cognito_resource_server",
		},
		{
			Factory:  ResourceRiskConfiguration,
			TypeName: "aws_cognito_risk_configuration",
		},
		{
			Factory:  ResourceUser,
			TypeName: "aws_cognito_user",
		},
		{
			Factory:  ResourceUserGroup,
			TypeName: "aws_cognito_user_group",
		},
		{
			Factory:  ResourceUserInGroup,
			TypeName: "aws_cognito_user_in_group",
		},
		{
			Factory:  ResourceUserPool,
			TypeName: "aws_cognito_user_pool",
		},
		{
			Factory:  ResourceUserPoolDomain,
			TypeName: "aws_cognito_user_pool_domain",
		},
		{
			Factory:  ResourceUserPoolUICustomization,
			TypeName: "aws_cognito_user_pool_ui_customization",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.CognitoIDP
}

var ServicePackage = &servicePackage{}
