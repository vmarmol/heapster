/*
Copyright 2014 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"regexp"
	"strings"
)

// IsDNSLabel tests for a string that conforms to the definition of a label in
// DNS (RFC 1123).
func IsDNSLabel(value string) bool {
	return IsDNS1123Label(value)
}

// IsDNSSubdomain tests for a string that conforms to the definition of a
// subdomain in DNS (RFC 1123).
func IsDNSSubdomain(value string) bool {
	return IsDNS1123Subdomain(value)
}

const dns1123LabelFmt string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"

var dns1123LabelRegexp = regexp.MustCompile("^" + dns1123LabelFmt + "$")

const dns1123LabelMaxLength int = 63

// IsDNS1123Label tests for a string that conforms to the definition of a label in
// DNS (RFC 1123).
func IsDNS1123Label(value string) bool {
	return len(value) <= dns1123LabelMaxLength && dns1123LabelRegexp.MatchString(value)
}

const dns1123SubdomainFmt string = dns1123LabelFmt + "(\\." + dns1123LabelFmt + ")*"

var dns1123SubdomainRegexp = regexp.MustCompile("^" + dns1123SubdomainFmt + "$")

const dns1123SubdomainMaxLength int = 253

// IsDNS1123Subdomain tests for a string that conforms to the definition of a
// subdomain in DNS (RFC 1123).
func IsDNS1123Subdomain(value string) bool {
	return len(value) <= dns1123SubdomainMaxLength && dns1123SubdomainRegexp.MatchString(value)
}

const dns952LabelFmt string = "[a-z]([-a-z0-9]*[a-z0-9])?"

var dns952LabelRegexp = regexp.MustCompile("^" + dns952LabelFmt + "$")

const dns952LabelMaxLength int = 24

// IsDNS952Label tests for a string that conforms to the definition of a label in
// DNS (RFC 952).
func IsDNS952Label(value string) bool {
	return len(value) <= dns952LabelMaxLength && dns952LabelRegexp.MatchString(value)
}

const cIdentifierFmt string = "[A-Za-z_][A-Za-z0-9_]*"

var cIdentifierRegexp = regexp.MustCompile("^" + cIdentifierFmt + "$")

// IsCIdentifier tests for a string that conforms the definition of an identifier
// in C. This checks the format, but not the length.
func IsCIdentifier(value string) bool {
	return cIdentifierRegexp.MatchString(value)
}

// IsValidPortNum tests that the argument is a valid, non-zero port number.
func IsValidPortNum(port int) bool {
	return 0 < port && port < 65536
}

// IsQualifiedName tests whether a string fits the "optionally-namespaced
// name" pattern: [ DNS_SUBDOMAIN "/" ] DNS_LABEL
func IsQualifiedName(value string) bool {
	var n, ns string
	parts := strings.Split(value, "/")
	switch len(parts) {
	case 1:
		n = parts[0]
	case 2:
		ns = parts[0]
		n = parts[1]
	default:
		return false
	}
	if (ns != "" && !IsDNSSubdomain(ns)) || !IsDNSLabel(n) {
		return false
	}
	return true
}
