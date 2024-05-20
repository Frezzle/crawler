package graph

import (
	"fmt"
	"os"

	"github.com/frezzle/web-crawler/utils"
)

// Outputs many URL->URL "connections" to a mermaid flowchart file,
// which can be rendered to an image e.g. using mermaid.js, mermaid CLI, mermaid.live.
func SaveMermaidFlowchart(connections [][2]string, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("unable to open or create file for mermaid graph: %w", err)
	}
	_, err = file.WriteString("flowchart LR\n")
	if err != nil {
		return fmt.Errorf("unable to write header to mermaid file: %w", err)
	}
	for _, conn := range connections {
		_, err = file.WriteString(fmt.Sprintf(
			"%s[%s]-->%s[%s]\n",
			utils.Truncate(utils.Hash(conn[0]), 10), // truncate to stay below char limit on mermaid.live
			utils.TruncateUrl(conn[0], 150),
			utils.Truncate(utils.Hash(conn[1]), 10), // truncate to stay below char limit on mermaid.live
			utils.TruncateUrl(conn[1], 150),
			// TODO: don't truncate now that i know how to increase char limit on mermaid.live / am using mermaid cli ?
		))
		if err != nil {
			return fmt.Errorf("unable to write url connection to mermaid file: %w", err)
		}
	}

	return nil
}
