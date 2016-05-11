package main

import(
	"fmt"
	"net/http"
	"io/ioutil"
)

func HealthCheck() error {
	fmt.Println("health check.")
	req, err := http.NewRequest("GET", "http://192.168.13.11:8000/Health", nil)
	if err != nil {
		panic(err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	robots, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    panic(err.Error())
	}
	fmt.Printf("%s", robots)
	//fmt.Println("health check resp : ", resp)
	return err
}

func main() {
	fmt.Println("test main start.")
	HealthCheck()

}