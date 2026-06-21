package calculationlink

type LinkSourceOutputValidator interface {
	EnsureValidSourceOutput(Source) error
}

type LinkTargetInputValidator interface {
	EnsureValidTargetInput(Target) error
}
