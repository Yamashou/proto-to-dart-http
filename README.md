# proto-to-dart-http

this repository is a command line tool to generate http api client from protofile.


# Usage
```
proto-to-dart-http \
 -i "proto/" \ 
 -o . 
 -p project-name \ 
 -pp /proto/v1/ \
  api/v1/*.proto

```


- *-i* : protfile root path
- *-o* : generated file path
- *-p* : project name
- *-pp*: dart project path



# Example

There are example in `example/`. 

example.proto
```example.proto

syntax = "proto3";

import "google/api/annotations.proto";

service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse){
        option (google.api.http) = {
       post: "/UserService/GetUser"
       body: "user"
       additional_bindings {
         post: "/UserService/GetUser2"
         body: "*"
       }
       additional_bindings {
         get: "/UserService/GetUser"
       }
    };
    };
}

message GetUserRequest {
    User user = 1;
}

message GetUserResponse {}

message User {
    string user_id= 1;
    string user_name = 2;
}
```


```
$ cd example
$ proto-to-dart-http -i "./,$GOPATH/src/github.com/googleapis/googleapis" -o . -p example -pp /proto/v1/ *.proto
```

```
import 'package:http/http.dart' as http;
import 'package:example/proto/v1/test.pb.dart'; // this path is decided path of `-pp`
class ExampleClient {
	String baseUrl;
	ExampleClient(String baseUrl) {this.baseUrl = baseUrl;}
	Future<GetUserResponse> getUser(GetUserRequest body, Map<String, String> headers) async {
		final response = await http.post(
			this.baseUrl + "/UserService/GetUser",
			body: body.writeToBuffer(),
			headers: headers);

		final GetUserResponse res = GetUserResponse.fromBuffer(response.bodyBytes);
		return res;
	}

	Future<GetUserResponse> getUser(GetUserRequest body, Map<String, String> headers) async {
		final response = await http.post(
			this.baseUrl + "/UserService/GetUser2",
			body: body.writeToBuffer(),
			headers: headers);

		final GetUserResponse res = GetUserResponse.fromBuffer(response.bodyBytes);
		return res;
	}

	Future<GetUserResponse> getUser(GetUserRequest body, Map<String, String> headers) async {
		final response = await http.get(
			this.baseUrl + "/UserService/GetUser",
			body: body.writeToBuffer(),
			headers: headers);

		final GetUserResponse res = GetUserResponse.fromBuffer(response.bodyBytes);
		return res;
	}

}
```


