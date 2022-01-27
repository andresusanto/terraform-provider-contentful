package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/the-urge-tech/terraform-provider-contentful/internal/utils"
	"github.com/the-urge-tech/terraform-provider-contentful/pkg/contentful"
)

func resourceContentfulContentType() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceContentTypeCreate,
		ReadContext:   resourceContentTypeRead,
		UpdateContext: resourceContentTypeUpdate,
		DeleteContext: resourceContentTypeDelete,

		Schema: map[string]*schema.Schema{
			"protected": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_field": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content_type_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"env_id": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}
					return old == new
				},
			},
			"field": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"link_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"default_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"items": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"link_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"validations": {
										Type:             schema.TypeList,
										Optional:         true,
										Elem:             &schema.Schema{Type: schema.TypeString},
										DiffSuppressFunc: validationDiff,
									},
								},
							},
						},
						"required": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"localized": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"omitted": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"validations": {
							Type:             schema.TypeList,
							Optional:         true,
							Elem:             &schema.Schema{Type: schema.TypeString},
							DiffSuppressFunc: validationDiff,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func validationDiff(k, old, new string, d *schema.ResourceData) bool {
	oldMap := make(map[string]interface{})
	newMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(old), &oldMap)

	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(new), &newMap)

	if err != nil {
		return false
	}

	return reflect.DeepEqual(oldMap, newMap)
}

func resourceContentTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*contentful.Client)

	spaceID := d.Get("space_id").(string)
	envID := d.Get("env_id").(string)
	id := d.Get("content_type_id").(string)

	body := make(map[string]interface{})

	if v, ok := d.GetOk("name"); ok {
		body["name"] = v.(string)
	}

	body["description"] = warningMessage
	if v, ok := d.GetOk("description"); ok {
		body["description"] = warningMessage + v.(string)
	}

	if v, ok := d.GetOk("display_field"); ok {
		body["displayField"] = v.(string)
	}

	fields, err := convertFieldsForWriting(d.Get("field"))

	if err != nil {
		return diag.Errorf("Unknown error when converting field: %s", err.Error())
	}

	body["fields"] = fields

	res, err := client.ContentType.Put(ctx, spaceID, envID, id, 1, body)
	if err != nil {
		return diag.Errorf("Unknown error when performing upsert: %s", err.Error())
	}

	res, err = client.ContentType.Activate(ctx, spaceID, envID, id, getVersion(res))
	if err != nil {
		return diag.Errorf("Unknown error when activating content type: %s", err.Error())
	}

	d.Set("version", getVersion(res))
	d.SetId(fmt.Sprintf("%s/%s/%s", spaceID, envID, id))

	return diags
}

func resourceContentTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*contentful.Client)

	ids := strings.Split(d.Id(), "/")
	if len(ids) != 3 {
		return diag.Errorf("Got invalid id: %s", d.Id())
	}
	spaceID := ids[0]
	envID := ids[1]
	id := ids[2]

	ct, err := client.ContentType.Read(ctx, spaceID, envID, id)

	if err != nil && strings.Contains(err.Error(), "status code 404") {
		d.SetId("")
		return diags
	}

	if err != nil {
		return diag.Errorf("Unknown error when getting content type with id:%s : %s", d.Id(), err.Error())
	}

	err = convertFieldsForReading(ct["fields"])

	if err != nil {
		return diag.Errorf("Unknown error when processing fields for content type:%s : %s", d.Id(), err.Error())
	}

	d.Set("content_type_id", id)
	d.Set("env_id", envID)
	d.Set("space_id", spaceID)
	d.Set("version", getVersion(ct))
	d.Set("name", ct["name"])
	d.Set("description", strings.TrimPrefix(ct["description"].(string), warningMessage))
	d.Set("display_field", ct["displayField"])
	d.Set("field", ct["fields"])

	return diags
}

func resourceContentTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*contentful.Client)

	ids := strings.Split(d.Id(), "/")
	if len(ids) != 3 {
		return diag.Errorf("Got invalid id: %s", d.Id())
	}
	spaceID := ids[0]
	envID := ids[1]
	id := ids[2]

	protected := d.Get("protected").(bool)
	version := d.Get("version").(int)

	oldFields, newFields := d.GetChange("field")
	deletedIDs := getDeletedFieldIDs(oldFields.([]interface{}), newFields.([]interface{}))

	if protected && len(deletedIDs) > 0 {
		d.Set("field", oldFields)
		return diag.Errorf("Protected is set to true and these field(s) will be removed: %v", deletedIDs)
	}

	body := make(map[string]interface{})

	if v, ok := d.GetOk("name"); ok {
		body["name"] = v.(string)
	}

	body["description"] = warningMessage
	if v, ok := d.GetOk("description"); ok {
		body["description"] = warningMessage + v.(string)
	}

	if v, ok := d.GetOk("display_field"); ok {
		body["displayField"] = v.(string)
	}

	if len(deletedIDs) > 0 {
		fields, err := convertFieldsForWriting(oldFields)

		if err != nil {
			return diag.Errorf("Unknown error when converting field: %s", err.Error())
		}

		for i := 0; i < len(fields.([]interface{})); i++ {
			field := fields.([]interface{})[i].(map[string]interface{})
			if deletedIDs[field["id"].(string)] {
				field["omitted"] = true
			}
		}

		body["fields"] = fields

		res, err := client.ContentType.Put(ctx, spaceID, envID, id, version, body)
		if err != nil {
			return diag.Errorf("Unknown error when performing upsert: %s", err.Error())
		}

		version = getVersion(res)

		res, err = client.ContentType.Activate(ctx, spaceID, envID, id, version)
		if err != nil {
			return diag.Errorf("Unknown error when activating content type: %s", err.Error())
		}

		version = getVersion(res)
	}

	fields, err := convertFieldsForWriting(newFields)
	if err != nil {
		return diag.Errorf("Unknown error when converting field: %s", err.Error())
	}
	body["fields"] = fields

	res, err := client.ContentType.Put(ctx, spaceID, envID, id, version, body)
	if err != nil {
		return diag.Errorf("Unknown error when performing upsert: %s", err.Error())
	}
	version = getVersion(res)

	res, err = client.ContentType.Activate(ctx, spaceID, envID, id, version)
	if err != nil {
		return diag.Errorf("Unknown error when activating content type: %s", err.Error())
	}

	d.Set("version", getVersion(res))
	return diags
}

func resourceContentTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("not implemented")
}

func getVersion(ct map[string]interface{}) int {
	if ct["sys"] != nil && ct["sys"].(map[string]interface{})["version"] != nil {
		return int(ct["sys"].(map[string]interface{})["version"].(float64))
	}

	return 1
}

func processValidationForReading(validations interface{}) error {
	v := validations.([]interface{})
	for i := 0; i < len(v); i++ {
		res, err := json.Marshal(v[i])

		if err != nil {
			return err
		}

		v[i] = string(res)
	}
	return nil
}

func processValidationForWriting(validations interface{}) error {
	v := validations.([]interface{})
	for i := 0; i < len(v); i++ {
		r := v[i]
		vMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(r.(string)), &vMap)

		if err != nil {
			return err
		}

		v[i] = vMap
	}
	return nil
}

func convertFieldsForReading(fields interface{}) error {
	for _, f := range fields.([]interface{}) {
		field := f.(map[string]interface{})

		if field["validations"] != nil {
			err := processValidationForReading(field["validations"])

			if err != nil {
				return fmt.Errorf("unknown error when processing validation: %s", err.Error())
			}
		}

		utils.ConvertStringField(field, "linkType", "link_type")
		utils.ConvertMapField(field, "defaultValue", "default_value")

		if field["default_value"] != nil {
			res, err := json.Marshal(field["default_value"])
			if err != nil {
				return err
			}
			field["default_value"] = string(res)
		}

		if field["items"] != nil {
			items := field["items"].(map[string]interface{})
			utils.ConvertStringField(items, "linkType", "link_type")
		}

		if field["items"] != nil && field["items"].(map[string]interface{})["validations"] != nil {
			err := processValidationForReading(field["items"].(map[string]interface{})["validations"])
			if err != nil {
				return fmt.Errorf("unknown error when processing item validation: %s", err.Error())
			}
		}

		if field["items"] != nil {
			field["items"] = []interface{}{field["items"]}
		}
	}
	return nil
}

func convertFieldsForWriting(original interface{}) (interface{}, error) {
	fields := make([]interface{}, 0, len(original.([]interface{})))
	for _, f := range original.([]interface{}) {
		fields = append(fields, utils.CopyMap(f.(map[string]interface{})))
	}

	for _, f := range fields {
		field := f.(map[string]interface{})

		if field["validations"] != nil {
			err := processValidationForWriting(field["validations"])

			if err != nil {
				return nil, fmt.Errorf("unknown error when processing validation: %s", err.Error())
			}
		}

		utils.ConvertStringField(field, "link_type", "linkType")
		utils.ConvertStringField(field, "default_value", "defaultValue")

		if field["defaultValue"] != nil {
			vMap := make(map[string]interface{})
			err := json.Unmarshal([]byte(field["defaultValue"].(string)), &vMap)
			if err != nil {
				return nil, err
			}
			field["defaultValue"] = vMap
		}

		if field["items"] != nil {
			if len(field["items"].([]interface{})) > 0 {
				field["items"] = field["items"].([]interface{})[0]
			} else {
				field["items"] = nil
			}
		}

		if field["items"] != nil {
			items := field["items"].(map[string]interface{})
			utils.ConvertStringField(items, "link_type", "linkType")
		}

		if field["items"] != nil && field["items"].(map[string]interface{})["validations"] != nil {
			err := processValidationForWriting(field["items"].(map[string]interface{})["validations"])
			if err != nil {
				return nil, fmt.Errorf("unknown error when processing item validation: %s", err.Error())
			}
		}
	}
	return fields, nil
}

func getDeletedFieldIDs(old, new []interface{}) map[string]bool {
	result := make(map[string]bool)

	newIDs := make(map[string]bool)
	for _, iField := range new {
		field := iField.(map[string]interface{})
		newIDs[field["id"].(string)] = true
	}

	for _, iField := range old {
		field := iField.(map[string]interface{})
		if !newIDs[field["id"].(string)] {
			result[field["id"].(string)] = true
		}
	}

	return result
}
