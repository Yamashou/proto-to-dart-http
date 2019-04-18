import 'package:http/http.dart' as http;
import 'package:example/proto/v1/test.pb.dart';
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
