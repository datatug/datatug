package models

import (
	"strings"
	"testing"
)

func TestDatasetDefinition_Validate(t *testing.T) {
	projectEntity := ProjectEntity{
		ID:    "dataset-id",
		Title: "Dataset title",
	}
	t.Run("should fail", func(t *testing.T) {
		t.Run("no project entity", func(t *testing.T) {
			datasetDef := DatasetDefinition{
				Type:       "json",
				JSONSchema: "{}",
			}
			if err := datasetDef.Validate(); err == nil {
				t.Error("expected to return error")
			}
		})
		t.Run("type field", func(t *testing.T) {
			datasetDef := DatasetDefinition{
				ProjectEntity: projectEntity,
			}
			if err := datasetDef.Validate(); err == nil {
				t.Error("expected to return error for empty struct")
			}
		})
		t.Run("type=json", func(t *testing.T) {
			t.Run("no schema", func(t *testing.T) {
				datasetDef := DatasetDefinition{
					ProjectEntity: projectEntity,
					Type:          "json",
				}

				if err := datasetDef.Validate(); err == nil {
					t.Error("expected to get an error if no schema defined")
					//} else if !validation.IsBadRecordError(err) { TODO: fix validation for IsBadRecordError(err)
					//	t.Errorf("expected BadRecordError, got: %T: %v", err, err)
				} else if !strings.Contains(err.Error(), "jsonSchema") {
					t.Error("name field 'jsonSchema' expected to be in error message, got: ", err.Error())
				}
			})
		})
	})
	t.Run("should pass", func(t *testing.T) {
		t.Run("type=recordset", func(t *testing.T) {
			datasetDef := DatasetDefinition{
				ProjectEntity: projectEntity,
				Type:          "recordset",
			}
			if err := datasetDef.Validate(); err != nil {
				t.Error(err)
			}
		})
		t.Run("type=json", func(t *testing.T) {
			datasetDef := DatasetDefinition{
				ProjectEntity: projectEntity,
				Type:          "json",
				JSONSchema:    "{}",
			}
			if err := datasetDef.Validate(); err != nil {
				t.Error(err)
			}
		})
	})
}
