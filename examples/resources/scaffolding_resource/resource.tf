resource "contentful_contenttype" "test" {
  space_id        = var.contentful_space_id
  content_type_id = "test"
  name            = "Test Test"
  description     = "Testing the Test"
  display_field   = "uniqueName"
  protected       = true # will prevent accidental field-id changes


  field {
    id          = "uniqueName"
    name        = "Unique Name"
    type        = "Symbol"
    localized   = false
    required    = false
    validations = []
    disabled    = false
    omitted     = false
  }

  field {
    id        = "url"
    name      = "URL"
    type      = "Symbol"
    localized = false
    required  = false
    validations = [
      jsonencode({
        regexp = {
          pattern = "((([A-Za-z]{3,9}:(?:\\/\\/)?)(?:[-;:&=\\+\\$,\\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=\\+\\$,\\w]+@)[A-Za-z0-9.-]+)((?:\\/[\\+~%\\/.\\w-_]*)?\\??(?:[-\\+=&;%@.\\w_]*)#?(?:[\\w]*))?)"
        }
        message = "URL is not valid"
      }),
      jsonencode({
        prohibitRegexp = {
          pattern = "(\\[…\\])"
        }
        message = "Brackets ([...]) not allowed in URL"
      })
    ]
    disabled = false
    omitted  = false
  }
}
