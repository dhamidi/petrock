package generator

// ComponentType represents the type of component to generate
type ComponentType string

const (
	ComponentTypeCommand ComponentType = "command"
	ComponentTypeQuery   ComponentType = "query"
	ComponentTypeWorker  ComponentType = "worker"
)

// String returns the string representation of ComponentType
func (ct ComponentType) String() string {
	return string(ct)
}
