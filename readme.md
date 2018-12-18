# Gbeta2

Based on [httprouter](https://github.com/julienschmidt/httprouter)


Example code:

```go
package main

import (
	"demo/gbeta2"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"time"
	"io/ioutil"
	"github.com/json-iterator/go"
)


func Authed(w http.ResponseWriter, r *http.Request, ctx map[string]interface{},next gbeta2.Next){
	token:= r.Header.Get("token")
	if token != "ok"{
		w.Write([]byte(`{"success":"Access denied"}`))
	}else{
		next()
	}
}


func Timecost(handle gbeta2.Handler) gbeta2.Handler{
      return func(w http.ResponseWriter, r *http.Request ,ctx map[string]interface{},next gbeta2.Next){
		   now:= time.Now().UnixNano()
		   handle(w,r,ctx,next)
		   fmt.Println("time cost",time.Now().UnixNano()-now)
	  }
}

func JSON_Parser(w http.ResponseWriter, r *http.Request ,ctx map[string]interface{},next gbeta2.Next){
	 body,err:= ioutil.ReadAll(r.Body)
	 
	 if nil != err{
		 w.Write([]byte(`{"success":false,"msg":"Internal error"}`))
	 }else{
		 var json = jsoniter.ConfigCompatibleWithStandardLibrary
		 ctx["body"] = json.Get([]byte(body))
		 next()
	 }
}

func Name(w http.ResponseWriter, r *http.Request ,ctx map[string]interface{},next gbeta2.Next){

    params:= ctx["params"].(httprouter.Params)
    
    body:= ctx["body"].(jsoniter.Any)
    
	w.Write([]byte(`{"success":true,"message":"`+params.ByName("name")+`","age":"`+body.Get("age").ToString()+`"}`))
}

func  main()  {
	
    router := gbeta2.New()
    
	router.Mw(Timecost)
	
	router.Use("/",Authed) // token is required

	router.POST("/:name",JSON_Parser,Name) 

	http.ListenAndServe(":8080",router.Build())

}
```