# sesam-firestoresink
Sesam.io sink to GCP cloud firestore

System setup
```json
{
  "_id": "test-firestore",
  "type": "system:microservice",
  "docker": {
    "environment": {
      "GCP_PROJECT_ID": "sesam-228011",
      "GOOGLE_APPLICATION_CREDENTIALS": "credentials.json",
      "GOOGLE_APPLICATION_CREDENTIALS_CONTENT": {
        //GCP credentials as json or as string
      }
    },
    "image": "ohuenno/firestoresink",
    "port": 8080
  },
  "verify_ssl": true
}

```

Pipe setup 
```json
{
  "_id": "test-pipe-firestore-sink",
  "type": "pipe",
  "source": {
    "type": "embedded",
    "entities": [{
      "_id": 1,
      "key": "value"
    }, {
      "_id": 2,
      "key2": "value2"
    }, {
      "_id": 3,
      "key3": "value3"
    }]
  },
  "sink": {
    "type": "json",
    "system": "test-firestore",
    "url": "/test-collection-a"
  },
  "pump": {
    "cron_expression": "0 0 1 1 ?"
  }
}

```
