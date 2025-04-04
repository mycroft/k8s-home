;<?php http_response_code(403); /*

[main]
name = "PrivateBin"
; basepath = "https://privatebin.example.com/"
discussion = true
opendiscussion = false
password = true
fileupload = false
burnafterreadingselected = false
defaultformatter = "plaintext"
sizelimit = 10485760
template = "bootstrap-dark"
languageselection = false

[expire]
default = "1week"

[expire_options]
5min = 300
10min = 600
1hour = 3600
1day = 86400
1week = 604800
1month = 2592000
1year = 31536000
never = 0

[formatter_options]
plaintext = "Plain Text"
syntaxhighlighting = "Source Code"
markdown = "Markdown"

[traffic]
limit = 10

[purge]
limit = 300
batchsize = 10

[model]
class = Database
[model_options]
dsn = "pgsql:host=postgres-instance.postgres;dbname=privatebin"
tbl = "privatebin_"
usr = "POSTGRESQL_USERNAME"
pwd = "POSTGRESQL_PASSWORD"
opt[12] = true

; class = S3Storage
; [model_options]
; region = "fr-north-1"
; version = "latest"
; endpoint = "https://minio-storage.services.mkz.me"
; use_path_style_endpoint = true
; bucket = "privatebin"
; accesskey = "AWS_ACCESS_KEY_ID"
;secretkey = "AWS_SECRET_ACCESS_KEY"

[yourls]
