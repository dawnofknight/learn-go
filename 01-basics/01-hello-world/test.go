package main 

import
	("fmt")


func test(data, data2 int){
	fmt.Println("Variable data", data)
	fmt.Println("Variable data2", data2)
}

func calculate(number1, number2 int) int{
	return(number1-number2)
}

func main(){
	var data int
	data = 10
	data2 := 10
	test(data, data2)
	result := calculate(data, data2)
	fmt.Println("Hasil Penjumlahan", result)
}
