package enum

type LocationType string

const (
	LocationTypeOnSite LocationType = "onsite"
	LocationTypeRemote LocationType = "remote"
	LocationTypeHybrid LocationType = "hybrid"
)

func (r LocationType) Valid() bool {
	switch r {
	case LocationTypeOnSite, LocationTypeRemote, LocationTypeHybrid:
		return true
	default:
		return false
	}
}
