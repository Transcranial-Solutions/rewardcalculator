package db

// Bucket
type Bucket interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) bool
	Set(key []byte, value []byte) error
	Delete(key []byte) error
}

type BucketID string

//	Bucket ID
const (
	// For I-Score DB
	// I-Score
	PrefixIScore BucketID             = ""

	// For global DB

	// DB information for management
	PrefixDBInfo BucketID             = "DI"

	// Block information for management
	PrefixBlockInfo BucketID          = "BI"

	// IISS governance variable
	PrefixGovernanceVariable BucketID = "GV"

	// P-Rep candidate list
	PrefixPrepCandidate BucketID      = "PC"

	// FOR IISS data DB
	// Header
	PrefixIISSHeader BucketID         = "HD"

	// Governance variable
	PrefixIISSGV BucketID             = "gv"

	// P-Rep list
	PrefixIISSPRep BucketID           = "prep"

	// TX
	PrefixIISSTX BucketID             = "TX"

)

// internalKey returns key prefixed with the bucket's id.
func internalKey(id BucketID, key []byte) []byte {
	buf := make([]byte, len(key)+len(id))
	copy(buf, id)
	copy(buf[len(id):], key)
	return buf
}

// nonNilBytes returns empty []byte if bz is nil
func nonNilBytes(bz []byte) []byte {
	if bz == nil {
		return []byte{}
	}
	return bz
}