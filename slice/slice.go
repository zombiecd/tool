package slice

import "reflect"

func Contain[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
}

func ContainBy[T any](slice []T, predicate func(item T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}

	return false
}

func ContainSubSlice[T comparable](slice, subSlice []T) bool {
	for _, v := range subSlice {
		if !Contain(slice, v) {
			return false
		}
	}

	return true
}

func Chunk[T any](slice []T, size int) [][]T {
	result := [][]T{}

	if len(slice) == 0 || size <= 0 {
		return result
	}

	for _, item := range slice {
		l := len(result)
		if l == 0 || len(result[l-1]) == size {
			result = append(result, []T{})
			l++
		}

		result[l-1] = append(result[l-1], item)
	}

	return result
}

func Compact[T comparable](slice []T) []T {
	var zero T

	result := []T{}

	for _, v := range slice {
		if v != zero {
			result = append(result, v)
		}
	}

	return result
}

func Concat[T any](slice []T, slices ...[]T) []T {
	result := append([]T{}, slice...)

	for _, v := range slices {
		result = append(result, v...)
	}

	return result
}

func Difference[T comparable](slice, comparedSlice []T) []T {
	result := []T{}

	for _, v := range slice {
		if !Contain(comparedSlice, v) {
			result = append(result, v)
		}
	}

	return result
}

func DifferenceBy[T comparable](slice []T, comparedSlice []T, iteratee func(index int, item T) T) []T {
	orginSliceAfterMap := Map(slice, iteratee)
	comparedSliceAfterMap := Map(comparedSlice, iteratee)

	result := make([]T, 0)
	for i, v := range orginSliceAfterMap {
		if !Contain(comparedSliceAfterMap, v) {
			result = append(result, slice[i])
		}
	}

	return result
}

func Map[T any, U any](slice []T, iteratee func(index int, item T) U) []U {
	result := make([]U, len(slice), cap(slice))

	for i := 0; i < len(slice); i++ {
		result[i] = iteratee(i, slice[i])
	}

	return result
}

func Equal[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

func EqualWith[T, U any](slice1 []T, slice2 []U, comparator func(T, U) bool) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i, v := range slice1 {
		if !comparator(v, slice2[i]) {
			return false
		}
	}

	return true
}

func Every[T any](slice []T, predicate func(index int, item T) bool) bool {
	for i, v := range slice {
		if !predicate(i, v) {
			return false
		}
	}

	return true
}

func None[T any](slice []T, predicate func(index int, item T) bool) bool {
	l := 0
	for i, v := range slice {
		if !predicate(i, v) {
			l++
		}
	}

	return l == len(slice)
}

func Some[T any](slice []T, predicate func(index int, item T) bool) bool {
	for i, v := range slice {
		if predicate(i, v) {
			return true
		}
	}

	return false
}

func Filter[T any](slice []T, predicate func(index int, item T) bool) []T {
	result := make([]T, 0)

	for i, v := range slice {
		if predicate(i, v) {
			result = append(result, v)
		}
	}

	return result
}

func Count[T comparable](slice []T, item T) int {
	count := 0

	for _, v := range slice {
		if item == v {
			count++
		}
	}

	return count
}

func CountBy[T any](slice []T, predicate func(index int, item T) bool) int {
	count := 0

	for i, v := range slice {
		if predicate(i, v) {
			count++
		}
	}

	return count
}

func GroupBy[T any](slice []T, groupFn func(index int, item T) bool) ([]T, []T) {
	if len(slice) == 0 {
		return make([]T, 0), make([]T, 0)
	}

	groupB := make([]T, 0)
	groupA := make([]T, 0)

	for i, v := range slice {
		ok := groupFn(i, v)
		if ok {
			groupA = append(groupA, v)
		} else {
			groupB = append(groupB, v)
		}
	}

	return groupA, groupB
}

func GroupWith[T any, U comparable](slice []T, iteratee func(item T) U) map[U][]T {
	result := make(map[U][]T)

	for _, v := range slice {
		key := iteratee(v)
		if _, ok := result[key]; !ok {
			result[key] = []T{}
		}
		result[key] = append(result[key], v)
	}

	return result
}

func Find[T any](slice []T, predicate func(index int, item T) bool) (v T, ok bool) {
	index := -1

	for i, v := range slice {
		if predicate(i, v) {
			index = i
			break
		}
	}

	if index == -1 {
		return v, false
	}

	return slice[index], true
}

func FindLast[T any](slice []T, predicate func(index int, item T) bool) (v T, ok bool) {
	index := -1

	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(i, slice[i]) {
			index = i
			break
		}
	}

	if index == -1 {
		return v, false
	}

	return slice[index], true
}

func Flatten(slice any) any {
	sv := sliceValue(slice)

	var result reflect.Value
	if sv.Type().Elem().Kind() == reflect.Interface {
		result = reflect.MakeSlice(reflect.TypeOf([]interface{}{}), 0, sv.Len())
	} else if sv.Type().Elem().Kind() == reflect.Slice {
		result = reflect.MakeSlice(sv.Type().Elem(), 0, sv.Len())
	} else {
		return result
	}

	for i := 0; i < sv.Len(); i++ {
		item := reflect.ValueOf(sv.Index(i).Interface())
		if item.Kind() == reflect.Slice {
			for j := 0; j < item.Len(); j++ {
				result = reflect.Append(result, item.Index(j))
			}
		} else {
			result = reflect.Append(result, item)
		}
	}

	return result.Interface()
}

func sliceElemType(reflectType reflect.Type) reflect.Type {
	for {
		if reflectType.Kind() != reflect.Slice {
			return reflectType
		}

		reflectType = reflectType.Elem()
	}
}

func FlattenDeep(slice any) any {
	sv := sliceValue(slice)
	st := sliceElemType(sv.Type())

	tmp := reflect.MakeSlice(reflect.SliceOf(st), 0, 0)

	result := flattenRecursive(sv, tmp)

	return result.Interface()
}

func flattenRecursive(value reflect.Value, result reflect.Value) reflect.Value {
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		kind := item.Kind()

		if kind == reflect.Slice {
			result = flattenRecursive(item, result)
		} else {
			result = reflect.Append(result, item)
		}
	}

	return result
}

func ForEach[T any](slice []T, iteratee func(index int, item T)) {
	for i := 0; i < len(slice); i++ {
		iteratee(i, slice[i])
	}
}

func ForEachWithBreak[T any](slice []T, iteratee func(index int, item T) bool) {
	for i := 0; i < len(slice); i++ {
		if !iteratee(i, slice[i]) {
			break
		}
	}
}

func FilterMap[T any, U any](slice []T, iteratee func(index int, item T) (U, bool)) []U {
	result := []U{}

	for i, v := range slice {
		if a, ok := iteratee(i, v); ok {
			result = append(result, a)
		}
	}

	return result
}

func FlatMap[T any, U any](slice []T, iteratee func(index int, item T) []U) []U {
	result := make([]U, 0, len(slice))

	for i, v := range slice {
		result = append(result, iteratee(i, v)...)
	}

	return result
}

func Reduce[T any](slice []T, iteratee func(index int, item1, item2 T) T, initial T) T {
	if len(slice) == 0 {
		return initial
	}

	result := iteratee(0, initial, slice[0])

	tmp := make([]T, 2)
	for i := 1; i < len(slice); i++ {
		tmp[0] = result
		tmp[1] = slice[i]
		result = iteratee(i, tmp[0], tmp[1])
	}

	return result
}

func ReduceBy[T any, U any](slice []T, initial U, reducer func(index int, item T, agg U) U) U {
	accumulator := initial

	for i, v := range slice {
		accumulator = reducer(i, v, accumulator)
	}

	return accumulator
}

func ReduceRight[T any, U any](slice []T, initial U, reducer func(index int, item T, agg U) U) U {
	accumulator := initial

	for i := len(slice) - 1; i >= 0; i-- {
		accumulator = reducer(i, slice[i], accumulator)
	}

	return accumulator
}

func Replace[T comparable](slice []T, old T, new T, n int) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	for i := range result {
		if result[i] == old && n != 0 {
			result[i] = new
			n--
		}
	}

	return result
}

func ReplaceAll[T comparable](slice []T, old T, new T) []T {
	return Replace(slice, old, new, -1)
}

func Repeat[T any](item T, n int) []T {
	result := make([]T, n)

	for i := range result {
		result[i] = item
	}

	return result
}

func InterfaceSlice(slice any) []any {
	sv := sliceValue(slice)
	if sv.IsNil() {
		return nil
	}

	result := make([]any, sv.Len())
	for i := 0; i < sv.Len(); i++ {
		result[i] = sv.Index(i).Interface()
	}

	return result
}

func StringSlice(slice any) []string {
	v := sliceValue(slice)

	result := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		v, ok := v.Index(i).Interface().(string)
		if !ok {
			panic("invalid element type")
		}
		result[i] = v
	}

	return result
}

func IntSlice(slice any) []int {
	sv := sliceValue(slice)

	result := make([]int, sv.Len())
	for i := 0; i < sv.Len(); i++ {
		v, ok := sv.Index(i).Interface().(int)
		if !ok {
			panic("invalid element type")
		}
		result[i] = v
	}

	return result
}

func DeleteAt[T any](slice []T, start int, end ...int) []T {
	size := len(slice)

	if start < 0 || start >= size {
		return slice
	}

	if len(end) > 0 {
		end := end[0]
		if end <= start {
			return slice
		}
		if end > size {
			end = size
		}

		slice = append(slice[:start], slice[end:]...)
		return slice
	}

	if start == size-1 {
		slice = slice[:start]
	} else {
		slice = append(slice[:start], slice[start+1:]...)
	}

	return slice
}

func Drop[T any](slice []T, n int) []T {
	size := len(slice)

	if size <= n {
		return []T{}
	}

	if n <= 0 {
		return slice
	}

	result := make([]T, 0, size-n)

	return append(result, slice[n:]...)
}

func DropRight[T any](slice []T, n int) []T {
	size := len(slice)

	if size <= n {
		return []T{}
	}

	if n <= 0 {
		return slice
	}

	result := make([]T, 0, size-n)

	return append(result, slice[:size-n]...)
}

func DropWhile[T any](slice []T, predicate func(item T) bool) []T {
	i := 0

	for ; i < len(slice); i++ {
		if !predicate(slice[i]) {
			break
		}
	}

	result := make([]T, 0, len(slice)-i)

	return append(result, slice[i:]...)
}

func DropRightWhile[T any](slice []T, predicate func(item T) bool) []T {
	i := len(slice) - 1

	for ; i >= 0; i-- {
		if !predicate(slice[i]) {
			break
		}
	}

	result := make([]T, 0, i+1)

	return append(result, slice[:i+1]...)
}
