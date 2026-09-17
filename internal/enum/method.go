package enum

type DBMethod string

const (
	DBMethodCreate = "create"
	DBMethodUpdate = "update"
	DBMethodDelete = "delete"
)

func (e DBMethod) Valid() bool {
	switch e {
	case DBMethodCreate, DBMethodUpdate, DBMethodDelete:
		return true
	default:
		return false
	}
}

type DBTable string

const (
	DBTableProfile    = "profile"
	DBTableAbout      = "about"
	DBTableEducation  = "education"
	DBTableExperience = "experience"
	DBTableIcon       = "icon"
	DBTableSocial     = "social"
	DBTableNameserver = "nameserver"
	DBTableDomain     = "domain"
)

func (e DBTable) Valid() bool {
	switch e {
	case DBTableProfile,
		DBTableAbout,
		DBTableEducation,
		DBTableExperience,
		DBTableIcon,
		DBTableSocial,
		DBTableNameserver,
		DBTableDomain:
		return true
	default:
		return false
	}
}
