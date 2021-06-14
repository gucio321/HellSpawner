package hsstate

// EditorState holds information about the state of an open editor
type EditorState struct {
	Path    []byte `json:"path"`
	Encoded []byte `json:"state"`
	WindowState
}
