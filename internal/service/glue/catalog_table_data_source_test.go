package glue_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/glue"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccGlueCatalogTableDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_glue_catalog_table.test"
	datasourceName := "data.aws_glue_catalog_table.test"

	dbName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	tName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCatalogTableDataSourceConfig_basic(dbName, tName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(datasourceName, "catalog_id", resourceName, "catalog_id"),
					resource.TestCheckResourceAttrPair(datasourceName, "database_name", resourceName, "database_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "owner", resourceName, "owner"),
					resource.TestCheckResourceAttrPair(datasourceName, "parameters", resourceName, "parameters"),
					resource.TestCheckResourceAttrPair(datasourceName, "partition_keys", resourceName, "partition_keys"),
					resource.TestCheckResourceAttrPair(datasourceName, "retention", resourceName, "retention"),
					resource.TestCheckResourceAttrPair(datasourceName, "storage_descriptor", resourceName, "storage_descriptor"),
					resource.TestCheckResourceAttrPair(datasourceName, "table_type", resourceName, "table_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "target_table", resourceName, "target_table"),
					resource.TestCheckResourceAttrPair(datasourceName, "view_original_text", resourceName, "view_original_text"),
					resource.TestCheckResourceAttrPair(datasourceName, "view_expanded_text", resourceName, "view_expanded_text"),
					resource.TestCheckResourceAttrPair(datasourceName, "partition_index", resourceName, "partition_index"),
				),
			},
		},
	})
}

func testAccCatalogTableDataSourceConfig_basic(dbName, tName string) string {
	return fmt.Sprintf(`
resource "aws_glue_catalog_database" "test" {
  name = %[1]q
}

resource "aws_glue_catalog_table" "test" {
  database_name = aws_glue_catalog_database.test.name
  name          = %[2]q

  description = "aws_glue_catalog_table datasource acc test"

  table_type = "EXTERNAL_TABLE"

  parameters = {
    EXTERNAL              = "TRUE"
    "parquet.compression" = "SNAPPY"
  }

  storage_descriptor {
    location      = "s3://my-bucket/event-streams/my-stream"
    input_format  = "org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat"

    ser_de_info {
      name                  = "my-stream"
      serialization_library = "org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe"

      parameters = {
        "serialization.format" = 1
      }
    }

    columns {
      name = "my_string"
      type = "string"
    }

    columns {
      name = "my_double"
      type = "double"
    }
  }

  partition_keys {
    name = "my_partition_key"
    type = "string"

    comment = "my_partition_key"
  }
}

data "aws_glue_catalog_table" "test" {
  database_name = aws_glue_catalog_table.test.database_name
  name          = aws_glue_catalog_table.test.name
}
`, dbName, tName)
}
