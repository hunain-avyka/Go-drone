// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package json

import (
	"regexp"
	"strings"

	"github.com/hunain-avyka/Go-drone/convert/harness"
)

func SanitizeForId(spanName string, spanId string) string {
	// Replace invalid characters with underscores
	invalidCharRegex := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	sanitized := invalidCharRegex.ReplaceAllString(spanName, "_")

	// Trim leading and trailing underscores
	sanitized = regexp.MustCompile("^_+|_+$").ReplaceAllString(sanitized, "")

	if sanitized == "" {
		sanitized = "unamed"
	}

	if len(sanitized) > 58 {
		sanitized = sanitized[:58]
	}

	return sanitized + spanId[:6]
}

func SanitizeForName(spanName string) string {
	// Replace invalid characters with underscores
	invalidCharRegex := regexp.MustCompile(`[^\w\- ]+`)
	sanitized := invalidCharRegex.ReplaceAllString(spanName, "_")

	// Replace repeating separators
	repeatingCharRegex := regexp.MustCompile(`([_ \-]){2,}`)
	sanitized = repeatingCharRegex.ReplaceAllString(sanitized, "$1")

	// Trim leading and trailing underscores
	sanitized = strings.TrimLeft(sanitized, "_ ")
	sanitized = strings.TrimRight(sanitized, "_ ")

	if sanitized == "" {
		sanitized = harness.DefaultName
	}

	if len(sanitized) > 128 {
		sanitized = sanitized[:128]
	}

	return sanitized
}
