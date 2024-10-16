package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyMap(t *testing.T) {
	var map1 map[string]string // nil map
	map2 := map[string]string{
		"go":   "lang",
		"lang": "go",
	}

	assert.Nil(t, CopyMap(map1))
	assert.Equal(t, map2, CopyMap(map2))
}

func TestMergeMap(t *testing.T) {
	// case 1: normal
	assert.Equal(t,
		map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"},
		MergeMap(map[string]string{"key1": "value1"}, map[string]string{"key2": "value2"}, map[string]string{"key3": "value3"}),
	)

	// case2: nil base
	assert.Equal(t,
		map[string]string{"key2": "value2", "key3": "value3"},
		MergeMap(nil, map[string]string{"key2": "value2"}, map[string]string{"key3": "value3"}),
	)
	assert.Equal(t, map[string]string{}, MergeMap[string, string](nil, nil))
	assert.Nil(t, MergeMap[string, int](nil))

	// case 3: nil overrides
	assert.Equal(t,
		map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"},
		MergeMap(map[string]string{"key1": "value1"}, map[string]string{"key2": "value2"}, nil, map[string]string{"key3": "value3"}),
	)

	// case 4: override values
	assert.Equal(t,
		map[string]string{"key1": "value3"},
		MergeMap(map[string]string{"key1": "value1"}, map[string]string{"key1": "value2"}, map[string]string{"key1": "value3"}),
	)
}

func TestMergeStrIFaceMaps(t *testing.T) {
	assert.Equal(t,
		map[string]any{},
		MergeStrIFaceMaps(nil, nil),
	)

	assert.Equal(t,
		map[string]any{"key1": 1, "key2": 2},
		MergeStrIFaceMaps(nil, map[string]any{"key1": 1, "key2": 2}),
	)

	assert.Equal(t,
		map[string]any{"key1": 1, "key2": 2},
		MergeStrIFaceMaps(map[string]any{"key1": 1, "key2": 2}, nil),
	)

	assert.Equal(t,
		map[string]any{"key1": 1, "key2": 2},
		MergeStrIFaceMaps(map[string]any{"key1": 3, "key2": 4}, map[string]any{"key1": 1, "key2": 2}),
	)

	assert.Equal(t,
		map[string]any{"key1": 1, "key2": 2},
		MergeStrIFaceMaps(map[string]any{"key1": 1, "key2": map[string]any{"key3": "123"}}, map[string]any{"key1": 1, "key2": 2}),
	)

	assert.Equal(t,
		map[string]any{"key1": 1, "key2": map[string]any{"key3": "456", "key4": 123, "key5": []int{1, 2, 3}}, "key3": 3},
		MergeStrIFaceMaps(
			map[string]any{"key1": 1, "key2": map[string]any{"key3": "123", "key5": []int{1, 2, 3}}},
			map[string]any{"key1": 1, "key2": map[string]any{"key3": "456", "key4": 123}, "key3": 3},
		),
	)
}
