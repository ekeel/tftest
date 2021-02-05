package model

type Test struct {
	Type         string            `hcl:"type,attr"`
	Name         string            `hcl:",label"`
	QueryBy      string            `hcl:"query_by,attr"`
	InstanceName string            `hcl:"instance_name,optional"`
	InstanceID   string            `hcl:"instance_id,optional"`
	Fields       map[string]string `hcl:"fields,attr"`
	Tags         map[string]string `hcl:"tags,attr"`

	ValidationResults []ValidationResult
}
