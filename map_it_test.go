/*
 * Copyright (c) 2008-2021, Hazelcast, Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hazelcast_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	hz "github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/aggregate"
	"github.com/hazelcast/hazelcast-go-client/hzerrors"
	"github.com/hazelcast/hazelcast-go-client/internal/it"
	"github.com/hazelcast/hazelcast-go-client/predicate"
	"github.com/hazelcast/hazelcast-go-client/serialization"
	"github.com/hazelcast/hazelcast-go-client/types"
)

func TestMap_Put(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if _, err := m.Put(context.Background(), "key", targetValue); err != nil {
			t.Fatal(err)
		}
		if value, err := m.Get(context.Background(), "key"); err != nil {
			t.Fatal(err)
		} else if targetValue != value {
			t.Fatalf("target %v != %v", targetValue, value)
		}
	})
}

func TestMap_PutWithTTL(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if _, err := m.PutWithTTL(context.Background(), "key", targetValue, 20*time.Second); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, targetValue, it.MustValue(m.Get(context.Background(), "key")))
		it.Eventually(t, func() bool {
			return it.MustValue(m.Get(context.Background(), "key")) == nil
		})
	})
}

func TestMap_PutWithMaxIdle(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if _, err := m.PutWithMaxIdle(context.Background(), "key", targetValue, 1*time.Second); err != nil {
			t.Fatal(err)
		}
		it.Eventually(t, func() bool {
			return it.MustValue(m.Get(context.Background(), "key")) == nil
		})
	})
}

func TestMap_PutWithTTLAndMaxIdle(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if _, err := m.PutWithTTLAndMaxIdle(context.Background(), "key", targetValue, 1*time.Second, 1*time.Second); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, targetValue, it.MustValue(m.Get(context.Background(), "key")))
		it.Eventually(t, func() bool {
			return it.MustValue(m.Get(context.Background(), "key")) == nil
		})
	})
}

func TestMap_PutIfAbsent(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if _, err := m.PutIfAbsent(context.Background(), "key", targetValue); err != nil {
			t.Fatal(err)
		}
		it.AssertEquals(t, targetValue, it.MustValue(m.Get(context.Background(), "key")))
		if _, err := m.PutIfAbsent(context.Background(), "key", "another-value"); err != nil {
			t.Fatal(err)
		}
		it.AssertEquals(t, targetValue, it.MustValue(m.Get(context.Background(), "key")))
	})
}

func TestMap_PutIfAbsentWithTTL(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if _, err := m.PutIfAbsentWithTTL(context.Background(), "key", targetValue, 1*time.Second); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, targetValue, it.MustValue(m.Get(context.Background(), "key")))
		it.Eventually(t, func() bool {
			return it.MustValue(m.Get(context.Background(), "key")) == nil
		})
	})
}

func TestMap_PutIfAbsentWithTTLAndMaxIdle(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if _, err := m.PutIfAbsentWithTTLAndMaxIdle(context.Background(), "key", targetValue, 1*time.Second, 1*time.Second); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, targetValue, it.MustValue(m.Get(context.Background(), "key")))
		it.Eventually(t, func() bool {
			return it.MustValue(m.Get(context.Background(), "key")) == nil
		})
	})
}

func TestMap_PutTransient(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if err := m.PutTransient(context.Background(), "key", targetValue); err != nil {
			t.Fatal(err)
		}
		it.AssertEquals(t, targetValue, it.MustValue(m.Get(context.Background(), "key")))
	})
}

func TestMap_PutTransientWithTTL(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if err := m.PutTransientWithTTL(context.Background(), "key", targetValue, 1*time.Second); err != nil {
			t.Fatal(err)
		}
		it.AssertEquals(t, targetValue, it.MustValue(m.Get(context.Background(), "key")))
		it.Eventually(t, func() bool {
			return it.MustValue(m.Get(context.Background(), "key")) == nil
		})
	})
}

func TestMap_PutTransientWithMaxIdle(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if err := m.PutTransientWithMaxIdle(context.Background(), "key", targetValue, 1*time.Second); err != nil {
			t.Fatal(err)
		}
		it.AssertEquals(t, targetValue, it.MustValue(m.Get(context.Background(), "key")))
		it.Eventually(t, func() bool {
			return it.MustValue(m.Get(context.Background(), "key")) == nil
		})

	})
}

func TestMap_PutTransientWithTTLAndMaxIdle(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		// TODO: better test
		if err := m.PutTransientWithTTLAndMaxIdle(context.Background(), "key", targetValue, 1*time.Second, 1*time.Second); err != nil {
			t.Fatal(err)
		}
		it.AssertEquals(t, targetValue, it.MustValue(m.Get(context.Background(), "key")))
		it.Eventually(t, func() bool {
			return it.MustValue(m.Get(context.Background(), "key")) == nil
		})
	})
}

func TestMap_Set(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if err := m.Set(context.Background(), "key", targetValue); err != nil {
			t.Fatal(err)
		}
		value := it.MustValue(m.Get(context.Background(), "key"))
		if targetValue != value {
			t.Fatalf("target %v != %v", targetValue, value)
		}
	})
}

func TestMap_SetWithTTL(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		targetValue := "value"
		if err := m.SetWithTTL(ctx, "key", targetValue, 1*time.Second); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, targetValue, it.MustValue(m.Get(ctx, "key")))
		it.Eventually(t, func() bool {
			return it.MustValue(m.Get(ctx, "key")) == nil
		})
	})
}

func TestMap_Delete(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		it.Must(m.Set(context.Background(), "key", targetValue))
		if value := it.MustValue(m.Get(context.Background(), "key")); targetValue != value {
			t.Fatalf("target %v != %v", targetValue, value)
		}
		if err := m.Delete(context.Background(), "key"); err != nil {
			t.Fatal(err)
		}
		if value := it.MustValue(m.Get(context.Background(), "key")); nil != value {
			t.Fatalf("target nil != %v", value)
		}
	})
}

func TestMap_Evict(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		if err := m.Set(context.Background(), "key", targetValue); err != nil {
			t.Fatal(err)
		}
		if value, err := m.Evict(context.Background(), "key"); err != nil {
			t.Fatal(err)
		} else if !value {
			t.Fatalf("target true != %v", value)
		}
	})
}

func TestMap_Clear(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		it.Must(m.Set(context.Background(), "key", targetValue))
		if ok := it.MustBool(m.ContainsKey(context.Background(), "key")); !ok {
			t.Fatalf("key not found")
		}
		if ok := it.MustBool(m.ContainsValue(context.Background(), "value")); !ok {
			t.Fatalf("value not found")
		}
		value := it.MustValue(m.Get(context.Background(), "key"))
		if targetValue != value {
			t.Fatalf("target %v != %v", targetValue, value)
		}
		if err := m.Clear(context.Background()); err != nil {
			t.Fatal(err)
		}
		value = it.MustValue(m.Get(context.Background(), "key"))
		if nil != value {
			t.Fatalf("target nil!= %v", value)
		}
	})
}

func TestMap_Remove(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValue := "value"
		it.Must(m.Set(context.Background(), "key", targetValue))
		if !it.MustBool(m.ContainsKey(context.Background(), "key")) {
			t.Fatalf("key not found")
		}
		if value, err := m.Remove(context.Background(), "key"); err != nil {
			t.Fatal(err)
		} else if targetValue != value {
			t.Fatalf("target nil != %v", value)
		}
		if it.MustBool(m.ContainsKey(context.Background(), "key")) {
			t.Fatalf("key found")
		}
	})
}

func TestMap_GetAll(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		const maxKeys = 100
		makeKey := func(id int) string {
			return fmt.Sprintf("k%d", id)
		}
		makeValue := func(id int) string {
			return fmt.Sprintf("v%d", id)
		}
		var allPairs []types.Entry
		for i := 0; i < maxKeys; i++ {
			allPairs = append(allPairs, types.NewEntry(makeKey(i), makeValue(i)))
		}
		var keys []interface{}
		var target []types.Entry
		for i, pair := range allPairs {
			if i%2 == 0 {
				keys = append(keys, pair.Key)
				target = append(target, pair)
			}
		}
		for _, pair := range allPairs {
			it.Must(m.Set(context.Background(), pair.Key, pair.Value))
		}
		for _, pair := range allPairs {
			it.AssertEquals(t, pair.Value, it.MustValue(m.Get(context.Background(), pair.Key)))
		}
		if kvs, err := m.GetAll(context.Background(), keys...); err != nil {
			t.Fatal(err)
		} else if !entriesEqualUnordered(target, kvs) {
			t.Fatalf("target: %#v != %#v", target, kvs)
		}
	})
}

func TestMap_GetKeySet(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetKeySet := []interface{}{"k1", "k2", "k3"}
		it.Must(m.Set(context.Background(), "k1", "v1"))
		it.Must(m.Set(context.Background(), "k2", "v2"))
		it.Must(m.Set(context.Background(), "k3", "v3"))
		it.AssertEquals(t, "v1", it.MustValue(m.Get(context.Background(), "k1")))
		it.AssertEquals(t, "v2", it.MustValue(m.Get(context.Background(), "k2")))
		it.AssertEquals(t, "v3", it.MustValue(m.Get(context.Background(), "k3")))
		if keys, err := m.GetKeySet(context.Background()); err != nil {
			t.Fatal(err)
		} else if !assert.ElementsMatch(t, targetKeySet, keys) {
			t.FailNow()
		}
	})
}

func TestMap_GetKeySetWithPredicate(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetKeySet := []interface{}{serialization.JSON(`{"a": 15}`)}
		it.Must(m.Set(context.Background(), serialization.JSON(`{"a": 5}`), "v1"))
		it.Must(m.Set(context.Background(), serialization.JSON(`{"a": 10}`), "v2"))
		it.Must(m.Set(context.Background(), serialization.JSON(`{"a": 15}`), "v3"))
		if keys, err := m.GetKeySetWithPredicate(context.Background(), predicate.GreaterOrEqual("__key.a", 11)); err != nil {
			t.Fatal(err)
		} else if !assert.ElementsMatch(t, targetKeySet, keys) {
			t.FailNow()
		}
	})
}

func TestMap_GetValues(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValues := []interface{}{"v1", "v2", "v3"}
		it.Must(m.Set(context.Background(), "k1", "v1"))
		it.Must(m.Set(context.Background(), "k2", "v2"))
		it.Must(m.Set(context.Background(), "k3", "v3"))
		it.AssertEquals(t, "v1", it.MustValue(m.Get(context.Background(), "k1")))
		it.AssertEquals(t, "v2", it.MustValue(m.Get(context.Background(), "k2")))
		it.AssertEquals(t, "v3", it.MustValue(m.Get(context.Background(), "k3")))
		if values, err := m.GetValues(context.Background()); err != nil {
			t.Fatal(err)
		} else if !assert.ElementsMatch(t, targetValues, values) {
			t.FailNow()
		}
	})
}

func TestMap_GetValuesWithPredicate(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		targetValues := []interface{}{serialization.JSON(`{"A": 10, "B": 200}`), serialization.JSON(`{"A": 10, "B": 30}`)}
		it.Must(m.Set(context.Background(), "k1", serialization.JSON(`{"A": 10, "B": 200}`)))
		it.Must(m.Set(context.Background(), "k2", serialization.JSON(`{"A": 10, "B": 30}`)))
		it.Must(m.Set(context.Background(), "k3", serialization.JSON(`{"A": 5, "B": 200}`)))

		it.AssertEquals(t, serialization.JSON(`{"A": 10, "B": 200}`), it.MustValue(m.Get(context.Background(), "k1")))
		it.AssertEquals(t, serialization.JSON(`{"A": 10, "B": 30}`), it.MustValue(m.Get(context.Background(), "k2")))
		it.AssertEquals(t, serialization.JSON(`{"A": 5, "B": 200}`), it.MustValue(m.Get(context.Background(), "k3")))
		if values, err := m.GetValuesWithPredicate(context.Background(), predicate.Equal("A", 10)); err != nil {
			t.Fatal(err)
		} else if !assert.ElementsMatch(t, targetValues, values) {
			t.FailNow()
		}
	})
}

func TestMap_PutAll(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		pairs := []types.Entry{
			types.NewEntry("k1", "v1"),
			types.NewEntry("k2", "v2"),
			types.NewEntry("k3", "v3"),
		}
		if err := m.PutAll(context.Background(), pairs...); err != nil {
			t.Fatal(err)
		}
		it.AssertEquals(t, "v1", it.MustValue(m.Get(context.Background(), "k1")))
		it.AssertEquals(t, "v2", it.MustValue(m.Get(context.Background(), "k2")))
		it.AssertEquals(t, "v3", it.MustValue(m.Get(context.Background(), "k3")))
	})
}

func TestMap_GetEntrySet(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		target := []types.Entry{
			types.NewEntry("k1", "v1"),
			types.NewEntry("k2", "v2"),
			types.NewEntry("k3", "v3"),
		}
		if err := m.PutAll(context.Background(), target...); err != nil {
			t.Fatal(err)
		}
		if entries, err := m.GetEntrySet(context.Background()); err != nil {
			t.Fatal(err)
		} else if !entriesEqualUnordered(target, entries) {
			t.Fatalf("target: %#v != %#v", target, entries)
		}
	})
}

func TestMap_GetEntrySetWithPredicateUsingPortable(t *testing.T) {
	cbCallback := func(config *hz.Config) {
		config.Serialization.SetPortableFactories(it.SamplePortableFactory{})
	}
	it.MapTesterWithConfig(t, cbCallback, func(t *testing.T, m *hz.Map) {
		okValue := "foo-Ğİ"
		noValue := "foo"
		entries := []types.Entry{
			types.NewEntry("k1", &it.SamplePortable{A: okValue, B: 10}),
			types.NewEntry("k2", &it.SamplePortable{A: noValue, B: 15}),
			types.NewEntry("k3", &it.SamplePortable{A: okValue, B: 10}),
		}
		if err := m.PutAll(context.Background(), entries...); err != nil {
			t.Fatal(err)
		}
		target := []types.Entry{
			types.NewEntry("k1", &it.SamplePortable{A: okValue, B: 10}),
			types.NewEntry("k3", &it.SamplePortable{A: okValue, B: 10}),
		}
		if entries, err := m.GetEntrySetWithPredicate(context.Background(), predicate.And(predicate.Equal("A", okValue), predicate.Equal("B", 10))); err != nil {
			t.Fatal(err)
		} else if !entriesEqualUnordered(target, entries) {
			t.Fatalf("target: %#v != %#v", target, entries)
		}
	})
}

func TestMap_GetEntrySetWithPredicateUsingJSON(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		entries := []types.Entry{
			types.NewEntry("k1", it.SamplePortable{A: "foo", B: 10}.Json()),
			types.NewEntry("k2", it.SamplePortable{A: "foo", B: 15}.Json()),
			types.NewEntry("k3", it.SamplePortable{A: "foo", B: 10}.Json()),
		}
		if err := m.PutAll(context.Background(), entries...); err != nil {
			t.Fatal(err)
		}
		target := []types.Entry{
			types.NewEntry("k1", it.SamplePortable{A: "foo", B: 10}.Json()),
			types.NewEntry("k3", it.SamplePortable{A: "foo", B: 10}.Json()),
		}
		if entries, err := m.GetEntrySetWithPredicate(context.Background(), predicate.And(predicate.Equal("A", "foo"), predicate.Equal("B", 10))); err != nil {
			t.Fatal(err)
		} else if !entriesEqualUnordered(target, entries) {
			t.Fatalf("target: %#v != %#v", target, entries)
		}
	})
}

func TestMap_GetEntryView(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		it.Must(m.Set(context.Background(), "k1", "v1"))
		if ev, err := m.GetEntryView(context.Background(), "k1"); err != nil {
			t.Fatal(err)
		} else {
			it.AssertEquals(t, "k1", ev.Key)
			it.AssertEquals(t, "v1", ev.Value)
		}
	})
}

func TestMap_GetEntryView_2(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		it.Must(m.Set(context.Background(), int64(1), int64(2)))
		if ev, err := m.GetEntryView(context.Background(), int64(1)); err != nil {
			t.Fatal(err)
		} else {
			assert.Equal(t, int64(1), ev.Key)
			assert.Equal(t, int64(2), ev.Value)
		}
	})
}

func TestMap_GetEntryView_KeyNotFound(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		if ev, err := m.GetEntryView(context.Background(), "k1"); err != nil {
			t.Fatal(err)
		} else if ev != nil {
			t.Fatalf("ev should be nil")
		}
	})
}

func TestMap_AddIndexWithConfig(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		it.Must(m.Set(context.Background(), "k1", serialization.JSON(`{"A": 10, "B": 40}`)))
		indexConfig := types.IndexConfig{
			Name:               "my-index",
			Type:               types.IndexTypeBitmap,
			Attributes:         []string{"A"},
			BitmapIndexOptions: types.BitmapIndexOptions{UniqueKey: "B", UniqueKeyTransformation: types.UniqueKeyTransformationLong},
		}
		if err := m.AddIndex(context.Background(), indexConfig); err != nil {
			t.Fatal(err)
		}
	})
}

func TestMap_AddIndexValidationError(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		indexConfig := types.IndexConfig{
			Name:               "my-index",
			Type:               types.IndexTypeBitmap,
			Attributes:         []string{"A", "B"},
			BitmapIndexOptions: types.BitmapIndexOptions{UniqueKey: "B", UniqueKeyTransformation: types.UniqueKeyTransformationLong},
		}
		if err := m.AddIndex(context.Background(), indexConfig); err == nil {
			t.Fatalf("should have failed")
		} else if !errors.Is(err, hzerrors.ErrIllegalArgument) {
			t.Fatalf("should have returned an illegal argument error")
		}
	})
}

func TestMap_Flush(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		it.Must(m.Set(context.Background(), "k1", "v1"))
		if err := m.Flush(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestMap_LoadAllWithoutReplacing(t *testing.T) {
	makeMapName := func(_ ...string) string {
		// the map name for this test should be static.
		return "test-map"
	}
	it.MapTesterWithConfigAndName(t, makeMapName, nil, func(t *testing.T, m *hz.Map) {
		putSampleKeyValues(m, 2)
		it.Must(m.EvictAll(context.Background()))
		it.Must(m.PutTransient(context.Background(), "k0", "new-v0"))
		it.Must(m.PutTransient(context.Background(), "k1", "new-v1"))
		it.Must(m.LoadAllWithoutReplacing(context.Background(), "k0", "k1"))
		targetEntrySet := []types.Entry{
			{Key: "k0", Value: "new-v0"},
			{Key: "k1", Value: "new-v1"},
		}
		entrySet := it.MustValue(m.GetAll(context.Background(), "k0", "k1")).([]types.Entry)
		if !entriesEqualUnordered(targetEntrySet, entrySet) {
			t.Fatalf("target %#v != %#v", targetEntrySet, entrySet)
		}
	})
}

func TestMap_LoadAllReplacing(t *testing.T) {
	makeMapName := func(_ ...string) string {
		// the map name for this test should be static.
		return "test-map"
	}
	it.MapTesterWithConfigAndName(t, makeMapName, nil, func(t *testing.T, m *hz.Map) {
		keys := putSampleKeyValues(m, 10)
		it.Must(m.EvictAll(context.Background()))
		it.Must(m.LoadAllReplacing(context.Background()))
		entrySet := it.MustValue(m.GetAll(context.Background(), keys...)).([]types.Entry)
		if len(keys) != len(entrySet) {
			t.Fatalf("target len: %d != %d", len(keys), len(entrySet))
		}
		it.Must(m.EvictAll(context.Background()))
		keys = keys[:5]
		it.Must(m.LoadAllReplacing(context.Background(), keys...))
		entrySet = it.MustValue(m.GetAll(context.Background(), keys...)).([]types.Entry)
		if len(keys) != len(entrySet) {
			t.Fatalf("target len: %d != %d", len(keys), len(entrySet))
		}
	})
}

func TestMap_Lock(t *testing.T) {
	it.MapTester(t, func(t *testing.T, cm *hz.Map) {
		const goroutineCount = 100
		const key = "counter"
		wg := &sync.WaitGroup{}
		wg.Add(goroutineCount)
		for i := 0; i < goroutineCount; i++ {
			go func() {
				defer wg.Done()
				intValue := int64(0)
				lockCtx := cm.NewLockContext(context.Background())
				if err := cm.Lock(lockCtx, key); err != nil {
					panic(err)
				}
				defer cm.Unlock(lockCtx, key)
				v, err := cm.Get(lockCtx, key)
				if err != nil {
					panic(err)
				}
				if v != nil {
					intValue = v.(int64)
				}
				intValue++
				if err = cm.Set(lockCtx, key, intValue); err != nil {
					panic(err)
				}
			}()
		}
		wg.Wait()
		if lastValue, err := cm.Get(context.Background(), key); err != nil {
			panic(err)
		} else {
			assert.Equal(t, int64(goroutineCount), lastValue)
		}
	})
}

func TestMap_ForceUnlock(t *testing.T) {
	it.MapTester(t, func(t *testing.T, cm *hz.Map) {
		lockCtx := cm.NewLockContext(context.Background())
		if err := cm.Lock(lockCtx, "k1"); err != nil {
			t.Fatal(err)
		}
		if locked, err := cm.IsLocked(lockCtx, "k1"); err != nil {
			t.Fatal(err)
		} else {
			it.AssertEquals(t, true, locked)
		}
		if err := cm.ForceUnlock(lockCtx, "k1"); err != nil {
			t.Fatal(err)
		}
		if locked, err := cm.IsLocked(lockCtx, "k1"); err != nil {
			t.Fatal(err)
		} else {
			assert.Equal(t, false, locked)
		}
	})
}

func TestMap_IsEmptySize(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		if value, err := m.IsEmpty(context.Background()); err != nil {
			t.Fatal(err)
		} else if !value {
			t.Fatalf("target: true != false")
		}
		targetSize := 0
		if value, err := m.Size(context.Background()); err != nil {
			t.Fatal(err)
		} else if targetSize != value {
			t.Fatalf("target: %d != %d", targetSize, value)
		}
		it.MustValue(m.Put(context.Background(), "k1", "v1"))
		it.MustValue(m.Put(context.Background(), "k2", "v2"))
		it.MustValue(m.Put(context.Background(), "k3", "v3"))
		it.Eventually(t, func() bool {
			value, err := m.IsEmpty(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			return !value
		})
		targetSize = 3
		it.Eventually(t, func() bool {
			value, err := m.Size(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			return targetSize == value
		})
	})
}

func TestMap_RemoveAll(t *testing.T) {
	cbCallback := func(config *hz.Config) {
		config.Serialization.SetPortableFactories(it.SamplePortableFactory{})
	}
	it.MapTesterWithConfig(t, cbCallback, func(t *testing.T, m *hz.Map) {
		entries := []types.Entry{
			types.NewEntry("k1", &it.SamplePortable{A: "foo", B: 10}),
			types.NewEntry("k2", &it.SamplePortable{A: "foo", B: 15}),
			types.NewEntry("k3", &it.SamplePortable{A: "foo", B: 10}),
		}
		it.Must(m.PutAll(context.Background(), entries...))
		if err := m.RemoveAll(context.Background(), predicate.Equal("B", 10)); err != nil {
			t.Fatal(err)
		}
		target := []types.Entry{
			types.NewEntry("k2", &it.SamplePortable{A: "foo", B: 15}),
		}
		if kvs, err := m.GetAll(context.Background(), "k1", "k2", "k3"); err != nil {
			t.Fatal(err)
		} else if !entriesEqualUnordered(target, kvs) {
			t.Fatalf("target: %#v != %#v", target, kvs)
		}

	})
}

func TestMap_RemoveIfSame(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		it.Must(m.Set(context.Background(), "k1", "v1"))
		it.Must(m.Set(context.Background(), "k2", "v2"))
		it.AssertEquals(t, "v1", it.MustValue(m.Get(context.Background(), "k1")))
		it.AssertEquals(t, "v2", it.MustValue(m.Get(context.Background(), "k2")))
		if removed, err := m.RemoveIfSame(context.Background(), "k1", "v1"); err != nil {
			t.Fatal(err)
		} else if !removed {
			t.Fatalf("not removed")
		}
		it.AssertEquals(t, false, it.MustValue(m.ContainsKey(context.Background(), "k1")))
		if removed, err := m.RemoveIfSame(context.Background(), "k2", "v1"); err != nil {
			t.Fatal(err)
		} else if removed {
			t.Fatalf("removed")
		}
		it.AssertEquals(t, true, it.MustValue(m.ContainsKey(context.Background(), "k2")))
	})
}

func TestMap_Replace(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		it.Must(m.Set(context.Background(), "k1", "v1"))
		it.AssertEquals(t, "v1", it.MustValue(m.Get(context.Background(), "k1")))
		if oldValue, err := m.Replace(context.Background(), "k1", "v2"); err != nil {
			t.Fatal(err)
		} else {
			it.AssertEquals(t, "v1", oldValue)
		}
		it.AssertEquals(t, "v2", it.MustValue(m.Get(context.Background(), "k1")))
	})
}

func TestMap_ReplaceIfSame(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		it.Must(m.Set(context.Background(), "k1", "v1"))
		it.AssertEquals(t, "v1", it.MustValue(m.Get(context.Background(), "k1")))
		if replaced, err := m.ReplaceIfSame(context.Background(), "k1", "v1", "v2"); err != nil {
			t.Fatal(err)
		} else {
			it.AssertEquals(t, true, replaced)
		}
		it.AssertEquals(t, "v2", it.MustValue(m.Get(context.Background(), "k1")))
	})
}

func TestMap_EntryNotifiedEvent(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		const totalCallCount = int32(100)
		callCount := int32(0)
		handler := func(event *hz.EntryNotified) {
			if event.EventType == hz.EntryAdded {
				atomic.AddInt32(&callCount, 1)
			}
		}
		listenerConfig := hz.MapEntryListenerConfig{
			IncludeValue: true,
		}
		listenerConfig.NotifyEntryAdded(true)
		subscriptionID, err := m.AddEntryListener(context.Background(), listenerConfig, handler)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < int(totalCallCount); i++ {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			it.MustValue(m.Put(context.Background(), key, value))
		}
		it.Eventually(t, func() bool {
			cc := atomic.LoadInt32(&callCount)
			t.Logf("call count target: %d, current: %d", totalCallCount, cc)
			return cc == totalCallCount
		})
		atomic.StoreInt32(&callCount, 0)
		if err := m.RemoveEntryListener(context.Background(), subscriptionID); err != nil {
			t.Fatal(err)
		}
		// Clear the map to get entry added events
		err = m.Clear(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < int(totalCallCount); i++ {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			it.MustValue(m.Put(context.Background(), key, value))
		}
		if !assert.Equal(t, int32(0), atomic.LoadInt32(&callCount)) {
			t.FailNow()
		}
	})
}

func TestMap_EntryNotifiedEventWithAddListener(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		const totalCallCount = int32(100)
		callCount := int32(0)
		subscriptionID, err := m.AddListener(context.Background(), hz.MapListener{
			EntryAdded: func(event *hz.EntryNotified) {
				if event.EventType != hz.EntryAdded {
					t.Fatalf("unexpected event type: %v", event.EventType)
				}
				atomic.AddInt32(&callCount, 1)
			},
		}, true)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < int(totalCallCount); i++ {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			it.MustValue(m.Put(context.Background(), key, value))
		}
		it.Eventually(t, func() bool {
			cc := atomic.LoadInt32(&callCount)
			t.Logf("call count target: %d, current: %d", totalCallCount, cc)
			return cc == totalCallCount
		})
		atomic.StoreInt32(&callCount, 0)
		// Remove the listener and test that entry added events are not received anymore
		if err := m.RemoveListener(context.Background(), subscriptionID); err != nil {
			t.Fatal(err)
		}
		// Clear the map to get entry added events
		err = m.Clear(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < int(totalCallCount); i++ {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			it.MustValue(m.Put(context.Background(), key, value))
		}
		cc := atomic.LoadInt32(&callCount)
		if !assert.Equal(t, int32(0), cc) {
			t.Fatalf("call count target: %d, current: %d", 0, cc)
		}
	})
}

func TestMap_EntryNotifiedEventToKey(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		callCount := int32(0)
		handler := func(event *hz.EntryNotified) {
			if event.EventType != hz.EntryAdded {
				t.Fatalf("unexpected event type: %v", event.EventType)
			}
			atomic.AddInt32(&callCount, 1)
		}
		listenerConfig := hz.MapEntryListenerConfig{
			IncludeValue: true,
			Key:          "k1",
		}
		listenerConfig.NotifyEntryAdded(true)
		if _, err := m.AddEntryListener(context.Background(), listenerConfig, handler); err != nil {
			t.Fatal(err)
		}
		it.MustValue(m.Put(context.Background(), "k1", "v1"))
		it.Eventually(t, func() bool {
			cc := atomic.LoadInt32(&callCount)
			totalCallCount := int32(1)
			t.Logf("call count target: %d, current: %d", totalCallCount, cc)
			return cc == totalCallCount
		})
	})
}

func TestMap_EntryNotifiedEventToKeyWithAddListenerWithKey(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		const totalCallCount = int32(1)
		callCount := int32(0)
		if _, err := m.AddListenerWithKey(context.Background(), hz.MapListener{
			EntryAdded: func(event *hz.EntryNotified) {
				if event.EventType != hz.EntryAdded {
					t.Fatalf("unexpected event type: %v", event.EventType)
				}
				atomic.AddInt32(&callCount, 1)
			},
		}, "k1", true); err != nil {
			t.Fatal(err)
		}
		it.MustValue(m.Put(context.Background(), "k1", "v1"))
		it.Eventually(t, func() bool {
			cc := atomic.LoadInt32(&callCount)
			t.Logf("call count target: %d, current: %d", totalCallCount, cc)
			return cc == totalCallCount
		})
	})
}

func TestMap_EntryNotifiedEventWithPredicate(t *testing.T) {
	cbCallback := func(config *hz.Config) {
		config.Serialization.SetPortableFactories(it.SamplePortableFactory{})
	}
	it.MapTesterWithConfig(t, cbCallback, func(t *testing.T, m *hz.Map) {
		const totalCallCount = int32(100)
		callCount := int32(0)
		handler := func(event *hz.EntryNotified) {
			if event.EventType != hz.EntryAdded {
				t.Fatalf("unexpected event type: %v", event.EventType)
			}
			atomic.AddInt32(&callCount, 1)
		}
		listenerConfig := hz.MapEntryListenerConfig{
			IncludeValue: true,
			Predicate:    predicate.Equal("A", "foo"),
		}
		listenerConfig.NotifyEntryAdded(true)
		subID, err := m.AddEntryListener(context.Background(), listenerConfig, handler)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("TestMap_EntryNotifiedEventWithPredicate subscriptionID: %s", subID)
		for i := 0; i < int(totalCallCount); i++ {
			key := fmt.Sprintf("key-%d", i)
			value := &it.SamplePortable{A: "foo", B: int32(i)}
			it.MustValue(m.Put(context.Background(), key, value))
		}
		it.Eventually(t, func() bool {
			cc := atomic.LoadInt32(&callCount)
			t.Logf("call count target: %d, current: %d", totalCallCount, cc)
			return cc == totalCallCount
		})
	})
}

func TestMap_EntryNotifiedEventWithPredicateWithAddListenerWithPredicate(t *testing.T) {
	cbCallback := func(config *hz.Config) {
		config.Serialization.SetPortableFactories(it.SamplePortableFactory{})
	}
	it.MapTesterWithConfig(t, cbCallback, func(t *testing.T, m *hz.Map) {
		const totalCallCount = int32(100)
		callCount := int32(0)
		subID, err := m.AddListenerWithPredicate(context.Background(), hz.MapListener{
			EntryAdded: func(event *hz.EntryNotified) {
				if event.EventType != hz.EntryAdded {
					t.Fatalf("unexpected event type: %v", event.EventType)
				}
				atomic.AddInt32(&callCount, 1)
			},
		}, predicate.Equal("A", "foo"), true)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("TestMap_EntryNotifiedEventWithPredicateWithAddListenerWithPredicate subscriptionID: %s", subID)
		for i := 0; i < int(totalCallCount); i++ {
			key := fmt.Sprintf("key-%d", i)
			value := &it.SamplePortable{A: "foo", B: int32(i)}
			it.MustValue(m.Put(context.Background(), key, value))
		}
		it.Eventually(t, func() bool {
			cc := atomic.LoadInt32(&callCount)
			t.Logf("call count target: %d, current: %d", totalCallCount, cc)
			return cc == totalCallCount
		})
	})
}

func TestMap_EntryNotifiedEventToKeyAndPredicate(t *testing.T) {
	cbCallback := func(config *hz.Config) {
		config.Serialization.SetPortableFactories(it.SamplePortableFactory{})
	}
	it.MapTesterWithConfig(t, cbCallback, func(t *testing.T, m *hz.Map) {
		callCount := int32(0)
		handler := func(event *hz.EntryNotified) {
			if event.EventType != hz.EntryAdded {
				t.Fatalf("unexpected event type: %v", event.EventType)
			}
			atomic.AddInt32(&callCount, 1)
		}
		listenerConfig := hz.MapEntryListenerConfig{
			IncludeValue: true,
			Key:          "k1",
			Predicate:    predicate.Equal("A", "foo"),
		}
		listenerConfig.NotifyEntryAdded(true)
		if _, err := m.AddEntryListener(context.Background(), listenerConfig, handler); err != nil {
			t.Fatal(err)
		}
		it.MustValue(m.Put(context.Background(), "k1", &it.SamplePortable{A: "foo", B: 10}))
		it.MustValue(m.Put(context.Background(), "k1", &it.SamplePortable{A: "bar", B: 10}))
		it.MustValue(m.Put(context.Background(), "k2", &it.SamplePortable{A: "foo", B: 10}))
		it.Eventually(t, func() bool {
			cc := atomic.LoadInt32(&callCount)
			totalCallCount := int32(1)
			t.Logf("call count target: %d, current: %d", totalCallCount, cc)
			return cc == totalCallCount
		})
	})
}

func TestMap_EntryNotifiedEventToKeyAndPredicateWithAddListenerWithPredicateAndKey(t *testing.T) {
	cbCallback := func(config *hz.Config) {
		config.Serialization.SetPortableFactories(it.SamplePortableFactory{})
	}
	it.MapTesterWithConfig(t, cbCallback, func(t *testing.T, m *hz.Map) {
		const totalCallCount = int32(1)
		callCount := int32(0)
		_, err := m.AddListenerWithPredicateAndKey(context.Background(), hz.MapListener{
			EntryAdded: func(event *hz.EntryNotified) {
				if event.EventType != hz.EntryAdded {
					t.Fatalf("unexpected event type: %v", event.EventType)
				}
				atomic.AddInt32(&callCount, 1)
			},
		}, predicate.Equal("A", "foo"), "k1", true)

		if err != nil {
			t.Fatal(err)
		}
		it.MustValue(m.Put(context.Background(), "k1", &it.SamplePortable{A: "foo", B: 10}))
		it.MustValue(m.Put(context.Background(), "k1", &it.SamplePortable{A: "bar", B: 10}))
		it.MustValue(m.Put(context.Background(), "k2", &it.SamplePortable{A: "foo", B: 10}))
		it.Eventually(t, func() bool {
			cc := atomic.LoadInt32(&callCount)
			t.Logf("call count target: %d, current: %d", totalCallCount, cc)
			return cc == totalCallCount
		})
	})
}

func TestMap_Destroy(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		if err := m.Destroy(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestMap_Aggregate(t *testing.T) {
	cbCallback := func(config *hz.Config) {
		config.Serialization.SetPortableFactories(it.SamplePortableFactory{})
	}
	it.MapTesterWithConfig(t, cbCallback, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		it.MustValue(m.Put(ctx, "k1", &it.SamplePortable{A: "foo", B: 10}))
		it.MustValue(m.Put(ctx, "k2", &it.SamplePortable{A: "bar", B: 30}))
		it.MustValue(m.Put(ctx, "k3", &it.SamplePortable{A: "zoo", B: 30}))
		result, err := m.Aggregate(ctx, aggregate.Count("B"))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int64(3), result)
	})
}

func TestMap_Aggregate_2(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		it.MustValue(m.Put(ctx, "k1", serialization.JSON(`{"A": "foo", "B": 10}`)))
		it.MustValue(m.Put(ctx, "k2", serialization.JSON(`{"A": "bar", "B": 30}`)))
		it.MustValue(m.Put(ctx, "k3", serialization.JSON(`{"A": "zoo", "B": 30}`)))
		result, err := m.Aggregate(ctx, aggregate.Count("B"))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int64(3), result)
	})
}

func TestMap_AggregateWithPredicate(t *testing.T) {
	cbCallback := func(config *hz.Config) {
		config.Serialization.SetPortableFactories(it.SamplePortableFactory{})
	}
	it.MapTesterWithConfig(t, cbCallback, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		it.MustValue(m.Put(ctx, "k1", &it.SamplePortable{A: "foo", B: 10}))
		it.MustValue(m.Put(ctx, "k2", &it.SamplePortable{A: "bar", B: 30}))
		it.MustValue(m.Put(ctx, "k2", &it.SamplePortable{A: "zoo", B: 30}))
		result, err := m.AggregateWithPredicate(ctx, aggregate.Count("B"), predicate.Equal("A", "foo"))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int64(1), result)
	})
}

func TestMap_SetWithTTLAndMaxIdle(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		targetValue := "value"
		if err := m.SetWithTTLAndMaxIdle(ctx, "key", targetValue, 20*time.Second, 5*time.Second); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, targetValue, it.MustValue(m.Get(ctx, "key")))
		time.Sleep(10 * time.Second)
		it.Eventually(t, func() bool {
			v := it.MustValue(m.Get(ctx, "key"))
			return v == nil
		})
	})
}

func TestMap_TryLock(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		const key = "foo"
		go func() {
			ctx := m.NewLockContext(context.Background())
			it.Must(m.Lock(ctx, key))
			wg.Done()
		}()
		wg.Wait()
		mainCtx := m.NewLockContext(context.Background())
		b, err := m.TryLock(mainCtx, key)
		if err != nil {
			t.Fatal(err)
		}
		assert.False(t, b)
	})
}

func TestMap_LockWithLease(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		const key = "foo"
		go func() {
			ctx := m.NewLockContext(context.Background())
			it.Must(m.LockWithLease(ctx, key, 100*time.Millisecond))
			wg.Done()
		}()
		wg.Wait()
		mainCtx := m.NewLockContext(context.Background())
		err := m.Lock(mainCtx, key)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestMap_TryLockWithLease(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		const key = "foo"
		go func() {
			ctx := m.NewLockContext(context.Background())
			if b := it.MustValue(m.TryLockWithLease(ctx, key, 100*time.Millisecond)); b != true {
				panic("unexpected value")
			}
			wg.Done()
		}()
		wg.Wait()
		mainCtx := m.NewLockContext(context.Background())
		err := m.Lock(mainCtx, key)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestMap_TryLockWithTimeout(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		const key = "foo"
		go func() {
			ctx := m.NewLockContext(context.Background())
			it.Must(m.Lock(ctx, key))
			wg.Done()
		}()
		wg.Wait()
		mainCtx := m.NewLockContext(context.Background())
		b, err := m.TryLockWithTimeout(mainCtx, key, 100*time.Millisecond)
		if err != nil {
			t.Fatal(err)
		}
		assert.False(t, b)
	})
}

func TestMap_TryLockWithLeaseAndTimeout(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		const key = "foo"
		go func() {
			ctx := m.NewLockContext(context.Background())
			it.Must(m.Lock(ctx, key))
			wg.Done()
		}()
		wg.Wait()
		mainCtx := m.NewLockContext(context.Background())
		b, err := m.TryLockWithLeaseAndTimeout(mainCtx, key, 100*time.Millisecond, 200*time.Millisecond)
		if err != nil {
			t.Fatal(err)
		}
		assert.False(t, b)
	})
}

func TestMap_SetTTL(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		targetValue := "value"
		it.Must(m.Set(ctx, "key", targetValue))
		if err := m.SetTTL(ctx, "key", 20*time.Second); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, targetValue, it.MustValue(m.Get(ctx, "key")))
		it.Eventually(t, func() bool {
			return it.MustValue(m.Get(ctx, "key")) == nil
		})
	})
}

func TestMap_SetTTLAffected(t *testing.T) {
	it.SkipIf(t, "hz < 4.2")
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		testCases := []struct {
			key           string
			isEffected    bool
			errorExpected bool
		}{
			// happy path
			{
				key:        "k1",
				isEffected: true,
			},
			// setTTL on non-existing key
			{
				key:           "k2",
				isEffected:    false,
				errorExpected: false,
			},
			// setTTL on already expired key
			{
				key:           "k3",
				isEffected:    false,
				errorExpected: false,
			},
		}
		_ = it.MustValue(m.Put(ctx, "k1", "someValue"))
		_ = it.MustValue(m.PutWithTTL(ctx, "k3", "someValue", time.Millisecond))
		time.Sleep(time.Millisecond)
		for _, tc := range testCases {
			affected, err := m.SetTTLAffected(ctx, tc.key, time.Second)
			assert.Equal(t, tc.errorExpected, err != nil)
			assert.Equal(t, tc.isEffected, affected)
		}
	})
}

func TestMap_ExecuteOnEntries(t *testing.T) {
	cb := func(c *hz.Config) {
		c.Serialization.SetIdentifiedDataSerializableFactories(&SimpleEntryProcessorFactory{})
	}
	it.MapTesterWithConfig(t, cb, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		it.MustValue(m.Put(ctx, "my-key", "my-value"))
		vs, err := m.ExecuteOnEntries(ctx, &SimpleEntryProcessor{value: "test"})
		if err != nil {
			t.Fatal(err)
		}
		target := []types.Entry{{Key: "my-key", Value: "test"}}
		assert.Equal(t, target, vs)
	})
}

func TestMap_ExecuteOnEntriesWithPredicate(t *testing.T) {
	cb := func(c *hz.Config) {
		c.Serialization.SetIdentifiedDataSerializableFactories(&SimpleEntryProcessorFactory{})
	}
	it.MapTesterWithConfig(t, cb, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		it.Must(m.Set(ctx, "k1", serialization.JSON(`{"A": 10, "B": 200}`)))
		it.Must(m.Set(ctx, "k2", serialization.JSON(`{"A": 10, "B": 30}`)))
		it.Must(m.Set(ctx, "k3", serialization.JSON(`{"A": 5, "B": 200}`)))
		vs, err := m.ExecuteOnEntriesWithPredicate(ctx, &SimpleEntryProcessor{value: "test"}, predicate.Equal("A", 10))
		if err != nil {
			t.Fatal(err)
		}
		updatedTo := []types.Entry{{Key: "k1", Value: "test"}, {Key: "k2", Value: "test"}}
		sort.Slice(vs, func(i, j int) bool {
			return vs[i].Key.(string) < vs[j].Key.(string)
		})
		assert.Equal(t, updatedTo, vs)
		it.AssertEquals(t, serialization.JSON(`{"A": 5, "B": 200}`), it.MustValue(m.Get(context.Background(), "k3"))) //ensure this should not change
	})
}

func TestMap_ExecuteOnKey(t *testing.T) {
	cb := func(c *hz.Config) {
		c.Serialization.SetIdentifiedDataSerializableFactories(&SimpleEntryProcessorFactory{})
	}
	it.MapTesterWithConfig(t, cb, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		it.MustValue(m.Put(ctx, "k1", "my-value"))
		it.MustValue(m.Put(ctx, "k2", "my-value"))
		vs, err := m.ExecuteOnKey(ctx, &SimpleEntryProcessor{value: "test"}, "k1")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "test", vs)
		assert.Equal(t, "test", it.MustValue(m.Get(ctx, "k1")))
		assert.Equal(t, "my-value", it.MustValue(m.Get(ctx, "k2")))
	})
}

func TestMap_ExecuteOnKeys(t *testing.T) {
	cb := func(c *hz.Config) {
		c.Serialization.SetIdentifiedDataSerializableFactories(&SimpleEntryProcessorFactory{})
	}
	it.MapTesterWithConfig(t, cb, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		var (
			k1     = "my-key1"            //stable key
			k2, k3 = "my-key2", "my-key3" //to be updated
		)
		it.MustValue(m.Put(ctx, k1, "my-value"))
		it.MustValue(m.Put(ctx, k2, "my-value"))
		it.MustValue(m.Put(ctx, k3, "my-value"))
		vs, err := m.ExecuteOnKeys(ctx, &SimpleEntryProcessor{value: "test"}, k2, k3)
		if err != nil {
			t.Fatal(err)
		}
		target := []interface{}{"test", "test"}
		assert.Equal(t, target, vs)
		assert.Equal(t, "my-value", it.MustValue(m.Get(ctx, k1)))

		// Zero length keys
		vs, err = m.ExecuteOnKeys(ctx, &SimpleEntryProcessor{value: "test"})
		assert.Nil(t, vs, "zero length key list should return nil value")
		assert.Nil(t, err, "error should be nil on zero length key list")
	})
}

func TestMap_AddInterceptor(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		ctx := context.Background()
		prefix := "My Prefix"
		id, err := m.AddInterceptor(ctx, &MapGetInterceptor{Prefix: prefix})
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if _, err := m.RemoveInterceptor(ctx, id); err != nil {
				t.Fatal(err)
			}
		}()
		v, err := m.Get(ctx, "foo")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, prefix, v)
		ok, err := m.RemoveInterceptor(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, ok, true)
		v, err = m.Get(ctx, "foo")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, nil, v)
	})
}

func TestMap_TryPut(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		ctx1 := m.NewLockContext(context.Background())
		if err := m.Lock(ctx1, "foo"); err != nil {
			t.Fatal(err)
		}
		defer m.Unlock(ctx1, "foo")
		// TryPut with a different lock context returns false
		ctx2 := m.NewLockContext(context.Background())
		ok, err := m.TryPut(ctx2, "foo", "bar")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, false, ok)
		// TryPut with the same lock context returns true
		ok, err = m.TryPut(ctx1, "foo", "bar")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, true, ok)
	})
}

func TestMap_TryPutWithTimeout(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		ctx1 := m.NewLockContext(context.Background())
		if err := m.Lock(ctx1, "foo"); err != nil {
			t.Fatal(err)
		}
		// TryPut with a different lock context returns false
		ctx2 := m.NewLockContext(context.Background())
		ok, err := m.TryPutWithTimeout(ctx2, "foo", "bar", 1*time.Second)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, false, ok)
		// unlock 5 seconds later
		go func() {
			time.Sleep(5 * time.Second)
			if err := m.Unlock(ctx1, "foo"); err != nil {
				panic(err)
			}
		}()
		// TryPut after the timeout
		ok, err = m.TryPutWithTimeout(ctx2, "foo", "bar", 2*time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, true, ok)
	})
}

func TestMap_TryRemove(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		ctx1 := m.NewLockContext(context.Background())
		if _, err := m.Put(ctx1, "foo", "bar"); err != nil {
			t.Fatal(err)
		}
		if err := m.Lock(ctx1, "foo"); err != nil {
			t.Fatal(err)
		}
		// TryRemove with a different lock context returns false
		ctx2 := m.NewLockContext(context.Background())
		ok, err := m.TryRemove(ctx2, "foo")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, false, ok)
		// TryPut with the same lock context returns true
		ok, err = m.TryRemove(ctx1, "foo")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, true, ok)
	})
}

func TestMap_TryRemoveWithTimeout(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		ctx1 := m.NewLockContext(context.Background())
		if _, err := m.Put(ctx1, "foo", "bar"); err != nil {
			t.Fatal(err)
		}
		if err := m.Lock(ctx1, "foo"); err != nil {
			t.Fatal(err)
		}
		// TryPut with a different lock context returns false
		ctx2 := m.NewLockContext(context.Background())
		ok, err := m.TryPutWithTimeout(ctx2, "foo", "bar", 1*time.Second)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, false, ok)
		// unlock 5 seconds later
		go func() {
			time.Sleep(5 * time.Second)
			if err := m.Unlock(ctx1, "foo"); err != nil {
				panic(err)
			}
		}()
		// TryPut after the timeout
		ok, err = m.TryPutWithTimeout(ctx2, "foo", "bar", 2*time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, true, ok)
	})
}

// ==== Reliability Tests ====

func TestMapSetGet1000(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		const setGetCount = 1000
		for i := 0; i < setGetCount; i++ {
			key := fmt.Sprintf("k%d", i)
			value := fmt.Sprintf("v%d", i)
			it.Must(m.Set(context.Background(), key, value))
		}
		for i := 0; i < setGetCount; i++ {
			key := fmt.Sprintf("k%d", i)
			targetValue := fmt.Sprintf("v%d", i)
			if !assert.Equal(t, targetValue, it.MustValue(m.Get(context.Background(), key))) {
				t.FailNow()
			}
		}
	})
}

func TestMapSetGetLargePayload(t *testing.T) {
	it.MapTester(t, func(t *testing.T, m *hz.Map) {
		const payloadSize = 1 * 1024 * 1024
		payload := make([]byte, payloadSize)
		for i := 0; i < len(payload); i++ {
			payload[i] = byte(i)
		}
		it.Must(m.Set(context.Background(), "k1", payload))
		v := it.MustValue(m.Get(context.Background(), "k1"))
		if !assert.Equal(t, payload, v) {
			t.FailNow()
		}
	})
}

func makeStringSet(items []interface{}) map[string]struct{} {
	result := map[string]struct{}{}
	for _, item := range items {
		result[item.(string)] = struct{}{}
	}
	return result
}

func entriesEqualUnordered(p1, p2 []types.Entry) bool {
	if len(p1) != len(p2) {
		return false
	}
	// inefficient
	for _, item1 := range p1 {
		if entriesIndex(item1, p2) < 0 {
			return false
		}
	}
	for _, item2 := range p2 {
		if entriesIndex(item2, p1) < 0 {
			return false
		}
	}
	return true
}

func entriesIndex(p types.Entry, ps []types.Entry) int {
	for i, item := range ps {
		if reflect.DeepEqual(p, item) {
			return i
		}
	}
	return -1
}

func putSampleKeyValues(m *hz.Map, count int) []interface{} {
	keys := []interface{}{}
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("k%d", i)
		value := fmt.Sprintf("v%d", i)
		it.MustValue(m.Put(context.Background(), key, value))
		keys = append(keys, key)
	}
	return keys
}

const simpleEntryProcessorFactoryID = 66
const simpleEntryProcessorClassID = 1

type SimpleEntryProcessor struct {
	value string
}

func (s SimpleEntryProcessor) FactoryID() int32 {
	return simpleEntryProcessorFactoryID
}

func (s SimpleEntryProcessor) ClassID() int32 {
	return simpleEntryProcessorClassID
}

func (s SimpleEntryProcessor) WriteData(output serialization.DataOutput) {
	output.WriteString(s.value)
}

func (s *SimpleEntryProcessor) ReadData(input serialization.DataInput) {
	s.value = input.ReadString()
}

type SimpleEntryProcessorFactory struct {
}

func (f SimpleEntryProcessorFactory) Create(id int32) serialization.IdentifiedDataSerializable {
	if id == simpleEntryProcessorClassID {
		return &SimpleEntryProcessor{}
	}
	panic(fmt.Sprintf("unknown class ID: %d", id))
}

func (f SimpleEntryProcessorFactory) FactoryID() int32 {
	return simpleEntryProcessorFactoryID
}

type MapGetInterceptor struct {
	Prefix string
}

func (m MapGetInterceptor) FactoryID() int32 {
	return 666
}

func (m MapGetInterceptor) ClassID() int32 {
	return 6
}

func (m MapGetInterceptor) WriteData(output serialization.DataOutput) {
	output.WriteString(m.Prefix)
}

func (m *MapGetInterceptor) ReadData(input serialization.DataInput) {
	m.Prefix = input.ReadString()
}
