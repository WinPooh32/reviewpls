package prompt

import (
	"errors"
	"fmt"
)

type PromptSet map[string]Prompt

func (ps PromptSet) Require(name string, names ...string) error {
	var errs []error

	if err := ps.checkLoaded(name); err != nil {
		errs = append(errs, err)
	}

	for _, name := range names {
		if err := ps.checkLoaded(name); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (ps PromptSet) checkLoaded(name string) error {
	_, ok := ps[name]
	if !ok {
		return fmt.Errorf("%q is not loaded", name)
	}

	return nil
}
