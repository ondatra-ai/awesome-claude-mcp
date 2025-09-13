package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Default output path when not provided as an argument
const defaultCurrentConversation = "tmp/CONV_CURRENT.json"

// Note: Keep the exact misspelling to match the spec
const sentinelJSON = "{ \"id\": \"No More Converations\" }\n"

type conversationHeader struct {
	ID string `json:"id"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <conversations_json> <processed_ids_txt> [output_json]\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	conversationsPath := os.Args[1]
	processedIDsPath := os.Args[2]
	outPath := defaultCurrentConversation
	if len(os.Args) >= 4 {
		outPath = os.Args[3]
	}

	rawItems, err := readConversations(conversationsPath)
	if err != nil {
		log.Fatalf("failed to read conversations: %v", err)
	}

	processed, err := readProcessedIDs(processedIDsPath)
	if err != nil {
		log.Fatalf("failed to read processed IDs: %v", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}

	// Find the first unprocessed conversation by order of appearance
	for _, raw := range rawItems {
		var hdr conversationHeader
		if err := json.Unmarshal(raw, &hdr); err != nil {
			// Skip items without readable id
			continue
		}
		if hdr.ID == "" {
			// Skip items without id
			continue
		}
		if _, seen := processed[hdr.ID]; seen {
			continue
		}
		// Write full object as JSON
		if err := os.WriteFile(outPath, append(raw, '\n'), 0o644); err != nil {
			log.Fatalf("failed to write current conversation: %v", err)
		}
		fmt.Printf("Selected conversation: %s\n", hdr.ID)
		return
	}

	// None remain: write sentinel
	if err := os.WriteFile(outPath, []byte(sentinelJSON), 0o644); err != nil {
		log.Fatalf("failed to write sentinel current conversation: %v", err)
	}
	fmt.Println("No unprocessed conversations remain. Wrote sentinel.")
}

func readConversations(path string) ([]json.RawMessage, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// Expect an array of objects
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return nil, fmt.Errorf("invalid conversations JSON: %w", err)
	}
	return arr, nil
}

func readProcessedIDs(path string) (map[string]struct{}, error) {
	result := make(map[string]struct{})

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Missing is fine; treat as empty
			return result, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if id == "" {
			continue
		}
		result[id] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}


