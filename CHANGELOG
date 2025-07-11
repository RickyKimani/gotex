# Changelog

## [v0.1.3] - 2025-07-11

### Fixed

- **Braced expressions in equation environments** - Superscripts and subscripts with braces (e.g., `x^{ab}`, `a_{xy}`) now parse correctly without errors in equation environments
- **Vinculum positioning** - Square root horizontal bars now align more precisely with the radical symbol

## [v0.1.2] - 2025-07-10

### Added

- **Vinculum rendering for square roots** - Square root expressions now display with proper horizontal bars above the radicand

### Fixed

- **Square root rendering** - Braces `{}` are now correctly parsed as groups instead of plain text
- **Fraction bar length** - Fraction bars now properly extend to cover the last operand in complex expressions

### Known Issues

- **Braced expressions in equation environments** - Superscripts and subscripts with braces (e.g., `x^{ab}`, `a_{xy}`) in equation environments still trigger parsing errors

### Removed

- **sqrt symbol** - Removed from symbols map to prevent conflicts with proper square root command handling

## [v0.1.1] - 2025-07-09

### Added

- **Rust-style colorized error reporting** with severity-based color coding
  - 🟡 Yellow warnings for recoverable issues (e.g., unmatched closing braces)
  - 🔴 Red errors for syntax problems (e.g., missing environment closures)
  - 🔴 Red+White fatal errors for critical parsing failures
  - Precise error location pointers with code context
- **Enhanced CLI filename resolution**
  - Auto-append `.tex` extension when omitted (e.g., `gotex assignment` finds `assignment.tex`)
  - Reject unsupported file extensions with clear error messages
  - Improved user experience for common workflows
- **Colorized scan output** with subtle visual improvements
  - 🟡 Yellow headers for scan operations
  - 🔵 Blue numbering for easier list navigation
  - 🟢 Green file sizes for quick reference
- **Comprehensive error classification**
  - Warning vs Error vs Fatal severity levels
  - Appropriate severity assignment based on error type
  - Better error recovery and user feedback

### Improved

- **Error reporting system** now provides professional, developer-friendly output
- **CLI help text** updated to reflect new filename resolution behavior
- **Parser error handling** with proper warning/error distinction
