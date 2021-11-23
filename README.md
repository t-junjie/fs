# fs

fs is a simple server program that provides support for the followng workflow.

**Data provider:**

- Upload a binary file containing a 2 dimensional array of integers in csv format to the server.
- Remove a previously uploaded binary file from the server.

**Data analyst:**

- Retrieve a list of all files available on the server.
- Retrieve the details of a specific file, such as the number of rows and columns.
- Accept a request to compute the summation of a list of integers given by the analyst. Each integer in the list is uniquely identifiable by a resource identifier, row and column index. The results of the computation must be made accessible to the analyst.

## Install

To be updated

## How to run and test

- Download and install [Postman](https://www.postman.com/downloads/) for testing.
- Run the server and initiate an appropriate request to http://localhost:8080/files via Postman.
- Test data is available in the data/ directory.

### Upload File

- Send a POST request to http://localhost:8080/files.
- The request body must be of the type multipart/form-data, where "file" is the key and "data/test.csv" is the value, the binary file to be sent.

<img src="./data/Uploading.png"  width="600" height="250">

### Test Cases

| Checked | Request Type |               Endpoint               |                                                                 Request Body                                                                  | Status Code |                                                         Response Body                                                         | Remarks                                                                                                                         |
| :-----: | :----------: | :----------------------------------: | :-------------------------------------------------------------------------------------------------------------------------------------------: | :---------: | :---------------------------------------------------------------------------------------------------------------------------: | ------------------------------------------------------------------------------------------------------------------------------- |
|  `[x]`  |     GET      |     http://localhost:8080/files/     |                                                                     empty                                                                     |     200     |  `{"files":[{"id": "cb6a5e12-d582-46df-94d8-97e3cfa64006", "name": "test.csv","size": "30.4 kB","rows": 100,"cols": 104}]}`   | if no files are present, files will be an empty array                                                                           |
|  `[x]`  |     POST     |     http://localhost:8080/files/     |                                                     multipart/form-data (with .csv file)                                                      |     200     | `{"files": {"id": "311ed6fe-0374-4bb5-9d4e-1b3166189a81", "name": "test.csv", "size": "30.4 kB", "rows": 100, "cols": 104 }}` |                                                                                                                                 |
|  `[x]`  |     POST     |     http://localhost:8080/files/     |                                                   multipart/form-data (other file formats)                                                    |   400/500   |                                             `{"error": "unable to process file"}`                                             |                                                                                                                                 |
|  `[x]`  |     GET      | http://localhost:8080/files/{fileID} |                                                                     empty                                                                     |     200     | `{"files": {"id": "311ed6fe-0374-4bb5-9d4e-1b3166189a81", "name": "test.csv", "size": "30.4 kB", "rows": 100, "cols": 104 }}` |                                                                                                                                 |
|  `[x]`  |     GET      | http://localhost:8080/files/{fileID} |                                                                     empty                                                                     |     400     |                                                  `{"error": "invalid uuid}`                                                   |                                                                                                                                 |
|  `[x]`  |     GET      | http://localhost:8080/files/{fileID} |                                                                     empty                                                                     |     404     |                                                 `{"error": "file not found"}`                                                 |                                                                                                                                 |
|  `[x]`  |    DELETE    | http://localhost:8080/files/{fileID} |                                                                     empty                                                                     |     200     | `{"files": {"id": "311ed6fe-0374-4bb5-9d4e-1b3166189a81", "name": "test.csv", "size": "30.4 kB", "rows": 100, "cols": 104 }}` |                                                                                                                                 |
|  `[x]`  |    DELETE    | http://localhost:8080/files/{fileID} |                                                                     empty                                                                     |     400     |                                                  `{"error": "invalid uuid}`                                                   |                                                                                                                                 |
|  `[x]`  |    DELETE    | http://localhost:8080/files/{fileID} |                                                                     empty                                                                     |     404     |                                                 `{"error": "file not found"}`                                                 |                                                                                                                                 |
|  `[x]`  |     POST     |   http://localhost:8080/files/sum    | raw/json : `[{"uuid":"3a3192fd-d5b0-468a-bbf1-a066b9f1b679","row":1,"col":2},{"uuid":"3a3192fd-d5b0-468a-bbf1-a066b9f1b63","row":4,"col":5}]` |     200     |                                                         `{"sum": 10}`                                                         | `{uuid,row,col}` represents a value from cell(row,col) in the file with fileID=uuid. Note that row, col are zero-based indices. |
|  `[x]`  |     POST     |   http://localhost:8080/files/sum    | raw/json : `[{"uuid":"3a3192fd-d5b0-468a-bbf1-a066b9f1b679","row":1,"col":2},{"uuid":"3a3192fd-d5b0-468a-bbf1-a066b9f1b63","row":4,"col":5}]` |     400     |                                                             `{}`                                                              | Invalid `{uuid,row,col}` specified in request body or wrong request content type                                                |
