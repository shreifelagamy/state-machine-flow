package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"os/exec"
)

// Status represents a state in the flow and its potential next statuses.
// It is used to define the nodes and edges of the status flow graph.
type Status struct {
	Name       string      // The name of the current status
	NextStatus []string    // A list of statuses that can be transitioned to from the current status
}

// GenerateDOT creates a DOT representation of the statuses and their flows.
// It takes a slice of Status objects and returns a string in DOT graph format.
// The generated graph is left-to-right oriented with rounded blue boxes for nodes.
func GenerateDOT(statuses []Status) string {
	dot := "digraph G {\n"  // Start of the DOT graph definition
	dot += "rankdir=LR;\n"  // Set the direction of the graph to left-to-right
	dot += "node [shape=box, style=rounded, color=blue, fontname=Helvetica];\n"  // Set node style to be more creative and visually appealing
	visited := make(map[string]struct{})  // A map to keep track of which statuses have already been visited using an empty struct for efficiency

	// Loop over each status and create the flow relationships.
	for _, status := range statuses {
		// Add the current status node if it hasn't been added yet
		if _, ok := visited[status.Name]; !ok {
			dot += fmt.Sprintf("  \"%s\";\n", status.Name)  // Add the status as a node in the graph
			visited[status.Name] = struct{}{}
		}
		// Loop over each of the next possible statuses
		for _, next := range status.NextStatus {
			// Add the next status node if it hasn't been added yet
			if _, ok := visited[next]; !ok {
				dot += fmt.Sprintf("  \"%s\";\n", next)  // Add the next status as a node in the graph
				visited[next] = struct{}{}
			}
			// Create an edge from the current status to the next status
			dot += fmt.Sprintf("  \"%s\" -> \"%s\";\n", status.Name, next)
		}
	}
	dot += "}\n"  // End of the DOT graph definition
	return dot
}

// WriteDOTToFile writes the DOT content to a file.
func WriteDOTToFile(dot string, filename string) error {
	// Create a new file with the given filename
	file, err := os.Create(filename)
	if err != nil {
		return err  // Return an error if the file could not be created
	}
	defer file.Close()  // Ensure the file is closed when the function exits

	// Write the DOT content to the file
	_, err = file.WriteString(dot)
	if err != nil {
		return err  // Return an error if writing to the file fails
	}
	return nil
}

// GenerateGraphImage generates an image using Graphviz.
func GenerateGraphImage(dotFile string, outputFile string) error {
	// Use the "dot" command from Graphviz to generate a PNG image from the DOT file
	cmd := exec.Command("dot", "-Tpng", dotFile, "-o", outputFile)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error generating graph image: %v", err)  // Return an error if the command fails
	}
	return nil
}

func main() {
	// Define command line flags
	outputPath := flag.String("path", ".", "Output directory path")
	outputName := flag.String("name", "status_flow", "Base name for output files (without extension)")
	useStaticData := flag.Bool("static", false, "Use static data for generating the image")
	flag.Parse()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputPath, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	var statuses []Status

	if *useStaticData {
		// Use predefined static data
		statuses = []Status{
			{Name: "Start", NextStatus: []string{"In Progress"}},
			{Name: "In Progress", NextStatus: []string{"Completed", "Failed"}},
			{Name: "Completed", NextStatus: []string{}},
			{Name: "Failed", NextStatus: []string{}},
		}
	} else {
		// Read input JSON from stdin
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}

		// Parse the input JSON into an array of statuses
		err = json.Unmarshal(input, &statuses)
		if err != nil {
			fmt.Printf("Error parsing input JSON: %v\n", err)
			return
		}
	}

	// Generate the DOT representation
	dotContent := GenerateDOT(statuses)

	// Create full file paths
	dotFilename := filepath.Join(*outputPath, *outputName+".dot")
	pngFilename := filepath.Join(*outputPath, *outputName+".png")

	// Write DOT content to a file
	if err := WriteDOTToFile(dotContent, dotFilename); err != nil {
		fmt.Printf("Error writing DOT file: %v\n", err)
		return
	}

	// Generate the status flow image using Graphviz
	if err := GenerateGraphImage(dotFilename, pngFilename); err != nil {
		fmt.Printf("Error generating graph image: %v\n", err)
		return
	}

	// Inform the user that the image was successfully generated
	fmt.Printf("Status flow image generated: %s\n", pngFilename)
}
