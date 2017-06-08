// Copyright 2017 Santhosh Kumar Tekuri. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dom_test

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/dom"
)

func TestIdentity(t *testing.T) {
	tests := []string{
		`<test a1="v1" a2="v2"/>`,
		`<x:test xmlns:x="ns1" x:a1="v1" a2="v2"/>`,
		`<test xmlns="ns1" a1="v1" a2="v2"/>`,
		`<x:one xmlns:x="ns1"><x:two/></x:one>`,
		`<x><!--ignore me--></x>`,
		`<x><?abcd hello world?></x>`,
	}
	for i, test := range tests {
		d, err := dom.Unmarshal(xml.NewDecoder(strings.NewReader(test)))
		if err != nil {
			t.Errorf("#%d: %s", i, err)
			continue
		}
		buf := new(bytes.Buffer)
		if err := dom.Marshal(d, buf); err != nil {
			t.Errorf("#%d: %s", i, err)
		}
		if s := buf.String(); s != test {
			t.Errorf("expected:\n%s\nbut got:\n%s\n", test, s)
		}
	}
}

func TestNormalized(t *testing.T) {
	tests := []struct {
		raw, normalized string
	}{
		{`<?xml version="1.0" encoding="UTF-8" standalone="no" ?><test/>`, `<test/>`},
		{` <a/>`, `<a/>`},
		{`<e a='v'/>`, `<e a="v"/>`},
		{`<e a='v"'/>`, `<e a="v&quot;"/>`},
		{`<a>one<![CDATA[two]]>three<![CDATA[four]]>five</a>`, `<a>onetwothreefourfive</a>`},
	}
	for i, test := range tests {
		d, err := dom.Unmarshal(xml.NewDecoder(strings.NewReader(test.raw)))
		if err != nil {
			t.Errorf("#%d: %s", i, err)
			continue
		}
		buf := new(bytes.Buffer)
		if err := dom.Marshal(d, buf); err != nil {
			t.Errorf("#%d: %s", i, err)
		}
		if s := buf.String(); s != test.normalized {
			t.Errorf("expected:\n%s\nbut got:\n%s\n", test.normalized, s)
		}
	}
}

func TestInvalidXML(t *testing.T) {
	tests := []string{
		``,                  // no root element
		`<e1`,               // incomplete start element
		`<e1>`,              // missing end element
		`<e1/><e2/>`,        // more than one root element
		`<ns1:e1/>`,         // unresolved element prefix
		`<e1 ns1:p1="v1"/>`, // unresolved attribute prefix
		`<e1>hai</e2>`,      // wrong end element
		`hai<e1/>`,          // text outside root element
	}

	for i, test := range tests {
		if _, err := dom.Unmarshal(xml.NewDecoder(strings.NewReader(test))); err == nil {
			t.Errorf("#%d: error expected", i)
			continue
		} else {
			t.Logf("#%d: %v", i, err)
		}
	}
}
