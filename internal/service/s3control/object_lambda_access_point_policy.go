package s3control

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceObjectLambdaAccessPointPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceObjectLambdaAccessPointPolicyCreate,
		Read:   resourceObjectLambdaAccessPointPolicyRead,
		Update: resourceObjectLambdaAccessPointPolicyUpdate,
		Delete: resourceObjectLambdaAccessPointPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"has_public_access_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: verify.SuppressEquivalentPolicyDiffs,
			},
		},
	}
}

func resourceObjectLambdaAccessPointPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	accountID := meta.(*conns.AWSClient).AccountID
	if v, ok := d.GetOk("account_id"); ok {
		accountID = v.(string)
	}
	name := d.Get("name").(string)
	resourceID := ObjectLambdaAccessPointCreateResourceID(accountID, name)

	input := &s3control.PutAccessPointPolicyForObjectLambdaInput{
		AccountId: aws.String(accountID),
		Name:      aws.String(name),
		Policy:    aws.String(d.Get("policy").(string)),
	}

	log.Printf("[DEBUG] Creating S3 Object Lambda Access Point Policy: %s", input)
	_, err := conn.PutAccessPointPolicyForObjectLambda(input)

	if err != nil {
		return fmt.Errorf("error creating S3 Object Lambda Access Point (%s) Policy: %w", resourceID, err)
	}

	d.SetId(resourceID)

	return resourceObjectLambdaAccessPointPolicyRead(d, meta)
}

func resourceObjectLambdaAccessPointPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	accountID, name, err := ObjectLambdaAccessPointParseResourceID(d.Id())

	if err != nil {
		return err
	}

	policy, status, err := FindObjectLambdaAccessPointPolicyAndStatusByAccountIDAndName(conn, accountID, name)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] S3 Object Lambda Access Point Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading S3 Object Lambda Access Point Policy (%s): %w", d.Id(), err)
	}

	d.Set("account_id", accountID)
	d.Set("has_public_access_policy", status.IsPublic)
	d.Set("name", name)
	d.Set("policy", policy)

	return nil
}

func resourceObjectLambdaAccessPointPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	accountID, name, err := ObjectLambdaAccessPointParseResourceID(d.Id())

	if err != nil {
		return err
	}

	input := &s3control.PutAccessPointPolicyForObjectLambdaInput{
		AccountId: aws.String(accountID),
		Name:      aws.String(name),
		Policy:    aws.String(d.Get("policy").(string)),
	}

	log.Printf("[DEBUG] Updating S3 Object Lambda Access Point Policy: %s", input)
	_, err = conn.PutAccessPointPolicyForObjectLambda(input)

	if err != nil {
		return fmt.Errorf("error updating S3 Object Lambda Access Point Policy (%s): %w", d.Id(), err)
	}

	return resourceObjectLambdaAccessPointPolicyRead(d, meta)
}

func resourceObjectLambdaAccessPointPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	accountID, name, err := ObjectLambdaAccessPointParseResourceID(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting S3 Object Lambda Access Point Policy: %s", d.Id())
	_, err = conn.DeleteAccessPointPolicyForObjectLambda(&s3control.DeleteAccessPointPolicyForObjectLambdaInput{
		AccountId: aws.String(accountID),
		Name:      aws.String(name),
	})

	if tfawserr.ErrCodeEquals(err, errCodeNoSuchAccessPoint) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting S3 Object Lambda Access Point Policy (%s): %w", d.Id(), err)
	}

	return nil
}
