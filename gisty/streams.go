package gisty

import "github.com/KEINOS/go-gisty/internal/ghcmd"

func (g *Gisty) streams() ghcmd.Streams {
	return ghcmd.Streams{
		Stdin:  g.Stdin,
		Stdout: g.Stdout,
		Stderr: g.Stderr,
	}
}
