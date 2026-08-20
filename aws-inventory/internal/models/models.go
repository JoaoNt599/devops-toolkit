package models

// Report is the top-level document: one key per AWS service.
type Report struct {
	Meta map[string]any            `json:"_meta"`
	Data map[string]map[string]any `json:"architecture"`
}

// Section returns (creating if needed) the sub-map for a given service key,
// e.g. "ec2", "ecs", "elbv2". Collectors write their results into it.
func (r *Report) Section(svc string) map[string]any {
	if r.Data == nil {
		r.Data = map[string]map[string]any{}
	}
	if r.Data[svc] == nil {
		r.Data[svc] = map[string]any{}
	}
	return r.Data[svc]
}
