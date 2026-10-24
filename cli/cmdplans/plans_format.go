package cmdplans

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/constants"
)

// RenderJSON serializes the plans report as pretty-printed JSON to standard output.
func RenderJSON(report PlansReport) error {
	data, err := json.MarshalIndent(report, "", constants.DefaultJSONIndent)
	if err != nil {
		return apperror.WrapSimple(err, "json.MarshalIndent plans")
	}

	fmt.Fprintln(os.Stdout, string(data))

	return nil
}

// RenderYAML serializes the plans report as YAML document to standard output.
func RenderYAML(report PlansReport) error {
	data, err := yaml.Marshal(report)
	if err != nil {
		return apperror.WrapSimple(err, "yaml.Marshal plans")
	}

	fmt.Fprint(os.Stdout, string(data))

	return nil
}
