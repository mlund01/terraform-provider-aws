package ses

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKResource("aws_ses_active_receipt_rule_set")
func ResourceActiveReceiptRuleSet() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceActiveReceiptRuleSetUpdate,
		UpdateWithoutTimeout: resourceActiveReceiptRuleSetUpdate,
		ReadWithoutTimeout:   resourceActiveReceiptRuleSetRead,
		DeleteWithoutTimeout: resourceActiveReceiptRuleSetDelete,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_set_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
		},
	}
}

func resourceActiveReceiptRuleSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SESConn()

	ruleSetName := d.Get("rule_set_name").(string)

	createOpts := &ses.SetActiveReceiptRuleSetInput{
		RuleSetName: aws.String(ruleSetName),
	}

	_, err := conn.SetActiveReceiptRuleSetWithContext(ctx, createOpts)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "setting active SES rule set: %s", err)
	}

	d.SetId(ruleSetName)

	return append(diags, resourceActiveReceiptRuleSetRead(ctx, d, meta)...)
}

func resourceActiveReceiptRuleSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SESConn()

	describeOpts := &ses.DescribeActiveReceiptRuleSetInput{}

	response, err := conn.DescribeActiveReceiptRuleSetWithContext(ctx, describeOpts)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, ses.ErrCodeRuleSetDoesNotExistException) {
			log.Printf("[WARN] SES Receipt Rule Set (%s) belonging to SES Active Receipt Rule Set not found, removing from state", d.Id())
			d.SetId("")
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "reading SES Active Receipt Rule Set: %s", err)
	}

	if response.Metadata == nil {
		log.Print("[WARN] No active Receipt Rule Set found")
		d.SetId("")
		return diags
	}

	d.Set("rule_set_name", response.Metadata.Name)

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "ses",
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("receipt-rule-set/%s", d.Id()),
	}.String()
	d.Set("arn", arn)

	return diags
}

func resourceActiveReceiptRuleSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SESConn()

	deleteOpts := &ses.SetActiveReceiptRuleSetInput{
		RuleSetName: nil,
	}

	_, err := conn.SetActiveReceiptRuleSetWithContext(ctx, deleteOpts)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting active SES rule set: %s", err)
	}

	return diags
}
