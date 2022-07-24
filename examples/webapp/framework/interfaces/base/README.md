~~~bash
python -m grpc_tools.protoc --proto_path=/generated --python_out=/generated --grpc_python_out=/generated base_model.proto
~~~

~~~protobuf
syntax = "proto3";
package mlserving;

option java_package = "mlserving.model";
option java_multiple_files = true;


/*
 * Request
 */
message ModelInterfaceRequest {
  map<string, string> details = 1;
  bytes data = 2;
}


/*
 * Response
 */
message ModelInterfaceResponse {
  bytes data = 1;
}


/*
 * Service
 */
service BaseModelInterface {

  // Make prediction based on user request
  rpc Predict(ModelInterfaceRequest) returns (ModelInterfaceResponse);
}
~~~