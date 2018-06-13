package domain

// Scope contains a list of permission scopes
type Scope []string

// HasSubscope returns whether this scope has the scope passed as argument
func (s Scope) HasSubscope(subScope Scope) bool {
OUTER:
	for tryScope := range subScope {
		for allowedScope := range s {
			if tryScope == allowedScope {
				continue OUTER
			}
		}
		return false
	}
	return true
}
