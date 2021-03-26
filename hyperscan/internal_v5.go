// +build !hyperscan_v4

package hyperscan

/*
#include <hs.h>
*/
import "C"

const (
	// Combination represents logical combination.
	Combination CompileFlag = C.HS_FLAG_COMBINATION
	// Quiet represents don't do any match reporting.
	Quiet CompileFlag = C.HS_FLAG_QUIET
)

func init() {
	compileFlags['C'] = Combination
	compileFlags['Q'] = Quiet
}
