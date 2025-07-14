package pagination

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

// SmartBind - binds query params and auto-applies defaults
func SmartBind(c *gin.Context, params *DefaultQueryParams) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return fmt.Errorf("invalid query parameters: %w", err)
	}
	params.SetDefaults()
	return nil
}

// BindAndSetDefaults - helper function to bind any struct and apply defaults
// Works with existing DTO structs
func BindAndSetDefaults(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return fmt.Errorf("invalid query parameters: %w", err)
	}

	ApplyDefaultsToStruct(req)
	return nil
}

// ApplyDefaultsToStruct uses reflection to apply defaults to any struct with Page/Limit fields
func ApplyDefaultsToStruct(req interface{}) {
	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	// Apply defaults to common pagination fields
	if pageField := val.FieldByName("Page"); pageField.IsValid() && pageField.CanSet() && pageField.Int() < 1 {
		pageField.SetInt(1)
	}

	if limitField := val.FieldByName("Limit"); limitField.IsValid() && limitField.CanSet() {
		if limitField.Int() < 1 {
			limitField.SetInt(10)
		} else if limitField.Int() > 100 {
			limitField.SetInt(100)
		}
	}

	if sortByField := val.FieldByName("SortBy"); sortByField.IsValid() && sortByField.CanSet() && sortByField.String() == "" {
		sortByField.SetString("created_at")
	}

	if sortOrderField := val.FieldByName("SortOrder"); sortOrderField.IsValid() && sortOrderField.CanSet() {
		if sortOrderField.String() == "" {
			sortOrderField.SetString("desc")
		} else if sortOrderField.String() != "asc" && sortOrderField.String() != "desc" {
			sortOrderField.SetString("desc")
		}
	}
}
