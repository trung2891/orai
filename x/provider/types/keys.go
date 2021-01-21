package types

const (
	// ModuleName is the name of the module
	ModuleName = "provider"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName

	// DefaultQueryPage sets the default page query value
	DefaultQueryPage = 1

	// DefaultQueryLimit sets the default query limit value
	DefaultQueryLimit = 5

	// DefaultValNum default validator number
	DefaultValNum = 1

	// ScriptPath is the path that stores all the script files (oracle scripts, data sources, test cases)
	ScriptPath = ".oraifiles/"

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"

	// Denom is the denominator of the currency
	Denom = "orai"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
