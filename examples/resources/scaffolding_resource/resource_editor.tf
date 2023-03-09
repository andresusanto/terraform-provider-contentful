resource "contentful_contenttype_editor" "test" {
    space_id = "test
    content_type_id = "test"

    control {
        field_id = "testField"
        widget_id = "singleLine"
        widget_namespace = "builtin"
        settings {
            bulk_editing = true,
            help_text = "Test help text",
            show_create_entity_action = true,
            show_link_entity_action = true
        }
    }

    control {
        field_id = "AnotherTestField"
        widget_id = "entryLinksEditor"
        widget_namespace = "builtin"
        settings {
            bulk_editing = true,
            help_text = "Test another help text",
            show_create_entity_action = true,
            show_link_entity_action = true
        }
    }
}