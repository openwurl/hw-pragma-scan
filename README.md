# HW Pragma Scan
Simple tool for fetching Pragmas and other useful headers from a Highwinds site

# Usage
```
go run main.go scan -u http://website.com/path/to/file
```

or

```
hw-pragma scan -u http://website.com/path/to/file
```

# Example
```
+-----------+---------------+---------------+-------------+
|   FIELD   | TTL (SECONDS) | TTL (MINUTES) | TTL (HOURS) |
+-----------+---------------+---------------+-------------+
| CDN Cache |         86382 |          1439 |          23 |
+-----------+---------------+---------------+-------------+
+-----------------------------+--------------------+
|            FIELD            |       VALUE        |
+-----------------------------+--------------------+
| File Size                   |             851264 |
+-----------------------------+--------------------+
| Access-Control-Allow-Origin | *                  |
+-----------------------------+--------------------+
| Cache-Control               | max-age=3600       |
+-----------------------------+--------------------+
| Content-Type                | video/MP2T         |
+-----------------------------+--------------------+
| X-HW-Cache-Compressed-Size  | [NA]               |
+-----------------------------+--------------------+
| X-HW-Cache-Behavior         | DEFAULT            |
+-----------------------------+--------------------+
| X-HW-Cache-Last-Modified    | 437925h6m7.088353s |
+-----------------------------+--------------------+
| X-HW-Cache-Originated       | 18.088353s         |
+-----------------------------+--------------------+
| X-HW-Cache-Last-Refresh     | 18.088353s         |
+-----------------------------+--------------------+
| X-HW-Cache-Last-Request     | 88.353ms           |
+-----------------------------+--------------------+
```