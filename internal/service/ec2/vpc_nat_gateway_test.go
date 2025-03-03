package ec2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccVPCNATGateway_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var natGateway ec2.NatGateway
	resourceName := "aws_nat_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNATGatewayDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNATGatewayConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNATGatewayExists(ctx, resourceName, &natGateway),
					resource.TestCheckResourceAttrSet(resourceName, "allocation_id"),
					resource.TestCheckResourceAttrSet(resourceName, "association_id"),
					resource.TestCheckResourceAttr(resourceName, "connectivity_type", "public"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
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

func TestAccVPCNATGateway_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var natGateway ec2.NatGateway
	resourceName := "aws_nat_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNATGatewayDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNATGatewayConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNATGatewayExists(ctx, resourceName, &natGateway),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfec2.ResourceNATGateway(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPCNATGateway_ConnectivityType_private(t *testing.T) {
	ctx := acctest.Context(t)
	var natGateway ec2.NatGateway
	resourceName := "aws_nat_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNATGatewayDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNATGatewayConfig_connectivityType(rName, "private"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNATGatewayExists(ctx, resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "allocation_id", ""),
					resource.TestCheckResourceAttr(resourceName, "association_id", ""),
					resource.TestCheckResourceAttr(resourceName, "connectivity_type", "private"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttr(resourceName, "public_ip", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", rName),
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

func TestAccVPCNATGateway_privateIP(t *testing.T) {
	ctx := acctest.Context(t)
	var natGateway ec2.NatGateway
	resourceName := "aws_nat_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNATGatewayDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNATGatewayConfig_privateIP(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNATGatewayExists(ctx, resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "allocation_id", ""),
					resource.TestCheckResourceAttr(resourceName, "association_id", ""),
					resource.TestCheckResourceAttr(resourceName, "connectivity_type", "private"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface_id"),
					resource.TestCheckResourceAttr(resourceName, "private_ip", "10.0.0.8"),
					resource.TestCheckResourceAttr(resourceName, "public_ip", ""),
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

func TestAccVPCNATGateway_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var natGateway ec2.NatGateway
	resourceName := "aws_nat_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNATGatewayDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNATGatewayConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNATGatewayExists(ctx, resourceName, &natGateway),
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
				Config: testAccVPCNATGatewayConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNATGatewayExists(ctx, resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVPCNATGatewayConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNATGatewayExists(ctx, resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckNATGatewayDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_nat_gateway" {
				continue
			}

			_, err := tfec2.FindNATGatewayByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("EC2 NAT Gateway %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckNATGatewayExists(ctx context.Context, n string, v *ec2.NatGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 NAT Gateway ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn()

		output, err := tfec2.FindNATGatewayByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccNATGatewayConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "private" {
  vpc_id                  = aws_vpc.test.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.test.id
  cidr_block              = "10.0.2.0/24"
  map_public_ip_on_launch = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_eip" "test" {
  domain = "vpc"

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVPCNATGatewayConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfig_base(rName), `
resource "aws_nat_gateway" "test" {
  allocation_id = aws_eip.test.id
  subnet_id     = aws_subnet.public.id

  depends_on = [aws_internet_gateway.test]
}
`)
}

func testAccVPCNATGatewayConfig_connectivityType(rName, connectivityType string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_nat_gateway" "test" {
  connectivity_type = %[2]q
  subnet_id         = aws_subnet.test[0].id

  tags = {
    Name = %[1]q
  }
}
`, rName, connectivityType))
}

func testAccVPCNATGatewayConfig_privateIP(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_nat_gateway" "test" {
  connectivity_type = "private"
  private_ip        = "10.0.0.8"
  subnet_id         = aws_subnet.test[0].id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccVPCNATGatewayConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfig_base(rName), fmt.Sprintf(`
resource "aws_nat_gateway" "test" {
  allocation_id = aws_eip.test.id
  subnet_id     = aws_subnet.public.id

  tags = {
    %[1]q = %[2]q
  }

  depends_on = [aws_internet_gateway.test]
}
`, tagKey1, tagValue1))
}

func testAccVPCNATGatewayConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfig_base(rName), fmt.Sprintf(`
resource "aws_nat_gateway" "test" {
  allocation_id = aws_eip.test.id
  subnet_id     = aws_subnet.public.id

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }

  depends_on = [aws_internet_gateway.test]
}
`, tagKey1, tagValue1, tagKey2, tagValue2))
}
