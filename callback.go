package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type CheckStruct struct {
	Val   int     `json:"val"`
	Str   string  `json:"str"`
	Float float64 `json:"float"`
}

func unmarshalJSON(src js.Value, dst interface{}) {
	str := js.Global().Get("JSON").Call("stringify", src).String()
	err := json.Unmarshal([]byte(str), &dst)
	if err != nil {
		fmt.Println(err)
	}
}

func CheckJSON(this js.Value, args []js.Value) interface{} {
	var checkStruct CheckStruct
	unmarshalJSON(args[0], &checkStruct)
	fmt.Println(checkStruct)
	return nil
}

type TestStruct struct {
	Val      int    `json:"val"`
	Callback string `json:"callback"`
}

func unmarshalJSONwithCallback(src js.Value, dst interface{}) {
	replacerString := `
		if (typeof v === 'function') {
			return v.toString();
		}
		return v;
	`
	replacer := js.Global().Get("Function").New(js.ValueOf("k"), js.ValueOf("v"), js.ValueOf(replacerString))
	str := js.Global().Get("JSON").Call("stringify", src, replacer).String()
	err := json.Unmarshal([]byte(str), &dst)
	if err != nil {
		fmt.Println(err)
	}
}

func getJSFuncFromString(function string) js.Value {
	this := js.Global().Get("this")
	return js.Global().Get("Function").Call("call", this, js.ValueOf("return"+function)).Invoke()
}

func Callback(this js.Value, args []js.Value) interface{} {
	var config TestStruct
	unmarshalJSONwithCallback(args[0], &config)
	callback := getJSFuncFromString(config.Callback)
	callback.Invoke(js.ValueOf(config.Val * 2))
	return nil
}

func registerCallbacks() {
	js.Global().Set("Callback", js.FuncOf(Callback))
	js.Global().Set("CheckJson", js.FuncOf(CheckJSON))
}

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}
