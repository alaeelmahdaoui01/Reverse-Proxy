package proxy

// from the backend in serverpool since it stores everything and 
func (serverPool *ServerPool) GetNextValidPeer() *Backend {
	if len(serverPool.Backends) == 0 {
		return nil
	}

	// i call the start backend then loop on the backends one after one in a cyclic one
	// i return the first available one/alive 

	for i:=0 ; i<len(serverPool.Backends) ; i++ {

	}

	// should i update serverPool.current ??
	
	return nil // if none found 
}



// TO IMPLEMENT 

// AddBackend(backend *Backend)
// SetBackendStatus(uri *url.URL, alive bool)