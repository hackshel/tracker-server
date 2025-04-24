package main

import (

  "fmt"
  "strconv"
)

func main(){
  s := "12"

  d, err := strconv.ParseInt( s, 10, 64)
  if err != nil{
     fmt.Print("errors %v", err)
  }

  fmt.Printf("value id %v, type %T", d, d)

}
