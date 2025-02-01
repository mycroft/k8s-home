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
class = S3Storage
[model_options]
region = ""
version = "latest"
endpoint = "https://minio-storage.services.mkz.me"
use_path_style_endpoint = true
bucket = "privatebin"
accesskey = "AWS_ACCESS_KEY_ID"
secretkey = "AWS_SECRET_ACCESS_KEY"

[yourls]
