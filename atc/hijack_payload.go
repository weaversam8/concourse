package atc

type HijackProcessSpec struct {
	Path string   `json:"path"`
	Args []string `json:"args"`
	Env  []string `json:"env"`
	Dir  string   `json:"dir"`

	Privileged bool   `json:"privileged"`
	User       string `json:"user"`

	TTY *HijackTTYSpec `json:"tty"`
}

type HijackTTYSpec struct {
	WindowSize HijackWindowSize `json:"window_size"`
}

type HijackWindowSize struct {
	Columns uint16 `json:"columns"`
	Rows    uint16 `json:"rows"`
}

type HijackInput struct {
	Closed  bool           `json:"closed,omitempty"`
	Stdin   []byte         `json:"stdin,omitempty"`
	TTYSpec *HijackTTYSpec `json:"tty,omitempty"`
}

type HijackOutput struct {
	Stdout             []byte `json:"stdout,omitempty"`
	Stderr             []byte `json:"stderr,omitempty"`
	Error              string `json:"error,omitempty"`
	ExitStatus         *int   `json:"exit_status,omitempty"`
	ExecutableNotFound bool   `json:"executable_not_found,omitempty"`
}
