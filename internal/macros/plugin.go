package macros

import (
	"fmt"
	"plugin"
)

type Plugin interface {
	RegisterMacros() map[string]Macro
}

func loadPlugin(path string) error {
	plug, err := plugin.Open(path)
	if err != nil {
		return err
	}

	symPlugin, err := plug.Lookup("Plugin")
	if err != nil {
		return err
	}

	pluginImpl, ok := symPlugin.(Plugin)
	if !ok {
		return fmt.Errorf("plugin does not implement Plugin interface")
	}

	for name, macro := range pluginImpl.RegisterMacros() {
		Macros[name] = macro
	}

	return nil
}
