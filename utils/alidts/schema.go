package alidts

const ALIYUN_DTS_SCHEMA = `
{
  "name": "com.alibaba.alidts.formats.avro.Record",
  "type": "record",
  "fields": [
    {
      "doc": "version information",
      "name": "version",
      "type": "int"
    },
    {
      "doc": "unique id of this record in the whole stream",
      "name": "id",
      "type": "long"
    },
    {
      "doc": "record log timestamp",
      "name": "sourceTimestamp",
      "type": "long"
    },
    {
      "doc": "record source location information",
      "name": "sourcePosition",
      "type": "string"
    },
    {
      "default": "",
      "doc": "safe record source location information, use to recovery.",
      "name": "safeSourcePosition",
      "type": "string"
    },
    {
      "default": "",
      "doc": "record transation id",
      "name": "sourceTxid",
      "type": "string"
    },
    {
      "doc": "source datasource",
      "name": "source",
      "type": {
        "fields": [
          {
            "name": "sourceType",
            "type": {
              "name": "SourceType",
              "namespace": "com.alibaba.alidts.formats.avro",
              "symbols": [
                "MySQL",
                "Oracle",
                "SQLServer",
                "PostgreSQL",
                "MongoDB",
                "Redis",
                "DB2",
                "PPAS",
                "DRDS",
                "HBASE",
                "HDFS",
                "FILE",
                "OTHER"
              ],
              "type": "enum"
            }
          },
          {
            "doc": "source datasource version information",
            "name": "version",
            "type": "string"
          }
        ],
        "name": "Source",
        "namespace": "com.alibaba.alidts.formats.avro",
        "type": "record"
      }
    },
    {
      "name": "operation",
      "namespace": "com.alibaba.alidts.formats.avro",
      "type": {
        "name": "Operation",
        "symbols": [
          "INSERT",
          "UPDATE",
          "DELETE",
          "DDL",
          "BEGIN",
          "COMMIT",
          "ROLLBACK",
          "ABORT",
          "HEARTBEAT",
          "CHECKPOINT",
          "COMMAND",
          "FILL",
          "FINISH",
          "CONTROL",
          "RDB",
          "NOOP",
          "INIT"
        ],
        "type": "enum"
      }
    },
    {
      "default": null,
      "name": "objectName",
      "type": [
        "null",
        "string"
      ]
    },
    {
      "default": null,
      "doc": "time when this record is processed along the stream dataflow",
      "name": "processTimestamps",
      "type": [
        "null",
        {
          "items": "long",
          "type": "array"
        }
      ]
    },
    {
      "default": {

      },
      "doc": "tags to identify properties of this record",
      "name": "tags",
      "type": {
        "type": "map",
        "values": "string"
      }
    },
    {
      "default": null,
      "name": "fields",
      "type": [
        "null",
        "string",
        {
          "items": {
            "fields": [
              {
                "name": "name",
                "type": "string"
              },
              {
                "name": "dataTypeNumber",
                "type": "int"
              }
            ],
            "name": "Field",
            "namespace": "com.alibaba.alidts.formats.avro",
            "type": "record"
          },
          "type": "array"
        }
      ]
    },
    {
      "default": null,
      "name": "beforeImages",
      "type": [
        "null",
        "string",
        {
          "items": [
            "null",
            {
              "fields": [
                {
                  "name": "precision",
                  "type": "int"
                },
                {
                  "name": "value",
                  "type": "string"
                }
              ],
              "name": "Integer",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "fields": [
                {
                  "name": "charset",
                  "type": "string"
                },
                {
                  "name": "value",
                  "type": "bytes"
                }
              ],
              "name": "Character",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "fields": [
                {
                  "name": "value",
                  "type": "string"
                },
                {
                  "name": "precision",
                  "type": "int"
                },
                {
                  "name": "scale",
                  "type": "int"
                }
              ],
              "name": "Decimal",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "fields": [
                {
                  "name": "value",
                  "type": "double"
                },
                {
                  "name": "precision",
                  "type": "int"
                },
                {
                  "name": "scale",
                  "type": "int"
                }
              ],
              "name": "Float",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "fields": [
                {
                  "name": "timestamp",
                  "type": "long"
                },
                {
                  "name": "millis",
                  "type": "int"
                }
              ],
              "name": "Timestamp",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "fields": [
                {
                  "default": null,
                  "name": "year",
                  "type": [
                    "null",
                    "int"
                  ]
                },
                {
                  "default": null,
                  "name": "month",
                  "type": [
                    "null",
                    "int"
                  ]
                },
                {
                  "default": null,
                  "name": "day",
                  "type": [
                    "null",
                    "int"
                  ]
                },
                {
                  "default": null,
                  "name": "hour",
                  "type": [
                    "null",
                    "int"
                  ]
                },
                {
                  "default": null,
                  "name": "minute",
                  "type": [
                    "null",
                    "int"
                  ]
                },
                {
                  "default": null,
                  "name": "second",
                  "type": [
                    "null",
                    "int"
                  ]
                },
                {
                  "default": null,
                  "name": "millis",
                  "type": [
                    "null",
                    "int"
                  ]
                }
              ],
              "name": "DateTime",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "fields": [
                {
                  "name": "value",
                  "type": "com.alibaba.alidts.formats.avro.DateTime"
                },
                {
                  "name": "timezone",
                  "type": "string"
                }
              ],
              "name": "TimestampWithTimeZone",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "fields": [
                {
                  "name": "type",
                  "type": "string"
                },
                {
                  "name": "value",
                  "type": "bytes"
                }
              ],
              "name": "BinaryGeometry",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "fields": [
                {
                  "name": "type",
                  "type": "string"
                },
                {
                  "name": "value",
                  "type": "string"
                }
              ],
              "name": "TextGeometry",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "fields": [
                {
                  "name": "type",
                  "type": "string"
                },
                {
                  "name": "value",
                  "type": "bytes"
                }
              ],
              "name": "BinaryObject",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "fields": [
                {
                  "name": "type",
                  "type": "string"
                },
                {
                  "name": "value",
                  "type": "string"
                }
              ],
              "name": "TextObject",
              "namespace": "com.alibaba.alidts.formats.avro",
              "type": "record"
            },
            {
              "name": "EmptyObject",
              "namespace": "com.alibaba.alidts.formats.avro",
              "symbols": [
                "NULL",
                "NONE"
              ],
              "type": "enum"
            }
          ],
          "type": "array"
        }
      ]
    },
    {
      "default": null,
      "name": "afterImages",
      "type": [
        "null",
        "string",
        {
          "items": [
            "null",
            "com.alibaba.alidts.formats.avro.Integer",
            "com.alibaba.alidts.formats.avro.Character",
            "com.alibaba.alidts.formats.avro.Decimal",
            "com.alibaba.alidts.formats.avro.Float",
            "com.alibaba.alidts.formats.avro.Timestamp",
            "com.alibaba.alidts.formats.avro.DateTime",
            "com.alibaba.alidts.formats.avro.TimestampWithTimeZone",
            "com.alibaba.alidts.formats.avro.BinaryGeometry",
            "com.alibaba.alidts.formats.avro.TextGeometry",
            "com.alibaba.alidts.formats.avro.BinaryObject",
            "com.alibaba.alidts.formats.avro.TextObject",
            "com.alibaba.alidts.formats.avro.EmptyObject"
          ],
          "type": "array"
        }
      ]
    }
  ]
}`
