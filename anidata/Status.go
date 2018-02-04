package anidata

type Status int

const (
	PlanToWatch Status = iota
	Watching
	Completed
	Dropped
)

func (s Status) String() string {
	switch s {
	case PlanToWatch:
		return "Plan To Watch"
	case Watching:
		return "Watching"
	case Completed:
		return "Completed"
	case Dropped:
		return "Dropped"
	default:
		return "Status Unknown"
	}
}
