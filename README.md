Author-Aditya Agarwal

# GO-TAMBOON ไปทำบุญ


Config file(Keys, Rate Limit etc)- src/custom/constants/constants.go, can create a yml file

Running the App
 1) go build app.go <fileName>
 OR
 2) go install . && ./go-tamboon <fileName>
 
 The Application uses clean architecture
 
 Start reading the k bytes from file at a time, decrypt user and push user struct to C channel of Cn length(src/repositries/user.go:70),
  when buffer is finished refill the buffer and push user to C chanel and so on
 
 The user channel is parsed in src/useCase/collect.go:24, for each user a goroutine is launched(src/useCase/collect.go:27)
  if tx fail, the goroutine try to retry the request and return result to Transaction channel
 
 The data is collected and processed in src/useCase/collect.go:38