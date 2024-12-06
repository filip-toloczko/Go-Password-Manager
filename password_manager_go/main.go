package main

import (
	"fmt"
	"os"
  "bufio"
  "strings"
)

//Declaration for type Entry
type Entry struct{
  website string
  user string
  password string
}

//Type declaration for EntrySlice
type EntrySlice []Entry

//Declaration for passwordMap
var passwordMap = make(map[string]EntrySlice)

//_______________________________________________________________________
// initialize before main()
//_______________________________________________________________________
func init () {
  //Initialize passwordMap
  pmRead()
}

//_______________________________________________________________________
// print the list in columns
//_______________________________________________________________________
func pmList() {
  //Loop through passwordMap
  for _, value := range passwordMap{
    //Loop through each individual EntrySlice
    for i := 0; i < len(value); i++ {
      fmt.Printf("%40s %20s %20s\n", value[i].website, value[i].user, value[i].password)
    }
  }
}

//_______________________________________________________________________
//  add an entry if the site, user is not already found
//_______________________________________________________________________
func pmAdd(site, user, password string) {
  _, ok := passwordMap[site]
  //If the site is already in the map
  if ok {
    present := false
    //Iterate over elements of each slice
    for _, entry := range passwordMap[site] {
      //If the user already exists give an error
      if entry.user == user{
        fmt.Println("add: duplicate entry")
        present = true
      }
    }
    //Add the entry if the user is not present
    if present == false {
      passwordMap[site] = append(passwordMap[site], Entry{
        website: site,
        user: user,
        password: password,
      })
      pmWrite()
    }
  //Add the site if the site is not in the map
  }else{
    entry := Entry{website: site, user: user, password: password}
    passwordMap[site] = append(passwordMap[site], entry)
    pmWrite()
  }
}

//_______________________________________________________________________
// remove by site and user
//_______________________________________________________________________
func pmRemove(site, user string) {
  exists := false
  //Loop through the map
  for key, value := range passwordMap {
    if key == site {
      for i := 0; i < len(value); i++ {
        if value[i].user == user {
          //Remove the entry containing the user from the entryslice
          exists = true
          passwordMap[key] = append(value[:i], value[i+1:]...)
          pmWrite()
        }
      }
    }
  }
  //If the site does not exist print an error
  if exists == false{
    fmt.Println("remove: site not found")
  }
}

//_______________________________________________________________________
// remove the whole site if there is a single user at that site
//_______________________________________________________________________
func pmRemoveSite(site string) {
  //If the site has only one user, delete the site
  if len(passwordMap[site]) == 1{
    delete(passwordMap, site)
    pmWrite()
  //If the site does not exist print an error
  } else if len(passwordMap[site]) == 0 {
    fmt.Println("remove: site not found")
  //If the site had multiple users print an error
  } else{
    fmt.Println("remove: attempted to remove multiple users")
  }
}

//_______________________________________________________________________
// read the passwordVault
//_______________________________________________________________________
func pmRead() {
  //Open the file
  f, err := os.Open("passwordVault.txt")
  if err != nil {
     os.Create("passwordVault.txt")
  }
  //Create a scanner to read the file
  fileScanner := bufio.NewScanner(f)
  //Split the input into lines
  fileScanner.Split(bufio.ScanLines)

  var lines []string
  //Add each line to lines
  for fileScanner.Scan(){
    lines = append(lines, fileScanner.Text())
  }

  //Close the file
  f.Close()

  //Loop through the lines
  for _, line := range lines {
    if line != ""{
      //Add each string from the line to passwordMap
      parsedLine := strings.Split(line, " ")
      pmAdd(parsedLine[0], parsedLine[1], parsedLine[2])
    }
  }
}

//_______________________________________________________________________
// write the passwordVault
//_______________________________________________________________________
func pmWrite() {
  //Truncate the file
  f, err := os.Create("passwordVault.txt")
  if err != nil {
     panic(err)
  }
  //Close the file
  defer f.Close()

  //Loop through passwordMap
  for _, values := range passwordMap{
    //Loop through each entry in the slice 
    for _, value := range values{
      //Write the values to the file
      fmt.Fprintf(f, "%s %s %s\n", value.website, value.user, value.password)
    }    
  }
}

//_______________________________________________________________________
// do forever loop reading the following commands
//    l
//    a s u p
//    r s
//    r s u
//    x
//  where l,a,r,x are list, add, remove, and exit
//  and s,u,p are site, user, and password
//_______________________________________________________________________
func loop() {
  //Main loop that runs the program
  for{
    //Create a Reader to read the input
    reader := bufio.NewReader(os.Stdin)
    line, err := reader.ReadString('\n')
    if err != nil {
      panic(err)
    }
    //Remove whitespace from the input
    inputs := strings.Fields(strings.TrimSpace(line))

    //The length of the input is 1
    if len(inputs) == 1 {
      //The input is l
      if inputs[0] == "l"{
        pmList()
      //The input is x
      }else if inputs[0] == "x"{
        break
      //Print an error
      }else{
        fmt.Println("invalid command")
      }
    }
    //The length of the input is 2
    if len(inputs) == 2 {
      //The first input is r
      if inputs[0] == "r"{
        pmRemoveSite(inputs[1])
      //Print an error
      }else{
        fmt.Println("invalid command")
      }
    }
    //The length of the input is 3
    if len(inputs) == 3 {
      //The first input is r
      if inputs[0] == "r"{
        pmRemove(inputs[1], inputs[2])
      //Print an error
      }else{
        fmt.Println("invalid command")
      }
    }
    //The length of the input is 4
    if len(inputs) == 4 {
      //The first input is a
      if inputs[0] == "a"{
        pmAdd(inputs[1], inputs[2], inputs[3])
      //Print an error
      }else{
        fmt.Println("invalid command")
      }
    }
    //If none of the above applies print an error
    if len(inputs) > 4 || len(inputs) == 0 {
      fmt.Println("invalid command")
    }
  }
}

//_______________________________________________________________________
//  let her rip
//_______________________________________________________________________
func main() {
  loop()
}
