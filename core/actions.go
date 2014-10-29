package core

import (
	"fmt"
	"github.com/sclevine/agouti/core/internal/types"
)

type actionsFunc func(types.Element) error

func (s *selection) forEachElement(actions actionsFunc) error {
	elements, err := s.getSelectedElements()
	if err != nil {
		return fmt.Errorf("failed to select '%s': %s", s, err)
	}

	for _, element := range elements {
		if err := actions(element); err != nil {
			return err
		}
	}
	return nil
}

func (s *selection) Click() error {
	return s.forEachElement(func(element types.Element) error {
		if err := element.Click(); err != nil {
			return fmt.Errorf("failed to click on '%s': %s", s, err)
		}
		return nil
	})
}

func (s *selection) DoubleClick() error {
	return s.forEachElement(func(element types.Element) error {
		if err := s.client.MoveTo(element, nil); err != nil {
			return fmt.Errorf("failed to move mouse to '%s': %s", s, err)
		}
		if err := s.client.DoubleClick(); err != nil {
			return fmt.Errorf("failed to double-click on '%s': %s", s, err)
		}
		return nil
	})
}

func (s *selection) Fill(text string) error {
	return s.forEachElement(func(element types.Element) error {
		if err := element.Clear(); err != nil {
			return fmt.Errorf("failed to clear '%s': %s", s, err)
		}
		if err := element.Value(text); err != nil {
			return fmt.Errorf("failed to enter text into '%s': %s", s, err)
		}
		return nil
	})
}

func (s *selection) Check() error {
	return s.setChecked(true)
}

func (s *selection) Uncheck() error {
	return s.setChecked(false)
}

func (s *selection) setChecked(checked bool) error {
	return s.forEachElement(func(element types.Element) error {
		elementType, err := element.GetAttribute("type")
		if err != nil {
			return fmt.Errorf("failed to retrieve type of '%s': %s", s, err)
		}

		if elementType != "checkbox" {
			return fmt.Errorf("'%s' does not refer to a checkbox", s)
		}

		selected, err := element.IsSelected()
		if err != nil {
			return fmt.Errorf("failed to retrieve state of '%s': %s", s, err)
		}

		if selected != checked {
			if err := element.Click(); err != nil {
				return fmt.Errorf("failed to click on '%s': %s", s, err)
			}
		}
		return nil
	})
}

func (s *selection) Select(text string) error {
	return s.forEachElement(func(element types.Element) error {
		optionXPath := fmt.Sprintf(`./option[normalize-space(text())="%s"]`, text)
		optionToSelect := types.Selector{Using: "xpath", Value: optionXPath}
		options, err := element.GetElements(optionToSelect)
		if err != nil {
			return fmt.Errorf("failed to select specified option for some '%s': %s", s, err)
		}

		if len(options) == 0 {
			return fmt.Errorf(`no options with text "%s" found for some '%s'`, text, s)
		}

		for _, option := range options {
			if err := option.Click(); err != nil {
				return fmt.Errorf(`failed to click on option with text "%s" for some '%s': %s`, text, s, err)
			}
		}
		return nil
	})
}

func (s *selection) Submit() error {
	return s.forEachElement(func(element types.Element) error {
		if err := element.Submit(); err != nil {
			return fmt.Errorf("failed to submit '%s': %s", s, err)
		}
		return nil
	})
}