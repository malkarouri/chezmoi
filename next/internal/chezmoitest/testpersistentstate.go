package chezmoitest

// A PersistentState is a persistent state for testing.
type PersistentState map[string]map[string][]byte

func NewPersistentState() PersistentState {
	return make(PersistentState)
}

func (s PersistentState) Delete(bucket, key []byte) error {
	bucketMap, ok := s[string(bucket)]
	if !ok {
		return nil
	}
	delete(bucketMap, string(key))
	return nil
}

func (s PersistentState) Get(bucket, key []byte) ([]byte, error) {
	bucketMap, ok := s[string(bucket)]
	if !ok {
		return nil, nil
	}
	return bucketMap[string(key)], nil
}

func (s PersistentState) Set(bucket, key, value []byte) error {
	bucketMap, ok := s[string(bucket)]
	if !ok {
		bucketMap = make(map[string][]byte)
		s[string(bucket)] = bucketMap
	}
	bucketMap[string(key)] = value
	return nil
}
