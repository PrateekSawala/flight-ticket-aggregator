## Overview

Flight ticket upgrader processes the flight Records to send an email, offering a discount on upgrade to a higher class, to all the passengers who have booked tickets on its flights.

## API

Below is a list of API endpoint with their respective input and output.

### Upload Flight-Records

The API processes the input file sent via the request over the endpoint. Then write it to the passed record file. The records that fail the validation, get stored in the failed record file. So, that someone can look at them and fix the problem. Both files get stored inside the uploads folder.

### Flight Record schema 
```
First_name, Last_name, PNR, Fare_class, Travel_date, Pax, Ticketing_date, Email, Mobile_phone, Booked_cabin
```

#### Endpoint

```
POST
/upload/flightRecord
```

#### Input
1. Header.
    - Content-Type = multipart/form-data 
2. Body.
    - formKey = flightRecord
    - formValue =  file in .CSV format  
    
#### Output

```json
[
    {
        "Upload":   "File_Upload_Status",
        "Filename": "Uploaded_Filename",
        "Records": {
            "PassedRecordFileName": "Passed_Record_FileName",
            "FailedRecordFileName": "Failed_Record_FileName",
            "PassedRecordFilePathUrl": "Passed_Record_FilePath_Url",
            "FailedRecordFilePathUrl": "Failed_Record_FilePath_Url"
        }
    },
]
```

`NOTE: Only CSV file formats are allowed for the file upload, File Upload Size limit has been set to 5 MB, The value can be adjusted via constant attribute UploadFileSizeLimit declared inside constant.go file.`

## Useful make commands

###### Run the linter for checking stylistic errors
- make lint

###### Integration-test the application 
- make test
     
###### Compile the application 
- make build

###### Create a docker image to package the application
- make image

###### Run the application 
- make run

###### Execute all operations
- make

## FileWatcher

The File watcher has been setupped on the import folder to process the flight records manually via the application which checks for the newly available CSV files in the import directory and processes them. Then writes the passed entries to the passed record file. The records that fail the validation, get stored as a failed record file. Both files get stored inside the uploads folder.

## GUI

After Starting the server via the **make run** command the GUI will be available over the path [Url](http://localhost:3002), A simple GUI accepting the files to be uploaded and sending those files over the API endpoint.

## Containerization

The docker file has been added with the application along with the docker-compose file, The docker image can be created using the **make image** command which later can get used by the docker-compose file to serve the application over the path [Url](http://localhost:3002) for the **GUI** access.


## Integration test

The Integration tests have been added with the application inside the service folder for testing out the upload of a file with a test flight-record **flightRecord.csv** present inside the template folder and compare the results obtained via upload with the expected results, The integration-tests also contains the tests over validation checks for detecting the correctness of the file record entries. The integration test can be executed using the **make test** command.