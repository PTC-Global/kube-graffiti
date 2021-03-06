/*
Copyright (C) 2018 Expedia Group.
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

package graffiti

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyObject(t *testing.T) {
	_, err := makeFieldMapFromRawObject([]byte{})
	require.Error(t, err)
	assert.Equal(t, "no fields found", err.Error())

}

func TestTopLevelObjectMustBeAMap(t *testing.T) {
	validJSON := `[ "apple", "orange", "banana" ]`
	_, err := makeFieldMapFromRawObject([]byte(validJSON))
	assert.Error(t, err)
	assert.Equal(t, "failed to unmarshal object: json: cannot unmarshal array into Go value of type map[string]interface {}", err.Error())
}

func TestBaseTypesAsStrings(t *testing.T) {
	// when creating a fieldmap the following json types are converted to strings

	// strings
	testJSON := `{ "test": "dave" }`
	fm, err := makeFieldMapFromRawObject([]byte(testJSON))
	require.NoError(t, err)
	assert.Equal(t, "dave", fm["test"])

	// ints
	testJSON = `{ "test": 100 }`
	fm, err = makeFieldMapFromRawObject([]byte(testJSON))
	require.NoError(t, err)
	assert.Equal(t, "100", fm["test"])

	// floats
	testJSON = `{ "test": 63.333392 }`
	fm, err = makeFieldMapFromRawObject([]byte(testJSON))
	require.NoError(t, err)
	assert.Equal(t, "63.333392", fm["test"])

	// bools
	testJSON = `{ "test": true }`
	fm, err = makeFieldMapFromRawObject([]byte(testJSON))
	require.NoError(t, err)
	assert.Equal(t, "true", fm["test"])
}

func TestSlicesAreReferencedByIndex(t *testing.T) {
	testJSON := `{ "test": [ "dave", 100, 63.49, true ] }`
	fm, err := makeFieldMapFromRawObject([]byte(testJSON))
	require.NoError(t, err)

	assert.Equal(t, "dave", fm["test.0"])
	assert.Equal(t, "100", fm["test.1"])
	assert.Equal(t, "63.49", fm["test.2"])
	assert.Equal(t, "true", fm["test.3"])
}

func TestMapsAreReferencedByKey(t *testing.T) {
	testJSON := `{ "test": { "band": "Queen", "singer": "Freddie Mercury", "status": "legend" }}`
	fm, err := makeFieldMapFromRawObject([]byte(testJSON))
	require.NoError(t, err)

	assert.Equal(t, "Queen", fm["test.band"])
	assert.Equal(t, "Freddie Mercury", fm["test.singer"])
	assert.Equal(t, "legend", fm["test.status"])
}

func TestComplexObject(t *testing.T) {
	var testJSON = `{
		"metadata":{
			"name":"test-namespace",
			"creationTimestamp":null,
			"labels":{
				"author": "david",
				"group": "runtime"
			},
			"annotations":{
				"level": "v.special",
				"prometheus.io/path": "/metrics"
			}
		},
		"spec":{},
		"status":{
			"phase":"Active"
		}
	 }`

	fm, err := makeFieldMapFromRawObject([]byte(testJSON))
	require.NoError(t, err)

	assert.Equal(t, "test-namespace", fm["metadata.name"])
	assert.Equal(t, "david", fm["metadata.labels.author"])
	assert.Equal(t, "v.special", fm["metadata.annotations.level"])
	assert.Equal(t, "Active", fm["status.phase"])
}
