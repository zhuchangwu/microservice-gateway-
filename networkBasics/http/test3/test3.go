package main

type Handle interface {
	Serve(string, int, string)
}

type HandleImpl struct {

}

func (h HandleImpl)Serve(string, int, string){

}

/*
func Entry(string,Handle){

}

func MyHandler(name string, age int, address string) {
	fmt.Println(name,age,address)
}*/

func main() {
	//Entry("123",MyHandler)

}
