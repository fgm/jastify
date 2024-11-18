package schemaindexer

import (
	"log"
	"os"
	"reflect"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"

	"github.com/fgm/jastify/cmd"
)

var (
	// discriminators is a map of the references of the oneOfs in the schema,
	// to the Schema.Properties map entry used to discriminate them, e.g. "type"
	// for a widget, or "data_source" for a query_value request.
	discriminators = map[string]string{
		// TODO complete the 24 missing ones as needed.
		// #/components/schemas/DistributionPoint
		// #/components/schemas/DistributionWidgetHistogramRequestQuery
		"#/components/schemas/FormulaAndFunctionQueryDefinition": "data_source",
		// #/components/schemas/LogsProcessor
		// #/components/schemas/MonitorFormulaAndFunctionQueryDefinition
		// #/components/schemas/NotebookCellCreateRequestAttributes
		// #/components/schemas/NotebookCellResponseAttributes
		// #/components/schemas/NotebookCellTime
		// #/components/schemas/NotebookCellUpdateRequestAttributes
		// #/components/schemas/NotebookGlobalTime
		// #/components/schemas/NotebookUpdateCell
		// #/components/schemas/SLODataSourceQueryDefinition
		// #/components/schemas/SLOSliSpec
		// #/components/schemas/SharedDashboardInvitesData
		// #/components/schemas/SplitGraphSourceWidgetDefinition
		// #/components/schemas/SunburstWidgetLegend
		// #/components/schemas/SyntheticsAPIStep
		// #/components/schemas/SyntheticsAssertion
		// #/components/schemas/SyntheticsBasicAuth
		// #/components/schemas/SyntheticsMobileStepParamsValue
		// #/components/schemas/SyntheticsTestRequestPort
		// #/components/schemas/TableWidgetTextFormatReplace
		// #/components/schemas/ToplistWidgetDisplay
		"#/components/schemas/WidgetDefinition": "type",
		// #/components/schemas/WidgetSortOrderBy
		// #/components/schemas/WidgetTime
	}
)

func selectArrayItem(cur *base.Schema) *base.Schema {
	if cur.Items == nil {
		log.Fatalf("expected items, got nil")
	}
	if !cur.Items.IsA() {
		log.Fatalf("expected items to be *base.SchemaProxy, got %T", cur.Items)
	}
	curSP := cur.Items.A
	cur, err := curSP.BuildSchema()
	if err != nil {
		log.Fatalf("building schema: %v", err)
	}
	return cur
}

func selectObjectProperty(cur *base.Schema, selector any) *base.Schema {
	sel, ok := selector.(string)
	if !ok {
		log.Fatalf("object selector must be string, got %v (%T)", selector, selector)
	}
	if cur.Properties == nil {
		log.Fatalf("expected properties, got nil")
	}

	curSP := cur.Properties.GetOrZero(sel)
	if curSP == nil {
		log.Fatalf("no such property %v", sel)
	}
	cur, err := curSP.BuildSchema()
	if err != nil {
		log.Fatalf("building schema: %v", err)
	}
	return cur
}

func selectOneOfVariant(cur *base.Schema, selector any) *base.Schema {
	var (
		ds                     *base.Schema      // Discriminator schema
		of                     *base.SchemaProxy // Each of the oneOf schema proxies in turn
		err                    error
		variant, discriminator string
		ok                     bool
	)
	if variant, ok = selector.(string); !ok {
		log.Fatalf("OneOf selector must be a string, got %v (%v)", selector, reflect.TypeOf(selector))
	}
	name := cur.ParentProxy.GetReference()
	if discriminator, ok = discriminators[name]; !ok {
		log.Fatalf("no discriminator found for %q", name)
	}

	ofs := cur.OneOf
	for _, of = range ofs {
		s := of.Schema()
		dsp := s.Properties.GetOrZero(discriminator)
		if ds, err = dsp.BuildSchema(); err != nil {
			log.Fatalf("building schema: %v", err)
		}
		ref := dsp.GetReference()
		switch {
		case ds.Type == nil:
			log.Fatalf("no type found for %q", ref)

		case len(ds.Type) == 0:
			log.Fatalf("type length 0 found for %q", ref)

		case ds.Type[0] == "string":
			if len(ds.Enum) == 0 {
				log.Fatalf("enum length 0 found for %q", ref)
			}
			for _, dv := range ds.Enum {
				if dv.Value == variant {
					return s
				}
			}
		default:
			log.Fatalf("unexpected type found for %q", ref)
		}
	}
	log.Fatalf("expected one OneOf to have type %q, found none", variant)
	return nil
}

// loadSchema loads the DataDog OpenAPI reference schema.
// It is copied from https://github.com/DataDog/datadog-api-client-go/blob/master/.generator/schemas/v1/openapi.yaml
// when updating the project.
func loadSchema() base.Schema {
	rawDoc, err := libopenapi.NewDocument(exp.OpenAPI)
	if err != nil {
		log.Fatal(err)
	}
	model, errs := rawDoc.BuildV3Model()
	if len(errs) != 0 {
		for _, err := range errs {
			log.Println(err)
		}
		os.Exit(1)
	}
	dashboardSchema, err := model.Model.Components.Schemas.
		GetOrZero("Dashboard").
		BuildSchema()
	if err != nil {
		log.Fatal(err)
	}
	return *dashboardSchema
}

// Index follows a path of selections into the DataDog client OpenAPI Schema,
// and returns the description of the last component.
//
// Caveats: for elements which are OneOf, like the WidgetDefinition in
// https://github.com/DataDog/datadog-api-client-go/blob/master/.generator/schemas/v1/openapi.yaml
// start by adding a path component for the variant using its discriminator,
// then a second one to index within it.
func Index(path []string) *base.Schema {
	ds := loadSchema()
	cur := &ds
	for _, selector := range path {
		var tl int

	ChooseAction:
		tl = len(cur.Type)
		if tl > 1 {
			log.Fatalf("expected 1 type, got %v", cur.Type)
		}
		switch {
		case cur == nil:
			log.Fatal("expected non-nil schema, got nil")

		case tl == 1 && cur.Type[0] == "array":
			// This is an array, so we need to re-scan for its items type,
			// without progressing into the path.
			cur = selectArrayItem(cur)
			goto ChooseAction

		case tl == 1 && cur.Type[0] == "object":
			// This is an object, so we select one of its properties
			cur = selectObjectProperty(cur, selector)

		case tl == 0 && cur.OneOf != nil:
			// This is a oneOf, so we need to auto-detect which one, because DataDog provided no discriminator.
			cur = selectOneOfVariant(cur, selector)

		default:
			log.Fatalf("unimplemented case: tl = %d, schema: %#v", tl, cur)
		}
	}
	return cur
}
