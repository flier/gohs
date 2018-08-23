// +build !hyperscan_v4

package hyperscan

/*
#include <hs.h>
*/
import "C"

const (
	Combination     CompileFlag = C.HS_FLAG_COMBINATION  // Logical combination.
	Quiet           CompileFlag = C.HS_FLAG_QUIET        // Don't do any match reporting.
)

func init() {
	compileFlags['C'] = Combination
	compileFlags['Q'] = Quiet
}