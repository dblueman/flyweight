package flyweight

import (
   "os"
   "testing"
)

type User struct {
   Counter int
}

const (
   filename = "test.yaml"
)

func Test(t *testing.T) {
   user := User{}

   os.Remove(filename)

   s, err := NewState(filename, &user, 8*1024*1024)
   if err != nil {
      t.Fatal("NewState:", err)
   }

   for i := 0; i < 100000; i++ {
      user.Counter = i

      err = s.Capture()
      if err != nil {
         t.Fatal("Capture:", err)
      }
   }

   localCounter := user.Counter

   err = s.Close()
   if err != nil {
      t.Fatal("Close:", err)
   }

   _, err = NewState(filename, &user, 8*1024*1024)
   if err != nil {
      t.Fatal("NewState:", err)
   }

   if user.Counter != localCounter {
      t.Fatal("counter disagree")
   }
}
