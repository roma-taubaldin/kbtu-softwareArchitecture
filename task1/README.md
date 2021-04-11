# kbtu-softwareArchitecture

to build task1 you need:
1. pull this project to your pc
2. install latest go version
3. build this project with this params "CGO_ENABLED=0 go build main.go"
4. do "docker build ."
5. do "docker run -p 8000:8000" with image
6. you can test it with another terminal doing curl "curl localhost:8000/health"
7. it should return {"status": "OK"}
