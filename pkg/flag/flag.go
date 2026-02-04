package flag

import (
	"os"
)

var NEO_LOG = os.Getenv("NEO_LOG")
var NEO_LOG_CHILDREN = isTrue("NEO_LOG_CHILDREN")
var NEO_PRINT_LOGS = isTrue("NEO_PRINT_LOGS")
var NEO_NO_CLEANUP = isTrue("NEO_NO_CLEANUP")
var NEO_PASSPHRASE = os.Getenv("NEO_PASSPHRASE")
var NEO_PULUMI_PATH = os.Getenv("NEO_PULUMI_PATH")

// NEO_BUILD_CONCURRENCY is deprecated, use NEO_FUNCTION_BUILD_CONCURRENCY instead
var NEO_BUILD_CONCURRENCY = os.Getenv("NEO_BUILD_CONCURRENCY")
var NEO_BUILD_CONCURRENCY_FUNCTION = os.Getenv("NEO_BUILD_CONCURRENCY_FUNCTION")
var NEO_BUILD_CONCURRENCY_SITE = os.Getenv("NEO_BUILD_CONCURRENCY_SITE")
var NEO_SKIP_DEPENDENCY_CHECK = isTrue("NEO_SKIP_DEPENDENCY_CHECK")
var NEO_TELEMETRY_DISABLED = isTrue("NEO_TELEMETRY_DISABLED") || isTrue("DO_NOT_TRACK")
var NEO_BUN_VERSION = os.Getenv("NEO_BUN_VERSION")
var NEO_VERBOSE = isTrue("NEO_VERBOSE")
var NEO_EXPERIMENTAL = isTrue("NEO_EXPERIMENTAL") || isTrue("NEO_EXPERIMENTAL_RUN")
var NEO_RUN_ID = os.Getenv("NEO_RUN_ID")
var NEO_SKIP_APPSYNC = isTrue("NEO_SKIP_APPSYNC")
var NEO_NO_BUN = isTrue("NO_BUN") || isTrue("NEO_NO_BUN")

func isTrue(name string) bool {
	val, ok := os.LookupEnv(name)
	if !ok {
		return false
	}
	if val == "1" {
		return true
	}
	if val == "true" {
		return true
	}
	return false
}
