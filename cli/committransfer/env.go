package committransfer

import "os"

// currentEnv snapshots os.Environ. Split out so tests can stub it.
var currentEnv = os.Environ
