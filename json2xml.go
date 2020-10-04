package json2xml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
)

// JSON2XML comment
type JSON2XML struct {
	encoder *xml.Encoder
}

// New comment
func New(w io.Writer) (*JSON2XML, error) {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return nil, fmt.Errorf("cannot write xml header: %w", err)
	}
	return &JSON2XML{
		encoder: xml.NewEncoder(w),
	}, nil
}

// encodeXML comment
func (c *JSON2XML) encodeXML(root interface{}, key string) error {
	tagName := xml.Name{Space: "", Local: key}
	switch val := root.(type) {
	// TODO nil is not handled
	case string, float64, bool:
		if err := c.element(val, tagName); err != nil {
			return err
		}
	case map[string]interface{}:
		if err := c.start(tagName); err != nil {
			return err
		}
		sortedEntries := sortMap(val)
		for _, k := range sortedEntries {
			if err := c.encodeXML(val[k], k); err != nil {
				return err
			}
		}
		if err := c.end(tagName); err != nil {
			return err
		}
	case []interface{}:
		if err := c.start(tagName); err != nil {
			return err
		}
		for _, element := range val {
			if err := c.encodeXML(element, ""); err != nil {
				return err
			}
		}
		if err := c.end(tagName); err != nil {
			return err
		}
	}
	return nil
}

// element comment
func (c *JSON2XML) element(val interface{}, tagName xml.Name) error {
	if err := c.encoder.EncodeElement(val, xml.StartElement{Name: tagName}); err != nil {
		return fmt.Errorf("cannot encode XML Element: %w", err)
	}
	return nil
}

// start comment
func (c *JSON2XML) start(tagName xml.Name) error {
	if err := c.encoder.EncodeToken(xml.StartElement{Name: tagName}); err != nil {
		return fmt.Errorf("cannot create start start: %w", err)
	}
	return nil
}

// end comment
func (c *JSON2XML) end(tagName xml.Name) error {
	if err := c.encoder.EncodeToken(xml.EndElement{Name: tagName}); err != nil {
		return fmt.Errorf("cannot create end tag: %w", err)
	}
	return nil
}

// Convert comment
func (c *JSON2XML) Convert(content []byte) error {
	var decodeJSON map[string]interface{}
	err := json.Unmarshal(content, &decodeJSON)
	if err != nil {
		return fmt.Errorf("cannot unmarshal from JSON: %w", err)
	}
	if err := c.encodeXML(decodeJSON, ""); err != nil {
		return fmt.Errorf("cannot encode XML: %w", err)
	}
	if err := c.encoder.Flush(); err != nil {
		return fmt.Errorf("cannot flush XML: %w", err)
	}
	return nil
}

// sortMap comment
func sortMap(val map[string]interface{}) []string {
	keys := make([]string, len(val))
	i := 0
	for k := range val {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
