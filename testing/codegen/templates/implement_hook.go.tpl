{{ printf "Service must be set by user code (in a _test file) to provide the service implementation being tested." | comment }}
// Example:
//   func init() {
//       {{ .VarName }}test.Service = func(t *testing.T) {{ .PkgName }}.Service { 
//           return NewMyService(...) 
//       }
//   }
var Service func(t *testing.T) {{ .PkgName }}.Service


