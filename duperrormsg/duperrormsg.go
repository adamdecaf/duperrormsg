package duperrormsg

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is the main analyzer for the duplicate-error checker
var Analyzer = &analysis.Analyzer{
	Name:     "duperror",
	Doc:      "Checks for duplicate error messages across different code paths",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// ErrorInfo stores information about an error message
type ErrorInfo struct {
	Pos       ast.Node // Position in source
	Construct string   // Which error construction method was used
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Map to store error messages and their locations
	errorMap := make(map[string][]ErrorInfo)

	// Get the inspector from the analyzer requirements
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Define the node filter for efficiently inspecting only relevant nodes
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	// Use Preorder to visit all call expressions
	inspector.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)

		// Check if this is a function call we're interested in
		construct, msg := extractErrorMessage(call)
		if construct == "" || msg == "" {
			return
		}

		// Add to our map
		info := ErrorInfo{
			Pos:       node,
			Construct: construct,
		}

		errorMap[msg] = append(errorMap[msg], info)
	})

	// Check for duplicates
	for msg, locations := range errorMap {
		if len(locations) > 1 {
			// Report the first occurrence
			firstLoc := locations[0]
			pass.Reportf(firstLoc.Pos.Pos(), "duplicate error message %q used in multiple locations", msg)

			// Report all subsequent occurrences with reference to the first
			for i := 1; i < len(locations); i++ {
				pass.Reportf(locations[i].Pos.Pos(), "duplicate error message %q also used at %v",
					msg, pass.Fset.Position(firstLoc.Pos.Pos()))
			}
		}
	}

	return nil, nil
}

func extractErrorMessage(call *ast.CallExpr) (string, string) {
	construct := getErrorConstructName(call)
	if construct == "" {
		return "", ""
	}

	var msgArg ast.Expr

	// Check if there are any arguments
	if len(call.Args) == 0 {
		return "", ""
	}

	switch construct {
	case "errors.New":
		// errors.New takes a single string argument
		if len(call.Args) != 1 {
			return "", ""
		}
		msgArg = call.Args[0]

	case "fmt.Errorf":
		// fmt.Errorf takes a format string and optional arguments
		msgArg = call.Args[0]

	case "log", "logger", "Log", "Logf", "LogError", "LogErrorf":
		// Log functions take format string as first argument
		msgArg = call.Args[0]

	default:
		// For custom error constructors that likely take a message as first arg
		// First, check if the first argument is a string
		if len(call.Args) > 0 {
			if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
				msgArg = lit
			} else {
				// If first arg isn't a string, try to find any string literal among arguments
				for _, arg := range call.Args {
					if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						msgArg = lit
						break
					}
				}
			}
		}

		if msgArg == nil {
			return "", ""
		}
	}

	msg := extractStringLiteral(msgArg)
	if msg == "" {
		return "", ""
	}

	return construct, msg
}

func getErrorConstructName(call *ast.CallExpr) string {
	// First, handle chained calls like logger.Info().Logf()
	if selExpr, ok := call.Fun.(*ast.SelectorExpr); ok {
		// Check if the selector's X is another call expression (method chaining)
		if _, ok := selExpr.X.(*ast.CallExpr); ok {
			// This handles chained methods like logger.Info().Logf()
			// For log methods specifically
			if selExpr.Sel.Name == "Logf" ||
				selExpr.Sel.Name == "LogErrorf" ||
				selExpr.Sel.Name == "LogError" ||
				selExpr.Sel.Name == "Log" {
				return selExpr.Sel.Name
			}
		}

		// Check for standard selector expressions (e.g., errors.New, fmt.Errorf)
		if pkgIdent, ok := selExpr.X.(*ast.Ident); ok {
			// Common error construction patterns
			if pkgIdent.Name == "errors" && selExpr.Sel.Name == "New" {
				return "errors.New"
			}
			if pkgIdent.Name == "fmt" && selExpr.Sel.Name == "Errorf" {
				return "fmt.Errorf"
			}

			// Check for logging functions
			if pkgIdent.Name == "log" || strings.Contains(strings.ToLower(pkgIdent.Name), "log") {
				logFuncSuffixes := []string{
					"", "f", "ln", // Log, Logf, Logln
					"Error", "Errorf", "Errorln",
					"Fatal", "Fatalf", "Fatalln",
					"Panic", "Panicf", "Panicln",
					"Warning", "Warningf", "Warningln",
					"Info", "Infof", "Infoln",
				}

				for _, suffix := range logFuncSuffixes {
					if selExpr.Sel.Name == suffix ||
						selExpr.Sel.Name == "Log"+suffix ||
						selExpr.Sel.Name == "Print"+suffix {
						return pkgIdent.Name
					}
				}
			}

			// Check for common error constructor patterns
			if strings.HasSuffix(selExpr.Sel.Name, "Error") ||
				strings.HasPrefix(selExpr.Sel.Name, "New") ||
				strings.Contains(selExpr.Sel.Name, "Error") ||
				strings.Contains(strings.ToLower(selExpr.Sel.Name), "fail") {
				return selExpr.Sel.Name
			}
		}
	}

	// Also check for direct function idents (not selector expressions)
	// This handles cases like NewUserError("message")
	if ident, ok := call.Fun.(*ast.Ident); ok {
		if strings.HasPrefix(ident.Name, "New") &&
			(strings.Contains(ident.Name, "Error") ||
				strings.Contains(ident.Name, "Err") ||
				strings.Contains(ident.Name, "Fail")) {
			return ident.Name
		}
	}

	return ""
}

func extractStringLiteral(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			// Remove quotes and process format strings
			raw := strings.Trim(e.Value, "`\"")

			// For format strings, we normalize format specifiers
			// This approach catches %s, %d, %v, etc.
			formatSpecifier := regexp.MustCompile(`%[a-zA-Z0-9\.\-\+#]*[a-zA-Z]`)
			normalized := formatSpecifier.ReplaceAllString(raw, "%x")

			return normalized
		}
	}
	return ""
}
