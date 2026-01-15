package csvutils

import (
	"encoding/csv"
	"io"
	"log"
)

type RowSet map[string]struct{}

func RemoveDuplicates(input io.Reader, output io.Writer, keyColumn int) error {
	reader := csv.NewReader(input)
	writer := csv.NewWriter(output)
	defer writer.Flush()

	seen := make(RowSet)
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Warning: skipping malformed row: %v", err)
			continue
		}
		if len(record) <= keyColumn {
			continue
		}
		key := record[keyColumn]
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}