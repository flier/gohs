// Chimera is a software regular expression matching engine that is a hybrid of Hyperscan and PCRE.
// The design goals of Chimera are to fully support PCRE syntax as well as to
// take advantage of the high performance nature of Hyperscan.
//
// Chimera inherits the design guideline of Hyperscan with C APIs for compilation and scanning.
//
// The Chimera API itself is composed of two major components:
//
// # Compilation
//
// These functions take a group of regular expressions, along with identifiers and option flags,
// and compile them into an immutable database that can be used by the Chimera scanning API.
// This compilation process performs considerable analysis and optimization work in order to build a database
// that will match the given expressions efficiently.
//
// See Compiling Patterns for more details (https://intel.github.io/hyperscan/dev-reference/chimera.html#chcompile)
//
// # Scanning
//
// Once a Chimera database has been created, it can be used to scan data in memory.
// Chimera only supports block mode in which we scan a single contiguous block in memory.
//
// Matches are delivered to the application via a user-supplied callback function
// that is called synchronously for each match.
//
// For a given database, Chimera provides several guarantees:
//
// 1 No memory allocations occur at runtime with the exception of scratch space allocation,
// it should be done ahead of time for performance-critical applications:
//
// 2 Scratch space: temporary memory used for internal data at scan time.
// Structures in scratch space do not persist beyond the end of a single scan call.
//
// 3 The size of the scratch space required for a given database is fixed and determined at database compile time.
// This means that the memory requirement of the application are known ahead of time,
// and the scratch space can be pre-allocated if required for performance reasons.
//
// 4 Any pattern that has successfully been compiled by the Chimera compiler can be scanned against any input.
// There could be internal resource limits or other limitations caused by PCRE at runtime
// that could cause a scan call to return an error.
//
// * Note
//
// Chimera is designed to have the same matching behavior as PCRE, including greedy/ungreedy, capturing, etc.
// Chimera reports both start offset and end offset for each match like PCRE.
// Different from the fashion of reporting all matches in Hyperscan, Chimera only reports non-overlapping matches.
// For example, the pattern /foofoo/ will match foofoofoofoo at offsets (0, 6) and (6, 12).
//
// * Note
//
// Since Chimera is a hybrid of Hyperscan and PCRE in order to support full PCRE syntax,
// there will be extra performance overhead compared to Hyperscan-only solution.
// Please always use Hyperscan for better performance unless you must need full PCRE syntax support.
//
// See Scanning for Patterns for more details (https://intel.github.io/hyperscan/dev-reference/chimera.html#chruntime)
package chimera
