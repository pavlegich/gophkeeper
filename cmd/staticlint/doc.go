// The staticlint package contains a set of static source code analyzers,
// combined into a multichecker.
//
// multichecker contains:
//
// - all analyzers from golang.org/x/tools/go/analysis/passes;
//
// - all SA analyzers of staticcheck.io package;
//
// - analyzers S1020, ST1003, QF1003 of staticcheck.io package;
//
// - analyzer of correct request body closing request.Body (github.com/timakin/bodyclose/passes/bodyclose);
//
// - analyzer of correct error wrap (github.com/fatih/errwrap);
//
// Running examples:
//
//	./staticlint -S1020 <path>
//	./staticlint -fieldalignment <path>
//	./staticlint -errwrap <path>
package main
