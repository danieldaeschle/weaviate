//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2022 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package sorter

import (
	"github.com/semi-technologies/weaviate/entities/schema"
)

type classHelper struct {
	schema schema.Schema
}

func newClassHelper(schema schema.Schema) *classHelper {
	return &classHelper{schema}
}

func (s *classHelper) getDataType(className, property string) []string {
	class := s.schema.GetClass(schema.ClassName(className))
	if property == "id" || property == "_id" {
		// handle special ID property
		return []string{string(schema.DataTypeString)}
	}
	if property == "_creationTimeUnix" || property == "_lastUpdateTimeUnix" {
		return []string{string(schema.DataTypeInt)}
	}
	for _, prop := range class.Properties {
		if prop.Name == property {
			return prop.DataType
		}
	}
	return nil
}

func (s *classHelper) getOrder(order string) string {
	switch order {
	case "asc", "desc":
		return order
	default:
		return "asc"
	}
}
