// deprecated you gotta add the DLP to the service account for this work and update the credentials section in the env with the new path to the creds

package sanitize

import (
	"context"
	"fmt"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

var infoTypesToRedact = []string{
	"SOCIAL_SECURITY_NUMBER",         // Social Security numbers
	"PHONE_NUMBER",                   // Telephone numbers
	"US_LICENSE_NUMBER",              // Certificate or license numbers (US specific - may need others)
	"MEDICAL_RECORD_NUMBER",          // Medical record numbers
	"ACCOUNT_NUMBER",                 // Account numbers
	"HEALTH_PLAN_BENEFICIARY_NUMBER", // Health plan beneficiary numbers
	"PERSON_NAME",                    // Names (Patient names)
	"BIOMETRIC_ID",                   // Biometric elements / identifiers
	"DEVICE_SERIAL_NUMBER",           // Device identifiers and serial numbers
	"EMAIL_ADDRESS",                  // Email addresses
	"FAX_NUMBER",                     // Fax numbers
	"VEHICLE_IDENTIFICATION_NUMBER",  // Vehicle id and serial numbers / Vehicle identifiers
	"IP_ADDRESS",                     // Internet protocol (ip) address numbers / IP addresses
	"URL",                            // Web universal resource locators (urls) / URLs
	// "UNIQUE_IDENTIFIER",   // For "Any other unique identifying number" - this is broad,
	// you might need to define custom InfoTypes for specific ones
	"FULL_FACE_IMAGE", // Full face photographs and comparable images (DLP might have limited support)
	// Consider these as well based on previous discussion:
	"DATE", // Dates (can be identifying in context)
	"DATE_OF_BIRTH",
	"AGE", // If specific and potentially identifying
	"STREET_ADDRESS",
	"ZIP_CODE",              // Depending on context
	"GEOGRAPHIC_COORDINATE", // If present and precise
	"BLOOD_TYPE",
	// And any custom InfoTypes you create for specific medical terms, codes, etc.
}

func SanitizeTranscript(ctx context.Context, projectID, inputString string) (string, error) {
	// Instantiate a client.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("Sanitizer: dlp.NewClient: %w", err)
	}
	defer client.Close()

	// Convert the project id into a full resource id.
	parent := fmt.Sprintf("projects/%s/locations/global", projectID)

	// Define the InfoTypes to redact for HIPAA compliance.
	// This list is based on the discussion and should be reviewed
	// against your specific data and HIPAA requirements.

	// Construct inspect configuration.
	var inspectInfoTypes []*dlppb.InfoType
	for _, it := range infoTypesToRedact {
		inspectInfoTypes = append(inspectInfoTypes, &dlppb.InfoType{Name: it})
	}
	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: inspectInfoTypes,
	}

	// Construct deidentify configuration to redact all identified InfoTypes.
	deidentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
			InfoTypeTransformations: &dlppb.InfoTypeTransformations{
				Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
					{
						PrimitiveTransformation: &dlppb.PrimitiveTransformation{
							Transformation: &dlppb.PrimitiveTransformation_RedactConfig{
								RedactConfig: &dlppb.RedactConfig{},
							},
						},
					},
				},
			},
		},
	}

	// Construct item with the input transcript.
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: inputString,
		},
	}

	// Call the API.
	req := &dlppb.DeidentifyContentRequest{
		Parent:           parent,
		DeidentifyConfig: deidentifyConfig,
		InspectConfig:    inspectConfig,
		Item:             item,
	}

	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return "", fmt.Errorf("Sanitizer: error DeidentifyContent: %w", err)
	}

	return resp.GetItem().GetValue(), nil
}
