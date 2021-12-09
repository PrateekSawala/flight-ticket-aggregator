# Flight ticket aggregator

## Overview

The project aims to provide a microservice solution to process flight records for sending out report emails to the respected airlines offering the discounts provided to users upgraded to a higher class who have booked tickets on their flights.

## API

Below is a list of API endpoint with their respective input and output.

### Flight-Record

The API processes the input file sent via the request over the endpoint. Then write it to the passed record file. The records that fail the validation, get stored in the failed record file. So, that someone can look at them and fix the problem. Both files get stored in AWS S3 bucket.

If the record entry is a valid entry the we need to add a new column called discount code to the passed records file whose value will be calculated based on the fare class field in the input record. Fare class A - E will have discount code OFFER_20, F - K will have discount code OFFER_30, L - R will have OFFER_25 and the rest will have no offer code.

### Flight-Record schema 

```
First_name, Last_name, PNR, Fare_class, Travel_date, Pax, Ticketing_date, Email, Mobile_phone, Booked_cabin
```

### Flight-Record validations 

- Email ID is valid
- The mobile phone is valid
- Ticketing date is before travel date
- PNR is six characters long, Is alphanumeric, and Is unique.
- The booked cabin is valid (one of Economy, Premium Economy, Business, First)


#### Endpoints

```
To upload flightRecord -
POST {Base_URL}/upload/flightRecord

#### Input Request 

1. Header.
    - Content-Type = multipart/form-data 
2. Body.
    - formKey = flightRecord
    - formValue =  file in .CSV format  

#### Output Response
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

```
To download flightRecord -
GET {Base_URL}/download/flightRecord/<FilePath>/<FileName>
```

`NOTE: Only CSV file formats are allowed for the file upload, File Upload Size limit has been set to 5 MB, The value can be adjusted via constant attribute UploadFileSizeLimit declared inside constant.go file.`


## Microservices

#### Websever Service
The webserver microservice is our gateway and link between frontend and backend services. It is the only microservice in direct connection with the website and the only way to process data. It gives functionality to our whole system.

#### Ticket Service
The ticket microservice is there to process the flight records. It includes the business logic for processing input flight records which can later be sent to the respective airline via email and stored in the AWS S3 bucket.

#### Space Service
The space microservice is our file management service. It is connected to AWS S3.

#### Mail Service
The mail microservice is there to send out emails.

## Useful make commands

##### Transpile protobuffer definitions for the complete application
- fta all build:protoc

##### Transpile protobuffer definitions for a specific service
- fta <service> build:protoc

##### Compile the complete application
- fta all build:build

##### Compile a specfic service
- fta <service> build:build

##### Create docker images for the complete application
- fta all docker:build

##### Create docker image for specfic service 
- fta <service> docker:build

##### unit-test the complete application, Unit-test the specfic service  
- fta <service> build:unittest

##### Unit-test a specfic service  
- fta <service> build:unittest

##### Integrtion test a specfic service  
- fta <service> build:test

## FileWatcher

The File watcher has been setupped on the import folder to process the flight records manually via the application which checks for the newly available CSV files in the import directory and processes them. Then writes the passed entries to the passed record file. The records that fail the validation, get stored as a failed record file. So, that someone can look at them and fix the problem. Both files get stored in AWS S3 bucket.

## GUI

After Starting the server the GUI will be available over the path [Url](http://localhost:3002), A simple GUI accepting the files to be uploaded and sending those files over the API endpoint.

## Containerization

The docker file has been added with the application along with the docker-compose file, The created docker images can be can  used by the docker-compose file to serve the application over the path [Url](http://localhost:3002) for the **GUI** access.

## Unit tests

The unit tests have been added with the application for performing unit tests on the business logic of the application.