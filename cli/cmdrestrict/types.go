package cmdrestrict

// Restriction is one named window where commits are not allowed.
type Restriction struct {
	Name          string `json:"Name"`
	BlockedAfter  string `json:"BlockedAfter"`
	Email         string `json:"Email"`
	Timezone      string `json:"Timezone"`
	IsWeekendOpen bool   `json:"IsWeekendOpen"`
	IsPermitStay  bool   `json:"IsPermitStay"`
}

// RestrictionFile is the on-disk list and the machine-wide enforcement mode.
type RestrictionFile struct {
	EnforcementMode string        `json:"EnforcementMode"`
	Items           []Restriction `json:"Items"`
}
