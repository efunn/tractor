package registry

import (
	"errors"
	"reflect"
	"sync"
)

type Registry struct {
	entries []*Entry

	mu sync.Mutex
}

func New() *Registry {
	return &Registry{}
}

func Ref(v interface{}) *Entry {
	// TODO: allow values instead of assuming pointer references
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		v = &v
	}
	return &Entry{
		Ref: v,
	}
}

func (r *Registry) Entries() []*Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	e := make([]*Entry, len(r.entries))
	copy(e, r.entries)
	return e
}

func (r *Registry) Register(entries ...*Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range entries {
		e.RefType = reflect.TypeOf(e.Ref)

		// if not a pointer, ignore
		if e.RefType.Kind() != reflect.Ptr {
			return errors.New("value reference must be a pointer")
		}

		e.Type = e.RefType.Elem()
		e.Value = reflect.ValueOf(e.Ref)
		e.PkgPath = e.Type.PkgPath()
		e.TypeName = e.Type.Name()

		// error if the object has no package path
		// if e.TypeName == "" && e.PkgPath == "" {
		// 	return errors.New("unable to register object without name when it has no package path")
		// }

		// append entry to registry list
		r.entries = append(r.entries, e)

	}
	return nil
}

func (r *Registry) AssignableTo(t reflect.Type) []*Entry {
	var entries []*Entry
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	for _, entry := range r.Entries() {
		if entry.RefType.AssignableTo(t) {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (r *Registry) Populate(v interface{}) {
	rv := reflect.ValueOf(v)
	// TODO: assert struct
	var fields []reflect.Value
	for i := 0; i < rv.Elem().NumField(); i++ {
		// filter out unexported fields
		if len(rv.Elem().Type().Field(i).PkgPath) > 0 {
			continue
		}
		fields = append(fields, rv.Elem().Field(i))
		// TODO: filtering with struct tags
	}
	for _, field := range fields {
		if !isNilOrZero(field, field.Type()) {
			continue
		}
		assignable := r.AssignableTo(field.Type())
		if len(assignable) == 0 {
			continue
		}
		if field.Type().Kind() == reflect.Slice {
			field.Set(reflect.MakeSlice(field.Type(), 0, len(assignable)))
			for _, entry := range assignable {
				field.Set(reflect.Append(field, entry.Value))
			}
		} else {
			field.Set(assignable[0].Value)
		}
	}
}

func isNilOrZero(v reflect.Value, t reflect.Type) bool {
	switch v.Kind() {
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(t).Interface())
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
}

// remember to use reflect.Indirect on rv after
func (r *Registry) ValueTo(rv reflect.Value) {
	for _, e := range r.Entries() {
		if rv.Elem().Type().Kind() == reflect.Struct {
			if e.Value.Elem().Type().AssignableTo(rv.Elem().Type()) {
				rv.Elem().Set(e.Value.Elem())
				return
			}
		} else {
			if e.Value.Type().Implements(rv.Elem().Type()) {
				rv.Elem().Set(e.Value)
				return
			}
		}
	}

}

type Entry struct {
	Ref      interface{}
	TypeName string
	PkgPath  string

	RefType reflect.Type
	Type    reflect.Type
	Value   reflect.Value
}
