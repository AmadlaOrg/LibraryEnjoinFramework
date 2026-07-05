package display

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Display interface {
	Display() error
	SetTableHeaders(tableHeaders []string) Display
}

type displayImpl struct {
	cmd          *cobra.Command
	tableHeaders []string
	data         map[string]any
	writer       io.Writer
}

var (
	jsonMarshalIndent = json.MarshalIndent
	yamlMarshal       = yaml.Marshal
)

// attachFlags adds the -o/--output flag (table|json|yaml).
func (s *displayImpl) attachFlags() {
	s.cmd.Flags().StringP("output", "o", "table", "Output format: table, json, or yaml")
}

// Display renders the data in the requested format.
func (s *displayImpl) Display() error {
	format, _ := s.cmd.Flags().GetString("output")

	switch format {
	case "json":
		return s.json()
	case "yaml":
		return s.yaml()
	case "table", "":
		s.table()
		return nil
	default:
		return fmt.Errorf("unknown output format %q (use table, json, or yaml)", format)
	}
}

// SetTableHeaders sets custom table headers.
func (s *displayImpl) SetTableHeaders(tableHeaders []string) Display {
	s.tableHeaders = tableHeaders
	return s
}

// getTableHeaders returns custom or default headers.
func (s *displayImpl) getTableHeaders() []string {
	if len(s.tableHeaders) == 0 {
		return []string{"Category", "Key", "Value"}
	}
	return s.tableHeaders
}

// table renders data as a key-value table.
func (s *displayImpl) table() {
	table := tablewriter.NewTable(s.writer)
	headers := make([]any, len(s.getTableHeaders()))
	for i, h := range s.getTableHeaders() {
		headers[i] = h
	}
	table.Header(headers...)

	for category, data := range s.data {
		switch v := data.(type) {
		case map[string]string:
			for key, value := range v {
				table.Append(category, key, value)
			}
		case map[string]any:
			for key, subData := range v {
				if subMap, ok := subData.(map[string]string); ok {
					for subKey, subValue := range subMap {
						table.Append(
							category,
							fmt.Sprintf("%s.%s", key, subKey),
							subValue,
						)
					}
				}
			}
		}
	}

	table.Render()
}

// json renders the data in JSON format.
func (s *displayImpl) json() error {
	jsonData, err := jsonMarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}
	_, err = fmt.Fprintln(s.writer, string(jsonData))
	return err
}

// yaml renders the data in YAML format.
func (s *displayImpl) yaml() error {
	yamlData, err := yamlMarshal(s.data)
	if err != nil {
		return fmt.Errorf("error encoding YAML: %w", err)
	}
	_, err = fmt.Fprint(s.writer, string(yamlData))
	return err
}
