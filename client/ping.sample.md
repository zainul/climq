
## ping
- FieldsRequest : ping(string)
- FieldsResponse : pong(string)

### Sample : 
- Request : ping,test_ping
- Result: [null,"pong"]

### Note :

- Request : 

before send format the string to multiple of 16 bit and if the string is below the value please added ||||

until match with multiple of 16bit 

for example test_ping|||||||

test_ping actually < 16 so it will be test_ping|||||||

and then encrypt like aes cbc format and the encrypted message will be

�`�w����{U��n��������3�d���n

- Response :

Response from server will be [null,"pong"]

first array is describe of error if null it mean success

Sample if error ["UCU00000","UnexpectedError","Sample of error",null]

3 index is described the error

1. "UCU00000" => error code
2. "UnexpectedError" => error message
3. "Sample of error" => error detail