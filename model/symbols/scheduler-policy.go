package symbols

type SchedulerPolicy int

const (
	SchedulerPolicyParent SchedulerPolicy = iota
	SchedulerPolicyAny
)
