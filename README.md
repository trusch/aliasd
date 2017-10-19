aliasd
======

a microservice for managing aliases

## Interface

* `GET /alias/{scope}/{key}`
  get an alias
* `POST /alias/{scope}/{key}`
  create an random alias
* `PUT /alias/{scope}/{key} <data>`
  set an alias
* `DELETE /alias/{scope}/{key}`
  delete an alias
* `GET /alias/{scope}`
  get all aliases of a given scope
