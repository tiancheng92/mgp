package mgp

type Group struct {
	GroupName string
	Routes    []*Route
	Groups    []*Group
}
