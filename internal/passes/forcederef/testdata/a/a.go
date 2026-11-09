package a // want package:"module example.com/local"

import (
	"example.com/local/external"
	"example.com/local/local"
	"strconv"
)

func consume(value int)                     {}
func pointer(value *int)                    {}
func anyValue(value any)                    {}
func variadic(prefix string, values ...int) {}
func generic[T any](value T)                {}

type consumer struct{}

func (consumer) consume(value int)        {}
func (*consumer) pointerMethod(value int) {}

type accepts interface {
	consume(int)
}

type stringer interface {
	String() string
}

type valueStringer struct{}

func (valueStringer) String() string { return "" }

func interfaceValue(value stringer) {}

func calls(value *int, values *[]int, receiver consumer, dynamic accepts, text *string, stringValue *valueStringer) {
	consume(*value)                              // want "explicit dereference to call consume"
	consume((*value))                            // want "explicit dereference to call consume"
	local.Consume(*value)                        // want "explicit dereference to call Consume"
	receiver.consume(*value)                     // want "explicit dereference to call consume"
	receiver.pointerMethod(*value)               // want "explicit dereference to call pointerMethod"
	consumer.consume(receiver, *value)           // want "explicit dereference to call consume"
	(*consumer).pointerMethod(&receiver, *value) // want "explicit dereference to call pointerMethod"
	local.Consumer{}.Consume(*value)             // want "explicit dereference to call Consume"
	generic(*value)                              // want "explicit dereference to call generic"
	generic[int](*value)                         // want "explicit dereference to call generic"
	local.Generic(*value)                        // want "explicit dereference to call Generic"
	variadic("values", *value, *value)           // want "explicit dereference to call variadic" "explicit dereference to call variadic"
	variadic("values", *values...)               // want "explicit dereference to call variadic"
	go consume(*value)                           // want "explicit dereference to call consume"
	defer consume(*value)                        // want "explicit dereference to call consume"

	consume(0)
	pointer(value)
	anyValue(*value)
	generic[any](*value)
	interfaceValue(*stringValue)
	external.Consume(*value)
	external.Consumer{}.Consume(*value)
	strconv.Atoi(*text)
	println(*value)
	_ = int(*value)
	_ = *value + 1
	indirect := consume
	indirect(*value)
	dynamic.consume(*value) // want "explicit dereference to call consume"
}
