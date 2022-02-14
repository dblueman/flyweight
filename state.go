package flyweight

import (
   "os"

   "golang.org/x/sys/unix"
   "gopkg.in/yaml.v3"
)

type State struct {
   file    *os.File
   mapping []byte
   length  int64
   user    interface{}
}

func NewState(fname string, user interface{}, maxSize int) (*State, error) {
   s := State{user: user}

   var err error
   s.file, err = os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0o600)
   if err != nil {
      return nil, err
   }

   s.mapping, err = unix.Mmap(int(s.file.Fd()), 0, maxSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
   if err != nil {
      return nil, err
   }

   stat, err := s.file.Stat()
   if err != nil {
      return nil, err
   }

   s.length = stat.Size()
   if s.length > 0 {
      err = yaml.Unmarshal(s.mapping[:s.length], user)
      if err != nil {
         return nil, err
      }
   }

   return &s, nil
}

func (s *State) Capture() error {
   out, err := yaml.Marshal(s.user)
   if err != nil {
      return err
   }

   length := int64(len(out))

   if length != s.length {
      err = s.file.Truncate(length)
      if err != nil {
         return err
      }

      s.length = length
   }

   copy(s.mapping, out)

   return nil
}

func (s *State) Close() error {
   err := unix.Munmap(s.mapping)
   if err != nil {
      return err
   }

   err = s.file.Close()
   if err != nil {
      return err
   }

   return nil
}
