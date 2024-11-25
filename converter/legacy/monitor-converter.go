package legacy

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fgm/jastify/converter"
)

// MONITOR and OPTIONS definitions
var MONITOR = map[string]stringFunc{
	"message": stringGen("message"),
	"name":    stringGen("name"),
	"query":   stringGen("query"),
	"type":    stringGen("type"),
	"options": func(v any) string {
		return "\n// Options" + mapContents(v.(converter.Jmap), func(k1 string, v1 any) string {
			return converter.Must(convertFromDefinition(OPTIONS, k1, v1))
		}) + "\n// /Options\n"
	},
	"id":               blankGen,
	"tags":             stringGen("tags"),
	"priority":         stringGen("priority"),
	"restricted_roles": stringGen("restricted_roles"),
}

var OPTIONS = map[string]stringFunc{
	"enable_logs_sample":     stringGen("enable_logs_sample"),
	"escalation_message":     stringGen("escalation_message"),
	"evaluation_delay":       stringGen("evaluation_delay"),
	"force_delete":           stringGen("force_delete"),
	"groupby_simple_monitor": stringGen("groupby_simple_monitor"),
	"include_tags":           stringGen("include_tags"),
	"locked":                 stringGen("locked"),
	"new_group_delay":        stringGen("new_group_delay"),
	"new_host_delay":         stringGen("new_host_delay"),
	"no_data_timeframe":      stringGen("no_data_timeframe"),
	"notify_audit":           stringGen("notify_audit"),
	"notify_no_data":         stringGen("notify_no_data"),
	"on_missing_data":        stringGen("on_missing_data"),
	"renotify_interval":      stringGen("renotify_interval"),
	"renotify_statuses":      stringGen("renotify_statuses"),
	"require_full_window":    stringGen("require_full_window"),
	"restricted_roles":       stringGen("restricted_roles"),
	"silenced":               blankGen, // Deprecated
	"threshold_windows": func(v any) string {
		return block("monitor_threshold_windows", v.(converter.Jmap), AssignmentString)
	},
	"thresholds": func(v any) string {
		return block("monitor_thresholds", v.(converter.Jmap), AssignmentString)
	},
	"timeout_h": stringGen("timeout_h"),
	"validate":  stringGen("timeout_h"),
}

func GenerateMonitorTerraformCode(resourceName string, data converter.Jmap) (string, error) {
	var (
		result strings.Builder
		keys   = make([]string, 0, len(data))
	)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s, err := convertFromDefinition(MONITOR, k, data[k])
		if err != nil {
			return "", err
		}
		result.WriteString(s)
	}
	return fmt.Sprintf("resource \"datadog_monitor\" \"%s\" {\n%s}\n", resourceName, result.String()), nil
}
