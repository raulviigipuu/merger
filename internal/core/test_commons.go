package core

import "path/filepath"

// BasePath points to the root-level testdata directory.
// Used by all unit tests in this package.
var BasePath = filepath.Join("..", "..", "testdata")
