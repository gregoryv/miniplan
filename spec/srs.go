package spec

// System requirements specification

var YRA05 = &R{
	ref:    "YRA05",
	Short:  "Delete a Change by uuid suffix reference",
	Reason: "UUIDs are long and not always visible",
	Rel:    []*R{AJE15},
}

var AJE15 = &R{
	ref:    "AJE15",
	Short:  "Change should have a short reference",
	Reason: "Simplifies talking about many changes, makes them unique in a group",
}
