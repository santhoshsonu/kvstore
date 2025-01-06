package store

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type FileTranscationLogger struct {
	events       chan<- Event // Write-only channel for sending events
	errors       <-chan error // Read-only channel for receiving errors
	lastSequence uint64       // last used sequence number
	file         *os.File     // the location of the transaction log
}

func (ftl *FileTranscationLogger) WriteDelete(key string) {
	ftl.events <- Event{EventType: EventDelete, Key: key, Value: ""}
}

func (ftl *FileTranscationLogger) WritePut(key, value string) {
	ftl.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (ftl *FileTranscationLogger) Err() <-chan error {
	return ftl.errors
}

// Playback Events
func (ftl *FileTranscationLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(ftl.file)
	outEvent := make(chan Event)
	outError := make(chan error)

	// Spawn a goroutine to read each line, process and send the transaction event to outEvent or outError (in case of error) channel
	go func() {
		log.Println("Initilizing transaction log reader...")
		defer log.Println("Exiting transaction log reader")

		defer close(outEvent)
		defer close(outError)

		var e Event
		for scanner.Scan() {
			line := scanner.Text()
			if n, err := fmt.Sscanf(line, LOG_FORMAT, &e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				if n <= 3 && err != io.EOF {
					outError <- fmt.Errorf("input parse error: %s %w", line, err)
					return
				}
			}
			// sanity check! Are the sequence numbers in increasing order ?
			if ftl.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction sequence out of order")
				return
			}
			ftl.lastSequence = e.Sequence
			outEvent <- e
		}
	}()

	return outEvent, outError
}

// Method to initialize a FileTransactionLogger
func NewFileTranscationLogger(filename string) (TransactionLogger, error) {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("could not open transaction log file: %s %w", filename, err)
	}
	return &FileTranscationLogger{file: fp}, nil
}

func (ftl *FileTranscationLogger) Run() {
	events := make(chan Event, 16) // Make a buffered event channel
	ftl.events = events

	errors := make(chan error, 1) // Make a buffered error channel
	ftl.errors = errors

	// Spawn a goroutine to continuously receive events (from WritePut and WriteDelete) and write to log file
	go func() {
		log.Println("Initializing transaction event listener...")
		defer log.Println("Exiting transaction event listener")
		for e := range events {
			ftl.lastSequence++ // increment the sequence number
			e.Sequence = ftl.lastSequence
			log.Printf("Writing event: %v to transaction log\n", e)
			_, err := fmt.Fprintf(ftl.file, LOG_FORMAT+"\n", ftl.lastSequence, e.EventType, e.Key, e.Value) // write to file
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}
