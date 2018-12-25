# Gbeta2

Based on [httprouter](https://github.com/julienschmidt/httprouter)


Example code:

```go
package main

import (
	"fmt"
	"time"
	"net/url"
	"net/http"
	"io/ioutil"
	"github.com/yyrdl/gbeta2"
	"github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
)


func Authed(w *gbeta2.Res, r *http.Request, ctx  *gbeta2.Ctx,next gbeta2.Next){
	token:= r.Header.Get("token")
	if token != "ok"{
		w.Write([]byte(`{"success":"Access denied"}`))
	}else{
		next()
	}
}


func Timecost(handle gbeta2.Handler) gbeta2.Handler{
      return func(w *gbeta2.Res, r *http.Request ,ctx  *gbeta2.Ctx,next gbeta2.Next){
		   now:= time.Now().UnixNano()
		   handle(w,r,ctx,next)
		   fmt.Println("time cost",time.Now().UnixNano()-now)
	  }
}

func JSON_Parser(w *gbeta2.Res, r *http.Request ,ctx  *gbeta2.Ctx,next gbeta2.Next){
	 body,err:= ioutil.ReadAll(r.Body)
	 
	 if nil != err{
		 w.Write([]byte(`{"success":false,"msg":"Internal error"}`))
	 }else{
		 var json = jsoniter.ConfigCompatibleWithStandardLibrary
		 ctx.Set("body" ,json.Get([]byte(body)))
		 next()
	 }
}

func Query_Parser(w *gbeta2.Res, r *http.Request ,ctx  *gbeta2.Ctx,next gbeta2.Next){
	query,err:= url.ParseQuery(r.URL.RawQuery)
	if nil != err{
		w.WriteHeader(500)
		w.Write([]byte(`{"success":false,"msg":"Failed to parse query:`+r.URL.Path+`"}`))
	}else{
		ctx.Set("query",query)
		next()
	}
}

func Name(w *gbeta2.Res, r *http.Request ,ctx *gbeta2.Ctx,next gbeta2.Next){

    params:= ctx.Get("params").(httprouter.Params)
    
    body:= ctx.Get("body").(jsoniter.Any)
    
	w.Write([]byte(`{"success":true,"message":"`+params.ByName("name")+`","age":"`+body.Get("age").ToString()+`"}`))
}

func Age(w *gbeta2.Res, r *http.Request ,ctx *gbeta2.Ctx,next gbeta2.Next){
    w.Write([]byte(`{"success":true,"age":18}`))
}

func UserRouter()*gbeta2.Router{
	router := gbeta2.New()
	router.GET("/age",Age)
	return router
}

func  main()  {
	
    router := gbeta2.New()
    
	router.Mw(Timecost)

	router.Use("/",Query_Parser)

	router.Use("/",Authed) // token is required

	router.POST("/hello/:name",JSON_Parser,Name) 

    router.SubRouter("/info",UserRouter)
	

	http.ListenAndServe(":8080",router.Build())

}
```