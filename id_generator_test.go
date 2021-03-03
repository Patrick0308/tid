package tid

import (
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
)

func TestDefaultGenerator(t *testing.T) {
	idg := NewDefaultGenerator("http://127.0.0.1:8080")
	wg := &sync.WaitGroup{}
	wg.Add(100)


	for i := 0; i < 100; i++ {
		go func() {
			id, err := idg.get()
			if err != nil {
				log.Printf("%v", err)
			}

			assert.NoError(t, err)
			log.Printf(id)
			wg.Done()
		}()
	}
	wg.Wait()
}
