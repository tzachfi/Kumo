package journey

import (
	"fmt"

	"github.com/tzachfi/kumo/server/internal/domain"
)

// Validate checks that enum fields on a Journey tree are recognized values.
// String-typed enums in Go are not enforced at compile time, so this runs at
// the package boundary after assembly and before persistence or API response.
func Validate(j *domain.Journey) error {
	if j == nil {
		return fmt.Errorf("journey: nil journey")
	}
	if !j.State.Valid() {
		return fmt.Errorf("journey: invalid state %q", j.State)
	}

	for i := range j.Milestones {
		m := &j.Milestones[i]
		if !m.State.Valid() {
			return fmt.Errorf("journey: milestone[%d]: invalid state %q", i, m.State)
		}
		for k := range m.Tasks {
			if !m.Tasks[k].State.Valid() {
				return fmt.Errorf("journey: milestone[%d].task[%d]: invalid state %q", i, k, m.Tasks[k].State)
			}
		}
	}

	return nil
}
