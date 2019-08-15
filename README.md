Provides HTTP-API for resources that are received as events and are modified and merged according to the configuration. Resources will be saved in elastic search.
Can be used to create materialized views in a cqrs environment.


# Config - Events ('events')
The 'events' field maps event topics to a list of action-groups. If the matviev instance receives an event from the event-broker, each corresponding action-group will be called.

## Action-Group
A Action-Group conditionally transforms and saves a event to specified elasticsearch documents.
May contain the fields `type`, `target`, `where`, `if`, `features`, `actions`, `init` and `id_feature`.

### Type ('type')
Valid values are `"root"` and `"child"`.

Root will ignore the `where` field and identifies which resources should be edited by the `id_feature` field.

Child will ignore the `id_feature` and `init` fields and identifies resources by the `where` field.

### Target ('target')
Target identifies which elasticsearch index will be used to find and save resources. It corresponds to the Keys of the Mapping in the Queries-Config.

### Features ('features')
Describes how the event should be transformed to a new map (`map[string]interface{}`) the keys of this map will be referenced by `feature` or `event_feature`.
Features consists of a list of descriptions, where each entry describes one field. These descriptions contain the fields

* `name`: (string) name of the feature
* `path`: (string) json-path, used on the event to get the value of the field (https://github.com/JumboInteractiveLimited/jsonpath)
* `temp`: (bool) if true this feature will be used for if-conditions and where-conditions but will not be included when saved
* `omitempty`: (bool) if true the feature will not be included if the jsonpath-result and default is empty (null, empty list etc). if omitempty is false the feature will exists with a null value.
* `default`: (anything) will be used if jsonpath-result is null. will evaluated before omitempty
* `default_ref`: (string) equivalent to default but refers to values like the current time. Implemented references are:
    * `"time.epoch_millis"`: current time as unix timestamp in milliseconds
    * `"time.epoch_second"`: current time as unix timestamp in seconds 
    
**Example:**
```
{
    ...
    "features": [
        {"name": "command", "path": "$.command+", "temp": true},
        {"name": "user", "path": "$.User+", "omitempty": true},
        {"name": "group", "path": "$.Group+", "omitempty": true},
        {"name": "right", "path": "$.Right+", "temp": true},
        {"name": "kind", "path": "$.Kind+", "temp": true},
        {"name": "resource", "path": "$.Resource+"},
        {"name": "date", "default_ref": "time.epoch_millis"}
    ],
    ...
}
```

### If-Condition ('if')
If-Conditions can be used in Action-Groups and Actions. If one of the listed conditions fails the Action-Group or Action will not be executed. A If-Condition consists of the fields

* `feature`: (string) reference to a event-feature
* `operation`: (string) operation that will be executed to determine the result of the condition
* `value`: (anything) value on which the operation will be executed. The type of the value is determined by the operation and feature.

Currently valid operations are:

* `==`: checks equality of the feature and value (both must be of the same type and have the same value)
* `!=`: not `==`
* `feature_str_contains_value`: feature must be string; value must be string; feature must contain value as substring

**Example:**
```
{
    ...
    "if": [
        {"feature": "command", "operation": "==", "value": "PUT"},
        {"feature": "right", "operation": "feature_str_contains_value", "value": "x"}
    ],
    ...
}
```

### Where-Condition ('where')
The `where` field contains a list of conditions which are used to find the document in the `target` which should be updated by the `actions`.
Used in Action-Groups with `type` = `"child"` and in the `init` section of Action-Groups with `type` = `"root"`.
Each condition can contain the following fields:

* `target_feature`: (string) reference to a feature saved in the elasticsearch document. May contain `'.'` to traverse (for example `device.name`)
* `operation`: (string) operation that will be executed to determine the result of the condition
* `value`: (anything) value on which the operation can be executed. The type of the value is determined by the operation and target_feature.
* `event_feature`: (string) reference to a event_feature (deviced from the `features` field) on which the operation can be executed.

Currently valid operations are:

* `==`:
    * checks equality with the feature.
    * Uses the `value` field if set, if not it tries to use the `event_feature` reference.
    * If neither is set the condition will check the non-existence of the `target_feature`. For example `{"target_feature": "kind", "operation":"==", "value":null}` searches for documents where the field `kind` does not exist.
    * `{"target_feature": "name", "operation":"==", "value":"foo"}` searches for documents where the field `name` is equal to `"foo"`.
* `!=`:
    * not `==`.
    * `{"target_feature": "kind", "operation":"!=", "value":null}` searches for documents where the field `kind` does exist.
    * `{"target_feature": "name", "operation":"!=", "event_feature":"name"}` searches for documents where the field `name` is equal to the event_feature with the name `name`.
* `any_target_in_value`
    * checks if any of the values matches the `target_feature`.
    * the `target_feature` may be a list but can also be a single value.
        * if list: any target matches any value
        * if single element: any value matches target
    * the `value` must be a list of values.
* `any_target_in_event`: same as `any_target_in_value` but with `event_feature` replacing `value`

**Example:**
```
{
    ...
    "where":[
        {"target_feature": "resource", "operation":"==", "event_feature":"id"},
        {"target_feature": "kind", "operation":"==", "value":"deviceinstance"}
    ],
    ...
}
```

#### Notes
* if target_feature is a number, elastic search tries to interpret string values as numbers. if it is unable to it will throw a error
* elasticsearch allows terms on lists as if they where elements -> `==` and `!=` work on lists as if there where none (list.a == "foo";  where_test.go line 479).
* elasticsearch is not able to compare objects (only primitives). It would be possible to rewrite a object to something like list.a but that would loose the correlation between the fields of the object.
Elasticsearch has some solutions for this problem but these would make this project more complex. https://www.elastic.co/blog/managing-relations-inside-elasticsearch https://www.elastic.co/guide/en/elasticsearch/guide/current/nested-objects.html

### Id-Feature ('id_feature')
Used to find the document in the `target` which should be updated by the `actions`.
Only used in Action-Groups with `type` = `"root"`. Uses given event-feature as id of the document. If no `id_feature` is in the root-group defined, a id will be generated.

### Actions ('actions')
Actions modify the found elasticsearch document if the if-conditions are met. A action may contain teh fields `type`, `if`, `fields` and `scale`.
`if` follows the same rules as the `if` in the Action-Group. `scale` helps `type` to differentiate its behavior between lists and single elements. Depending on `type` is may not exists, but if it does is must have the value `"one"` ore `"many"`.
`fields` describes on which document (target) fields the action should be performed.
`fields` does not allow the `'.'` notation that is working in where-conditions.

#### Action: insert-one
`type` is `"insert"`, `scale` is `"one"`

Writes event-features as object to all fields listed in `fields`.

If `fields` is `[""]` the event-features will be written to the root of the document.

**Example 1:**

features:
```
"features": [
    {"name": "command", "path": "$.command+", "temp":true},
    {"name": "name", "path": "$.device_instance.name+"},
],
```

features result (event-features):
```
{
    "name": "devicename",
    "command":"PUT"
}
```

action:
```
"actions": [{
    "type": "insert",
    "if": [{"feature": "command", "operation": "==", "value": "PUT"}],
    "fields": ["device"],
    "scale": "one"
}]
```
actions result:
```
{
    "device": {
        "name": "devicename"
    }
}
```

**Example 2:**

features as in example 1

action:
```
"actions": [{
    "type": "insert",
    "if": [{"feature": "command", "operation": "==", "value": "PUT"}],
    "fields": [""],
    "scale": "one"
}]
```
actions result:
```
{
   "name": "devicename"
}
```


#### Action: remove-one
`type` is `"remove"`, `scale` is `"one"`

Removes all fields listed in `fields` from document.

If `fields` is `[""]` all fields contained in event-features will removed from the document.


**Example 1:**

features:
```
"features": [
    {"name": "command", "path": "$.command+", "temp":true},
    {"name": "name", "path": "$.device_instance.name+"},
],
```

features result (event-features):
```
{
    "name": "devicename",
    "command":"DELETE"
}
```

action:
```
"actions": [{
    "type": "remove",
    "if": [{"feature": "command", "operation": "==", "value": "DELETE"}],
    "fields": ["device"],
    "scale": "one"
}]
```
before actions (document):
```
{
    "foo": "bar",
    "device": {
        "name": "something"
    }
}
```

after actions (document):
```
{
    "foo": "bar"
}
```

**Example 2:**

features as in example 1

action:
```
"actions": [{
    "type": "remove",
    "if": [{"feature": "command", "operation": "==", "value": "DELETE"}],
    "fields": [""],
    "scale": "one"
}]
```
before actions (document):
```
{
    "foo": "bar",
    "name": "something"
}
```

after actions (document):
```
{
    "foo": "bar"
}
```

**Example 3:**

features are irrelevant for this example

action:
```
"actions": [{
    "type": "remove",
    "fields": ["foo", "device"],
    "scale": "one"
}]
```
before actions (document):
```
{
    "foo": "bar",
    "name": "something",
    "bla": "bla"
}
```

after actions (document):
```
{
    "bla": "bla"
}
```


#### Action: insert-many
`type` is `"insert"`, `scale` is `"many"`

Adds event-features as object to all fields listed in `fields`.
Fields listed in `fields` are interpreted as lists.

`fields` may not be  `[""]`.

#### Action: remove-many
`type` is `"remove"`, `scale` is `"many"`

Removes all elements from fields listed in `fields` which completely match the event-features.
Fields listed in `fields` are interpreted as lists.

`fields` may not be  `[""]`.

**Example:**

features:
```
"features": [
    {"name": "command", "path": "$.command+", "temp":true},
    {"name": "name", "path": "$.device_instance.name+"},
],
```

features result (event-features):
```
{
    "name": "devicename",
    "command":"DELETE"
}
```

action:
```
"actions": [{
    "type": "remove",
    "if": [{"feature": "command", "operation": "==", "value": "DELETE"}],
    "fields": ["device"],
    "scale": "many"
}]
```
before actions (document):
```
{
    "foo": "bar",
    "device": [
        {
            "name": "devicename"
        },
        {
            "name": "foo"
        },
        {
            "name": "devicename",
            "foo": "bar"
        }
    ]
}
```

after actions (document):
```
{
    "foo": "bar",
    "device": [
        {
            "name": "foo"
        },
        {
            "name": "devicename",
            "foo": "bar"
        }
    ]
}
```

#### Action: remove_target
`type` is `"remove_target"`

removes the whole document.


### Init ('init')
Only used in Action-Groups with type = "root". 
This section will be executed after the actions of a root-actions-group if no existing document with the id_feature was found and a new document will be created.
The init section is used to read documents that should be included into this document. If a dependency document is changed those changes will be included by actions for the corresponding event.

init is structured similar to a child action group and consists of the following fields:

* `target`: ame as action group target
* `where`: same structure and working as action group where. But event_feature references the event features from the parent root action group and target_feature references the fields of the target of this init segment.
* `sorting`: optional. used to limit and sort `where` results. structure will be explained in sub-chapter.
* `default`: optional. will be used if `where` found nothing. structure will be explained in sub-chapter.
* `transform`: same structure and implementation as action group `features`. transforms the results of `where` or `default`. allows access to deeply nested fields, declaration of defaults and setting if the temp/omitempty flags. will be used by the init actions.
* `actions`: same as action group actions. uses the result of `transform` as features

#### Init-Sorting
Is optional. additional informations for `where`. consists of the fields:
* `by`: (string) reference to elasticsearch document field. allows `field.subfield` syntax.
* `asc`: (bool) default is false; defines the sorting direction.
* `limit`: (int) defines how many elements should be searched as maximum. elasticsearch default is 10 (used if no sorting segment is used). matview default is 1000 (used if sorting segment is used without no limit).

**Example:**
```
{
    ...
    "sorting": {"by": "date", "asc": true, "limit": 20},
    ...
}
```  

#### Init-Default
Is optional. Default will provide replacement results for `where` if nothing is found. 
It is a list of maps, where each map represents one where equivalent result. Each key of a map represents a field of a result.
Each value of a map describes what value result field should have. This can be done by:
* `feature`: (string) reference to a event feature of the parent root action group.
* `value`: (anything) value that should be used

**Example:**
```
{
    ...
    "default": [
        {"user":{"feature":"owner"}, "right":{"value":"rwxa"}, "command":{"value":"PUT"}},
        {"group":{"value":"admin"}, "right":{"value":"rwxa"}, "command":{"value":"PUT"}}
    ],
    ...
}
```

# Config - Queries
The queries section describes additional selections and projections for http-requests. It has the following structure:

```
"queries": {
    "{{resource}}":{
        "{{endpoint}}": {
            "selection": {{selection}},
            "projection": {{projection}}
        }
    }
}
```
`resource` corresponds with a target from a action group in the events-config. It also will be referenced in the ElasticMapping.
`endpoint` is a identifier for this selection/projection combination. `resource` and `endpint` will be referenced together in the http api.

## Selection ('selection')
The selection segment adds additional filters to a elasticsearch query that is triggered by a received http-api-request.

There are 4 kinds of selection:

### Selection-All
Used if no additional filter should be added to the query.

**Example:**
```
{
    ...
    "selection": {"all": true},
    ...
}
```

### Selection-Or
Used to combine a list of other Selections/Conditions. At least one condition must apply.

**Example:**
```
{
    ...
    "selection": {
        "or": [
            {"condition": {"feature": "write.user", "operation": "==", "ref": "jwt.user"}},
            {"condition": {"feature": "write.group", "operation": "any_value_in_feature", "ref": "jwt.groups"}}
        ]
    },
    ...
}
```

### Selection-And
Used to combine a list of other Selections/Conditions. All conditions must apply.

**Example:**
```
{
    ...
    "selection":{
        "and": [
            {
                "or": [
                    {"condition": {"feature": "read.user", "operation": "==", "ref": "jwt.user"}},
                    {"condition": {"feature": "read.group", "operation": "any_value_in_feature", "ref": "jwt.groups"}}
                ]
            },
            {
                "or": [
                    {"condition": {"feature": "write.user", "operation": "==", "ref": "jwt.user"}},
                    {"condition": {"feature": "write.group", "operation": "any_value_in_feature", "ref": "jwt.groups"}}
                ]
            }
        ]
    },
    ...
}
```

### Selection-Condition
Adds a Filter/Condition to the elasticsearch query. A `condition` has the following fields:

* `feature`: (string) reference to a feature saved in the elasticsearch document. May contain `'.'` to traverse (for example `device.name`)
* `operation`: (string) operation that will be executed to determine the result of the condition.
* `value`: (anything) value on which the operation can be executed. The type of the value is determined by the operation and target_feature.
* `ref`: (string) uses predefined references as value.

Currently valid operations are:

* `==`:
    * checks equality with the feature.
    * Uses the `value` field if set, if not it tries to use the `ref` reference.
    * If neither is set the condition will check the non-existence of the `feature`. For example `{"feature": "kind", "operation":"==", "value":null}` searches for documents where the field `kind` does not exist.
    * `{"feature": "name", "operation":"==", "value":"foo"}` searches for documents where the field `name` is equal to `"foo"`.
* `!=`:
    * not `==`.
    * `{"feature": "kind", "operation":"!=", "value":null}` searches for documents where the field `kind` does exist.
    * `{"feature": "name", "operation":"!=", "value":"foo"}` searches for documents where the field `name` is equal to `"foo"`.
* `any_value_in_feature`:
    * interprets the `value` or `ref` as list 
    * checks if any of the list-entries matches the `feature`.
    * the `feature` may be a list but can also be a single value.
        * if list: any target matches any value
        * if single element: any value matches target

Currently valid `ref` values are:

* `"jwt.user"`: (string) uses the user-id that was transmitted by the JWT-Authorisation-Token in the HTTP-Request.
* `"jwt.groups"`: ([]string) uses the groups that where transmitted by the JWT-Authorisation-Token in the HTTP-Request.
* other ref values will be interpretet as query-parameter of the http-request (usage shown in integration_test.go fom line 608, with queries definition in integration_test.go from line 196)


## Projection ('projection')
The projection segment defines which parts of a elasticsearch document the user may receive by a http-api-request. is consists of a simple list of field-names.
A special variation is the projection `["*"]`, which means that all fields will be used for the result.

**Example:**
```
{
    "projection": ["device", "gw"]
}
```

# Confit-ElasticMapping
This section will be used for the Mapping in elasticsearch https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping.html.
The configuration for each resource will be placed under `mapping.doc.properties`.
Matview prepares a index for searches, if you want a field to be searchable in the http-api use `"copy_to": "feature_search"`.
Types other then "Keyword" may influence results of `where` and `queries` by running elasticsearch analysis on this field (for example stemming).

**Example:**

```
"elastic_mapping": {
    "simple_resource": {
      "name":    {"type": "keyword"},
      "devices": {"type": "keyword"},
      "id":      {"type": "keyword"},
      "command": {"type": "keyword"}
    },
    "complex_resource":{
      "device":{
        "properties": {
          "name":         {"type": "keyword", "copy_to": "feature_search"},
          "description":  {"type": "text",    "copy_to": "feature_search"},
          "usertag":      {"type": "keyword", "copy_to": "feature_search"},
          "tag":          {"type": "keyword", "copy_to": "feature_search"},
          "devicetype":   {"type": "keyword"},
          "uri":          {"type": "keyword"},
          "img":          {"type": "keyword"}
        }
      },
      "gw":{
        "properties": {
          "name":         {"type": "keyword", "copy_to": "feature_search"}
        }
      }
    }
  }
```

# HTTP

## Routes

* `GET /search/:target/:searchtext/:endpoint`:
    * searches for partial matches of the `searchtext` in the index `target` on features with `"copy_to": "feature_search"`.
    * `endpoint` is a reference to a queries section which defines additional selection and projection criteria.
    * returns maximal 10 results (elasticsearch default, can be changed with postfix-route).
* `GET /get/:target/:endpoint`
    * uses only `endpoint` queries-configuration to find results.
    * returns maximal 10 results (elasticsearch default, can be changed with postfix-route).
* `GET /select/field/:target/:endpoint/:field/:value`
    * finds `target` documents where the `field` has the `value`.
    * field should have `"type": "keyword"` in `elastic_mapping` to prevent elasticsearch analysis to influence result.
    * `endpoint` is a reference to a queries section which defines additional selection and projection criteria.
    * returns maximal 10 results (elasticsearch default, can be changed with postfix-route).
* `POST /select/field/:target/:endpoint/:field`
    * gets list of values from post body.
    * finds `target` documents where the `field` has the any of the given values.
    * field should have `"type": "keyword"` in `elastic_mapping` to prevent elasticsearch analysis to influence result.
    * `endpoint` is a reference to a queries section which defines additional selection and projection criteria.
    * returns maximal 10 results (elasticsearch default, can be changed with postfix-route).

## Postfix-Routes

These routes can be appended on all routes to define sorting and paging.

* `/:limit/:offset` 
    * returns maximal `limit` results.
    * skips `offset` documents.
* `/:limit/:offset/:order_by/asc` 
    * returns maximal `limit` results.
    * skips `offset` documents.
    * orders by field `order_by` ascending.
    * `order_by` may have `field.subfield` syntax.
    * `order_by` must be descibed in ElasticMapping.
* `/:limit/:offset/:order_by/desc` 
    * returns maximal `limit` results
    * skips `offset` documents
    * orders by field `order_by` descending.
    * `order_by` may have `field.subfield` syntax.
    * `order_by` must be descibed in ElasticMapping.