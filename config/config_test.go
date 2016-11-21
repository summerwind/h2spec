package config

import "testing"

func TestRunMode(t *testing.T) {
	tests := []struct {
		sections []string
		target   string
		mode     int
	}{
		{sections: []string{"http2"}, target: "http2", mode: RunModeAll},
		{sections: []string{"http2"}, target: "http2/5", mode: RunModeAll},
		{sections: []string{"http2"}, target: "http2/5.1", mode: RunModeAll},
		{sections: []string{"http2"}, target: "http2/5.1.2", mode: RunModeAll},
		{sections: []string{"http2"}, target: "http2/5.1.2/1", mode: RunModeAll},
		{sections: []string{"http2"}, target: "hpack", mode: RunModeNone},

		{sections: []string{"http2/5"}, target: "http2", mode: RunModeGroup},
		{sections: []string{"http2/5"}, target: "http2/5", mode: RunModeAll},
		{sections: []string{"http2/5"}, target: "http2/5.1", mode: RunModeAll},
		{sections: []string{"http2/5"}, target: "http2/5.1.2", mode: RunModeAll},
		{sections: []string{"http2/5"}, target: "http2/5.1.2/1", mode: RunModeAll},
		{sections: []string{"http2/5"}, target: "hpack", mode: RunModeNone},
		{sections: []string{"http2/5"}, target: "http2/4", mode: RunModeNone},

		{sections: []string{"http2/5.1"}, target: "http2", mode: RunModeGroup},
		{sections: []string{"http2/5.1"}, target: "http2/5", mode: RunModeGroup},
		{sections: []string{"http2/5.1"}, target: "http2/5.1", mode: RunModeAll},
		{sections: []string{"http2/5.1"}, target: "http2/5.1.2", mode: RunModeAll},
		{sections: []string{"http2/5.1"}, target: "http2/5.1.2/1", mode: RunModeAll},
		{sections: []string{"http2/5.1"}, target: "hpack", mode: RunModeNone},
		{sections: []string{"http2/5.1"}, target: "http2/4", mode: RunModeNone},
		{sections: []string{"http2/5.1"}, target: "http2/5.2", mode: RunModeNone},

		{sections: []string{"http2/5.1.2"}, target: "http2", mode: RunModeGroup},
		{sections: []string{"http2/5.1.2"}, target: "http2/5", mode: RunModeGroup},
		{sections: []string{"http2/5.1.2"}, target: "http2/5.1", mode: RunModeGroup},
		{sections: []string{"http2/5.1.2"}, target: "http2/5.1.2", mode: RunModeAll},
		{sections: []string{"http2/5.1.2"}, target: "http2/5.1.2/1", mode: RunModeAll},
		{sections: []string{"http2/5.1.2"}, target: "hpack", mode: RunModeNone},
		{sections: []string{"http2/5.1.2"}, target: "http2/4", mode: RunModeNone},
		{sections: []string{"http2/5.1.2"}, target: "http2/5.2", mode: RunModeNone},
		{sections: []string{"http2/5.1.2"}, target: "http2/5.1.3", mode: RunModeNone},

		{sections: []string{"http2/5.1.2/1"}, target: "http2", mode: RunModeGroup},
		{sections: []string{"http2/5.1.2/1"}, target: "http2/5", mode: RunModeGroup},
		{sections: []string{"http2/5.1.2/1"}, target: "http2/5.1", mode: RunModeGroup},
		{sections: []string{"http2/5.1.2/1"}, target: "http2/5.1.2", mode: RunModeGroup},
		{sections: []string{"http2/5.1.2/1"}, target: "http2/5.1.2/1", mode: RunModeAll},
		{sections: []string{"http2/5.1.2/1"}, target: "hpack", mode: RunModeNone},
		{sections: []string{"http2/5.1.2/1"}, target: "http2/4", mode: RunModeNone},
		{sections: []string{"http2/5.1.2/1"}, target: "http2/5.2", mode: RunModeNone},
		{sections: []string{"http2/5.1.2/1"}, target: "http2/5.1.3", mode: RunModeNone},
		{sections: []string{"http2/5.1.2/1"}, target: "http2/5.1.2/2", mode: RunModeNone},

		{sections: []string{"http2", "http2/5.1.2/1"}, target: "http2", mode: RunModeAll},
		{sections: []string{"http2", "http2/5.1.2/1"}, target: "http2/5", mode: RunModeAll},
		{sections: []string{"http2", "http2/5.1.2/1"}, target: "http2/5.1", mode: RunModeAll},
		{sections: []string{"http2", "http2/5.1.2/1"}, target: "http2/5.1.2", mode: RunModeAll},
		{sections: []string{"http2", "http2/5.1.2/1"}, target: "http2/5.1.2/1", mode: RunModeAll},
		{sections: []string{"http2", "http2/5.1.2/1"}, target: "hpack", mode: RunModeNone},
		{sections: []string{"http2", "http2/5.1.2/1"}, target: "http2/4", mode: RunModeAll},
		{sections: []string{"http2", "http2/5.1.2/1"}, target: "http2/5.2", mode: RunModeAll},
		{sections: []string{"http2", "http2/5.1.2/1"}, target: "http2/5.1.3", mode: RunModeAll},
		{sections: []string{"http2", "http2/5.1.2/1"}, target: "http2/5.1.2/2", mode: RunModeAll},

		{sections: []string{"http2/5.1.2/1", "http2"}, target: "http2", mode: RunModeAll},
		{sections: []string{"http2/5.1.2/1", "http2"}, target: "http2/5", mode: RunModeAll},
		{sections: []string{"http2/5.1.2/1", "http2"}, target: "http2/5.1", mode: RunModeAll},
		{sections: []string{"http2/5.1.2/1", "http2"}, target: "http2/5.1.2", mode: RunModeAll},
		{sections: []string{"http2/5.1.2/1", "http2"}, target: "http2/5.1.2/1", mode: RunModeAll},
		{sections: []string{"http2/5.1.2/1", "http2"}, target: "hpack", mode: RunModeNone},
		{sections: []string{"http2/5.1.2/1", "http2"}, target: "http2/4", mode: RunModeAll},
		{sections: []string{"http2/5.1.2/1", "http2"}, target: "http2/5.2", mode: RunModeAll},
		{sections: []string{"http2/5.1.2/1", "http2"}, target: "http2/5.1.3", mode: RunModeAll},
		{sections: []string{"http2/5.1.2/1", "http2"}, target: "http2/5.1.2/2", mode: RunModeAll},
	}

	for i, tt := range tests {
		c := Config{
			Sections: tt.sections,
		}

		mode := c.RunMode(tt.target)
		if tt.mode != mode {
			t.Errorf("#%d mode - expect: %d, got: %d (%v / %v)", i, tt.mode, mode, tt.target, c.targetMap)
		}
	}
}
