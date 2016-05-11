package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "html/template"
    "io"
    "time"
    "crypto/md5"
    "strconv"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "log"
    "strings"
    "os"
)

const (
	URL = "192.168.13.14:27017"
)

type Person struct {
	User string  
	Pass string 
}

type JsonReturnHandler func(http.ResponseWriter, *http.Request) error
func (f JsonReturnHandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
	f(w, r)
}

type Routes []Route

type Route struct {
	Name string

	Method string
	Pattern string
	Handler http.Handler
}

var routes = Routes{ 
	Route{Name:"Images", Method:"GET", Pattern:"/list", Handler: JsonReturnHandler(ListImages)},
	Route{Name:"Images", Method:"GET", Pattern:"/pull", Handler: JsonReturnHandler(PullImage)},
	Route{Name:"login",  Method:"GET", Pattern:"/", Handler: JsonReturnHandler(Login)},
	Route{Name:"login",  Method:"POST", Pattern:"/login", Handler: JsonReturnHandler(LoginJoin)},
	Route{Name:"upload",  Method:"POST", Pattern:"/upload", Handler: JsonReturnHandler(Upload)},
	Route{Name:"upload",  Method:"GET", Pattern:"/Health", Handler: JsonReturnHandler(HealthCheck)},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		router.Methods(route.Method).
		       Path(route.Pattern).
		       Name(route.Name).
			   Handler(route.Handler)
	}
	return router
}

func ListImages(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("this is listImages func")
	w.Write([]byte("this is listImages func\n"))
	return nil
}

func PullImage(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("this is pullImages func")
	w.Write([]byte("this is pullImages func\n"))
	return nil
}

func Login(w http.ResponseWriter, r *http.Request) error{
	fmt.Println("this is Login func")
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime,10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	t,_ :=template.ParseFiles("login/login.gtpl")
	t.Execute(w, token)
	return nil
}

func LoginJoin(w http.ResponseWriter, r *http.Request) error{
	fmt.Println("this is LoginJoin func")
	r.ParseForm()
	fmt.Println("username : ", r.Form["username"])
	fmt.Println("password : ", r.Form["password"])
	user := ""
	pass := ""
	if len(r.Form["username"]) != 0 {
		user = r.Form["username"][0]
	}
	if len(r.Form["password"]) != 0 {
		pass = r.Form["password"][0]
	}
	if checklogin(user, pass) ==  true {
		fmt.Printf("login success.\n")
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime,10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		
		t,_ :=template.ParseFiles("upload/upload.gtpl")
		t.Execute(w, token)
	}else {
		fmt.Printf("login failed. \n")
	}
	return nil
}

func Upload(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("this is Upload func")
	r.ParseMultipartForm(1024*10)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./test/" + handler.Filename, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	io.Copy(f,file)
	return nil
}

func checklogin(user, pass string) bool {
	session, err := mgo.Dial(URL)
	if (err != nil) {
		panic(err)
		return false
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("test")
	collection := db.C("login")
	countNum, err := collection.Count()
	if err != nil {
		panic(err)
		return false
	}
	if countNum == 0 {
		err = collection.Insert(&Person{"admin","admin"})
		if err != nil {
			log.Fatal(err)
			return false
		}
	}
	fmt.Println("things objects count :", countNum)
	fmt.Printf("user : %s \n", user)
	doc := Person{}
	err = collection.Find(bson.M{"user":user}).One(&doc)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("username : %s \n", doc.User)
	fmt.Printf("password : %s \n", doc.Pass)
	if strings.EqualFold(user, doc.User) == true &&
	strings.EqualFold(pass, doc.Pass){
		return true
	}else {
		return false
	}
	return false
}


func HealthCheck(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Health check .")
	w.Write([]byte("health check success.\n"))
	return nil
}
/*func (p *MyMux)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		sayhelloName(w, r)
		return 
	}
	http.NotFound(w, r)
	return
}*/

/*func  sayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello myroute!")
}*/

/*func YourHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Gorilla!\n"))
}*/

func main() {
    //r := mux.NewRouter()
    // Routes consist of a path and a handler function.
    //r.HandleFunc("/", YourHandler)

   // mux := &MyMux{}
	router := NewRouter()

    // Bind to a port and pass our router in
    http.ListenAndServe(":8000", router)
}