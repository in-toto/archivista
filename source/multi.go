// Copyright 2022 The Witness Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package source

import (
	"context"
	"sync"
)

type MultiSource struct {
	sources []Sourcer
}

func NewMultiSource(sources ...Sourcer) *MultiSource {
	return &MultiSource{sources}
}

// Search concurrently queries all sources and returns the combined results.
func (s *MultiSource) Search(ctx context.Context, collectionName string, subjectDigests, attestations []string) ([]CollectionEnvelope, error) {
	results := []CollectionEnvelope{}
	errors := []error{}

	errs := make(chan error)                   // Channel for collecting errors from each source
	resChan := make(chan []CollectionEnvelope) // Channel for collecting results from each source

	errdone := make(chan bool)    // Signal channel indicating when error collection is done
	readerDone := make(chan bool) // Signal channel indicating when result collection is done

	// Goroutine for collecting results from the result channel
	go func() {
		for item := range resChan {
			results = append(results, item...)
		}
		readerDone <- true
	}()

	// Goroutine for collecting errors from the error channel
	go func() {
		for err := range errs {
			errors = append(errors, err)
		}
		errdone <- true
	}()

	var wg sync.WaitGroup // WaitGroup for waiting on all source queries to finish
	for _, source := range s.sources {
		source := source
		wg.Add(1)
		// Goroutine for querying a source and collecting the results or error
		go func(src Sourcer) {
			defer wg.Done()
			res, err := src.Search(ctx, collectionName, subjectDigests, attestations)
			if err != nil {
				errs <- err
			} else {
				resChan <- res
			}
		}(source)
	}
	wg.Wait()      // Wait for all source queries to finish
	close(resChan) // Close the result channel
	close(errs)    // Close the error channel

	<-errdone    // Wait for error collection to finish
	<-readerDone // Wait for result collection to finish

	// If any errors occurred, return the first error and discard the results
	if len(errors) > 0 {
		return nil, errors[0]
	}
	// Return the combined results from all sources
	return results, nil
}
