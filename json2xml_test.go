package json2xml

import (
	"bytes"
	"testing"
)

// TestJSONToXML comment
func TestJSONToXML(t *testing.T) {
	body := []byte(`{"root":{"a":{"a1":"value-a1"},"b":"value-b","c":"value-c","d":200,"e":[{"e1":"value-e1"},{"e2":{"e21":"value-e21"}}],"f":false}}`)
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<root><a><a1>value-a1</a1></a><b>value-b</b><c>value-c</c><d>200</d><e><e1>value-e1</e1><e2><e21>value-e21</e21></e2></e><f>false</f></root>`
	buf := new(bytes.Buffer)
	converter, err := New(buf)
	if err != nil {
		t.Fatalf("cannot instantiate newConverter: %v", err)
	}
	err = converter.Convert(body)
	if err != nil {
		t.Fatalf("cannot comnvert xml to json: %v", err)
	}
	if expected != buf.String() {
		t.Fatalf("expected \n\t%v but got \n\t%v", expected, buf.String())
	}
}
