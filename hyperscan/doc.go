/*

	Hyperscan (https://github.com/01org/hyperscan) is a software regular expression matching engine designed with high performance and flexibility in mind.	It is implemented as a library that exposes a straightforward C API.

	Hyperscan uses hybrid automata techniques to allow simultaneous matching of large numbers (up to tens of thousands) of regular expressions and for the matching of regular expressions across streams of data.

	Hyperscan is typically used in a DPI library stack.

	The Hyperscan API itself is composed of two major components:

	Compilation

	These functions take a group of regular expressions, along with identifiers and option flags, and compile them into an immutable database that can be used by the Hyperscan scanning API. This compilation process performs considerable analysis and optimization work in order to build a database that will match the given expressions efficiently.

	If a pattern cannot be built into a database for any reason (such as the use of an unsupported expression construct, or the overflowing of a resource limit), an error will be returned by the pattern compiler.

	Compiled databases can be serialized and relocated, so that they can be stored to disk or moved between hosts. They can also be targeted to particular platform features (for example, the use of Intel® Advanced Vector Extensions 2 (Intel® AVX2) instructions).

	See Compiling Patterns for more detail. (http://01org.github.io/hyperscan/dev-reference/compilation.html)

	Scanning

	Once a Hyperscan database has been created, it can be used to scan data in memory. Hyperscan provides several scanning modes, depending on whether the data to be scanned is available as a single contiguous block, whether it is distributed amongst several blocks in memory at the same time, or whether it is to be scanned as a sequence of blocks in a stream.

	Matches are delivered to the application via a user-supplied callback function that is called synchronously for each match.

	For a given database, Hyperscan provides several guarantees:

	1. No memory allocations occur at runtime with the exception of two fixed-size allocations, both of which should be done ahead of time for performance-critical applications:
		- Scratch space: temporary memory used for internal data at scan time.
		  Structures in scratch space do not persist beyond the end of a single scan call.
		- Stream state: in streaming mode only, some state space is required to store
		  data that persists between scan calls for each stream. This allows Hyperscan to
		  track matches that span multiple blocks of data.

	2. The sizes of the scratch space and stream state (in streaming mode) required for a given database are fixed and determined at database compile time. This means that the memory requirements of the application are known ahead of time, and these structures can be pre-allocated if required for performance reasons.

	3. Any pattern that has successfully been compiled by the Hyperscan compiler can be scanned against any input. There are no internal resource limits or other limitations at runtime that could cause a scan call to return an error.

	See Scanning for Patterns for more detail. (http://01org.github.io/hyperscan/dev-reference/runtime.html)

	Building a Database

	The Hyperscan compiler API accepts regular expressions and converts them into a compiled pattern database that can then be used to scan data.

	Compilation allows the Hyperscan library to analyze the given pattern(s) and pre-determine how to scan for these patterns in an optimized fashion that would be far too expensive to compute at run-time.

	When compiling expressions, a decision needs to be made whether the resulting compiled patterns are to be used in a streaming, block or vectored mode:

		- Streaming mode: the target data to be scanned is a continuous stream, not all of
		  which is available at once; blocks of data are scanned in sequence and matches may
		  span multiple blocks in a stream. In streaming mode, each stream requires a block
		  of memory to store its state between scan calls.
		- Block mode: the target data is a discrete, contiguous block which can be scanned
		  in one call and does not require state to be retained.
		- Vectored mode: the target data consists of a list of non-contiguous blocks that are
		  available all at once. As for block mode, no retention of state is required.
*/
package hyperscan
