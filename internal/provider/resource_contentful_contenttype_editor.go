package provider

import (
    "context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/the-urge-tech/terraform-provider-contentful/internal/utils"
	"github.com/the-urge-tech/terraform-provider-contentful/pkg/contentful"
)

func resourceContentfulContentTypeEditor() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

        CreateContext: resourceContentTypeEditorCreate,
        ReadContext:   resourceContentTypeEditorRead,
        UpdateContext: resourceContentTypeEditorUpdate,
        DeleteContext: resourceContentTypeEditorDelete,

        Schema: map[string]*schema.Schema{
            "space_id": {
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
            "version": {
                Type:     schema.TypeInt,
                Computed: true,
            },
            "content_type_id": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
			"control": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
						"field_id": {
							Type:     schema.TypeString,
							Required: true,
                        },
                        "widget_id": {
                            Type:     schema.TypeString,
                            Required: true,
                        },
                        "widget_namespace": {
                            Type:     schema.TypeString,
                            Required: true,
                        },
						"settings": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"help_text": {
										Type:     schema.TypeString,
										Optional: true,
                                    },
                                    "bulk_editing": {
                                        Type:     schema.TypeBool,
                                        Optional: true,
                                        Default: true,
                                    },
                                    "show_link_entity_action": {
                                        Type:     schema.TypeBool,
                                        Optional: true,
                                        Default: true,
                                    },
                                    "show_create_entity_action": {
										Type:     schema.TypeBool,
                                        Optional: true,
                                        Default: true,
									},
								},
							},
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

func resourceContentTypeEditorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*contentful.Client)

	spaceID := d.Get("space_id").(string)
	envID := d.Get("env_id").(string)
	id := d.Get("content_type_id").(string)

	body := make(map[string]interface{})

    controls, err := convertControlsForWriting(d.Get("control"))

	if err != nil {
		return diag.Errorf("Unknown error when converting controls: %s", err.Error())
	}

    body["controls"] = controls

	res, err := client.ContentTypeEditor.Put(ctx, spaceID, envID, id, 18, body)
	if err != nil {
		return diag.Errorf("Unknown error when performing upsert: %s", err.Error())
	}

	d.Set("version", getVersion(res))
	d.SetId(fmt.Sprintf("%s/%s/%s/editor_interface", spaceID, envID, id))

	return diags
}

func resourceContentTypeEditorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*contentful.Client)

	ids := strings.Split(d.Id(), "/")
	if len(ids) != 4 {
		return diag.Errorf("Got invalid id: %s", d.Id())
	}
	spaceID := ids[0]
	envID := ids[1]
	id := ids[2]

	ct, err := client.ContentTypeEditor.Read(ctx, spaceID, envID, id)

	if err != nil && strings.Contains(err.Error(), "status code 404") {
		d.SetId("")
		return diags
	}

	if err != nil {
        return diag.Errorf("Unknown error when getting content type editor with id:%s : %s", d.Id(), err.Error())
	}

    err = convertControlsForReading(ct["controls"])

	if err != nil {
		return diag.Errorf("Unknown error when processing fields for content type editor:%s : %s", d.Id(), err.Error())
	}

	d.Set("content_type_id", id)
	d.Set("env_id", envID)
	d.Set("space_id", spaceID)
	d.Set("version", getVersion(ct))
    d.Set("control", ct["controls"])

	return diags
}

func resourceContentTypeEditorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*contentful.Client)

	ids := strings.Split(d.Id(), "/")
	if len(ids) != 4 {
		return diag.Errorf("Got invalid id: %s", d.Id())
	}
	spaceID := ids[0]
	envID := ids[1]
	id := ids[2]

	version := d.Get("version").(int)
	body := make(map[string]interface{})

    controls, err := convertControlsForWriting(d.Get("control"))
	if err != nil {
		return diag.Errorf("Unknown error when converting controls: %s", err.Error())
	}
    body["controls"] = controls

	res, err := client.ContentTypeEditor.Put(ctx, spaceID, envID, id, version, body)
	if err != nil {
		return diag.Errorf("Unknown error when performing upsert: %s", err.Error())
	}
	version = getVersion(res)

//	res, err = client.ContentType.Activate(ctx, spaceID, envID, id, version)
//	if err != nil {
//		return diag.Errorf("Unknown error when activating content type: %s", err.Error())
//	}

	d.Set("version", getVersion(res))
	return diags
}

func resourceContentTypeEditorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("not implemented")
}

func convertControlsForWriting(original interface{}) (interface{}, error) {
    controls := make([]interface{}, 0, len(original.([]interface{})))

    for _, f := range original.([]interface{}) {
        controls = append(controls, utils.CopyMap(f.(map[string]interface{})))
    }

    for _, f := range controls {
        control := f.(map[string]interface{})

        utils.ConvertStringField(control, "field_id", "fieldId")
        utils.ConvertStringField(control, "widget_id", "widgetId")
        utils.ConvertStringField(control, "widget_namespace", "widgetNamespace")


        if control["settings"] != nil {
            if len(control["settings"].([]interface{})) > 0 {
                control["settings"] = control["settings"].([]interface{})[0]
            } else {
                control["settings"] = nil
            }
        }

        if control["settings"] != nil {
            settings := control["settings"].(map[string]interface{})
            utils.ConvertStringField(settings, "help_text", "helpText")
            utils.ConvertBoolField(settings, "bulk_editing", "bulkEditing")
            utils.ConvertBoolField(settings, "show_link_entity_action", "showLinkEntityAction")
            utils.ConvertBoolField(settings, "show_create_entity_action", "showCreateEntityAction")
        }
    }
    return controls, nil
}

func convertControlsForReading(controls interface{}) error {
    for _, f := range controls.([]interface{}) {
        control := f.(map[string]interface{})

        utils.ConvertStringField(control, "fieldId", "field_id")
        utils.ConvertStringField(control, "widgetId", "widget_id")
        utils.ConvertStringField(control, "widgetNamespace", "widget_namespace")


        if control["settings"] != nil {
            settings := control["settings"].(map[string]interface{})
            utils.ConvertStringField(settings, "helpText", "help_text")
            utils.ConvertBoolField(settings, "bulkEditing", "bulk_editing")
            utils.ConvertBoolField(settings, "showLinkEntityAction", "show_link_entity_action")
            utils.ConvertBoolField(settings, "showCreateEntityAction", "show_create_entity_action")
        }
    }
    return nil
}