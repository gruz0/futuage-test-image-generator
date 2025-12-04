package generator

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// GenerationResult represents the result of generating a single image
type GenerationResult struct {
	Spec     ImageSpec
	FileSize int64
	Error    error
}

// GenerationStats tracks statistics during generation
type GenerationStats struct {
	Total     int
	Completed int32
	Failed    int32
	StartTime time.Time
	EndTime   time.Time
}

// ProgressCallback is called periodically with generation progress
type ProgressCallback func(completed, total int, elapsed time.Duration)

// Orchestrator manages parallel image generation
type Orchestrator struct {
	MaxConcurrency int
	Stats          GenerationStats
	progressCb     ProgressCallback
}

// NewOrchestrator creates a new Orchestrator
func NewOrchestrator(maxConcurrency int) *Orchestrator {
	if maxConcurrency <= 0 {
		maxConcurrency = 10 // Default concurrency
	}

	return &Orchestrator{
		MaxConcurrency: maxConcurrency,
	}
}

// SetProgressCallback sets a callback for progress updates
func (o *Orchestrator) SetProgressCallback(cb ProgressCallback) {
	o.progressCb = cb
}

// GenerateAll generates all images in parallel
func (o *Orchestrator) GenerateAll(specs []ImageSpec) ([]GenerationResult, error) {
	o.Stats = GenerationStats{
		Total:     len(specs),
		StartTime: time.Now(),
	}

	if len(specs) == 0 {
		return nil, fmt.Errorf("no specs to generate")
	}

	// Channel for results
	resultsChan := make(chan GenerationResult, len(specs))

	// Create semaphore for concurrency control
	sem := make(chan struct{}, o.MaxConcurrency)

	// Wait group for all goroutines
	var wg sync.WaitGroup

	// Start progress reporter
	stopProgress := make(chan struct{})
	go o.reportProgress(stopProgress)

	// Launch generation goroutines
	for _, spec := range specs {
		wg.Add(1)
		go func(s ImageSpec) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			// Generate image
			result := o.generateOne(s)

			// Send result
			resultsChan <- result

			// Update stats
			if result.Error == nil {
				atomic.AddInt32(&o.Stats.Completed, 1)
			} else {
				atomic.AddInt32(&o.Stats.Failed, 1)
			}
		}(spec)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(resultsChan)
	close(stopProgress)

	o.Stats.EndTime = time.Now()

	// Collect results
	results := make([]GenerationResult, 0, len(specs))
	for result := range resultsChan {
		results = append(results, result)
	}

	// Check if any failed
	if o.Stats.Failed > 0 {
		return results, fmt.Errorf("%d image(s) failed to generate", o.Stats.Failed)
	}

	return results, nil
}

// generateOne generates a single image
func (o *Orchestrator) generateOne(spec ImageSpec) GenerationResult {
	result := GenerationResult{
		Spec: spec,
	}

	// Generate the image
	if err := Generate(spec); err != nil {
		result.Error = fmt.Errorf("failed to generate %s: %w", spec.Filename, err)
		return result
	}

	// Get file size
	fileSize, err := getFileSize(spec.OutputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to get file size for %s: %w", spec.Filename, err)
		return result
	}

	result.FileSize = fileSize
	return result
}

// reportProgress reports progress periodically
func (o *Orchestrator) reportProgress(stop chan struct{}) {
	if o.progressCb == nil {
		return
	}

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			// Final update
			completed := int(atomic.LoadInt32(&o.Stats.Completed))
			elapsed := time.Since(o.Stats.StartTime)
			o.progressCb(completed, o.Stats.Total, elapsed)
			return
		case <-ticker.C:
			completed := int(atomic.LoadInt32(&o.Stats.Completed))
			elapsed := time.Since(o.Stats.StartTime)
			o.progressCb(completed, o.Stats.Total, elapsed)
		}
	}
}

// Duration returns the total generation duration
func (s *GenerationStats) Duration() time.Duration {
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// ImagesPerSecond returns the generation rate
func (s *GenerationStats) ImagesPerSecond() float64 {
	duration := s.Duration()
	if duration == 0 {
		return 0
	}
	return float64(s.Completed) / duration.Seconds()
}

// getFileSize returns the size of a file in bytes
func getFileSize(path string) (int64, error) {
	// Import moved to avoid circular dependency
	// Use the filesystem package helper instead
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
