package editor

// Saveable denotes a struct that has data that can be saved to a file.
type Saveable interface {
	// GenerateSaveData is called by the underlying interface (namely editor.Editor) to retrieve the data
	// the editor wants written to the file.
	GenerateSaveData() []byte
}
