package template

// SprigFunctions documents commonly used Sprig template functions available in BMAD CLI templates.
// This file serves as reference documentation for template authors.
//
// The Sprig library provides 200+ template functions for Go's text/template and html/template.
// BMAD CLI integrates Sprig alongside custom functions like toYaml and nindent.
//
// Common Sprig Function Categories:
//
// String Functions:
//   - trim, trimAll, trimSuffix, trimPrefix: Remove whitespace or specific characters
//   - upper, lower, title, untitle: Change case
//   - repeat, replace, regexMatch, regexReplaceAll: String manipulation
//   - split, splitList, join: String/list operations
//   - contains, hasPrefix, hasSuffix: String testing
//   - quote, squote: Add quotes around strings
//
// Default/Conditional:
//   - default: Provide default value if empty
//   - empty, coalesce, ternary: Conditional logic
//   - required: Fail if value is empty
//   - fail: Explicitly fail template execution with message
//
// Date Functions:
//   - now: Current time
//   - date: Format date using Go layout (e.g., "2006-01-02")
//   - dateModify: Modify date by duration (e.g., "+24h")
//   - dateInZone: Format date in specific timezone
//   - ago: Time duration since given time
//
// Data Structure:
//   - list, dict: Create lists and dictionaries
//   - keys, values, has: Work with maps
//   - merge, mergeOverwrite: Combine maps
//   - append, prepend: Modify lists
//   - first, last, rest, initial: List operations
//   - compact, uniq: Filter lists
//
// Type Conversion:
//   - toString, toInt, toFloat, toBool: Type conversions
//   - kindOf, kindIs: Type checking
//   - typeOf, typeIsLike: Advanced type operations
//
// Math Functions:
//   - add, sub, mul, div, mod: Basic arithmetic
//   - min, max, ceil, floor, round: Math operations
//   - seq: Generate number sequences
//
// Encoding:
//   - b64enc, b64dec: Base64 encoding/decoding
//   - sha256sum, sha1sum: Hashing functions
//   - urlquery: URL query encoding
//
// Control Flow:
//   - with: Change context scope
//   - range: Iterate over collections
//   - if/else: Conditional blocks
//
// Examples for BMAD CLI Templates:
//
// 1. String cleaning and formatting:
//   {{.Title | trim | upper}}
//   {{.Description | default "No description provided"}}
//
// 2. Date formatting:
//   {{now | date "2006-01-02 15:04:05"}}
//   {{.CreatedAt | date "Jan 2, 2006"}}
//
// 3. Conditional logic:
//   {{.Value | default "unknown" | lower}}
//   {{if .Optional}}{{.Optional}}{{else}}N/A{{end}}
//
// 4. List operations:
//   {{.Tags | join ", " | quote}}
//   {{range .Items | compact}}{{.}}{{end}}
//
// 5. Map operations:
//   {{range $key, $value := .Config}}{{$key}}: {{$value}}{{end}}
//   {{.Metadata | keys | sortAlpha}}
//
// For the complete function reference, see:
// https://masterminds.github.io/sprig/
//
// Custom BMAD Functions (override Sprig if conflicts exist):
//   - toYaml: Convert value to YAML string
//   - nindent: Indent text by specified number of spaces
