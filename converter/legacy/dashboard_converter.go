package legacy

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/fgm/jastify/converter"
)

type (
	stringFunc func(any) string
)

var (
	blankGen = stringFunc(func(_ any) string { return "" })
)

func stringGen(name string) stringFunc {
	return func(v any) string { return AssignmentString(name, v) }
}

var WIDGET_DEFINITION map[string]stringFunc

var DASHBOARD = map[string]stringFunc{
	"dashboard_lists_removed": blankGen,
	"description":             stringGen("description"),
	"id":                      blankGen,
	"is_read_only":            stringGen("is_read_only"),
	"layout_type":             stringGen("layout_type"),
	"notify_list":             stringGen("notify_list"),
	"reflow_type":             stringGen("reflow_type"),
	"restricted_roles":        stringGen("restricted_roles"),
	"template_variables": func(v any) string {
		slice, ok := v.([]any)
		if !ok {
			log.Fatalf("template_variables expected as []any but got %T: %#v\n", v, v)
		}
		tvs := make(converter.Jmaps, len(slice))
		for i, tv := range slice {
			tvs[i], ok = tv.(converter.Jmap)
			if !ok {
				log.Fatalf("template_variables[%d] expected as Jmaps but got %T: %#v\n", i, tv, tv)
			}
		}
		return blockList(tvs, "template_variable", AssignmentString)
	},
	"template_variable_presets": func(v any) string {
		presets := converter.Must(converter.JmapsFromAny(v))
		return blockList(presets, "template_variable_preset", func(k1 string, v1 any) string {
			return converter.Must(convertFromDefinition(TEMPLATE_VARIABLE_PRESET, k1, v1))
		})
	},
	"tags":  stringGen("tags"),
	"title": stringGen("title"),
	"url":   stringGen("url"),
	"widgets": func(v any) string {
		widgets := converter.Must(converter.JmapsFromAny(v))
		return convertWidgets(widgets)
	},
}

var EVENT_QUERY = map[string]stringFunc{
	"aggregator":  stringGen("aggregator"),
	"compute":     func(v any) string { return block("compute", v.(converter.Jmap), AssignmentString) },
	"data_source": stringGen("data_source"),
	"group_by": func(v any) string {
		groups := converter.Must(converter.JmapsFromAny(v))
		return blockList(groups, "group_by", func(k1 string, v1 any) string {
			return converter.Must(convertFromDefinition(EVENT_QUERY_GROUP_BY, k1, v1))
		})
	},
	"indexes": stringGen("indexes"),
	"name":    stringGen("name"),
	"search":  func(v any) string { return block("search", v.(converter.Jmap), AssignmentString) },
}

var EVENT_QUERY_GROUP_BY = map[string]stringFunc{
	"facet": stringGen("facet"),
	"limit": stringGen("limit"),
	"sort":  func(v any) string { return block("sort", v.(converter.Jmap), AssignmentString) },
}

var FORMULA = map[string]stringFunc{
	"alias":             stringGen("alias"),
	"cell_display_mode": stringGen("cell_display_mode"),
	"formula":           stringGen("formula_expression "),
	"limit": func(v any) string {
		return block("limit", v.(converter.Jmap), func(k1 string, v1 any) string {
			return converter.Must(convertFromDefinition(FORMULA_LIMIT, k1, v1))
		})
	},
	//"number_format": func(v any) string {
	//	return block("number_format", v.(converter.Jmap), func(k1 string, v1 any) string {
	//		return converter.Must(convertFromDefinition(NUMBER_FORMAT, k1, v1))
	//	})
	//},
	//"style": func(v any) string { return blockList(Jmaps{v.(converter.Jmap)}, "style", AssignmentString) },
}

var FORMULA_LIMIT = map[string]stringFunc{
	"count": stringGen("count"),
	"order": stringGen("order"),
}

var GROUP_BY = map[string]stringFunc{
	"facet":      stringGen("facet"),
	"limit":      stringGen("limit"),
	"sort":       func(v any) string { return block("sort_query", v.(converter.Jmap), AssignmentString) },
	"sort_query": func(v any) string { return block("sort_query", v.(converter.Jmap), AssignmentString) },
}

var LOG_QUERY = map[string]stringFunc{
	"compute": func(v any) string {
		return block("compute_query", v.(converter.Jmap), AssignmentString)
	},
	"group_by": func(v any) string {
		groups := converter.Must(converter.JmapsFromAny(v))
		return blockList(groups, "group_by", func(k1 string, v1 any) string {
			return converter.Must(convertFromDefinition(GROUP_BY, k1, v1))
		})
	},
	"index": stringGen("index"),
	"multi_compute": func(v any) string {
		comps := converter.Must(converter.JmapsFromAny(v))
		return blockList(comps, "multi_compute", AssignmentString)
	},
	"search": func(v any) string {
		return AssignmentString("search_query", v.(converter.Jmap)["query"])
	},
	"search_query": stringGen("search_query"),
}

var NUMBER_FORMAT = map[string]stringFunc{
	"unit": func(v any) string {
		return block("unit", v.(converter.Jmap), AssignmentString)
	},
}

var QUERY = map[string]stringFunc{
	"name": stringGen("name"),
	// "indexes":      stringGen("indexes"),
	"data_source": stringGen("data_source"),
	"query":       stringGen("query"),
	//"query_string": stringGen("query_string"),
	//"sort":         stringGen("sort"),
	//"storage":      stringGen("storage"), // TODO validate value "hot"
}

var REQUEST = map[string]stringFunc{
	"aggregator":        stringGen("aggregator"),
	"alias":             stringGen("alias"),
	"apm_query":         stringGen("apm_query"),
	"apm_stats_query":   stringGen("apm_stats_query"),
	"cell_display_mode": stringGen("cell_display_mode"),
	"change_type":       stringGen("change_type"),
	//"columns": func(v any) string {
	//	values := converter.Must(JmapsFromAny(v))
	//	return blockList(values, "columns", func(k string, v any) string {
	//		return converter.Must(convertFromDefinition(REQUEST_COLUMNS, k, v))
	//	})
	//},
	"compare_to": stringGen("compare_to"),
	"conditional_formats": func(v any) string {
		formats := converter.Must(converter.JmapsFromAny(v))
		return blockList(formats, "conditional_formats", AssignmentString)
	},
	"display_type": stringGen("display_type"),
	"fill":         func(v any) string { return block("fill", v.(converter.Jmap), AssignmentString) },
	"formulas": func(v any) string {
		fs := converter.Must(converter.JmapsFromAny(v))
		return blockList(fs, "formula", func(k1 string, v1 any) string {
			return converter.Must(convertFromDefinition(FORMULA, k1, v1))
		})
	},
	"increase_good": stringGen("increase_good"),
	"limit":         stringGen("limit"),
	"log_query": func(v any) string {
		return block("log_query", v.(converter.Jmap), func(k1 string, v1 any) string {
			return converter.Must(convertFromDefinition(LOG_QUERY, k1, v1))
		})
	},
	"metadata": func(v any) string {
		meta := converter.Must(converter.JmapsFromAny(v))
		return blockList(meta, "metadata", AssignmentString)
	},
	"network_query":  stringGen("network_query"),
	"on_right_yaxis": stringGen("on_right_yaxis"),
	"order":          stringGen("order"),
	"order_by":       stringGen("order_by"),
	"order_dir":      stringGen("order_dir"),
	"process_query":  stringGen("process_query"),
	"q":              stringGen("q"),
	"queries": func(v any) string {
		queries := converter.Must(converter.JmapsFromAny(v))
		return queryBlockList(queries, AssignmentString)
	},
	//"query": func(v any) string {
	//	return block("query", v.(converter.Jmap), func(k string, v any) string {
	//		return converter.Must(convertFromDefinition(QUERY, k, v))
	//	})
	//},
	"response_format": func(v any) string {
		if v == "scalar" || v == "timeseries" {
			return ""
		}
		return AssignmentString("response_format", v)
	},
	"rum_query":      stringGen("rum_query"),
	"security_query": stringGen("security_query"),
	"show_present":   stringGen("show_present"),
	//"sort": func(v any) string {
	//	return block("sort", v.(converter.Jmap), func(k1 string, v1 any) string {
	//		return converter.Must(convertFromDefinition(REQUEST_SORT, k1, v1))
	//	})
	//},
	"style": func(v any) string { return blockList(converter.Jmaps{v.(converter.Jmap)}, "style", AssignmentString) },
	//"text_formats": func(v any) string {
	//	formats := converter.Must(JmapsFromAny(v))
	//	return blockList(formats, "text_formats", AssignmentString)
	//},
}

var REQUEST_COLUMNS = map[string]stringFunc{
	"field": stringGen("field"),
	"width": stringGen("width"),
}

var REQUEST_SORT = map[string]stringFunc{
	"count": stringGen("count"),
	"order_by": func(v any) string {
		return "" // FIXME
		orders := converter.Must(converter.JmapsFromAny(v))
		return blockList(orders, "order", AssignmentString)
	},
}

var TEMPLATE_VARIABLE_PRESET = map[string]stringFunc{
	"name": stringGen("name"),
	"template_variables": func(v any) string {
		vars := converter.Must(converter.JmapsFromAny(v))
		return blockList(vars, "template_variable", AssignmentString)
	},
}

var WIDGET = map[string]stringFunc{
	"definition": func(v any) string { return widgetDefinition(v.(converter.Jmap)) },
	"id":         blankGen,
	"layout": func(v any) string {
		return block("widget_layout", v.(converter.Jmap), AssignmentString)
	},
}

func init() {
	WIDGET_DEFINITION = map[string]stringFunc{
		"alert_id":         stringGen("alert_id"),
		"autoscale":        stringGen("autoscale"),
		"background_color": stringGen("background_color"),
		"check":            stringGen("check"),
		"color":            stringGen("color"),
		"color_by_groups":  stringGen("color_by_groups"),
		"color_preference": stringGen("color_preference"),
		"columns":          stringGen("columns"),
		"content":          stringGen("content"),
		"count":            blankGen,
		"custom_links": func(v any) string {
			links := converter.Must(converter.JmapsFromAny(v))
			return blockList(links, "custom_link", AssignmentString)
		},
		"custom_unit":    stringGen("custom_unit"),
		"display_format": stringGen("display_format"),
		"env":            stringGen("env"),
		"event":          func(v any) string { return block("event", v.(converter.Jmap), AssignmentString) },
		"events": func(v any) string {
			events := converter.Must(converter.JmapsFromAny(v))
			return blockList(events, "event", AssignmentString)
		},
		"event_size":            stringGen("event_size"),
		"filters":               stringGen("filters"),
		"font_size":             stringGen("font_size"),
		"global_time_target":    stringGen("global_time_target"),
		"group":                 stringGen("group"),
		"group_by":              stringGen("group_by"),
		"grouping":              stringGen("grouping"),
		"has_padding":           blankGen,
		"has_search_bar":        stringGen("has_search_bar"),
		"hide_zero_counts":      stringGen("hide_zero_counts"),
		"indexes":               stringGen("indexes"),
		"last_triggered_format": stringGen("last_triggered_format"),
		"layout_type":           stringGen("layout_type"),
		"legend_columns":        stringGen("legend_columns"),
		"legend_layout":         stringGen("legend_layout"),
		"legend_size":           stringGen("legend_size"),
		"live_span":             stringGen("live_span"),
		"logset":                blankGen,
		"margin":                stringGen("margin"),
		"markers": func(v any) string {
			markers := converter.Must(converter.JmapsFromAny(v))
			return blockList(markers, "marker", AssignmentString)
		},
		"message_display":     stringGen("message_display"),
		"no_group_hosts":      stringGen("no_group_hosts"),
		"no_metric_hosts":     stringGen("no_metric_hosts"),
		"node_type":           stringGen("node_type"),
		"precision":           stringGen("precision"),
		"query":               stringGen("query"),
		"requests":            convertRequests,
		"right_yaxis":         func(v any) string { return block("right_yaxis", v.(converter.Jmap), AssignmentString) },
		"scope":               stringGen("scope"),
		"service":             stringGen("service"),
		"show_breakdown":      stringGen("show_breakdown"),
		"show_date_column":    stringGen("show_date_column"),
		"show_distribution":   stringGen("show_distribution"),
		"show_error_budget":   stringGen("show_error_budget"),
		"show_errors":         stringGen("show_errors"),
		"show_hits":           stringGen("show_hits"),
		"show_last_triggered": stringGen("show_last_triggered"),
		"show_latency":        stringGen("show_latency"),
		"show_legend":         stringGen("show_legend"),
		"show_message_column": stringGen("show_message_column"),
		"show_priority":       stringGen("show_priority"),
		"show_resource_list":  stringGen("show_resource_list"),
		"show_status":         stringGen("show_status"),
		"show_tick":           stringGen("show_tick"),
		"show_title":          stringGen("show_title"),
		"size_format":         stringGen("size_format"),
		"sizing":              stringGen("sizing"),
		"slo_id":              stringGen("slo_id"),
		"sort":                convertSort,
		"span_name":           stringGen("span_name"),
		"start":               blankGen,
		"style":               func(v any) string { return block("style", v.(converter.Jmap), AssignmentString) },
		"summary_type":        stringGen("summary_type"),
		"tags":                stringGen("tags"),
		"tags_execution":      stringGen("tags_execution"),
		"text":                stringGen("text"),
		"text_align":          stringGen("text_align"),
		"tick_edge":           stringGen("tick_edge"),
		"tick_pos":            stringGen("tick_pos"),
		"time": func(v any) string {
			if liveSpan, ok := v.(converter.Jmap)["live_span"]; ok {
				return AssignmentString("live_span", liveSpan)
			}
			return ""
		},
		"time_windows": stringGen("time_windows"),
		//"timeseries_background": func(v any) string {
		//	// TODO Validate constraint: Max block length == 1
		//	return block("timeseries_background", v.(converter.Jmap), AssignmentString)
		//},
		"title":          stringGen("title"),
		"title_align":    stringGen("title_align"),
		"title_size":     stringGen("title_size"),
		"type":           blankGen,
		"unit":           stringGen("unit"),
		"url":            stringGen("url"),
		"vertical_align": blankGen,
		"view_mode":      stringGen("view_mode"),
		"view_type":      stringGen("view_type"),
		"viz_type":       stringGen("viz_type"),
		"widget_layout": func(v any) string {
			return block("widget_layout", v.(converter.Jmap), AssignmentString)
		},
		"widgets": func(v any) string {
			return convertWidgets(converter.Must(converter.JmapsFromAny(v)))
		},
		"xaxis": func(v any) string { return block("xaxis", v.(converter.Jmap), AssignmentString) },
		"yaxis": func(v any) string { return block("yaxis", v.(converter.Jmap), AssignmentString) },
	}
}

func ConvertEventQuery(value converter.Jmap) string {
	return block("query", value, func(_ string, _ any) string {
		return blockList(converter.Jmaps{value}, "event_query", func(k1 string, v1 any) string {
			return converter.Must(convertFromDefinition(EVENT_QUERY, k1, v1))
		})
	})
}

// convertRequests accepts either a single request as a Jmap or a requests Jmaps.
func convertRequests(value any) string {
	if reflect.ValueOf(value).Kind() == reflect.Slice {
		values := converter.Must(converter.JmapsFromAny(value))
		return blockList(values, "request", func(k string, v any) string {
			return converter.Must(convertFromDefinition(REQUEST, k, v))
		})
	}
	return block("request", value.(converter.Jmap), func(k string, v any) string {
		return converter.Must(convertFromDefinition(REQUEST, k, v))
	})
}

func convertSort(v any) string {
	if sortStr, ok := v.(string); ok {
		return AssignmentString("sort", sortStr)
	}
	return block("sort", v.(converter.Jmap), AssignmentString)
}

func convertWidgets(value converter.Jmaps) string {
	return blockList(value, "widget", func(k1 string, v1 any) string {
		return converter.Must(convertFromDefinition(WIDGET, k1, v1))
	})
}

func widgetDefinition(contents converter.Jmap) string {
	definitionType := contents["type"].(string)
	if definitionType == "slo" {
		definitionType = "service_level_objective"
	}
	return block(fmt.Sprintf("%s_definition", definitionType), contents, func(k string, v any) string {
		return converter.Must(convertFromDefinition(WIDGET_DEFINITION, k, v))
	})
}

func GenerateDashboardTerraformCode(resourceName string, data converter.Jmap) (string, error) {
	var (
		result strings.Builder
		keys   = make([]string, 0, len(data))
	)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s, err := convertFromDefinition(DASHBOARD, k, data[k])
		if err != nil {
			return "", err
		}
		result.WriteString(s)
	}
	return fmt.Sprintf("resource \"datadog_dashboard\" \"%s\" {\n%s}\n", resourceName, result.String()), nil
}
