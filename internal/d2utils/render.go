package d2utils

import (
	"context"
	"os"
	"path/filepath"

	"oss.terrastruct.com/d2/d2format"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/lib/textmeasure"
)

func RenderSVG(g *d2graph.Graph) error {
	// Turn the graph into a script
	script := d2format.Format(g.AST)

	// Initialize a ruler to measure font glyphs
	ruler, _ := textmeasure.NewRuler()

	// Compile the script into a diagram
	ctx := context.Background()
	diagram, _, _ := d2lib.Compile(ctx, script, &d2lib.CompileOptions{
		Layout: d2dagrelayout.DefaultLayout,
		Ruler:  ruler,
	})

	// Render to SVG
	diagramImage, _ := d2svg.Render(diagram, &d2svg.RenderOpts{
		Pad: d2svg.DEFAULT_PADDING,
	})

	// Write to disk the script and the SVG image
	_ = os.WriteFile(filepath.Join("out.svg"), diagramImage, 0600)
	_ = os.WriteFile(filepath.Join("out.d2"), []byte(script), 0600)

	return nil
}
