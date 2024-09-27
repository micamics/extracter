## Data Extraction Services
Service for extracting data from a file. 

This service facilitates the upload, extraction of data and basic processing of data from files.

*NOTE: This currently only supports excel files.*

The following are the instructions to use the service.

1. Build the project.

```bash
    ./scripts/build.sh
```

2. Run the server. This will be hosted locally.

```bash
    ./scripts/serve.sh
```

3. In another terminal, send a request to the server.

```bash
    curl http://localhost:8080/file/ -F file=@path/to/file
```

4. The extracted data will be saved in a csv file named `extracted_data.csv`.
