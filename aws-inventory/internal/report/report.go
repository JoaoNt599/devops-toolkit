package report

import (
	"encoding/json"
	"fmt"
	"os"

	"aws-inventory/internal/models"
)

// Jsonify turns any SDK output struct into a generic map so the raw shape
// (all fields) is preserved in the final document.
func Jsonify(v any) any {
	b, err := json.Marshal(v)
	if err != nil {
		return map[string]any{"_error": err.Error()}
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return map[string]any{"_error": err.Error()}
	}
	return out
}

// SetErr records a per-section error without aborting the rest of the run.
func SetErr(sec map[string]any, key string, err error) {
	sec[key] = map[string]any{"_error": err.Error()}
}

// Write serializes the report to an indented JSON file at path.
func Write(rep *models.Report, path string) error {
	buf, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling report: %w", err)
	}
	if err := os.WriteFile(path, buf, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	fmt.Printf("Architecture report written to %s (%d bytes)\n", path, len(buf))
	return nil
}
