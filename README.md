# ECC-API
A lightweight API with endpoints for doing operations using the bn256 elliptic curve.

## Installation
1. Install go 1.10+.
2. Make sure your `GOPATH` environment variable is pointing to an appropriate location
3. Get the repo along with its dependencies: `go get -d -u github.com/rynobey/ECC-API.git`
4. Compile by browsing to `$GOPATH/src/github.com/rynobey/ECC-API` and running `go build`

## Running and testing
1. Start the API server by browsing to `$GOPATH/src/github.com/rynobey/ECC-API` and running `./ECC-API`. The server runs on port 8083.
2. Once the API server is started, test it by running `GOCACHE=off go test -v .` (while in `$GOPATH/src/github.com/rynobey/ECC-API`)

## Routes
A list of available routes follows, with examples using `curl`.

### Routes implementing cryptographic algorithms
#### `/generate/commitment`
* Description: Generate Pedersen commitment: `result = v * g + b * h`, where `g` and `h` are ec curve points, `v` is the value being comitted to and `b` is the blinding factor    
* Method: `POST`  
*	Input: JSON object containing two integers, b and v, and two curve points, h and g, in hex: For ex. 
	```json
	{
	  "b":"0x010644e7fe131b029b85045b48181885d978163916871cffd3c208c16d87cfd3",
	  "v":"0x0ade",
	  "h":{
	    "x":"0x0d4826f08fe82224dfebd536358a1c0b3cd499b8dabec6e49abc37e78be1037a",
	    "y":"0x19e129957f1b471f2bb563bb32b3836412adbcc943362c896c143a47438aa518"
	  },
	  "g":{
	    "x":"0x0000000000000000000000000000000000000000000000000000000000000001",
	    "y":"0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd45"
	  }
	}
	```
* Output: JSON object containing the result in hex:  
	```json
	{
	  "curvepoint":{
	    "x":"0x086298450940b58d7f132a1765439e28658128fadcd790016466b3e67d4e2350",
	    "y":"0x00188d6bbfd592a7ef9e133fc1324c4f7a2eb89ccd6f70a4d1ccd304cbade81c"
	  }
	}
	```
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"b":"0x010644e7fe131b029b85045b48181885d978163916871cffd3c208c16d87cfd3","v":"0x0ade","h":{"x":"0x0d4826f08fe82224dfebd536358a1c0b3cd499b8dabec6e49abc37e78be1037a","y":"0x19e129957f1b471f2bb563bb32b3836412adbcc943362c896c143a47438aa518"},"g":{"x":"0x0000000000000000000000000000000000000000000000000000000000000001","y":"0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd45"}}' http://localhost:8083/generate/commitment/
	```
 
#### `/generate/keccak256/`
* Description: Generate hash of input text: `result = Hash(text)`  
* Method: `POST`  
*	Input: JSON object containing an input string to the hash function: For ex. `{"t":"Input to hash function"}`
* Output: JSON object containing the resulting hash in hex: For ex. 
	```
	{"number":{"v":"0xce567e7482a2e206d4337ff13f5c0b8e13bb5467844e2209d87e9f1484477799"}}
	```
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"t":"Input to hash function"}' http://localhost:8083/generate/keccak256/
	```

#### `/generate/schnorr/`
TODO

#### `/verify/schnorr/`
TODO

### Routes for math using elliptic curve points
#### `/ec/order`  
* Description: Returns bn256 EC order q: `result = q`  
* Method: `GET`  
* Output: JSON object containing the result in hex:  
	`{"number":{"v":"0x30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001"}}`
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request GET http://localhost:8083/ec/order
	```
	
#### `/ec/add/`  
* Description: Addition of two elliptic curve points: `result = a + b`  
* Method: `POST`  
*	Input: JSON object containing two curve points, a and b in hex: For ex. 
	```
	{"a":{"x":"0x030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd3","y":"0x1a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83"},"b":{"x":"0x0000000000000000000000000000000000000000000000000000000000000001","y":"0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd45"}}
	```
* Output: JSON object containing the resulting curve point in hex: For ex. 
	```
	{"curvepoint":{"x":"0x0769bf9ac56bea3ff40232bcb1b6bd159315d84715b8e679f2d355961915abf0","y":"0x05acb4b400e90c0063006a39f478f3e865e306dd5cd56f356e2e8cd8fe7edae6"}}
	```
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"a":{"x":"0x030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd3","y":"0x1a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83"},"b":{"x":"0x0000000000000000000000000000000000000000000000000000000000000001","y":"0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd45"}}' http://localhost:8083/ec/add/
	```
	
#### `/ec/sub/`  
* Description: Subtraction of one elliptic curve point from another: `result = a - b`  
* Method: `POST`  
*	Input: JSON object containing two curve points, a and b in hex: For ex. 
	```
	{"a":{"x":"0x030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd3","y":"0x1a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83"},"b":{"x":"0x0000000000000000000000000000000000000000000000000000000000000001","y":"0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd45"}}
	```
* Output: JSON object containing the resulting curve point in hex: For ex. 
	```
	{"curvepoint":{"x":"0x0000000000000000000000000000000000000000000000000000000000000001","y":"0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd45"}}
	```
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"a":{"x":"0x030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd3","y":"0x1a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83"},"b":{"x":"0x0000000000000000000000000000000000000000000000000000000000000001","y":"0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd45"}}' http://localhost:8083/ec/sub/
	```
	
#### `/ec/mul/`  
* Description: Multiplication of an elliptic curve point by a scalar: `result = s * a`  
* Method: `POST`  
*	Input: JSON object containing one integer, s, and one curve point, a, in hex: For ex. 
	```
	{"s":{"v":"0x01"},"a":{"x":"0x0000000000000000000000000000000000000000000000000000000000000001","y":"0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd45"}}
	```
* Output: JSON object containing the resulting curve point in hex: For ex. 
	```
	{"curvepoint":{"x":"0x030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd3","y":"0x1a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83"}}
	```
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"s":{"v":"0x02"},"a":{"x":"0x0000000000000000000000000000000000000000000000000000000000000001","y":"0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd45"}}' http://localhost:8083/ec/mul/
	```

#### `/ec/basemul/`  
* Description: Multiplication of the base generator (g) elliptic curve point by a scalar: `result = v * g`  
* Method: `POST`  
*	Input: JSON object containing one integer v in hex: For ex. `{"v":"0x02"}`
* Output: JSON object containing the resulting curve point in hex: For ex. 
	```
	{"curvepoint":{"x":"0x030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd3","y":"0x1a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83"}}
	```
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"v":"0x02"}' http://localhost:8083/ec/basemul/
	```

#### `/ec/hashtopoint/`  
* Description: Hash to a point on the elliptic curve (with unknown private key): `result = HashToPoint(text)`  
* Method: `POST`  
*	Input: JSON object containing an input string to the hash function: For ex. `{"t":"Input to hash function"}`
* Output: JSON object containing the resulting curve point in hex: For ex. 
	```
	{"curvepoint":{"x":"0x0d4826f08fe82224dfebd536358a1c0b3cd499b8dabec6e49abc37e78be1037a","y":"0x19e129957f1b471f2bb563bb32b3836412adbcc943362c896c143a47438aa518"}}
	```
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"t":"Input to hash function"}' http://localhost:8083/ec/hashtopoint/
	```

### Routes for math using big integers
#### `/big/add/`  
* Description: Addition of two big integers: `result = a + b`  
* Method: `POST`  
* Input: JSON object containing two numbers, a and b in hex: For ex. `{"a":"0x01", "b":"0x0adef342"}`  
* Output: JSON object containing the result in hex: For ex. `{"number":{"v":"0xadef343"}}`  
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"a":"0x01", "b":"0x0adef342"}' http://localhost:8083/big/add/
	```
	
#### `/big/submod/`  
* Description: Modular subtraction of one big integer from another: `result = (a - b) mod c`  
* Method: `POST`  
* Input: JSON object containing three numbers, a, b and c in hex: `{"a":"0x01", "b":"0x0adef342", "c":"0xffffffff}`  
* Output: JSON object containing the result in hex: `{"number":{"v":"0xf5210cbe"}}`  
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"a":"0x01", "b":"0x0adef342", "c":"0xffffffff"}' http://localhost:8083/big/submod/  
	```

#### `/big/mul/`  
* Description: Multiplication of two big integers: `result = a * b`  
* Method: `POST`  
* Input: JSON object containing two numbers, a and b in hex: `{"a":"0x01", "b":"0x0adef342"}`  
* Output: JSON object containing the result in hex: `{"number":{"v":"0xadef342"}}`  
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"a":"0x01", "b":"0x0adef342"}' http://localhost:8083/big/mul/
	```

#### `/big/mod/`  
* Description: Mod of a big integer: `result = a mod b`  
* Method: `POST`  
* Input: JSON object containing two numbers, a and b in hex: `{"a":"0x07", "b":"0x05"}`  
* Output: JSON object containing the result in hex: `{"number":{"v":"0x2"}}`  
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"a":"0x07", "b":"0x05"}' http://localhost:8083/big/mod/
	```

#### `/big/invmod/`  
* Description: Modular multiplicative inverse of a big integer: `(result * a) mod b = 1`  
* Method: `POST`  
* Input: JSON object containing two numbers, a and b in hex: `{"a":"0x07", "b":"0x05"}`  
* Output: JSON object containing the result in hex: `{"number":{"v":"0x3"}}`  
* Example usage: 
	```
	curl --header "Content-Type: application/json" --request POST --data '{"a":"0x07", "b":"0x05"}' http://localhost:8083/big/invmod/
	```
