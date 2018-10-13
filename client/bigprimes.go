package main

import "os"
import "math/big"

// lb is included, ub isn't
func checkdivisibility(p, lb, ub *big.Int) bool {
  z := new (big.Int)
  for i := new(big.Int).Set(lb); i.Cmp(ub) == -1; i.Add(i, big.NewInt(1)) {
    z.Mod(p, i)
    if z.Cmp(big.NewInt(0)) == 0 {
      return true
    }
  }
  return false
}

func main() {
  lb := new (big.Int)
  ub := new (big.Int)
  lb.SetString(os.Args[1], 10)
  ub.SetString(os.Args[2], 10)
  println(checkdivisibility(big.NewInt(81), lb, ub))
}

