package gisty

// ReadFile retrieves the content of a file from a gist.
//
// This method is a short hand of `Gisty.Read(gistID).Files[file].Content`.
func (g *Gisty) ReadFile(gist string, file string) ([]byte, error) {
	gists, err := g.read(gist, g.AltFunctions.Read)
	if err != nil {
		return nil, WrapIfErr(err, "failed to read gist info")
	}

	gistFile, ok := gists.Files[file]
	if !ok {
		return nil, NewErr("file not found")
	}

	return []byte(gistFile.Content), nil
}
