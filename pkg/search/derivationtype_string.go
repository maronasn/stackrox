// Code generated by "stringer -type=DerivationType"; DO NOT EDIT.

package search

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CountDerivationType-0]
	_ = x[SimpleReverseSortDerivationType-1]
}

const _DerivationType_name = "CountDerivationTypeSimpleReverseSortDerivationType"

var _DerivationType_index = [...]uint8{0, 19, 50}

func (i DerivationType) String() string {
	if i < 0 || i >= DerivationType(len(_DerivationType_index)-1) {
		return "DerivationType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DerivationType_name[_DerivationType_index[i]:_DerivationType_index[i+1]]
}