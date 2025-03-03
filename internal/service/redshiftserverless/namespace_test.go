package redshiftserverless_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/redshiftserverless"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfredshiftserverless "github.com/hashicorp/terraform-provider-aws/internal/service/redshiftserverless"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccRedshiftServerlessNamespace_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_redshiftserverless_namespace.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshiftserverless.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNamespaceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamespaceExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "redshift-serverless", regexp.MustCompile("namespace/.+$")),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_id"),
					resource.TestCheckResourceAttr(resourceName, "log_exports.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "iam_roles.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNamespaceConfig_updated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamespaceExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "redshift-serverless", regexp.MustCompile("namespace/.+$")),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_id"),
					resource.TestCheckResourceAttr(resourceName, "log_exports.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "iam_roles.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "iam_roles.*", "aws_iam_role.test.0", "arn"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "iam_roles.*", "aws_iam_role.test.1", "arn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccRedshiftServerlessNamespace_defaultIAMRole(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_redshiftserverless_namespace.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshiftserverless.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNamespaceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfig_defaultIAMRole(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamespaceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "default_iam_role_arn", "aws_iam_role.test.0", "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedshiftServerlessNamespace_user(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_redshiftserverless_namespace.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshiftserverless.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNamespaceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamespaceExists(ctx, resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNamespaceConfig_user(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamespaceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "admin_user_password", "Password123"),
				),
			},
		},
	})
}

func TestAccRedshiftServerlessNamespace_tags(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_redshiftserverless_namespace.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshiftserverless.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNamespaceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamespaceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNamespaceConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccNamespaceConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamespaceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccRedshiftServerlessNamespace_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_redshiftserverless_namespace.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshiftserverless.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNamespaceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamespaceExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfredshiftserverless.ResourceNamespace(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRedshiftServerlessNamespace_withWorkgroup(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_redshiftserverless_namespace.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshiftserverless.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNamespaceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfig_withWorkgroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNamespaceExists(ctx, resourceName),
				),
			},
			{
				Config:   testAccNamespaceConfig_withWorkgroup(rName),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckNamespaceDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).RedshiftServerlessConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_redshiftserverless_namespace" {
				continue
			}
			_, err := tfredshiftserverless.FindNamespaceByName(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Redshift Serverless Namespace %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckNamespaceExists(ctx context.Context, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Redshift Serverless Namespace is not set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RedshiftServerlessConn()

		_, err := tfredshiftserverless.FindNamespaceByName(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccNamespaceConfig_baseIAMRole(rName string, n int) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  count = %[2]d

  name = "%[1]s-${count.index}"
  path = "/service-role/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": [
                    "redshift-serverless.amazonaws.com",
                    "redshift.amazonaws.com",
                    "sagemaker.amazonaws.com"
                ]
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
EOF
}

data "aws_partition" "current" {}

resource "aws_iam_role_policy_attachment" "test" {
  count = %[2]d

  role       = aws_iam_role.test[count.index].name
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonRedshiftAllCommandsFullAccess"
}

`, rName, n)
}

func testAccNamespaceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_redshiftserverless_namespace" "test" {
  namespace_name = %[1]q
}
`, rName)
}

func testAccNamespaceConfig_user(rName string) string {
	return fmt.Sprintf(`
resource "aws_redshiftserverless_namespace" "test" {
  namespace_name      = %[1]q
  admin_user_password = "Password123"
}
`, rName)
}

func testAccNamespaceConfig_updated(rName string) string {
	return acctest.ConfigCompose(testAccNamespaceConfig_baseIAMRole(rName, 2), fmt.Sprintf(`
resource "aws_redshiftserverless_namespace" "test" {
  namespace_name = %[1]q
  iam_roles      = aws_iam_role.test[*].arn
}
`, rName))
}

func testAccNamespaceConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_redshiftserverless_namespace" "test" {
  namespace_name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccNamespaceConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_redshiftserverless_namespace" "test" {
  namespace_name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccNamespaceConfig_defaultIAMRole(rName string) string {
	return acctest.ConfigCompose(testAccNamespaceConfig_baseIAMRole(rName, 1), fmt.Sprintf(`
resource "aws_redshiftserverless_namespace" "test" {
  namespace_name       = %[1]q
  default_iam_role_arn = aws_iam_role.test[0].arn
  iam_roles            = aws_iam_role.test[*].arn
}
`, rName))
}

func testAccNamespaceConfig_withWorkgroup(rName string) string {
	return acctest.ConfigCompose(testAccNamespaceConfig_baseIAMRole(rName, 2), fmt.Sprintf(`
resource "aws_redshiftserverless_namespace" "test" {
  namespace_name       = %[1]q
  default_iam_role_arn = aws_iam_role.test[0].arn
  iam_roles            = aws_iam_role.test[*].arn
}

resource "aws_redshiftserverless_workgroup" "test" {
  namespace_name = aws_redshiftserverless_namespace.test.namespace_name
  workgroup_name = %[1]q
}
`, rName))
}
