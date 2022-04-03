# simple-search-es

> A simple search using Golang and Elasticsearch


## How to setup:
- Clone project
- Create **.env** file following **.env.example**
- Run `docker-compose up -d` to set up Elasticsearch
- Run `go run main.go`

## API
**Create an advert**
```
curl --location --request POST 'http://localhost:4041/api/ads' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Gree sofa",
    "text": "A beautiful green sofa",
    "tags": ["furniture"]
}'
```


**Delete an advert**
```
curl --location --request DELETE 'http://localhost:4041/api/ads/{id}'
```

**Find one advert**
```
curl --location --request GET 'http://localhost:4041/api/ads/{id}'
```

**Search adverts**

- Return all
```
curl --location --request GET 'http://localhost:4040/api/ads'
```

- Search by specific properties
```
curl --location --request GET 'http://localhost:4040/api/ads?title=Green sofa'
```

- Return results that match a provided text (...In progress...)
```
curl --location --request GET 'http://localhost:4040/api/ads?search=Green'
```