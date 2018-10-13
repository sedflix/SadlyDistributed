package main

import "os"
import "strconv"

// lb is included, ub isn't
func checkdivisibility(p, lb, ub int32) bool {
  for i := lb; i < ub; i++ {
    if (p % i) == 0 {
      return true
    }
  }
  return false
}

func main() {
  lb, _ := strconv.ParseInt(os.Args[1], 10, 32)
  ub, _ := strconv.ParseInt(os.Args[2], 10, 32)
  println(checkdivisibility(81, int32(lb), int32(ub)))
}
