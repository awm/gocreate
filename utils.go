package gocreate

import (
    "math"
    "reflect"
)

// Sequence generates a slice of regularly increasing (or decreasing) values.
//  values := gocreate.Sequence(0, 10, 2).([]int) // values is []int{0, 2, 4, 6, 8}
//  values = gocreate.Sequence(0, -10, -2).([]int) // values is []int{0, -2, -4, -6, -8}
//
// Note: The return type must be asserted as a slice of the parameter type.  For example:
//  values := gocreate.Sequence(uint32(0), uint32(10), uint32(1)).([]uint32) // values is []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
func Sequence(start interface{}, stop interface{}, step interface{}) interface{} {
    vstart := reflect.ValueOf(start)
    vstop := reflect.ValueOf(stop)
    vstep := reflect.ValueOf(step)

    if vstart.Kind() != vstop.Kind() || vstart.Kind() != vstep.Kind() {
        return nil
    }

    var test func(vi *reflect.Value) bool
    var incr func(vi *reflect.Value)
    var capacity int = 0
    switch vstart.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        test = func(vi *reflect.Value) bool {
            if vstep.Int() < 0 {
                return vi.Int() > vstop.Int()
            } else {
                return vi.Int() < vstop.Int()
            }
        }
        incr = func(vi *reflect.Value) {
            vi.SetInt(vi.Int() + vstep.Int())
        }

        c := (vstop.Int() - vstart.Int()) / vstep.Int()
        if c < 0 {
            c = -c
        }
        capacity = int(c)
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
        test = func(vi *reflect.Value) bool {
            return vi.Uint() < vstop.Uint()
        }
        incr = func(vi *reflect.Value) {
            vi.SetUint(vi.Uint() + vstep.Uint())
        }

        capacity = int((vstop.Uint() - vstart.Uint()) / vstep.Uint())
    case reflect.Float32, reflect.Float64:
        test = func(vi *reflect.Value) bool {
            if vstep.Float() < 0 {
                return vi.Float() > vstop.Float()
            } else {
                return vi.Float() < vstop.Float()
            }
        }
        incr = func(vi *reflect.Value) {
            vi.SetFloat(vi.Float() + vstep.Float())
        }

        c := (vstop.Float() - vstart.Float()) / vstep.Float()
        if c < 0.0 {
            c = -c
        }
        capacity = int(math.Ceil(c))
    default:
        return nil
    }
    tresult := reflect.SliceOf(vstart.Type())
    vresult := reflect.MakeSlice(tresult, 0, capacity)

    i := start
    var vi reflect.Value
    switch i := i.(type) {
    case int:
        vi = reflect.ValueOf(&i).Elem()
    case uint:
        vi = reflect.ValueOf(&i).Elem()
    case uintptr:
        vi = reflect.ValueOf(&i).Elem()
    case int8:
        vi = reflect.ValueOf(&i).Elem()
    case uint8:
        vi = reflect.ValueOf(&i).Elem()
    case int16:
        vi = reflect.ValueOf(&i).Elem()
    case uint16:
        vi = reflect.ValueOf(&i).Elem()
    case int32:
        vi = reflect.ValueOf(&i).Elem()
    case uint32:
        vi = reflect.ValueOf(&i).Elem()
    case int64:
        vi = reflect.ValueOf(&i).Elem()
    case uint64:
        vi = reflect.ValueOf(&i).Elem()
    case float32:
        vi = reflect.ValueOf(&i).Elem()
    case float64:
        vi = reflect.ValueOf(&i).Elem()
    default:
        return nil
    }

    for ; test(&vi); incr(&vi) {
        vresult = reflect.Append(vresult, vi)
    }

    return vresult.Interface()
}

func Sum(slice interface{}) interface{} {
    vslice := reflect.ValueOf(slice)
    if vslice.Kind() != reflect.Slice || vslice.Len() == 0 {
        return nil
    }

    titems := vslice.Index(0).Type()
    vresult := reflect.New(titems).Elem()

    var incr func(vi reflect.Value)
    switch titems.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        incr = func(vi reflect.Value) {
            vresult.SetInt(vresult.Int() + vi.Int())
        }
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
        incr = func(vi reflect.Value) {
            vresult.SetUint(vresult.Uint() + vi.Uint())
        }
    case reflect.Float32, reflect.Float64:
        incr = func(vi reflect.Value) {
            vresult.SetFloat(vresult.Float() + vi.Float())
        }
    default:
        return nil
    }

    for i := 0; i < vslice.Len(); i++ {
        vi := vslice.Index(i)
        if vi.Kind() != titems.Kind() {
            return nil
        }
        incr(vi)
    }

    return vresult.Interface()
}
