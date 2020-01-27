package main
import(
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"io/ioutil"
	"context"
	"os"
	"os/signal"
	"net/http"
	//"reflect"
	"strconv"
	"strings"
	// "encoding/json"
)
func main(){
	handlefunc()
	ctx:=context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()
}
func handlefunc(){
	myRouter:=mux.NewRouter()
	myRouter.HandleFunc("/{c1}/{c2}",index).Methods("GET")
	//myRouter.HandleFunc("/a/b/c",apicall).Methods("GET")
	log.Println("Server started on http://localhost:3000")
	http.ListenAndServe(":3000",myRouter)
}
func index(w http.ResponseWriter,r *http.Request){
	//c1:=mux.Vars(r)["c1"]
	c2:=mux.Vars(r)["c2"]
	c3,_:=r.URL.Query()["amount"]
	l:=strings.Join(c3,"")
	c4,_:=strconv.ParseFloat(l,7)
	str:="http://api.currencylayer.com/live?access_key=74097efca3ec5a3c10c8a5865fe7847f&currencies="+c2+"&format=1"
	res,err:=http.Get(str)
	if err!=nil{log.Println(err)}else{
		result,_:=ioutil.ReadAll(res.Body)
		body:=string(result)
		bodyArray:=strings.Fields(body)
		status:=strings.Split(bodyArray[1],":")[1]
		value:=strings.Split(bodyArray[7],":")[1]
		if status=="true,"{
			finalVal,err:=strconv.ParseFloat(value,7)
			if err!=nil{log.Println(err);}else{fmt.Fprintf(w,"%v", finalVal*c4)}
		}else{}
	}
}
