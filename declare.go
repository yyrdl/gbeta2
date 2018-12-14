package gbeta2


import (
	"net/http"
	"strings"
	"github.com/julienschmidt/httprouter"
)

var (
	
	_T_Middleware int = 0
	_T_HTTP_Handler int = 1
	_T_Filter int = 2
)

type Filter func (w http.ResponseWriter, r *http.Request, ctx map[string]interface{}) (bool) 

type Handler func (w http.ResponseWriter, r *http.Request ,ctx map[string]interface{})

type Middleware func(handle Handler) Handler

type pkg struct {
	p_type int
	path string
	method string
	mw Middleware
	filter Filter
	handle Handler
}

type  Router struct{
	pkg []*pkg
}

func (r *Router )Use(args ... interface{})*Router{
    if len(args) == 0{
		return r
	}

	if len(args) == 1{
      if mw,ok:= args[0].(Middleware);ok{
		 _pkg:= new(pkg)
		 _pkg.p_type = _T_Middleware
		 _pkg.mw = mw
		 r.pkg = append(r.pkg,_pkg)
	  }
	  if fit,ok:= args[0].(Filter);ok{
		  _pkg:= new(pkg)
		  _pkg.p_type = _T_Filter
		  _pkg.filter = fit
		  r.pkg = append(r.pkg,_pkg)
	  }
	  return r
	}

	if len(args) == 2{
		var path string
		var sub_router *Router
		var ok bool = true
		path,ok = args[0].(string)
		if ok == false{
			panic("The type of path must be string")
		}
		sub_router ,ok= args[1].(*Router)

		if ok == false{
			panic("The type of sub_router must be gebeta2.Router")
		}

		r.link(path,sub_router)

	}

	return  r;
}

func mergePath(p1,p2 string) string{
    if p1 == ""{
		return p2
	}
	if p2 == ""{
		return p1
	}
	if p1[len(p1)-1:] == "/" ||  p1[len(p1)-1:] == "\\"{
		p1 = p1[0:len(p1)-1]
	}

	if p1[len(p1)-1:] == "/" ||  p1[len(p1)-1:] == "\\"{
		p1 = p1[0:len(p1)-1]
	}
	return p1+"/"+p2
}
// link subrouter

func (r*Router)link(path string,router *Router){
	for i:=0;i<len(router.pkg);i++{
		if router.pkg[i].p_type != _T_Middleware{
		   router.pkg[i].path = mergePath(path,router.pkg[i].path)
		}
		r.pkg = append(r.pkg,router.pkg[i])
	}
}

func (r *Router)handle(method,path string ,args ... interface{}){
	if len(args) == 0{
		panic("Http handler is required!")
	}

    _pkg:= new(pkg)
	_pkg.p_type = _T_HTTP_Handler
	_pkg.path = path
	_pkg.method = method

	var filters []Filter
	var ft Filter
	var _handle Handler
	var ok bool

    // type check
	if len(args) > 1{
		filters  = make([]Filter, len(args)-1)
		for i:=0;i<len(args) -1 ;i++{
			ft,ok = args[i].(Filter)
			if ok == false{
				panic("The args before http handler must be Filter")
			}
			filters[i] = ft
		}
	}

	_handle ,ok =  args[len(args) -1].(Handler)
	

	if ok ==false{
		panic("The last argument must be a http handler!")
	}

	if len(args) > 1{
		_pkg.handle = func (w http.ResponseWriter, r *http.Request ,ctx map[string]interface{}){

			 var next bool = true

			 for i:=0;i<len(filters);i++{
				 next = filters[i](w,r,ctx)
				 if next == false{
					 break
				 }
			 }

			 if true == next{
				 _handle(w,r,ctx)
			 }

		}
	}else{
		_pkg.handle = _handle
	}


	r.pkg = append(r.pkg,_pkg)
}


func (r*Router)POST(path string ,args ... interface{}){
	r.handle("post",path,args...)
}

func (r*Router)GET(path string ,args ... interface{}){
	r.handle("get",path,args...)
}

func (r*Router)HEAD(path string ,args ... interface{}){
	r.handle("head",path,args...)
}

func (r*Router)OPTIONS(path string ,args ... interface{}){
	r.handle("option",path,args...)
}

func (r*Router)PUT(path string ,args ... interface{}){
	r.handle("put",path,args...)
}

func (r*Router)PATCH(path string ,args ... interface{}){
	r.handle("patch",path,args...)
}

func (r*Router)DELETE(path string ,args ... interface{}){
	r.handle("delete",path,args...)
}


func pathMatch(src, dst string) bool{
	if len(src) < len(dst){
		return src == dst[0:len(src)]
	}
	return src[0:len(dst)] == dst
}

func linkMiddlewareAndFilter(pkgs []*pkg,end int, path string ,handle Handler)Handler{
    for i:=0;i<end;i++{
		if pkgs[i].p_type == _T_Middleware {
			handle = pkgs[i].mw(handle)
		}
		if pkgs[i].p_type == _T_Filter {
            if pathMatch(pkgs[i].path,path){
				var _handle Handler = handle
				handle = func(ft Filter,hd Handler)Handler{
                       return func(w http.ResponseWriter, r *http.Request ,ctx map[string]interface{}){
							next:= ft(w,r,ctx)
							if true == next{
                                 hd(w,r,ctx)
							}
					   }
				}(pkgs[i].filter,_handle)
			}
		}
	}
	return handle
}


func (r*Router) Build()*httprouter.Router{

	 router:= httprouter.New()

	
	 for i:=0;i<len(r.pkg);i++{
         if r.pkg[i].path != ""{
			 r.pkg[i].path = strings.Replace(r.pkg[i].path,"\\","/",-1)
		 }
		
		 if r.pkg[i].p_type == _T_HTTP_Handler {
			 handle := linkMiddlewareAndFilter(r.pkg,i,r.pkg[i].path,r.pkg[i].handle)

			 md := r.pkg[i].method
			 pt := r.pkg[i].path

			 if "get" == md {
				 router.GET(pt,handle)
			 }else if "post" == md{
				 router.POST(pt,handle)
			 }else if "put" == md{
				 router.PUT(pt,handle)
			 }
		 }

	 }

	 return router
}


func New()*Router{
	return new(Router)
}