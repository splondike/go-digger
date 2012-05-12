package main

// The built in one wants float64s
func min(one int, two int) (rtn int) {
   if one < two {
      rtn = one
   } else {
      rtn = two
   }

   return
}
func max(one int, two int) (rtn int) {
   if one > two {
      rtn = one
   } else {
      rtn = two
   }

   return
}
