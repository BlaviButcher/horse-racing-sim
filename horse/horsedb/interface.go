package horsedb

// This will allow us to have dynamic methods for the database. Ability to set and read multiple types
// with the same method (accepts a Value)
type Value interface {
	// Obtain the key for this value. Every value should be self-describing in that way
	Key() Key
	// Some values omit storing the data in the key, so this provides a way of updating the value with the raw key bytes
	updateWithKey([]byte)

	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
